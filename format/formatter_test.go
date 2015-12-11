package format // import "eriol.xyz/piken/format"
import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseFormatterFields(t *testing.T) {
	bf := baseFormatter{}

	fields := []string{"f1", "f2"}
	bf.SetFields(fields)
	assert.Equal(t, bf.Fields(), fields)
}

func TestBaseFormatterShowGlyph(t *testing.T) {
	bf := baseFormatter{}

	assert.Equal(t, bf.ShowGlyph(), false)
	bf.SetShowGlyph(true)
	assert.Equal(t, bf.ShowGlyph(), true)
}
