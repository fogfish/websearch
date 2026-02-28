//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package arxiv

import (
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/fogfish/gurl/v2/http"
)

type ArXiv struct {
	http.Stack
	html2md *converter.Converter
}

type Config struct {
	http.Stack
}

// Creates new instance of HTTP client
func New(cfg Config) (*ArXiv, error) {
	stack := &ArXiv{
		Stack: cfg.Stack,
	}

	if stack.Stack == nil {
		stack.Stack = http.New()
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
