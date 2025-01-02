package test

import (
	"bytes"
	"testing"

	"github.com/moby/term"
)

func TestEscapeProxyRead(t *testing.T) {
	t.Run("no escape keys, keys [a]", func(t *testing.T) {
		escapeKeys, _ := term.ToBytes("")
		keys, _ := term.ToBytes("a")
		reader := term.NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, len(keys))
		nr, err := reader.Read(buf)
		if err != nil {
			t.Error(err)
		}
		if expected := len(keys); nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}
		if !bytes.Equal(buf, keys) {
			t.Errorf("expected: %+v, got: %+v", keys, buf)
		}
	})

	t.Run("no escape keys, keys [a,b,c]", func(t *testing.T) {
		escapeKeys, _ := term.ToBytes("")
		keys, _ := term.ToBytes("a,b,c")
		reader := term.NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, len(keys))
		nr, err := reader.Read(buf)
		if err != nil {
			t.Error(err)
		}
		if expected := len(keys); nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}
		if !bytes.Equal(buf, keys) {
			t.Errorf("expected: %+v, got: %+v", keys, buf)
		}
	})

	t.Run("no escape keys, no keys", func(t *testing.T) {
		escapeKeys, _ := term.ToBytes("")
		keys, _ := term.ToBytes("")
		reader := term.NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, len(keys))
		nr, err := reader.Read(buf)
		if err == nil {
			t.Error("expected an error when there are no keys are to read")
		}
		if expected := 0; len(keys) != expected {
			t.Errorf("expected: %d, got: %d", expected, len(keys))
		}
		if expected := len(keys); nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}
		if expected := len(keys); len(buf) != expected {
			t.Errorf("expected: %d, got: %d", expected, len(buf))
		}
	})

	t.Run("DEL escape key, keys [a,b,c,+]", func(t *testing.T) {
		escapeKeys, _ := term.ToBytes("DEL")
		keys, _ := term.ToBytes("a,b,c,+")
		reader := term.NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, len(keys))
		nr, err := reader.Read(buf)
		if err != nil {
			t.Error(err)
		}
		if expected := len(keys); nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}
		if !bytes.Equal(buf, keys) {
			t.Errorf("expected: %+v, got: %+v", keys, buf)
		}
	})

	t.Run("DEL escape key, no keys", func(t *testing.T) {
		escapeKeys, _ := term.ToBytes("DEL")
		keys, _ := term.ToBytes("")
		reader := term.NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, len(keys))
		nr, err := reader.Read(buf)
		if err == nil {
			t.Error("expected an error when there are no keys are to read")
		}
		if expected := 0; len(keys) != expected {
			t.Errorf("expected: %d, got: %d", expected, len(keys))
		}
		if expected := len(keys); nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}
		if expected := len(keys); len(buf) != expected {
			t.Errorf("expected: %d, got: %d", expected, len(buf))
		}
	})

	t.Run("ctrl-x,ctrl-@ escape key, keys [DEL]", func(t *testing.T) {
		escapeKeys, _ := term.ToBytes("ctrl-x,ctrl-@")
		keys, _ := term.ToBytes("DEL")
		reader := term.NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, len(keys))
		nr, err := reader.Read(buf)
		if err != nil {
			t.Error(err)
		}
		if expected := 1; nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}
		if !bytes.Equal(buf, keys) {
			t.Errorf("expected: %+v, got: %+v", keys, buf)
		}
	})

	t.Run("ctrl-c escape key, keys [ctrl-c]", func(t *testing.T) {
		escapeKeys, _ := term.ToBytes("ctrl-c")
		keys, _ := term.ToBytes("ctrl-c")
		reader := term.NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, len(keys))
		nr, err := reader.Read(buf)
		if expected := "read escape sequence"; err == nil || err.Error() != expected {
			t.Errorf("expected: %v, got: %v", expected, err)
		}
		if expected := 0; nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}
		if !bytes.Equal(buf, keys) {
			t.Errorf("expected: %+v, got: %+v", keys, buf)
		}
	})

	t.Run("ctrl-c,ctrl-z escape key, keys [ctrl-c],[ctrl-z]", func(t *testing.T) {
		escapeKeys, _ := term.ToBytes("ctrl-c,ctrl-z")
		keys, _ := term.ToBytes("ctrl-c,ctrl-z")
		reader := term.NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, 1)
		nr, err := reader.Read(buf)
		if err != nil {
			t.Error(err)
		}
		if expected := 0; nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}
		if !bytes.Equal(buf, keys[0:1]) {
			t.Errorf("expected: %+v, got: %+v", keys, buf)
		}

		nr, err = reader.Read(buf)
		if expected := "read escape sequence"; err == nil || err.Error() != expected {
			t.Errorf("expected: %v, got: %v", expected, err)
		}
		if expected := 0; nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}
		if !bytes.Equal(buf, keys[1:]) {
			t.Errorf("expected: %+v, got: %+v", keys, buf)
		}
	})

	t.Run("ctrl-c,ctrl-z escape key, keys [ctrl-c,ctrl-z]", func(t *testing.T) {
		escapeKeys, _ := term.ToBytes("ctrl-c,ctrl-z")
		keys, _ := term.ToBytes("ctrl-c,ctrl-z")
		reader := term.NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, 2)
		nr, err := reader.Read(buf)
		if expected := "read escape sequence"; err == nil || err.Error() != expected {
			t.Errorf("expected: %v, got: %v", expected, err)
		}
		if expected := 0; nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}
		if !bytes.Equal(buf, keys) {
			t.Errorf("expected: %+v, got: %+v", keys, buf)
		}
	})

	t.Run("ctrl-c,ctrl-z escape key, keys [ctrl-c],[DEL,+]", func(t *testing.T) {
		escapeKeys, _ := term.ToBytes("ctrl-c,ctrl-z")
		keys, _ := term.ToBytes("ctrl-c,DEL,+")
		reader := term.NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, 1)
		nr, err := reader.Read(buf)
		if err != nil {
			t.Error(err)
		}
		if expected := 0; nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}
		if !bytes.Equal(buf, keys[0:1]) {
			t.Errorf("expected: %+v, got: %+v", keys, buf)
		}

		buf = make([]byte, len(keys))
		nr, err = reader.Read(buf)
		if err != nil {
			t.Error(err)
		}
		if expected := len(keys); nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}
		if !bytes.Equal(buf, keys) {
			t.Errorf("expected: %+v, got: %+v", keys, buf)
		}
	})

	t.Run("ctrl-c,ctrl-z escape key, keys [ctrl-c],[DEL]", func(t *testing.T) {
		escapeKeys, _ := term.ToBytes("ctrl-c,ctrl-z")
		keys, _ := term.ToBytes("ctrl-c,DEL")
		reader := term.NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, 1)
		nr, err := reader.Read(buf)
		if err != nil {
			t.Error(err)
		}
		if expected := 0; nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}
		if !bytes.Equal(buf, keys[0:1]) {
			t.Errorf("expected: %+v, got: %+v", keys, buf)
		}

		buf = make([]byte, len(keys))
		nr, err = reader.Read(buf)
		if err != nil {
			t.Error(err)
		}
		if expected := len(keys); nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}
		if !bytes.Equal(buf, keys) {
			t.Errorf("expected: %+v, got: %+v", keys, buf)
		}
	})

	t.Run("a,b,c,d escape key, keys [a,b],[c,d]", func(t *testing.T) {
		escapeKeys, _ := term.ToBytes("a,b,c,d")
		keys, _ := term.ToBytes("a,b,c,d")
		reader := term.NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, 2)
		nr, err := reader.Read(buf)
		if err != nil {
			t.Error(err)
		}
		if expected := 0; nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}
		if !bytes.Equal(buf, keys[0:2]) {
			t.Errorf("expected: %+v, got: %+v", keys, buf)
		}

		buf = make([]byte, 2)
		nr, err = reader.Read(buf)
		if expected := "read escape sequence"; err == nil || err.Error() != expected {
			t.Errorf("expected: %v, got: %v", expected, err)
		}
		if expected := 0; nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}
		if !bytes.Equal(buf, keys[2:4]) {
			t.Errorf("expected: %+v, got: %+v", keys, buf)
		}
	})

	t.Run("ctrl-p,ctrl-q escape key, keys [ctrl-p],[a],[ctrl-p,ctrl-q]", func(t *testing.T) {
		escapeKeys, _ := term.ToBytes("ctrl-p,ctrl-q")
		keys, _ := term.ToBytes("ctrl-p,a,ctrl-p,ctrl-q")
		reader := term.NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, 1)
		nr, err := reader.Read(buf)
		if err != nil {
			t.Error(err)
		}
		if expected := 0; nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}

		buf = make([]byte, 1)
		nr, err = reader.Read(buf)
		if err != nil {
			t.Error(err)
		}
		if expected := 1; nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}
		if !bytes.Equal(buf, keys[:1]) {
			t.Errorf("expected: %+v, got: %+v", keys, buf)
		}

		buf = make([]byte, 2)
		nr, err = reader.Read(buf)
		if err != nil {
			t.Error(err)
		}
		if expected := 1; nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}
		if !bytes.Equal(buf, keys[1:3]) {
			t.Errorf("expected: %+v, got: %+v", keys, buf)
		}

		buf = make([]byte, 2)
		nr, err = reader.Read(buf)
		if expected := "read escape sequence"; err == nil || err.Error() != expected {
			t.Errorf("expected: %v, got: %v", expected, err)
		}
		if expected := 0; nr != expected {
			t.Errorf("expected: %d, got: %d", expected, nr)
		}
	})
}
