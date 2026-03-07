//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package opts

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/fogfish/websearch"
	_ "github.com/fogfish/websearch/pkg/arxiv/adapter"
	_ "github.com/fogfish/websearch/pkg/duckduckgo/adapter"
	_ "github.com/fogfish/websearch/pkg/hackernews/adapter"
	_ "github.com/fogfish/websearch/pkg/webkit/adapter"
	service "github.com/fogfish/websearch/pkg/websearch"
	_ "github.com/fogfish/websearch/pkg/wikipedia/adapter"
	"github.com/spf13/cobra"
)

func Execute(vsn string) {
	rootCmd.Version = vsn

	if err := rootCmd.Execute(); err != nil {
		e := err.Error()
		fmt.Fprintf(os.Stderr, "\n ❌ Something went wrong. Check the error below for details.\n   Run `websearch help` for guidance.\n\n   %s\n\n", strings.ToUpper(e[:1])+e[1:])
		os.Exit(1)
	}
}

var (
	provider string
	latest   string
	format   string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&provider, "provider", "p", "", "search provider to use")

	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(extractCmd)

	searchCmd.Flags().StringVar(&latest, "latest", "", "filter only latest results (1w, 1m, 1y)")
	searchCmd.Flags().StringVar(&format, "format", "json", "output format (json, text)")
}

var rootCmd = &cobra.Command{
	Use:          "websearch",
	Short:        "search and extract web content",
	Long:         `search and extract web content`,
	SilenceUsage: true,
	Run:          func(cmd *cobra.Command, args []string) { cmd.Help() },
}

// Print JSON to stdout
func PrintJSON(data any) (err error) {
	cookie, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var obj any //map[string]any
	json.Unmarshal(cookie, &obj)

	f := colorjson.NewFormatter()
	f.Indent = 2

	encoded, err := f.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(encoded)
	return err
}

func PrintText(data []websearch.Fact) {
	for _, fact := range data {
		var sb strings.Builder
		fmt.Fprintf(&sb, "## %s\n", fact.Title)
		if fact.Date != nil {
			fmt.Fprintf(&sb, "Date: %s\n", fact.Date.Format(time.RFC1123))
		}
		fmt.Fprintf(&sb, "URL: %s\n\n", fact.Url)
		if len(fact.Snippet) > 0 {
			sb.WriteString(fact.Snippet)
			sb.WriteString("\n\n")
		}
		sb.WriteString("\n")
		os.Stdout.Write([]byte(sb.String()))
	}
}

//------------------------------------------------------------------------------

var searchCmd = &cobra.Command{
	Use:          "search",
	Short:        "search web content",
	Long:         `search web content`,
	SilenceUsage: true,
	RunE:         search,
}

func search(cmd *cobra.Command, args []string) error {
	srv, err := service.New(provider)
	if err != nil {
		return err
	}
	defer srv.Close()

	_, reply, err := srv.Search(cmd.Context(), nil, service.SearchInput{
		Query: strings.Join(args, " "),
	})
	if err != nil {
		return err
	}
	reply.Facts = websearch.OnlyLatest(latest, reply.Facts)

	switch format {
	case "json":
		PrintJSON(reply)
	case "text":
		PrintText(reply.Facts)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	return nil
}

//------------------------------------------------------------------------------

var extractCmd = &cobra.Command{
	Use:          "extract",
	Short:        "extract web content",
	Long:         `extract web content`,
	SilenceUsage: true,
	RunE:         extract,
}

func extract(cmd *cobra.Command, args []string) error {
	srv, err := service.New(provider)
	if err != nil {
		return err
	}
	defer srv.Close()

	for _, url := range args {
		_, reply, err := srv.Extracts(cmd.Context(), nil, service.ExtractInput{
			Url: url,
		})
		if err != nil {
			return err
		}
		os.Stdout.Write([]byte(reply.Content))
	}

	return nil
}
