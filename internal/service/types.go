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

	"github.com/fogfish/websearch"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Provider interface {
	ID() string
	Close() error
}

type Searcher interface {
	Search(ctx context.Context, req *mcp.CallToolRequest, input SearchInput) (*mcp.CallToolResult, SearchReply, error)
}

type Extractor interface {
	Extracts(ctx context.Context, req *mcp.CallToolRequest, input ExtractInput) (*mcp.CallToolResult, ExtractReply, error)
}

type Sort string

const (
	SortRelevance Sort = "relevance"
	SortDate      Sort = "date"
)

type SearchInput struct {
	Query   string   `json:"query" jsonschema:"the query to search for"`
	Hashtag []string `json:"hashtag,omitempty" jsonschema:"hashtag to scope search results"`
	Sort    Sort     `json:"sort,omitempty" jsonschema:"the sort order of the search results, 'relevance' or 'date'"`
}

type SearchReply struct {
	Facts []websearch.Fact `json:"facts" jsonschema:"the facts to return to the client"`
}

type ExtractInput struct {
	Url string `json:"url" jsonschema:"the url the web page to extract content from"`
}

type ExtractReply struct {
	Content string `json:"content" jsonschema:"the content extracted from the web page"`
}
