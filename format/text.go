// Copyright © 2015 Daniele Tricoli <eriol@mornie.org>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package format // import "eriol.xyz/piken/format"

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"eriol.xyz/piken/sql"
)

type TextFormatter struct {
	baseFormatter
	Separator string
}

func NewTextFormatter(fields []string, separator string, glyph bool) *TextFormatter {

	return &TextFormatter{
		baseFormatter: baseFormatter{fields: fields, showGlyph: glyph},
		Separator:     separator}
}

func (t *TextFormatter) Format(s *sql.UnicodeData) (string, error) {

	var buffer []string

	glyph, err := CodePointToGlyph(s.CodePoint)
	if err != nil {
		return "", err
	}

	for _, field := range t.Fields() {
		r := reflect.ValueOf(s)
		f := reflect.Indirect(r).FieldByName(field)
		buffer = append(buffer, f.String())
	}

	if t.ShowGlyph() {
		buffer = append(buffer, glyph)
	}

	return strings.Join(buffer, t.Separator), nil
}

// Convert an unicode codepoint into a string.
func CodePointToGlyph(codepoint string) (string, error) {

	s, err := strconv.ParseInt(codepoint, 16, 32)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%c", s), nil
}
