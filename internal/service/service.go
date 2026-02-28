//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	registry = map[string]Provider{}
	lock     sync.Mutex
)

// Register search providers
func Register(p Provider) {
	lock.Lock()
	defer lock.Unlock()

	registry[p.ID()] = p
}

//------------------------------------------------------------------------------

type Server struct {
	srv  *mcp.Server
	api  Provider
	sapi Searcher
	eapi Extractor
}

func New(id string) (server *Server, err error) {
	lock.Lock()
	defer lock.Unlock()

	srv := &Server{}

	api, ok := registry[id]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", id)
	}

	srv.api = api

	if sapi, ok := api.(Searcher); ok {
		srv.sapi = sapi
	}

	if eapi, ok := api.(Extractor); ok {
		srv.eapi = eapi
	}

	return srv, nil
}

func (srv *Server) Run(ctx context.Context) error {
	srv.srv = mcp.NewServer(&mcp.Implementation{Name: fmt.Sprintf("%s search and extract", srv.api.ID())}, nil)

	if srv.sapi != nil {
		mcp.AddTool(srv.srv, &mcp.Tool{Name: "search", Description: fmt.Sprintf("Search %s content", srv.api.ID())}, srv.Search)
	}

	if srv.eapi != nil {
		mcp.AddTool(srv.srv, &mcp.Tool{Name: "extracts", Description: fmt.Sprintf("Extract content from %s", srv.api.ID())}, srv.Extracts)
	}

	return srv.srv.Run(ctx, &mcp.StdioTransport{})
}

func (srv *Server) Close() error { return srv.api.Close() }

func (srv *Server) Search(ctx context.Context, req *mcp.CallToolRequest, input SearchInput) (*mcp.CallToolResult, SearchReply, error) {
	if srv.sapi == nil {
		return nil, SearchReply{}, fmt.Errorf("search not supported by provider: %s", srv.api.ID())
	}

	return srv.sapi.Search(ctx, req, input)
}

func (srv *Server) Extracts(ctx context.Context, req *mcp.CallToolRequest, input ExtractInput) (*mcp.CallToolResult, ExtractReply, error) {
	if srv.eapi == nil {
		return nil, ExtractReply{}, fmt.Errorf("extract not supported by provider: %s", srv.api.ID())
	}

	return srv.eapi.Extracts(ctx, req, input)
}
