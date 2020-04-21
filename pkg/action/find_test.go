package action

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gopasspw/gopass/pkg/ctxutil"
	"github.com/gopasspw/gopass/pkg/out"
	"github.com/gopasspw/gopass/pkg/store/secret"
	"github.com/gopasspw/gopass/tests/gptest"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestFind(t *testing.T) {
	u := gptest.NewUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	ctx = ctxutil.WithTerminal(ctx, false)
	act, err := newMock(ctx, u)
	require.NoError(t, err)
	require.NotNil(t, act)

	buf := &bytes.Buffer{}
	out.Stdout = buf
	stdout = buf
	defer func() {
		stdout = os.Stdout
		out.Stdout = os.Stdout
	}()
	color.NoColor = true

	app := cli.NewApp()

	actName := "action.test"

	if runtime.GOOS == "windows" {
		actName = "action.test.exe"
	}

	// find
	c := cli.NewContext(app, flag.NewFlagSet("default", flag.ContinueOnError), nil)
	if err := act.Find(ctx, c); err == nil || err.Error() != fmt.Sprintf("Usage: %s find <NEEDLE>", actName) {
		t.Errorf("Should fail: %s", err)
	}

	// find fo
	fs := flag.NewFlagSet("default", flag.ContinueOnError)
	assert.NoError(t, fs.Parse([]string{"fo"}))
	c = cli.NewContext(app, fs, nil)

	assert.NoError(t, act.Find(ctx, c))
	assert.Equal(t, "Found exact match in 'foo'\nsecret", strings.TrimSpace(buf.String()))
	buf.Reset()

	// find yo
	fs = flag.NewFlagSet("default", flag.ContinueOnError)
	assert.NoError(t, fs.Parse([]string{"yo"}))
	c = cli.NewContext(app, fs, nil)

	assert.Error(t, act.Find(ctx, c))
	buf.Reset()

	// add some secrets
	assert.NoError(t, act.Store.Set(ctx, filepath.Join("bar", "baz"), secret.New("foo", "bar")))
	assert.NoError(t, act.Store.Set(ctx, filepath.Join("bar", "zab"), secret.New("foo", "bar")))
	buf.Reset()

	// find bar
	fs = flag.NewFlagSet("default", flag.ContinueOnError)
	assert.NoError(t, fs.Parse([]string{"bar"}))
	c = cli.NewContext(app, fs, nil)

	assert.NoError(t, act.Find(ctx, c))
	assert.Equal(t, filepath.Join("bar", "baz")+"\n"+filepath.Join("bar", "zab"), strings.TrimSpace(buf.String()))
	buf.Reset()
}
