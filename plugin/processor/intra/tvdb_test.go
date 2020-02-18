/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package intra

// package tvdb

// import (
// 	"context"
// 	"sort"
// 	"testing"

// 	"github.com/lithammer/fuzzysearch/fuzzy"
// )

// func TestClient_FetchMetadata(t *testing.T) {
// 	c, _ := NewTVDBDecorator()
// 	name := matcher.ReplaceAllLiteralString("Mr-Robot -", "")
// 	res, err := c.SearchSeries(context.TODO(), map[string]string{"name": name})
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Logf("res %+v\n", res)

// 	names := []string{}
// 	for _, res := range res {
// 		names = append(names, res.SeriesName)
// 	}
// 	t.Logf("names %v\n", names)
// 	matches := fuzzy.RankFindFold(name, names)
// 	t.Logf("matches %v\n", matches)
// 	sort.Sort(matches)
// 	t.Logf("sorted %v", matches)
// 	if matches[0].Target != "Mr. Robot" {
// 		t.Errorf("wanted 'Mr. Robot', got %s", matches[0].Target)
// 	}
// }

// func TestFuzzy(t *testing.T) {
// 	name := matcher.ReplaceAllLiteralString("Mr-Robot", "")
// 	matches := fuzzy.RankFindNormalizedFold(name, []string{"Mr. Robot"})
// 	t.Logf("matches %v\n", matches)
// }
