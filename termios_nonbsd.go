//go:build !darwin && !freebsd && !netbsd && !openbsd && !js && !windows
// +build !darwin,!freebsd,!netbsd,!openbsd,!js,!windows

package term

import (
	"golang.org/x/sys/unix"
)

const (
	getTermios = unix.TCGETS
	setTermios = unix.TCSETS
)
