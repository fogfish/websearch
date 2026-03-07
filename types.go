//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package websearch

import "time"

// Fact represents a single search result fact.
// The type acts as denominator for various search engines.
type Fact struct {
	// Unique identifier for the fact, document, or resource.
	ID string `json:"id,omitempty"`

	// Category or group the fact belongs to.
	Category string `json:"category,omitempty"`

	// Human-readable headline, title about the fact.
	Title string `json:"title,omitempty"`

	// Short description or abstract of the fact.
	Snippet string `json:"snippet,omitempty"`

	// Url to the resource represented by the fact.
	Url string `json:"url,omitempty"`

	// Date when the fact was created, published, or indexed.
	Date *time.Time `json:"date,omitempty"`
}

func OnlyLatest(wnd string, facts []Fact) []Fact {
	if len(wnd) == 0 {
		return facts
	}

	var dur time.Duration
	switch wnd[len(wnd)-1] {
	case 'w':
		dur = 7 * 24 * time.Hour
	case 'm':
		dur = 30 * 24 * time.Hour
	case 'y':
		dur = 365 * 24 * time.Hour
	default:
		return facts
	}

	cutoff := time.Now().Add(-dur)
	var latest []Fact
	for _, fact := range facts {
		if fact.Date != nil && fact.Date.After(cutoff) {
			latest = append(latest, fact)
		}
	}
	return latest
}
