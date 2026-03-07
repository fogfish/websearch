//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package wikipedia

import (
	"context"
	"fmt"

	"github.com/fogfish/gurl/v2/http"
	"github.com/fogfish/websearch"
)

// The request type for wikipedia query search module
// https://en.wikipedia.org/w/api.php?action=help&modules=query%2Bsearch
type Search struct {
	_       int           `wiki:"list=search"`
	Lang    string        `wiki:"-"`
	Type    SearchType    `wiki:"srwhat"`      // The type of search to perform.
	Query   string        `wiki:"srsearch"`    // The search query string.
	Limit   int           `wiki:"srlimit"`     // The maximum number of search results to return.
	Offset  int           `wiki:"sroffset"`    // The offset for the search results.
	Ranking SearchRanking `wiki:"srqiprofile"` // The ranking method to use.
}

// Ranking methods for wikipedia search
type SearchRanking string

const (
	// Ranking based on the number of incoming links, some templates, page language and recency (templates/language/recency may not be activated on this wiki).
	RankingClassic = SearchRanking("classic")

	// Ranking based on some templates, page language and recency when activated on this wiki.
	RankingClassicNoBoost = SearchRanking("classic_noboostlinks")

	// Ranking based solely on query dependent features (for debug only).
	RankingEmpty = SearchRanking("empty")

	// Weighted sum based on incoming links
	RankingIncomingLinks = SearchRanking("wsum_inclinks")

	//Weighted sum based on incoming links and weekly pageviews
	RankingIncomingLinksWeeklyViews = SearchRanking("wsum_inclinks_pv")

	// TODO: add more rankings as needed
	// * popular_inclinks_pv Ranking based primarily on page views
	// * popular_inclinks Ranking based primarily on incoming link counts
	// * mlr-1024rs Weighted sum based on incoming links and weekly pageviews
	// * mlr-1024rs-next Weighted sum based on incoming links and weekly pageviews
	// * growth_underlinked Internal rescore profile used in GrowthExperiments link recommendations for prioritizing articles which do not yet have enough links. This is a no-op when Link Recommendations are disabled.
	// * engine_autoselect Let the search engine decide on the best profile to use.
)

// Type of search to perform
type SearchType string

const (
	SearchTypeText      = SearchType("text")      // Search the text of pages.
	SearchTypeTitle     = SearchType("title")     // Search only page titles.
	SearchTypeNearMatch = SearchType("nearmatch") // Search only page titles that begin with the search term.
)

type searchBag struct {
	BatchComplete bool           `json:"batchcomplete"`
	Continue      map[string]any `json:"continue,omitempty"`
	Query         searchQuery    `json:"query"`
}

type searchQuery struct {
	Info searchInfo  `json:"searchinfo,omitempty"`
	Hits []searchHit `json:"search,omitempty"`
}

type searchInfo struct {
	TotalHits  int    `json:"totalhits,omitempty"`
	Suggestion string `json:"suggestion,omitempty"`
}

type searchHit struct {
	ID      int    `json:"pageid,omitempty"`
	Title   string `json:"title,omitempty"`
	Snippet string `json:"snippet,omitempty"`
}

// Extracts wikipedia article content
// https://en.wikipedia.org/w/api.php?action=help&modules=query%2Bimageinfo
func (api *Wikipedia) Search(ctx context.Context, req Search) ([]websearch.Fact, error) {
	bag, err := http.IO[searchBag](api.WithContext(ctx), api.query(req.Lang, req))
	if err != nil {
		return nil, err
	}

	facts := make([]websearch.Fact, 0, len(bag.Query.Hits))
	for _, hit := range bag.Query.Hits {
		fact := websearch.Fact{
			Title:   hit.Title,
			Snippet: hit.Snippet,
			Url:     fmt.Sprintf("https://en.wikipedia.org/?curid=%d", hit.ID),
		}
		facts = append(facts, fact)
	}

	return facts, nil
}
