//
// Copyright (C) 2025 - 2026 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package adapter

import (
	"context"

	"github.com/fogfish/websearch/internal/hackernews"
	"github.com/fogfish/websearch/internal/service"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func init() {
	srv, err := hackernews.New(hackernews.Config{})
	if err != nil {
		panic(err)
	}

	service.Register(&hn{srv})
}

type hn struct {
	*hackernews.HackerNews
}

func (api *hn) ID() string { return "hackernews" }

func (api *hn) Close() error { return nil }

func (api *hn) Search(ctx context.Context, req *mcp.CallToolRequest, input service.SearchInput) (*mcp.CallToolResult, service.SearchReply, error) {
	request := hackernews.Search{
		Query: input.Query,
		Tags:  input.Hashtag,
	}
	if input.Sort == service.SortRelevance {
		request.SortBy = "relevance"
	}

	facts, err := api.HackerNews.Search(ctx, request)
	if err != nil {
		return nil, service.SearchReply{}, err
	}

	return nil, service.SearchReply{Facts: facts}, nil
}
