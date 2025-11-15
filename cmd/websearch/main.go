//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package main

import (
	"fmt"

	"github.com/fogfish/websearch/cmd/websearch/opts"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	opts.Execute(fmt.Sprintf("websearch/%s (%s), %s", version, commit[:7], date))
}
