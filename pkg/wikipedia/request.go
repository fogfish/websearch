//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/websearch
//

package wikipedia

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/fogfish/gold"
)

var _ = gold.Register[Wikipedia]("wiki")

type ID = gold.IRI[Wikipedia]

const EN = "en"

var ErrNotFound = errors.New("not found")

//------------------------------------------------------------------------------

// Wikipedia api request object.
type request map[string]string

func marshalQuery(opts any) request {
	v := reflect.ValueOf(opts)
	t := v.Type()

	r := request{"action": "query"}
	for i := 0; i < t.NumField(); i++ {
		r.injectValue(t.Field(i), v.Field(i))
	}

	switch v := opts.(type) {
	case interface{ Encode(request) }:
		v.Encode(r)
	}

	return r
}

func (r request) injectValue(t reflect.StructField, v reflect.Value) {
	tag := t.Tag.Get("wiki")
	if tag == "" {
		return
	}

	// This is a special field that defines wikipedia module
	if t.Name == "_" {
		seq := strings.Split(tag, "=")
		if len(seq) == 2 {
			r[seq[0]] = seq[1]
		} else {
			r["prop"] = tag
		}

		if ext := t.Tag.Get("wext"); ext != "" {
			for x := range strings.SplitSeq(ext, ",") {
				seq := strings.Split(x, "=")
				if len(seq) == 2 {
					r[seq[0]] = seq[1]
				}
			}
		}

		return
	}

	// This is a special object identity (article, files, etc).
	if tag == "id" {
		iri, err := gold.AsIRI[Wikipedia](v.String())
		if err != nil {
			return
		}

		ref := iri.Reference()
		if strings.HasPrefix(ref, "pageid/") {
			r["pageids"] = ref[7:]
		} else {
			r["titles"] = ref
		}
		return
	}

	switch v.Kind() {
	case reflect.String:
		s := v.String()
		if len(s) > 0 {
			r[tag] = s
		}
	case reflect.Int:
		i := int(v.Int())
		if i > 0 {
			r[tag] = strconv.Itoa(i)
		}
	case reflect.Bool:
		b := v.Bool()
		if b {
			r[tag] = ""
		}
	case reflect.Float64:
		f := v.Float()
		if f > 0.0 {
			r[tag] = fmt.Sprintf("%f", f)
		}
	}
}

//------------------------------------------------------------------------------

// bag is a container of Wikipedia reply
type bag[T any] struct {
	BatchComplete bool              `json:"batchcomplete"`
	Continue      map[string]string `json:"continue,omitempty"`
	Query         query[T]          `json:"query,omitempty"`
}

// query is a container of Query reply.
type query[T any] struct {
	Pages []T `json:"pages,omitempty"`
}
