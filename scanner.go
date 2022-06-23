package sjson

import (
	"io"
	"unicode/utf8"
)

type PositionalRuneScanner interface {
	io.RuneScanner
	BytePos() int // the current byte position of the the underlying data
	RunePos() int // the current rune number
}

// wraps a RuneScanner so that the runes are stored into a slice of bytes as they are read
type StoringRuneScanner struct {
	scanner io.RuneScanner
	bytes   []byte
	runePos int
	unreads int
}

func (r *StoringRuneScanner) ReadRune() (rune, int, error) {
	if c, w, err := r.scanner.ReadRune(); err != nil {
		return c, w, err
	} else {
		r.bytes = utf8.AppendRune(r.bytes, c)
		r.runePos += 1
		return c, w, err
	}
}

func (r *StoringRuneScanner) UnreadRune() error {
	if len(r.bytes) == 0 {
		return io.EOF
	}
	_, w := utf8.DecodeLastRune(r.bytes)
	r.bytes = r.bytes[:len(r.bytes)-w]
	r.runePos -= 1
	return r.scanner.UnreadRune()
}

func (r *StoringRuneScanner) BytePos() int {
	return len(r.bytes)
}

func (r *StoringRuneScanner) RunePos() int {
	return r.runePos
}

type UTF8RuneScanner struct {
	bytes   []byte
	bytePos int
	runePos int
}

func (r *UTF8RuneScanner) ReadRune() (rune, int, error) {
	if r.bytePos >= len(r.bytes) {
		return utf8.RuneError, 0, io.EOF
	}
	c, w := utf8.DecodeRune(r.bytes[r.bytePos:])
	if c == utf8.RuneError {
		return c, w, RuneDecodingError
	}
	r.bytePos = r.bytePos + w
	r.runePos = r.runePos + 1
	return c, w, nil
}

func (r *UTF8RuneScanner) UnreadRune() error {
	if r.bytePos == 0 {
		return io.EOF
	}
	_, w := utf8.DecodeLastRune(r.bytes)
	r.bytePos -= w
	r.runePos -= 1
	return nil
}

func (r *UTF8RuneScanner) BytePos() int {
	return r.bytePos
}

func (r *UTF8RuneScanner) RunePos() int {
	return r.runePos
}

func readString(r io.RuneScanner, l int) (string, error) {
	runes := []rune{}
	for i := 0; i < l; i++ {
		if r, _, err := r.ReadRune(); err != nil {
			return "", err
		} else {
			runes = append(runes, r)
		}
	}
	return string(runes), nil
}
