/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package path

// import (
// 	"fmt"
// 	"path/filepath"
// 	"testing"

// 	"github.com/rbtr/pachinko/internal/config"
// )

// func Test_walkDir(t *testing.T) {
// 	f, _ := filepath.Abs("testdata")
// 	input, _ := NewPathInput(config.Config{
// 		SrcDir: f,
// 	})
// 	files := []string{}
// 	for out := range input.Consume() {
// 		files = append(files, out.SourcePath)
// 	}
// 	if len(files) != 5 {
// 		t.Errorf("expected %d files, got %d", 5, len(files))
// 	}
// 	for _, t := range files {
// 		fmt.Println(t)
// 	}
// }
