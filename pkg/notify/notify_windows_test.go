package notify

import (
	"image/png"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIcon(t *testing.T) {
	fn := strings.TrimPrefix(iconURI(), "file:///")
	_ = os.Remove(fn)
	_ = iconURI()
	fh, err := os.Open(fn)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, fh.Close())
	}()
	require.NotNil(t, fh)
	_, err = png.Decode(fh)
	assert.NoError(t, err)
}
