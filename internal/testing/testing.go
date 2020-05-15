package testing

import (
	"github.com/rbtr/pachinko/types"
	"github.com/rbtr/pachinko/types/metadata/movie"
	"github.com/rbtr/pachinko/types/metadata/tv"
)

type Test struct {
	Name   string
	Inputs []string
	Want   types.Item
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
		types.Item{
			Category:  types.Video,
			MediaType: tv.TV,
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
		"dir multimatch SSxEE",
		[]string{
			"/src/Broad City (2014) S01-03 Season 01-03 (1080p HEVC AAC 2.0)/Broad City 03x09 Getting There.mkv",
		},
		types.Item{
			Category:  types.Video,
			MediaType: tv.TV,
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
		types.Item{
			Category:  types.Video,
			MediaType: tv.TV,
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
	{
		"special characters in title",
		[]string{
			"/src/Tom Clancy's Jack Ryan S01E01.mkv",
		},
		types.Item{
			Category:  types.Video,
			MediaType: tv.TV,
			TVMetadata: tv.Metadata{
				Name: "Tom Clancy's Jack Ryan",
				Episode: tv.Episode{
					Number: 1,
					Season: tv.Season{
						Number: 1,
					},
				},
			},
		},
	},
	{
		"special characters in title",
		[]string{
			"/src/Handmaid's Tale S03E01.mkv",
		},
		types.Item{
			Category:  types.Video,
			MediaType: tv.TV,
			TVMetadata: tv.Metadata{
				Name: "Handmaid's Tale",
				Episode: tv.Episode{
					Number: 1,
					Season: tv.Season{
						Number: 3,
					},
				},
			},
		},
	},
	{
		"special characters in title",
		[]string{
			"src/Marvel's Runaways S01E01.mkv",
		},
		types.Item{
			Category:  types.Video,
			MediaType: tv.TV,
			TVMetadata: tv.Metadata{
				Name: "Marvel's Runaways",
				Episode: tv.Episode{
					Number: 1,
					Season: tv.Season{
						Number: 1,
					},
				},
			},
		},
	},
}

var Movies []*Test = []*Test{
	{
		"year in title",
		[]string{
			"/src/Blade Runner 2049 (2017)/Blade Runner 2049.mkv",
			"/src/Blade Runner 2049 (2017).mkv",
		},
		types.Item{
			Category:  types.Video,
			MediaType: movie.Movie,
			MovieMetadata: movie.Metadata{
				Title:       "Blade Runner 2049",
				ReleaseYear: 2017,
			},
		},
	},
	{
		"year as title",
		[]string{
			"/movies/1917 (2020).mkv",
		},
		types.Item{
			Category:  types.Video,
			MediaType: movie.Movie,
			MovieMetadata: movie.Metadata{
				Title:       "1917",
				ReleaseYear: 2020,
			},
		},
	},
	{
		"typical",
		[]string{
			"/src/Finding Nemo (2003).mkv",
		},
		types.Item{
			Category:  types.Video,
			MediaType: movie.Movie,
			MovieMetadata: movie.Metadata{
				Title:       "Finding Nemo",
				ReleaseYear: 2003,
			},
		},
	},
	{
		"metadata",
		[]string{
			"/src/Frozen 2 (2019) [1080p x265 10bit FS93].mkv",
		},
		types.Item{
			Category:  types.Video,
			MediaType: movie.Movie,
			MovieMetadata: movie.Metadata{
				Title:       "Frozen 2",
				ReleaseYear: 2019,
			},
		},
	},
	{
		"subtitled",
		[]string{
			"TRON - Legacy (2010) (1080p BluRay).mkv",
		},
		types.Item{
			Category:  types.Video,
			MediaType: movie.Movie,
			MovieMetadata: movie.Metadata{
				Title:       "TRON Legacy",
				ReleaseYear: 2010,
			},
		},
	},
	{
		"special character :",
		[]string{
			"3:10 To Yuma (2007).mkv",
		},
		types.Item{
			Category:  types.Video,
			MediaType: movie.Movie,
			MovieMetadata: movie.Metadata{
				Title:       "3:10 To Yuma",
				ReleaseYear: 2007,
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
