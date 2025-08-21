//go:build !darwin && !freebsd && !js && !netbsd && !openbsd && !wasip1 && !windows

package term

import (
	"golang.org/x/sys/unix"
)

const (
	getTermios = unix.TCGETS
	setTermios = unix.TCSETS
)
