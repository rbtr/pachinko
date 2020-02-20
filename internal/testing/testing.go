package testing

import (
	"github.com/rbtr/pachinko/types"
	"github.com/rbtr/pachinko/types/metadata/movie"
	"github.com/rbtr/pachinko/types/metadata/tv"
)

type Test struct {
	Name   string
	Inputs []string
	Want   types.Media
}

var TV []*Test = []*Test{
	{
		"SssEee",
		[]string{
			"Mr-Robot-2x5",
			"Mr.Robot.2x5",
			"Mr Robot 2x5",
			"Mr-Robot-2-5",
			"Mr.Robot.2-5",
			"Mr Robot 2-5",
			"Mr Robot S02E05",
			"Mr-Robot/Season-2/5",
			"Mr.Robot/Season.2/5",
			"Mr Robot/Season 2/5",
			"Mr-Robot-Season-02-Episode-05",
			"Mr.Robot.Season.02.Episode.05",
			"Mr Robot Season 02 Episode 05",
			"Mr Robot/Season 02/Episode 05",
		},
		types.Media{
			Category: types.Video,
			Type:     tv.TV,
			TVMetadata: tv.Metadata{
				Name: "Mr Robot",
				Episode: tv.Episode{
					Number: 5,
					Season: tv.Season{
						Number: 2,
					},
				},
			},
		},
	},
	{
		"dir multimathc SSxEE",
		[]string{
			"/src/Broad City (2014) S01-03 Season 01-03 (1080p HEVC AAC 2.0)/Broad City 03x09 Getting There.mkv",
		},
		types.Media{
			Category: types.Video,
			Type:     tv.TV,
			TVMetadata: tv.Metadata{
				Name: "Broad City",
				Episode: tv.Episode{
					Number: 9,
					Season: tv.Season{
						Number: 3,
					},
				},
			},
		},
	},
	{
		"dotted year",
		[]string{
			"/src/Doctor Who 2005 S12E03 1080p HEVC x265/Doctor.Who.2005.S12E03.1080p.HEVC.x265.mkv",
		},
		types.Media{
			Category: types.Video,
			Type:     tv.TV,
			TVMetadata: tv.Metadata{
				Name:        "Doctor Who",
				ReleaseYear: 2005,
				Episode: tv.Episode{
					Number: 3,
					Season: tv.Season{
						Number: 12,
					},
				},
			},
		},
	},
}

var Movies []*Test = []*Test{
	{
		"year as title",
		[]string{
			"/src/Blade Runner 2049 (2017)/Blade Runner 2049.mkv",
		},
		types.Media{
			Category: types.Video,
			Type:     movie.Movie,
			MovieMetadata: movie.Metadata{
				Title:       "Blade Runner 2049",
				ReleaseYear: 2017,
			},
		},
	},
	{
		"typical",
		[]string{
			"/src/Finding Nemo (2003).mkv",
		},
		types.Media{
			Category: types.Video,
			Type:     movie.Movie,
			MovieMetadata: movie.Metadata{
				Title:       "Finding Nemo",
				ReleaseYear: 2003,
			},
		},
	},
}

var NotTV []*Test = []*Test{}
var NotMovies []*Test = []*Test{}

func init() {
	NotTV = append(NotTV, Movies...)
	NotMovies = append(NotMovies, TV...)
}
