/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package util

import "testing"

func TestStringSliceContains(t *testing.T) {
	testSlice := []string{
		"a",
		"d",
		"b",
	}
	tests := []struct {
		name string
		x    string
		want bool
	}{
		{
			name: "contains",
			x:    "d",
			want: true,
		},
		{
			name: "should insert",
			x:    "c",
			want: false,
		},
		{
			name: "should append",
			x:    "z",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSliceContains(testSlice, tt.x); got != tt.want {
				t.Errorf("StringSliceContains() = %v, want %v", got, tt.want)
			}
		})
	}
}
