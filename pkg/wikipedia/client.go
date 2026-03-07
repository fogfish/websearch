//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package wikipedia

import (
	"fmt"
	"strconv"

	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
)

// Wikipedia HTTP client
type Wikipedia struct {
	http.Stack
	host   string
	maxlag int
	arrow  http.Arrow
}

type Config struct {
	http.Stack
}

// Creates new instance of HTTP client
func New(cfg Config) (*Wikipedia, error) {
	stack := &Wikipedia{
		Stack:  cfg.Stack,
		host:   "https://%s.wikipedia.org/w/api.php",
		maxlag: 5,
	}

	if stack.Stack == nil {
		stack.Stack = http.New()
	}

	stack.arrow = http.Join(
		ø.Params(map[string]string{
			"format":        "json",
			"formatversion": "2",
			"maxlag":        strconv.Itoa(stack.maxlag),
		}),
		ø.Accept.ApplicationJSON,
		ø.UserAgent.Set("gurl/v2 (https://github.com/fogfish/gurl)"),
		// TODO: ø.AcceptEncoding.Set("gzip"),
	)

	return stack, nil
}

// defines the host of wikipedia
func (api *Wikipedia) url(lang string) string {
	if len(lang) == 0 {
		lang = EN
	}

	return fmt.Sprintf(api.host, lang)
}

func (api *Wikipedia) query(lang string, req any) http.Arrow {
	return http.GET(
		ø.URI(api.url(lang)),
		api.arrow,
		ø.Params(marshalQuery(req)),

		ƒ.Status.OK,
	)
}
