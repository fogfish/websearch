//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package websearch

import (
	"strconv"
	"time"
)

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

// toDuration parses a duration string where the last character is the
// time unit (h, d, w, m, y) and an optional leading integer is the multiplier.
// For example: "1h", "5w", "2m". Returns 0, false if the string is invalid.
func toDuration(wnd string) (time.Duration, bool) {
	if len(wnd) == 0 {
		return 0, false
	}

	dim := wnd[len(wnd)-1]
	numStr := wnd[:len(wnd)-1]

	n := 1
	if len(numStr) > 0 {
		var err error
		n, err = strconv.Atoi(numStr)
		if err != nil || n <= 0 {
			return 0, false
		}
	}

	var unit time.Duration
	switch dim {
	case 'h':
		unit = time.Hour
	case 'd':
		unit = 24 * time.Hour
	case 'w':
		unit = 7 * 24 * time.Hour
	case 'm':
		unit = 30 * 24 * time.Hour
	case 'y':
		unit = 365 * 24 * time.Hour
	default:
		return 0, false
	}

	return time.Duration(n) * unit, true
}

func OnlyLatest(wnd string, facts []Fact) []Fact {
	if len(wnd) == 0 {
		return facts
	}

	dur, ok := toDuration(wnd)
	if !ok {
		return facts
	}

	cutoff := time.Now().Add(-dur)
	return LatestAfter(cutoff, facts)
}

func LatestAfter(cutoff time.Time, facts []Fact) []Fact {
	var latest []Fact
	for _, fact := range facts {
		if fact.Date != nil && fact.Date.After(cutoff) {
			latest = append(latest, fact)
		}
	}
	return latest
}
