//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package opts

import (
	"context"

	"github.com/fogfish/websearch/internal/service"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start websearch as a server",
	Long: `
Start websearch as an MCP server exposing it as a tool to agents.
`,
	Example: `
	websearch serve 
	`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          serve,
}

func serve(cmd *cobra.Command, args []string) error {
	srv, err := service.New(service.Provider(provider))
	if err != nil {
		return err
	}
	defer srv.Close()

	if err := srv.Run(context.Background()); err != nil {
		return err
	}

	return nil
}
