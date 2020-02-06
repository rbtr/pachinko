/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package util

import "sort"

// StringSliceContains sorts and performs a binary search on a slice s to
// determine if x exists in the slice
func StringSliceContains(s []string, x string) bool {
	i := sort.SearchStrings(s, x)
	return i < len(s) && s[i] == x
}
