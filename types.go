//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package websearch

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
}
