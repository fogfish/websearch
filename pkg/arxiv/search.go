//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package arxiv

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
	"github.com/fogfish/websearch"
	"github.com/mmcdole/gofeed"
)

type Search struct {
	Query   string // The search query string.
	Size    int
	SortBy  string
	OrderBy string
}

func (api *ArXiv) Search(ctx context.Context, req Search) ([]websearch.Fact, error) {
	if req.Size == 0 {
		req.Size = 20
	}

	if req.SortBy == "" {
		// allowed: relevance, lastUpdatedDate, submittedDate
		req.SortBy = "submittedDate"
	}

	if req.OrderBy == "" {
		// allowed: "ascending" or "descending"
		req.OrderBy = "descending"
	}

	var buf bytes.Buffer
	err := api.Stack.IO(ctx,
		http.GET(
			ø.URI("https://export.arxiv.org/api/query"),
			ø.Param("search_query", req.Query),
			ø.Param("max_results", fmt.Sprintf("%d", req.Size)),
			ø.Param("sortBy", req.SortBy),
			ø.Param("sortOrder", req.OrderBy),

			ƒ.Status.OK,
			ƒ.Bytes(&buf),
		),
	)
	if err != nil {
		return nil, err
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseString(buf.String())
	if err != nil {
		return nil, err
	}

	var facts []websearch.Fact
	for _, item := range feed.Items {
		published, _ := time.Parse("2006-01-02T15:04:05Z", item.Published)

		facts = append(facts, websearch.Fact{
			Title:   item.Title,
			Snippet: item.Description,
			Url:     item.Link,
			Date:    &published,
		})
	}

	return facts, nil
}
