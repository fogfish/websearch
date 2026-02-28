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

	"github.com/fogfish/websearch/internal/duckduckgo"
	"github.com/fogfish/websearch/internal/service"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func init() {
	srv, err := duckduckgo.New(duckduckgo.Config{})
	if err != nil {
		panic(err)
	}

	service.Register(&ddg{srv})
}

type ddg struct {
	*duckduckgo.DuckDuckGo
}

func (api *ddg) ID() string { return "duckduckgo" }

func (api *ddg) Close() error { return nil }

func (api *ddg) Search(ctx context.Context, req *mcp.CallToolRequest, input service.SearchInput) (*mcp.CallToolResult, service.SearchReply, error) {
	facts, err := api.DuckDuckGo.Search(ctx, duckduckgo.Search{
		Query: input.Query,
	})
	if err != nil {
		return nil, service.SearchReply{}, err
	}

	return nil, service.SearchReply{Facts: facts}, nil
}
