//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package hackernews

import (
	"context"
	"fmt"
	"strings"

	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
	"github.com/fogfish/websearch"
)

type Search struct {
	Query string
	Size  int
	Tags  []string
}

type bag struct {
	Hits []hit `json:"hits,omitempty"`
}

type hit struct {
	Title string `json:"title,omitempty"`
	Url   string `json:"url,omitempty"`
}

func (api *HackerNews) Search(ctx context.Context, req Search) ([]websearch.Fact, error) {
	if req.Size == 0 {
		req.Size = 20
	}
	if len(req.Tags) == 0 {
		req.Tags = api.tags
	}

	bag, err := http.IO[bag](api.WithContext(ctx),
		http.GET(
			ø.URI(api.host),
			ø.Param("query", req.Query),
			ø.Param("tags", fmt.Sprintf("(%s)", strings.Join(req.Tags, ","))),
			ø.Param("hitsPerPage", req.Size),

			ƒ.Status.OK,
		),
	)
	if err != nil {
		return nil, err
	}

	facts := make([]websearch.Fact, 0, len(bag.Hits))
	for _, hit := range bag.Hits {
		fact := websearch.Fact{
			Title: hit.Title,
			Url:   hit.Url,
		}
		facts = append(facts, fact)
	}

	return facts, nil
}
