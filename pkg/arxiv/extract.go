//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package arxiv

import (
	"bytes"
	"context"
	"fmt"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
)

type Extracts struct {
	ID string
}

func (api *ArXiv) Extracts(ctx context.Context, req Extracts) (string, error) {

	var buf bytes.Buffer
	err := api.Stack.IO(ctx,
		http.GET(
			ø.URI(fmt.Sprintf("https://ar5iv.labs.arxiv.org/html/%s", req.ID)),

			ƒ.Status.OK,
			ƒ.Bytes(&buf),
		),
	)
	if err != nil {
		return "", err
	}

	md, err := api.html2md.ConvertString(buf.String(),
		converter.WithDomain("ar5iv.labs.arxiv.org"),
	)
	if err != nil {
		return "", err
	}

	return md, nil
}
