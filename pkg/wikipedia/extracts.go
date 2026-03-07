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

	"github.com/fogfish/gurl/v2/http"
)

// The request type for wikipedia extracts module
// https://en.wikipedia.org/w/api.php?action=help&modules=query%2Bextracts
type Extracts struct {
	_          int        `wiki:"extracts"`
	Lang       string     `wiki:"-"`
	ID         ID         `wiki:"id"`
	TextOnly   bool       `wiki:"explaintext"`
	TextFormat TextFormat `wiki:"exsectionformat"`
	Chars      int        `wiki:"exchars"`     // How many characters to return. Actual text returned might be slightly longer.
	Sentences  int        `wiki:"exsentences"` // How many sentences to return.
	OnlyIntro  bool       `wiki:"exintro"`     // Return only content before the first section
}

type TextFormat string

const (
	FORMAT_WIKI = TextFormat("wiki")
	FORMAT_TEXT = TextFormat("plain")
	FORMAT_RAW  = TextFormat("raw")
)

type extracts struct {
	Title   string `json:"title,omitempty"`
	Extract string `json:"extract,omitempty"`
}

// Extracts wikipedia article content
// https://en.wikipedia.org/w/api.php?action=help&modules=query%2Bimageinfo
func (api *Wikipedia) Extracts(ctx context.Context, req Extracts) (string, string, error) {
	bag, err := http.IO[bag[extracts]](api.WithContext(ctx), api.query(req.Lang, req))
	if err != nil {
		return "", "", err
	}

	if len(bag.Query.Pages) == 0 {
		return "", "", ErrNotFound
	}

	page := bag.Query.Pages[0]

	return page.Title, page.Extract, nil
}
