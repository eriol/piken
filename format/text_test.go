package format // import "eriol.xyz/piken/format"
import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodePointToGlyph(t *testing.T) {
	glyph, err := CodePointToGlyph("1F602")
	assert.Nil(t, err)
	assert.Equal(t, glyph, "😂")

	glyph, err = CodePointToGlyph("1000000000")
	assert.Equal(t, glyph, "")
	assert.Error(t, err)
}
