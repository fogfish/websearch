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

	"github.com/TylerBrock/colorjson"
	"github.com/fogfish/websearch/internal/service"
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

var provider string

func init() {
	rootCmd.PersistentFlags().StringVarP(&provider, "provider", "p", string(service.ProviderDuckDuckGo), "search provider to use (wikipedia, duckduckgo)")

	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(extractCmd)
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

//------------------------------------------------------------------------------

var searchCmd = &cobra.Command{
	Use:          "search",
	Short:        "search web content",
	Long:         `search web content`,
	SilenceUsage: true,
	RunE:         search,
}

func search(cmd *cobra.Command, args []string) error {
	srv, err := service.New(service.Provider(provider))
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
	PrintJSON(reply)

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
	srv, err := service.New(service.Provider(provider))
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
