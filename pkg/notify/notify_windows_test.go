package notify

import (
	"image/png"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIcon(t *testing.T) {
	tmp, err := ioutil.TempDir("", "gopass-ico")
	assert.NoError(t, err)
	icon, err := iconURI(tmp)
	assert.NoError(t, err)
	fn := strings.TrimPrefix(icon, "file:///")
	t.Logf("icon path: %s", fn)
	_ = os.Remove(fn)
	icon, err = iconURI(tmp)
	assert.NoError(t, err)
	fh, err := os.Open(fn)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, fh.Close())
	}()
	require.NotNil(t, fh)
	_, err = png.Decode(fh)
	assert.NoError(t, err)
}
