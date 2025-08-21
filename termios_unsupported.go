//go:build js || wasip1

package term

import (
	"fmt"
	"runtime"
)

func makeRaw(fd uintptr) (*State, error) {
	return nil, fmt.Errorf("terminal: MakeRaw not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
}
