// Copyright © 2015 Daniele Tricoli <eriol@mornie.org>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package format // import "eriol.xyz/piken/format"

import (
	"strconv"
	"testing"

	"eriol.xyz/piken/sql"

	"github.com/stretchr/testify/assert"
)

func TestCodePointToGlyph(t *testing.T) {
	glyph, err := CodePointToGlyph("1F602")
	assert.Nil(t, err)
	assert.Equal(t, glyph, "😂")

	glyph, err = CodePointToGlyph("1000000000")
	assert.Equal(t, glyph, "")
	if assert.Error(t, err) {
		assert.Equal(t, err.(*strconv.NumError).Err, strconv.ErrRange)
	}
}

func TestFormat(t *testing.T) {

	s := sql.UnicodeData{CodePoint: "1F602",
		Name: "FACE WITH TEARS OF JOY"}

	formatter := NewTextFormatter(
		[]string{"CodePoint", "Name"}, " -- ", true)
	b, _ := formatter.Format(&s)
	assert.Equal(t, b, "1F602 -- FACE WITH TEARS OF JOY -- 😂")

	formatter = NewTextFormatter(
		[]string{"CodePoint", "Name"}, " ## ", true)
	b, _ = formatter.Format(&s)
	assert.Equal(t, b, "1F602 ## FACE WITH TEARS OF JOY ## 😂")

	formatter = NewTextFormatter(
		[]string{"Name"}, " -- ", true)
	b, _ = formatter.Format(&s)
	assert.Equal(t, b, "FACE WITH TEARS OF JOY -- 😂")

	formatter = NewTextFormatter(
		[]string{"Name"}, " -- ", false)
	b, _ = formatter.Format(&s)
	assert.Equal(t, b, "FACE WITH TEARS OF JOY")
}
