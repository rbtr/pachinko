/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package intra

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/pkg/errors"
	api "github.com/rbtr/go-tvdb"
	"github.com/rbtr/go-tvdb/generated/models"
	"github.com/rbtr/pachinko/plugin/processor"
	"github.com/rbtr/pachinko/types"
	"github.com/rbtr/pachinko/types/metadata/tv"
	log "github.com/sirupsen/logrus"
)

// WordMatcher regex
var matcher *regexp.Regexp = regexp.MustCompile(`[^\w]`)

// TVDbClient adds metadata from the TVDb
type TVDbClient struct {
	APIKey       string `mapstructure:"api-key"`
	RequestLimit int64  `mapstructure:"request-limit"`

	client  *api.Client
	limiter *time.Ticker
}

func (c *TVDbClient) Init(context.Context) error {
	authn := &models.Auth{
		Apikey: c.APIKey,
	}
	c.client = api.DefaultClient(authn)
	c.limiter = time.NewTicker((time.Second / time.Duration(c.RequestLimit)))
	return nil
}

func (c *TVDbClient) identify(m types.Media) (*models.Episode, *models.SeriesSearchResult, error) {
	cleanName := matcher.ReplaceAllLiteralString(m.TVMetadata.Name, " ")
	log.Debugf("tvdb_decorator: identifying %s", cleanName)

	// note: when TV has a (YEAR), it is because there's multiple series with the same
	// name (i.e. it's been rebooted) and the year is part of the name that is used to
	// disambiguate (at least that's how thetvdb does it)
	param := map[string]string{"name": cleanName}
	if m.TVMetadata.ReleaseYear > 0 {
		log.Tracef("tvdb_decorator: show has year %d", m.TVMetadata.ReleaseYear)
		param["name"] = fmt.Sprintf("%s (%d)", param["name"], m.TVMetadata.ReleaseYear)
	}

	res, err := c.client.SearchSeries(context.TODO(), param)
	if err != nil {
		return nil, nil, err
	}

	resMap := map[string]*models.SeriesSearchResult{}
	log.Tracef("tvdb_decorator: search series found: %#v", resMap)
	resKeys := []string{}
	for _, res := range res {
		resMap[res.SeriesName] = res
		resKeys = append(resKeys, res.SeriesName)
	}
	matches := fuzzy.RankFindFold(cleanName, resKeys)
	sort.Sort(matches)
	if len(matches) == 0 {
		return nil, nil, errors.Errorf("tvdb_decorator: no matches for %s", param["name"])
	}
	name := matches[0].Target
	series := resMap[name]
	log.Debugf("tvdb_decorator: search for %s found %s", param["name"], name)

	eps, _, jsonErr, err := c.client.GetSeriesEpisode(context.TODO(), series.ID, 0, map[string]string{"airedSeason": strconv.Itoa(m.TVMetadata.Season.Number), "airedEpisode": strconv.Itoa(m.TVMetadata.Episode.Number)})
	if err != nil {
		return nil, nil, err
	}
	if jsonErr != nil {
		return nil, nil, err
	}

	if len(eps) == 0 {
		return nil, nil, errors.New("no matching episode found")
	}
	return eps[0], series, nil
}

func (c *TVDbClient) addTVDBMetadata(m types.Media) types.Media {
	ep, series, err := c.identify(m)
	if err != nil || ep == nil || series == nil {
		log.Errorf("tvdb_decorator: error identifying episode: %s", err)
		return m
	}
	log.Debugf("tvdb_decorator: got episode from tvdb: %v", ep)
	m.Identifiers["tvdb"] = strconv.FormatInt(ep.ID, 10)
	m.TVMetadata.Name = series.SeriesName
	m.TVMetadata.AbsoluteNumber = int(ep.AbsoluteNumber)
	m.TVMetadata.Episode.Title = ep.EpisodeName
	log.Tracef("tvdb_decorator: populated %v from tvdb", m)
	return m
}

func (c *TVDbClient) Process(in <-chan types.Media, out chan<- types.Media) {
	log.Trace("started tvdb_decorator processor")
	for m := range in {
		log.Tracef("tvdb_decorator: received input: %#v", m)
		if m.Type == tv.TV {
			log.Infof("tvdb_decorator: looking up %s in tvdb", m.SourcePath)
			<-c.limiter.C // rate limiting on tvdb api calls
			m = c.addTVDBMetadata(m)
		} else {
			log.Debugf("tvdb_decorator: %s type [%s] != TV, skipping", m.SourcePath, m.Type)
		}
		out <- m
	}
}

func init() {
	processor.Register(processor.Intra, "tvdb", func() processor.Processor {
		return &TVDbClient{
			RequestLimit: 10,
		}
	})
}
