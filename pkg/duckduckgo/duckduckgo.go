//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package duckduckgo

import "github.com/fogfish/gurl/v2/http"

type DuckDuckGo struct {
	http.Stack
	host string
}

type Config struct {
	http.Stack
}

func New(cfg Config) (*DuckDuckGo, error) {
	stack := &DuckDuckGo{
		Stack: cfg.Stack,
		host:  "https://api.duckduckgo.com/",
	}

	if stack.Stack == nil {
		stack.Stack = http.New()
	}

	return stack, nil
}
