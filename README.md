<p align="center">
  <p align="center"><strong>Web Search Model-Context-Protocol (MCP) Server</strong></p>
</p>

<p align="center">
  <a href="https://github.com/fogfish/websearch/releases">
    <img src="https://img.shields.io/github/v/tag/fogfish/websearch?label=version" />
  </a>
  <a href="https://github.com/fogfish/websearch/actions/">
    <img src="https://github.com/fogfish/websearch/workflows/build/badge.svg" />
  </a>
  <a href="http://github.com/fogfish/websearch">
    <img src="https://img.shields.io/github/last-commit/fogfish/websearch.svg" />
  </a>
  <a href="https://coveralls.io/github/fogfish/websearch?branch=main">
    <img src="https://coveralls.io/repos/github/fogfish/websearch/badge.svg?branch=main" />
  </a>
  <a href="LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" />
  </a>
</p>

`websearch` is Model Context Protocol (MCP) server to search web content within agentic workflows. It focuse on search within content platform and providers rather than generic search.  


## ✨ Features

- **Search accross content platforms**: Query across Wikipedia, DuckDuckGo, Hacker News, and others
- **MCP Integration**: Native support for Model Context Protocol servers
- **CLI Tools**: Command-line interface for direct usage
- **Markdown Output**: Convert web content to clean markdown format


## Installation

<details>
<summary>Build from sources</summary>

Requires Go [installed](https://go.dev/doc/install).

```bash
go install github.com/fogfish/websearch/cmd/websearch@latest
```

</details>


## Quick Start

Enable discovery of HackerNews in VS Code using MCP 

```json
{
  "servers": {
    "hackernews": {
      "command": "websearch",
      "args": [
        "serve",
        "--provider",
        "hackernews"
      ],
      "env": {}
    }
  }
}
```

The utlity is also enables search in command cline

```bash
websearch search --provider wikipedia "artificial intelligence"
```

## Supported Content Platforms

* [Hacker News](https://hn.algolia.com/api)
* [Wikipedia](https://en.wikipedia.org/w/api.php?action=help&modules=query%2Bsearch)
* DuckDuckGo Instant Answer API

### On our future development list

- [arxiv.org](https://info.arxiv.org/help/api/index.html)
- [OpenLibrary](https://openlibrary.org/developers/api)
- [CrossRef](https://api.crossref.org)
- [Semantic Scholar API](https://www.semanticscholar.org/product/api)
- [PubMed](https://www.ncbi.nlm.nih.gov/books/NBK25497/)
- [Nonprofit Explorer API](https://projects.propublica.org/nonprofits/api)
- [SERP API](https://serper.dev)
- [Tavily MCP](https://docs.tavily.com/documentation/api-reference/endpoint/search)
- Substack


#### Extract Web Content

Extract content from a URL

```bash 
websearch extract https://example.com
```

## Contributing

`websearch` is [MIT](LICENSE) licensed and accepts contributions via GitHub pull requests:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

### Adding New Providers

1. Implement the provider interface in `internal/{provider}/`
2. Add provider constant to `internal/service/types.go`
3. Update CLI flags and help text


## License

[![See LICENSE](https://img.shields.io/github/license/fogfish/iq.svg?style=for-the-badge)](LICENSE)
