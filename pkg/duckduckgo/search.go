//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package duckduckgo

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
	"github.com/fogfish/websearch"
)

type Search struct {
	Query string // The search query string.
}

type bag struct {
	Heading        string           `json:"Heading"`
	Abstract       string           `json:"Abstract"`
	AbstractSource string           `json:"AbstractSource"`
	AbstractText   string           `json:"AbstractText"`
	AbstractURL    string           `json:"AbstractURL"`
	RelatedTopics  []map[string]any `json:"RelatedTopics"`
}

func (api *DuckDuckGo) Search(ctx context.Context, req Search) ([]websearch.Fact, error) {
	var buf bytes.Buffer
	err := api.Stack.IO(ctx,
		http.GET(
			ø.URI(api.host),
			ø.Param("q", req.Query),
			ø.Param("format", "json"),

			ƒ.Status.Accepted,
			ƒ.Bytes(&buf),
		),
	)
	if err != nil {
		return nil, err
	}

	var bag bag
	if err := json.Unmarshal(buf.Bytes(), &bag); err != nil {
		return nil, err
	}

	facts := make([]websearch.Fact, 0)

	if len(bag.AbstractURL) > 0 {
		fact := websearch.Fact{
			Title:   bag.Heading,
			Snippet: bag.AbstractText,
			Url:     bag.AbstractURL,
		}
		facts = append(facts, fact)
	}

	for _, item := range bag.RelatedTopics {
		subitems, ok := item["Topics"].([]any)
		if subitems == nil || !ok {
			fact := websearch.Fact{}
			if snippet, ok := item["Text"].(string); ok {
				fact.Snippet = snippet
			}
			if url, ok := item["FirstURL"].(string); ok {
				fact.Url = url
			}
			if len(fact.Snippet) > 0 {
				facts = append(facts, fact)
			}
			continue
		}

		for _, subitem := range subitems {
			subsubitem, ok := subitem.(map[string]any)
			if !ok {
				continue
			}
			fact := websearch.Fact{}
			if snippet, ok := subsubitem["Text"].(string); ok {
				fact.Snippet = snippet
			}
			if url, ok := subsubitem["FirstURL"].(string); ok {
				fact.Url = url
			}
			if cat, ok := item["Name"].(string); ok {
				fact.Category = cat
			}
			if len(fact.Snippet) > 0 {
				facts = append(facts, fact)
			}
		}
	}

	return facts, nil
}
