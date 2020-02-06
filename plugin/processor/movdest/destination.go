/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package movdest

import (
	"fmt"
	"path"

	"github.com/rbtr/pachinko/internal/types"
)

type DestinationSolver struct {
	DestDir      string
	TVPrefix     string
	SeasonDirs   bool
	MoviesPrefix string
	MovieDirs    bool
	OutputFormat string
}

func (d *DestinationSolver) Process(m types.Media) {
	if d.MovieDirs {
		m.DestinationPath = path.Join(
			d.DestDir,
			d.MoviesPrefix,
			fmt.Sprintf("%s (%d)", m.MovieMetadata.Title, m.MovieMetadata.ReleaseYear),
			fmt.Sprintf("%s (%d)%s", m.MovieMetadata.Title, m.MovieMetadata.ReleaseYear, path.Ext(m.SourcePath)),
		)
	} else {
		m.DestinationPath = path.Join(
			d.DestDir,
			d.MoviesPrefix,
			fmt.Sprintf("%s (%d)%s", m.MovieMetadata.Title, m.MovieMetadata.ReleaseYear, path.Ext(m.SourcePath)),
		)
	}
}
