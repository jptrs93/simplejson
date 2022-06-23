package sjson

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// tokens for parsing
const (
	BEGIN_OBJ = iota
	BEGIN_ARR = iota
	COLON     = iota
	SEP       = iota
	END_OBJ   = iota
	END_ARR   = iota
	NULL_T    = iota
	BOOL_T    = iota
	STRING_T  = iota
	NUMBER_T  = iota
	EOF       = iota
)

// Errors
var (
	BadNullTokenError      = errors.New("Bad NULL token")
	BadBoolTokenError      = errors.New("Bad BOOL token")
	BadNumberTokenError    = errors.New("Bad NUMBER token")
	UnclosedJsonError      = errors.New("Unclosed json")
	InvalidTokenStartError = errors.New("Invalid start of token")
	ControlCharacterError  = errors.New("Illegal control character in json string")
	EscapeSequenceError    = errors.New("Invalid escape sequence")
	RuneDecodingError      = errors.New("Error decoding rune")
	KeyPathError           = errors.New("No element at path")
	BoolValueError         = errors.New("Element cannot be converted to bool")
	NumberValueError       = errors.New("Element cannot be converted to number")
	StringValueError       = errors.New("Element cannot be converted to string")

	RuneReadingError = func(p int, e error) error { return fmt.Errorf("Error reading char %v: %w", p, e) }
	ParseError       = func(p int, e error) error { return fmt.Errorf("Char %v: %w", p, e) }
)

// Parse json from an io.RuneScanner. Note that the runes are copied into a backing byte slice as they are read.
func Parse(rs io.RuneScanner) (*Json, error) {
	scanner := StoringRuneScanner{scanner: rs, bytes: make([]byte, 0)}
	return parse(&scanner.bytes, &scanner)
}

// Parse json directly from an existing (UTF8) bytes slice. The input byte slice will be used directly as the backing
// slice for the json. No copy will be made.
func ParseUTF8(b []byte) (*Json, error) {
	return parse(&b, &UTF8RuneScanner{bytes: b})
}

func parse(content *[]byte, reader PositionalRuneScanner) (*Json, error) {
	parentStack, l := []*Json{nil}, 0
	prevToken, key, keyAva := -1, "", false
	var root *Json
	var j *Json
	for prevToken != EOF {
		start, end, t, jtype, err := nextToken(reader)
		if err != nil {
			return nil, err
		}
		if err := validateTokenOrder(t, prevToken, parentStack); err != nil {
			return nil, err
		}
		// if we reach here the json must be valid
		switch {
		case t == STRING_T && l > 0 && !keyAva && parentStack[l].jtype == OBJECT:
			// fmt.Println("Content length ", len(*content))
			key, keyAva = string((*content)[start+1:end-1]), true
		case t == EOF:
		case t == COLON || t == SEP:
		case t == END_ARR || t == END_OBJ:
			j = parentStack[l]
			j.end = end
			parentStack, l = parentStack[:l], l-1
		default:
			// todo: only initialise objectItems map and arrayItems slice for OBJECT and ARRAY json types
			j = &Json{parent: parentStack[l], bytes: content, start: start, end: end, jtype: jtype, objectItems: make(map[string]*Json, 0), arrayItems: make([]*Json, 0)}
			if l > 0 && keyAva {
				parentStack[l].objectItems[key], keyAva = j, false
			} else if l > 0 {
				parentStack[l].arrayItems = append(parentStack[l].arrayItems, j)
			}
			if t == BEGIN_ARR || t == BEGIN_OBJ {
				parentStack, l = append(parentStack, j), l+1
			}
		}
		if l == 0 {
			root = j
		}
		prevToken = t
	}
	return root, nil
}

func moveToNullTokenEnd(r PositionalRuneScanner) error {
	for _, expected := range [3]rune{'u', 'l', 'l'} {
		rp := r.RunePos()
		if c, _, err := r.ReadRune(); err != nil {
			return RuneReadingError(rp, err)
		} else if c != expected {
			return ParseError(rp, BadNullTokenError)
		}
	}
	return nil
}

func moveToBoolTokenEnd(c rune, r PositionalRuneScanner) error {
	if c == 't' {
		for _, expected := range [3]rune{'r', 'u', 'e'} {
			rp := r.RunePos()
			if c, _, err := r.ReadRune(); err != nil {
				return RuneReadingError(rp, err)
			} else if c != expected {
				return ParseError(rp, BadBoolTokenError)
			}
		}
	} else {
		for _, expected := range [4]rune{'a', 'l', 's', 'e'} {
			rp := r.RunePos()
			if c, _, err := r.ReadRune(); err != nil {
				return RuneReadingError(rp, err)
			} else if c != expected {
				return ParseError(rp, BadBoolTokenError)
			}
		}
	}
	return nil
}

func moveToStringTokenEnd(r PositionalRuneScanner) error {
	var err error
	for c, prev := '"', '\\'; !(c == '"' && !(prev == '\\')); {
		prev = c
		if c, _, err = r.ReadRune(); err != nil {
			return RuneReadingError(r.RunePos(), err)
		}
	}
	return nil
}

func moveToNumberTokenEnd(prev rune, r PositionalRuneScanner) error {
	dotCount := 0
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if !unicode.IsDigit(c) {
			if c == '+' || c == '-' {
				return ParseError(r.RunePos(), BadNumberTokenError)
			} else if c == '.' {
				if dotCount > 0 || !unicode.IsDigit(prev) {
					return ParseError(r.RunePos(), BadNumberTokenError)
				}
				dotCount++
			} else if !unicode.IsDigit(prev) {
				return ParseError(r.RunePos(), BadNumberTokenError)
			} else {
				break
			}
		}
		prev = c
	}
	if !unicode.IsSpace(prev) {
		r.UnreadRune()
	}
	return nil
}

// returns (byte start position, byte end poistion, token, json type, error)
func nextToken(r PositionalRuneScanner) (int, int, int, int, error) {
	var err error
	start, c := 0, ' '
	// skip white space
	for unicode.IsSpace(c) {
		start = r.BytePos()
		c, _, err = r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return start, r.BytePos(), EOF, 0, nil
			} else {
				return start, r.BytePos(), EOF, 0, RuneReadingError(r.RunePos(), err)
			}
		}
	}
	token, jtype := -1, -1
	switch {
	case c == '{':
		token, jtype = BEGIN_OBJ, OBJECT
	case c == '}':
		token, jtype = END_OBJ, OBJECT
	case c == '[':
		token, jtype = BEGIN_ARR, ARRAY
	case c == ']':
		token, jtype = END_ARR, ARRAY
	case c == ':':
		token, jtype = COLON, -1
	case c == ',':
		token, jtype = SEP, -1
	case c == 'n':
		err = moveToNullTokenEnd(r)
		token, jtype = NULL_T, NULL
	case c == 't' || c == 'f':
		err = moveToBoolTokenEnd(c, r)
		token, jtype = BOOL_T, BOOL
	case c == '"':
		err = moveToStringTokenEnd(r)
		token, jtype = STRING_T, STRING
	case strings.ContainsRune("+-0123456789", c):
		err = moveToNumberTokenEnd(c, r)
		token, jtype = NUMBER_T, NUMBER
	default:
		err = fmt.Errorf("Char %v (%c): %w", start, c, InvalidTokenStartError)
	}
	return start, r.BytePos(), token, jtype, err
}

func validateTokenOrder(token int, prevToken int, parentStack []*Json) error {
	// we can validate by only considering consecutive tokens and the parent token
	l := len(parentStack) - 1
	if token == EOF && l != 0 {
		return UnclosedJsonError
	}
	return nil
}
