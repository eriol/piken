package format // import "eriol.xyz/piken/format"

import "eriol.xyz/piken/sql"

type Formatter interface {
	Format(*sql.UnicodeData) (string, error)
}

type baseFormatter struct {
	fields    []string
	showGlyph bool
}

func (df *baseFormatter) SetFields(fields []string) {
	for _, field := range fields {
		df.fields = append(df.fields, field)
	}
}

func (df baseFormatter) Fields() []string {
	return df.fields
}

func (df *baseFormatter) SetShowGlyph(value bool) {
	df.showGlyph = value
}

func (df baseFormatter) ShowGlyph() bool {
	return df.showGlyph
}
