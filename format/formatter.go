package format // import "eriol.xyz/piken/format"

import "eriol.xyz/piken/sql"

type Formatter interface {
	Format(*sql.UnicodeData) (string, error)
}

type baseFormatter struct {
	fields    []string
	showGlyph bool
}

func (bf *baseFormatter) SetFields(fields []string) {
	for _, field := range fields {
		bf.fields = append(bf.fields, field)
	}
}

func (bf baseFormatter) Fields() []string {
	return bf.fields
}

func (bf *baseFormatter) SetShowGlyph(value bool) {
	bf.showGlyph = value
}

func (bf baseFormatter) ShowGlyph() bool {
	return bf.showGlyph
}
