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
	"net/url"

	"github.com/fogfish/websearch/internal/duckduckgo"
	"github.com/fogfish/websearch/internal/hackernews"
	"github.com/fogfish/websearch/internal/webkit"
	"github.com/fogfish/websearch/internal/wikipedia"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Server struct {
	srv      *mcp.Server
	api      provider
	provider Provider
}

type provider interface {
	Close() error
	Search(ctx context.Context, req *mcp.CallToolRequest, input SearchInput) (*mcp.CallToolResult, SearchReply, error)
	Extracts(ctx context.Context, req *mcp.CallToolRequest, input ExtractInput) (*mcp.CallToolResult, ExtractReply, error)
}

func New(provider Provider) (server *Server, err error) {
	server = &Server{provider: provider}

	switch provider {
	case ProviderWikipedia:
		api, err := wikipedia.New(wikipedia.Config{})
		if err != nil {
			return nil, err
		}
		server.api = &wiki{api}
	case ProviderDuckDuckGo:
		api, err := duckduckgo.New(duckduckgo.Config{})
		if err != nil {
			return nil, err
		}
		server.api = &ddg{api}
	case ProviderHackerNews:
		api, err := hackernews.New(hackernews.Config{})
		if err != nil {
			return nil, err
		}
		server.api = &hn{api}
	case ProviderWebkit:
		api, err := webkit.New(webkit.Config{
			AutoConfig: true,
			DriverDir:  "/tmp/websearch",
		})
		if err != nil {
			return nil, err
		}
		server.api = &web{api}
	default:
		api, err := wikipedia.New(wikipedia.Config{})
		if err != nil {
			return nil, err
		}
		server.api = &wiki{api}
	}

	return server, nil
}

func (srv *Server) Run(ctx context.Context) error {
	srv.srv = mcp.NewServer(&mcp.Implementation{Name: fmt.Sprintf("%s search and extract", srv.provider)}, nil)

	mcp.AddTool(srv.srv, &mcp.Tool{Name: "search", Description: fmt.Sprintf("Search %s content", srv.provider)}, srv.Search)
	switch srv.provider {
	case ProviderWikipedia:
		mcp.AddTool(srv.srv, &mcp.Tool{Name: "extracts", Description: fmt.Sprintf("Extract content from %s", srv.provider)}, srv.api.Extracts)
	}

	return srv.srv.Run(ctx, &mcp.StdioTransport{})
}

func (srv *Server) Close() error { return srv.api.Close() }

//------------------------------------------------------------------------------

func (srv *Server) Search(ctx context.Context, req *mcp.CallToolRequest, input SearchInput) (*mcp.CallToolResult, SearchReply, error) {
	return srv.api.Search(ctx, req, input)
}

func (srv *Server) Extracts(ctx context.Context, req *mcp.CallToolRequest, input ExtractInput) (*mcp.CallToolResult, ExtractReply, error) {
	return srv.api.Extracts(ctx, req, input)
}

//------------------------------------------------------------------------------

type wiki struct {
	*wikipedia.Wikipedia
}

func (api *wiki) Close() error { return nil }

func (api *wiki) Search(ctx context.Context, req *mcp.CallToolRequest, input SearchInput) (*mcp.CallToolResult, SearchReply, error) {
	facts, err := api.Wikipedia.Search(ctx, wikipedia.Search{
		Query: input.Query,
	})
	if err != nil {
		return nil, SearchReply{}, err
	}

	return nil, SearchReply{Facts: facts}, nil
}

func (api *wiki) Extracts(ctx context.Context, req *mcp.CallToolRequest, input ExtractInput) (*mcp.CallToolResult, ExtractReply, error) {
	u, err := url.Parse(input.Url)
	if err != nil {
		return nil, ExtractReply{}, err
	}
	pageid := u.Query().Get("curid")

	title, body, err := api.Wikipedia.Extracts(ctx, wikipedia.Extracts{
		ID:         wikipedia.ID("pageid/" + pageid).Norm(),
		TextOnly:   true,
		TextFormat: wikipedia.FORMAT_TEXT,
		OnlyIntro:  true,
	})
	if err != nil {
		return nil, ExtractReply{}, err
	}

	content := title + "\n\n" + body
	return nil, ExtractReply{Content: content}, nil
}

//------------------------------------------------------------------------------

type ddg struct {
	*duckduckgo.DuckDuckGo
}

func (api *ddg) Close() error { return nil }

func (api *ddg) Search(ctx context.Context, req *mcp.CallToolRequest, input SearchInput) (*mcp.CallToolResult, SearchReply, error) {
	facts, err := api.DuckDuckGo.Search(ctx, duckduckgo.Search{
		Query: input.Query,
	})
	if err != nil {
		return nil, SearchReply{}, err
	}

	return nil, SearchReply{Facts: facts}, nil
}

func (api *ddg) Extracts(ctx context.Context, req *mcp.CallToolRequest, input ExtractInput) (*mcp.CallToolResult, ExtractReply, error) {
	return nil, ExtractReply{}, fmt.Errorf("not implemented")
}

//------------------------------------------------------------------------------

type hn struct {
	*hackernews.HackerNews
}

func (api *hn) Close() error { return nil }

func (api *hn) Search(ctx context.Context, req *mcp.CallToolRequest, input SearchInput) (*mcp.CallToolResult, SearchReply, error) {
	facts, err := api.HackerNews.Search(ctx, hackernews.Search{
		Query: input.Query,
		Tags:  input.Hashtag,
	})
	if err != nil {
		return nil, SearchReply{}, err
	}

	return nil, SearchReply{Facts: facts}, nil
}

func (api *hn) Extracts(ctx context.Context, req *mcp.CallToolRequest, input ExtractInput) (*mcp.CallToolResult, ExtractReply, error) {
	return nil, ExtractReply{}, fmt.Errorf("not implemented")
}

//------------------------------------------------------------------------------

type web struct {
	*webkit.WebKit
}

func (api *web) Close() error { return api.WebKit.Close() }

func (api *web) Search(ctx context.Context, req *mcp.CallToolRequest, input SearchInput) (*mcp.CallToolResult, SearchReply, error) {
	return nil, SearchReply{}, fmt.Errorf("not implemented")
}

func (api *web) Extracts(ctx context.Context, req *mcp.CallToolRequest, input ExtractInput) (*mcp.CallToolResult, ExtractReply, error) {
	content, err := api.WebKit.Extract(input.Url)
	if err != nil {
		return nil, ExtractReply{}, err
	}

	return nil, ExtractReply{Content: content}, nil
}
