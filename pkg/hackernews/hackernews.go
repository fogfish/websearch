//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package hackernews

import (
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/fogfish/gurl/v2/http"
)

type HackerNews struct {
	http.Stack
	host    string
	tags    []string
	html2md *converter.Converter
}

type Config struct {
	http.Stack
	Tags []string
}

func New(cfg Config) (*HackerNews, error) {
	stack := &HackerNews{
		Stack: cfg.Stack,
		host:  "https://hn.algolia.com/api/v1",
		tags:  cfg.Tags,
	}

	if stack.Stack == nil {
		stack.Stack = http.New()
	}

	if len(stack.tags) == 0 {
		stack.tags = []string{"story"}
	}

	stack.html2md = converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
		),
	)
	stack.html2md.Register.TagType("img", converter.TagTypeRemove, converter.PriorityStandard)

	return stack, nil
}
