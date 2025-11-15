//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package service

import "github.com/fogfish/websearch"

type Provider string

const (
	ProviderWikipedia  Provider = "wikipedia"
	ProviderDuckDuckGo Provider = "duckduckgo"
	ProviderHackerNews Provider = "hackernews"
	ProviderWebkit     Provider = "webkit"
)

type SearchInput struct {
	Query   string   `json:"query" jsonschema:"the query to search for"`
	Hashtag []string `json:"hashtag,omitempty" jsonschema:"hashtag to scope search results"`
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
