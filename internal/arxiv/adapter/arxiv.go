//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package adapter

import (
	"context"
	"fmt"
	"strings"

	"github.com/fogfish/websearch/internal/arxiv"
	"github.com/fogfish/websearch/internal/service"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func init() {
	srv, err := arxiv.New(arxiv.Config{})
	if err != nil {
		panic(err)
	}

	service.Register(&arXiv{srv})
}

type arXiv struct {
	*arxiv.ArXiv
}

func (api *arXiv) ID() string { return "arxiv" }

func (api *arXiv) Close() error { return nil }

func (api *arXiv) Search(ctx context.Context, req *mcp.CallToolRequest, input service.SearchInput) (*mcp.CallToolResult, service.SearchReply, error) {
	request := arxiv.Search{Query: input.Query}
	if input.Sort == service.SortRelevance {
		request.SortBy = "relevance"
		request.OrderBy = "ascending"
	}

	facts, err := api.ArXiv.Search(ctx, request)
	if err != nil {
		return nil, service.SearchReply{}, err
	}

	return nil, service.SearchReply{Facts: facts}, nil
}

func (api *arXiv) Extracts(ctx context.Context, req *mcp.CallToolRequest, input service.ExtractInput) (*mcp.CallToolResult, service.ExtractReply, error) {
	if !strings.HasPrefix(input.Url, "https://arxiv.org/abs/") {
		return nil, service.ExtractReply{}, fmt.Errorf("invalid arxiv.org url: %s", input.Url)
	}

	id := strings.TrimPrefix(input.Url, "https://arxiv.org/abs/")

	content, err := api.ArXiv.Extracts(ctx, arxiv.Extracts{
		ID: id,
	})
	if err != nil {
		return nil, service.ExtractReply{}, err
	}

	return nil, service.ExtractReply{Content: content}, nil
}
