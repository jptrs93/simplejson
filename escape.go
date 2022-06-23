package sjson

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

var escapeChars = map[rune]rune{
	'"':  '"',
	'\\': '\\',
	'/':  '/',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'u':  '-'}

func isEscapeChar(x rune) bool {
	_, ok := escapeChars[x]
	return ok
}

func EscapeUTF8(bytes []byte) (string, error) {
	n := len(bytes)
	escaped := make([]rune, 0)
	escapeOpen := false
	var pos int
	var c rune
	var err error
	for pos < n {
		c, pos, err = nextRune(bytes, pos)
		if err != nil {
			return "", err
		}
		if escapeOpen {
			if isEscapeChar(c) {
				if c == 'u' {
					if c, pos, err = hexDigitsToRune(bytes, pos); err != nil {
						return "", err
					} else {
						escaped = append(escaped, c)
						escapeOpen = false
					}
				} else {
					v, _ := escapeChars[c]
					escaped = append(escaped, v)
					escapeOpen = false
				}
			} else {
				return "", fmt.Errorf("%w: \\%v", EscapeSequenceError, string(c))
			}
		} else if c == '\\' {
			escapeOpen = true
		} else if unicode.IsControl(c) {
			return "", fmt.Errorf("%w: %c", ControlCharacterError, c)
		} else {
			escaped = append(escaped, c)
		}
	}
	return string(escaped), nil
}

func hexDigitsToRune(bytes []byte, pos int) (rune, int, error) {
	if pos+3 < len(bytes) {
		hexdigits := string(bytes[pos : pos+4])
		if val, err := strconv.ParseInt(hexdigits, 16, 32); err != nil {
			return 0, pos + 4, fmt.Errorf("%w: \\%v", EscapeSequenceError, "u"+hexdigits)
		} else {
			return rune(val), pos + 4, nil
		}
	} else {
		return 0, pos + 4, fmt.Errorf("%w: \\%v", EscapeSequenceError, string(bytes[pos:]))
	}
}

func nextRune(b []byte, BytePos int) (rune, int, error) {
	r, w := utf8.DecodeRune(b[BytePos:])
	if r == utf8.RuneError {
		return r, BytePos, RuneDecodingError
	}
	return r, BytePos + w, nil
}
