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

	"github.com/fogfish/websearch/pkg/webkit"
	service "github.com/fogfish/websearch/pkg/websearch"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func init() {
	srv, err := webkit.New(webkit.Config{
		AutoConfig: true,
		DriverDir:  "/tmp/websearch",
	})
	if err != nil {
		panic(err)
	}

	service.Register(&web{srv})
}

type web struct {
	*webkit.WebKit
}

func (api *web) ID() string { return "webkit" }

func (api *web) Close() error { return api.WebKit.Close() }

func (api *web) Extracts(ctx context.Context, req *mcp.CallToolRequest, input service.ExtractInput) (*mcp.CallToolResult, service.ExtractReply, error) {
	content, err := api.WebKit.Extract(input.Url)
	if err != nil {
		return nil, service.ExtractReply{}, err
	}

	return nil, service.ExtractReply{Content: content}, nil
}
