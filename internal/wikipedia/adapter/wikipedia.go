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
	"net/url"

	"github.com/fogfish/websearch/internal/service"
	"github.com/fogfish/websearch/internal/wikipedia"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func init() {
	srv, err := wikipedia.New(wikipedia.Config{})
	if err != nil {
		panic(err)
	}

	service.Register(&wiki{srv})
}

type wiki struct {
	*wikipedia.Wikipedia
}

func (api *wiki) ID() string { return "wikipedia" }

func (api *wiki) Close() error { return nil }

func (api *wiki) Search(ctx context.Context, req *mcp.CallToolRequest, input service.SearchInput) (*mcp.CallToolResult, service.SearchReply, error) {
	facts, err := api.Wikipedia.Search(ctx, wikipedia.Search{
		Query: input.Query,
	})
	if err != nil {
		return nil, service.SearchReply{}, err
	}

	return nil, service.SearchReply{Facts: facts}, nil
}

func (api *wiki) Extracts(ctx context.Context, req *mcp.CallToolRequest, input service.ExtractInput) (*mcp.CallToolResult, service.ExtractReply, error) {
	u, err := url.Parse(input.Url)
	if err != nil {
		return nil, service.ExtractReply{}, err
	}
	pageid := u.Query().Get("curid")

	title, body, err := api.Wikipedia.Extracts(ctx, wikipedia.Extracts{
		ID:         wikipedia.ID("pageid/" + pageid).Norm(),
		TextOnly:   true,
		TextFormat: wikipedia.FORMAT_TEXT,
		OnlyIntro:  true,
	})
	if err != nil {
		return nil, service.ExtractReply{}, err
	}

	content := title + "\n\n" + body
	return nil, service.ExtractReply{Content: content}, nil
}
