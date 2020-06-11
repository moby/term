// +build !windows

package term

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/creack/pty"
	"github.com/google/go-cmp/cmp"
	"gotest.tools/v3/assert"
)

func newTtyForTest(t *testing.T) (*os.File, *os.File) {
	t.Helper()
	pty, tty, err := pty.Open()
	if err != nil {
		t.Fatalf("error creating pty: %v", err)
	}

	return pty, tty
}

func newTempFile() (*os.File, error) {
	return ioutil.TempFile(os.TempDir(), "temp")
}

func TestGetWinsize(t *testing.T) {
	pty, tty := newTtyForTest(t)
	defer pty.Close()
	defer tty.Close()
	winSize, err := GetWinsize(tty.Fd())
	assert.NilError(t, err)
	assert.Assert(t, winSize != nil)

	newSize := Winsize{Width: 200, Height: 200, x: winSize.x, y: winSize.y}
	err = SetWinsize(tty.Fd(), &newSize)
	assert.NilError(t, err)
	winSize, err = GetWinsize(tty.Fd())
	assert.NilError(t, err)
	assert.DeepEqual(t, *winSize, newSize, cmpWinsize)
}

var cmpWinsize = cmp.AllowUnexported(Winsize{})

func TestSetWinsize(t *testing.T) {
	pty, tty := newTtyForTest(t)
	defer pty.Close()
	defer tty.Close()
	winSize, err := GetWinsize(tty.Fd())
	assert.NilError(t, err)
	assert.Assert(t, winSize != nil)
	newSize := Winsize{Width: 200, Height: 200, x: winSize.x, y: winSize.y}
	err = SetWinsize(tty.Fd(), &newSize)
	assert.NilError(t, err)
	winSize, err = GetWinsize(tty.Fd())
	assert.NilError(t, err)
	assert.DeepEqual(t, *winSize, newSize, cmpWinsize)
}

func TestGetFdInfo(t *testing.T) {
	pty, tty := newTtyForTest(t)
	defer pty.Close()
	defer tty.Close()
	inFd, isTerminal := GetFdInfo(tty)
	assert.Equal(t, inFd, tty.Fd())
	assert.Equal(t, isTerminal, true)
	tmpFile, err := newTempFile()
	assert.NilError(t, err)
	defer tmpFile.Close()
	inFd, isTerminal = GetFdInfo(tmpFile)
	assert.Equal(t, inFd, tmpFile.Fd())
	assert.Equal(t, isTerminal, false)
}

func TestIsTerminal(t *testing.T) {
	pty, tty := newTtyForTest(t)
	defer pty.Close()
	defer tty.Close()
	isTerminal := IsTerminal(tty.Fd())
	assert.Equal(t, isTerminal, true)
	tmpFile, err := newTempFile()
	assert.NilError(t, err)
	defer tmpFile.Close()
	isTerminal = IsTerminal(tmpFile.Fd())
	assert.Equal(t, isTerminal, false)
}

func TestSaveState(t *testing.T) {
	pty, tty := newTtyForTest(t)
	defer pty.Close()
	defer tty.Close()
	state, err := SaveState(tty.Fd())
	assert.NilError(t, err)
	assert.Assert(t, state != nil)
	pty, tty = newTtyForTest(t)
	defer pty.Close()
	defer tty.Close()
	err = RestoreTerminal(tty.Fd(), state)
	assert.NilError(t, err)
}

func TestDisableEcho(t *testing.T) {
	pty, tty := newTtyForTest(t)
	defer pty.Close()
	defer tty.Close()
	state, err := SetRawTerminal(tty.Fd())
	defer RestoreTerminal(tty.Fd(), state)
	assert.NilError(t, err)
	assert.Assert(t, state != nil)
	err = DisableEcho(tty.Fd(), state)
	assert.NilError(t, err)
}
