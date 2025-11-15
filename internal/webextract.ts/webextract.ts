#!/usr/bin/env node
import { Server } from '@modelcontextprotocol/sdk/server/index.js';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';
import {
  CallToolRequestSchema,
  ErrorCode,
  ListToolsRequestSchema,
  McpError,
} from '@modelcontextprotocol/sdk/types.js';
import { load } from 'cheerio'
import type { Browser, BrowserContext } from 'playwright'
import { webkit as web } from 'playwright'
import { NodeHtmlMarkdown } from 'node-html-markdown'


class WebExtractServer {
  private server: Server;
  private pageReader: PageReader;

  constructor() {
    this.server = new Server(
      {
        name: 'webextract',
        version: '0.0.1',
      },
      {
        capabilities: {
          tools: {},
        },
      }
    );

    this.pageReader = new PageReader();
    this.setupToolHandlers();
  }

  private setupToolHandlers() {
    // List available tools
    this.server.setRequestHandler(ListToolsRequestSchema, async () => {
      return {
        tools: [
          {
            name: 'extract',
            description: 'Extract a webpage content to markdown format',
            inputSchema: {
              type: 'object',
              properties: {
                url: {
                  type: 'string',
                  description: 'The URL of the webpage to fetch',
                },
              },
              required: ['url'],
            },
          },
        ],
      };
    });

    // Handle tool calls
    this.server.setRequestHandler(CallToolRequestSchema, async (request) => {
      const { name, arguments: args } = request.params;

      try {
        switch (name) {
          case 'extract':
            return await this.handleExtract(args);
          default:
            throw new McpError(
              ErrorCode.MethodNotFound,
              `Unknown tool: ${name}`
            );
        }
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : String(error);
        throw new McpError(ErrorCode.InternalError, `Tool execution failed: ${errorMessage}`);
      }
    });
  }

  private async handleExtract(args: any) {
    const { url } = args;

    if (!url || typeof url !== 'string') {
      throw new McpError(ErrorCode.InvalidParams, 'URL is required and must be a string');
    }

    const webpage = await this.pageReader.read(url);

    const result: any = {
      url: webpage.url,
      format: 'markdown',
      content: webpage.markdown,
    };

    return {
      content: [
        {
          type: 'text',
          text: JSON.stringify(result, null, 2),
        },
      ],
    };
  }


  async run() {
    await this.pageReader.init();

    const transport = new StdioServerTransport();
    await this.server.connect(transport);

    console.error('WebSearch MCP server running on stdio');

    // Handle cleanup on exit
    process.on('SIGINT', async () => {
      await this.cleanup();
      process.exit(0);
    });

    process.on('SIGTERM', async () => {
      await this.cleanup();
      process.exit(0);
    });
  }

  private async cleanup() {
    console.error('Cleaning up WebSearch MCP server...');
    try {
      await this.pageReader.dispose();
    } catch (error) {
      console.error('Error during cleanup:', error);
    }
  }
}

// //------------------------------------------------------------------------------

// class WebSearchServer {
//   private server: Server;
//   private pageReader: PageReader;

//   constructor() {
//     this.server = new Server(
//       {
//         name: 'websearch',
//         version: '0.0.1',
//       },
//       {
//         capabilities: {
//           tools: {},
//         },
//       }
//     );

//     this.pageReader = new PageReader();
//     this.setupToolHandlers();
//   }

//   private setupToolHandlers() {
//     // List available tools
//     this.server.setRequestHandler(ListToolsRequestSchema, async () => {
//       return {
//         tools: [
//           {
//             name: 'search',
//             description: 'Search for a webpage and extract its content to markdown format',
//             inputSchema: {
//               type: 'object',
//               properties: {
//                 query: {
//                   type: 'string',
//                   description: 'The search query to find the webpage',
//                 },
//               },
//               required: ['query'],
//             },
//           },
//         ],
//       };
//     });

//     // Handle tool calls
//     this.server.setRequestHandler(CallToolRequestSchema, async (request) => {
//       const { name, arguments: args } = request.params;

//       try {
//         switch (name) {
//           case 'search':
//             return await this.handleExtract(args);
//           default:
//             throw new McpError(
//               ErrorCode.MethodNotFound,
//               `Unknown tool: ${name}`
//             );
//         }
//       } catch (error) {
//         const errorMessage = error instanceof Error ? error.message : String(error);
//         throw new McpError(ErrorCode.InternalError, `Tool execution failed: ${errorMessage}`);
//       }
//     });
//   }

//   private async handleExtract(args: any) {
//     const { query } = args;

//     if (!query || typeof query !== 'string') {
//       throw new McpError(ErrorCode.InvalidParams, 'Query is required and must be a string');
//     }


//     const url = `https://duckduckgo.com/?origin=funnel_home_website&ia=web&q=${encodeURIComponent(query)}`;
//     const webpage = await this.pageReader.read(url);

//     const result: any = {
//       url: webpage.url,
//       format: 'markdown',
//       content: webpage.markdown,
//     };

//     return {
//       content: [
//         {
//           type: 'text',
//           text: JSON.stringify(result, null, 2),
//         },
//       ],
//     };
//   }


//   async run() {
//     await this.pageReader.init();

//     const transport = new StdioServerTransport();
//     await this.server.connect(transport);

//     console.error('WebSearch MCP server running on stdio');

//     // Handle cleanup on exit
//     process.on('SIGINT', async () => {
//       await this.cleanup();
//       process.exit(0);
//     });

//     process.on('SIGTERM', async () => {
//       await this.cleanup();
//       process.exit(0);
//     });
//   }

//   private async cleanup() {
//     console.error('Cleaning up WebSearch MCP server...');
//     try {
//       await this.pageReader.dispose();
//     } catch (error) {
//       console.error('Error during cleanup:', error);
//     }
//   }
// }


//------------------------------------------------------------------------------

class PageReader {
  private browser?: Browser
  private context?: BrowserContext

  async init() {
    this.browser = await web.launch({
      headless: true,
    })

    this.context = await this.browser.newContext()
  }

  async read(pageUrl: string, selector?: string) {
    if (!this.context) {
      throw new Error('Browser context is not initialized. Call init() first.')
    }

    const page = await this.context.newPage()

    try {
      await page.goto(pageUrl)

      const pageHtml = await page.evaluate(() => {
        return globalThis.document.documentElement.outerHTML
      })

      const contentHtml = this.sanitizeHtml(pageHtml, selector)

      return {
        url: pageUrl,
        html: contentHtml,
        markdown: NodeHtmlMarkdown.translate(contentHtml),
      }
    } finally {
      await page.close()
    }
  }

  async dispose() {
    if (this.context) {
      await this.context.close()
    }

    if (this.browser) {
      await this.browser.close()
    }
  }

  private sanitizeHtml(html: string, selector?: string) {
    const $ = load(html)

    if (selector) {
      const selectedHtml = $(selector).html()

      if (!selectedHtml || !selectedHtml.trim()) {
        throw new Error(`No content found for selector: ${selector}`)
      }

      return selectedHtml
    }

    $('script, style, path, footer, header, head').remove()

    return $.html()
  }
}


// CLI mode for testing
async function cli() {
  var url = process.argv[2];
  if (!url) {
    console.error('Usage: tsx mcp-server.ts <url> [selector]');
    console.error('       tsx mcp-server.ts search query');
    console.error('       tsx mcp-server.ts serve [search] (to run as MCP server)');
    process.exit(1);
  }

  const selector = process.argv[3];

  const pageReader = new PageReader();

  try {
    await pageReader.init();
    const result = await pageReader.read(url, selector);

    console.log('# Webpage Content');
    console.log(`**URL:** ${result.url}`);
    console.log(`**Fetched:** ${new Date().toISOString()}`);
    console.log('---');
    console.log(result.markdown);
  } catch (error) {
    console.error('Error:', error);
    process.exit(1);
  } finally {
    await pageReader.dispose();
  }
}

// Main entry point
async function main() {
  const args = process.argv.slice(2);

  if (args.includes('serve') || process.env.MCP_SERVER === 'true') {
    const server = new WebExtractServer();
    await server.run();
  } else {
    await cli();
  }
}

if (import.meta.url === `file://${process.argv[1]}`) {
  main().catch((error) => {
    console.error('Fatal error:', error);
    process.exit(1);
  });
}