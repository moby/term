package term_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/moby/term"
)

func TestEscapeProxyRead(t *testing.T) {
	type readStep struct {
		size   int
		n      int
		expErr error
		buf    []byte
	}

	tests := []struct {
		name       string
		escapeKeys string
		keys       string
		reads      []readStep
	}{
		{
			name: "no escape keys, single key",
			keys: "a",
			reads: []readStep{
				{size: 1, n: 1, buf: []byte("a")},
			},
		},
		{
			name: "no escape keys, multiple keys",
			keys: "a,b,c",
			reads: []readStep{
				{size: 3, n: 3, buf: []byte("abc")},
			},
		},
		{
			name: "no escape keys, no input",
			reads: []readStep{
				{size: 0, n: 0, expErr: io.EOF},
			},
		},
		{
			name:       "DEL escape key, unmatched input",
			escapeKeys: "DEL",
			keys:       "a,b,c,+",
			reads: []readStep{
				{size: 4, n: 4, buf: []byte("abc+")},
			},
		},
		{
			name:       "DEL escape key, no input",
			escapeKeys: "DEL",
			reads: []readStep{
				{size: 0, n: 0, expErr: io.EOF},
			},
		},
		{
			name:       "ctrl-x,ctrl-@ escape sequence, unmatched DEL",
			escapeKeys: "ctrl-x,ctrl-@",
			keys:       "DEL",
			reads: []readStep{
				{size: 1, n: 1, buf: []byte{'\x7f'}},
			},
		},
		{
			name:       "ctrl-c escape key, complete match",
			escapeKeys: "ctrl-c",
			keys:       "ctrl-c",
			reads: []readStep{
				{size: 1, n: 0, expErr: term.EscapeError{}, buf: []byte{'\x03'}},
			},
		},
		{
			name:       "ctrl-c,ctrl-z escape sequence, match across reads",
			escapeKeys: "ctrl-c,ctrl-z",
			keys:       "ctrl-c,ctrl-z",
			reads: []readStep{
				{size: 1, n: 0, buf: []byte{'\x03'}},
				{size: 1, n: 0, expErr: term.EscapeError{}, buf: []byte{'\x1a'}},
			},
		},
		{
			name:       "ctrl-c,ctrl-z escape sequence, match in one read",
			escapeKeys: "ctrl-c,ctrl-z",
			keys:       "ctrl-c,ctrl-z",
			reads: []readStep{
				{size: 2, n: 0, expErr: term.EscapeError{}, buf: []byte{'\x03', '\x1a'}},
			},
		},
		{
			name:       "ctrl-c,ctrl-z escape sequence, partial match followed by DEL and plus",
			escapeKeys: "ctrl-c,ctrl-z",
			keys:       "ctrl-c,DEL,+",
			reads: []readStep{
				{size: 1, n: 0, buf: []byte{'\x03'}},
				{size: 3, n: 3, buf: []byte{'\x03', '\x7f', '+'}},
			},
		},
		{
			name:       "ctrl-c,ctrl-z escape sequence, partial match followed by DEL",
			escapeKeys: "ctrl-c,ctrl-z",
			keys:       "ctrl-c,DEL",
			reads: []readStep{
				{size: 1, n: 0, buf: []byte{'\x03'}},
				{size: 2, n: 2, buf: []byte{'\x03', '\x7f'}},
			},
		},
		{
			name:       "four-key escape sequence, match across reads",
			escapeKeys: "a,b,c,d",
			keys:       "a,b,c,d",
			reads: []readStep{
				{size: 2, n: 0, buf: []byte("ab")},
				{size: 2, n: 0, expErr: term.EscapeError{}, buf: []byte("cd")},
			},
		},
		{
			name:       "ctrl-p,ctrl-q escape sequence, match after mismatch",
			escapeKeys: "ctrl-p,ctrl-q",
			keys:       "ctrl-p,a,ctrl-p,ctrl-q",
			reads: []readStep{
				{size: 1, n: 0},
				{size: 1, n: 1, buf: []byte{'\x10'}},
				{size: 2, n: 1, buf: []byte{'a', '\x10'}},
				{size: 2, n: 0, expErr: term.EscapeError{}},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var escapeKeys []byte
			var err error
			if tc.escapeKeys != "" {
				escapeKeys, err = term.ToBytes(tc.escapeKeys)
				if err != nil {
					t.Fatalf("parse escape keys: %v", err)
				}
			}
			var keys []byte
			if tc.keys != "" {
				keys, err = term.ToBytes(tc.keys)
				if err != nil {
					t.Fatalf("parse input keys: %v", err)
				}
			}

			reader := term.NewEscapeProxy(bytes.NewReader(keys), escapeKeys)
			for i, step := range tc.reads {
				buf := make([]byte, step.size)
				n, err := reader.Read(buf)
				if !errors.Is(err, step.expErr) {
					t.Errorf("read %d: expected error %v, got %v", i+1, step.expErr, err)
				}
				if n != step.n {
					t.Errorf("read %d: expected %d bytes, got %d", i+1, step.n, n)
				}
				// Compare the entire buffer, including when Read returns n == 0.
				if step.buf != nil && !bytes.Equal(buf, step.buf) {
					t.Errorf("read %d: expected buffer %v, got %v", i+1, step.buf, buf)
				}
			}
		})
	}
}
