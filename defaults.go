// Copyright (c) 2023, The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grease

import "goki.dev/laser"

// SetFromDefaults sets the values of the given config object
// from `def:` field tag values. Parsing errors are automatically logged.
func SetFromDefaults(cfg any) error {
	return laser.SetFromDefaultTags(cfg)
}
