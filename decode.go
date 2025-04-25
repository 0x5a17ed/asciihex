package asciihex

import (
	"encoding/hex"
	"fmt"
	"iter"
	"unicode/utf8"
)

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF

	tokenByte
)

type token struct {
	typ tokenType // The type of this token.
	pos int       // The starting position of this token in the input.
	val string    // The value of this token.
}

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*scanner) stateFn

type scanner struct {
	input string // The string being scanned.

	yield func(token) bool // Yield callback.
	done  bool             // Set to true if yield returns false.
	prev  int              // Previous position (for undo).
	pos   int              // Current position in the input.
}

// next returns the next rune in the input and updates the lexer's position.
// It saves the current position to allow undo.
func (l *scanner) next() rune {
	// Save the current state for undo.
	l.prev = l.pos

	// Check if we are at the end of the input.
	if l.pos >= len(l.input) {
		return eof
	}

	// Advance to the next rune.
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += w
	return r
}

// undo reverts the lexer by one rune using the saved position.
func (l *scanner) undo() {
	l.pos = l.prev
}

// peek returns but does not consume the next rune.
func (l *scanner) peek() rune {
	r := l.next()
	l.undo()
	return r
}

// emit creates a Token from the current input and calls the yield callback.
func (l *scanner) emit(typ tokenType, val string) {
	if l.done || !l.yield(token{typ: typ, pos: l.pos, val: val}) {
		l.done = true
	}
}

// errorf emits an error Token and stops lexing.
func (l *scanner) errorf(format string, args ...any) stateFn {
	msg := fmt.Sprintf(format, args...)
	l.yield(token{typ: tokenError, pos: l.pos, val: msg})
	l.done = true
	return nil
}

// scanCaret is a state function that handles the '^' character.
func scanCaret(s *scanner) stateFn {
	switch ch := s.next(); {
	case ch == eof:
		return s.errorf("unexpected end of input")
	case ch >= '?' && ch <= '_':
		s.emit(tokenByte, string((ch-'@')&0x7F))
		return scanTop
	default:
		return s.errorf("invalid control character '^%c'", ch)
	}
}

// scanTildeHex is a state function that handles the '~' character followed by a hex sequence.
func scanTildeHex(s *scanner) stateFn {
	var buf [2]rune
	for i := 0; i < 2; i++ {
		if buf[i] = s.next(); buf[i] == eof {
			return s.errorf("unexpected end of input")
		}
	}

	b, err := hex.DecodeString(string(buf[:]))
	if err != nil {
		return s.errorf("invalid hex sequence '~%s'", string(buf[:]))
	}

	s.emit(tokenByte, string(b))
	return scanTop
}

// scanTilde is a state function that handles the '~' character.
func scanTilde(s *scanner) stateFn {
	switch ch := s.peek(); ch {
	case eof:
		return s.errorf("unexpected end of input")
	case '~', '^':
		s.next() // Consume the second '~' or '^'
		s.emit(tokenByte, string(ch))
		return scanTop
	default:
		return scanTildeHex(s)
	}
}

// scanTop is a state function that handles the string content.
func scanTop(l *scanner) stateFn {
	switch ch := l.next(); {
	case ch == eof:
		l.emit(tokenEOF, "")
		return nil
	case ch == '^':
		return scanCaret
	case ch == '~':
		return scanTilde
	case ch >= 0x20 && ch <= 0x7E:
		l.emit(tokenByte, string(ch))
		return scanTop
	default:
		return l.errorf("unexpected character")
	}
}

// runPattern is the main function that runs the scanner state machine.
func runPattern(l *scanner) {
	for state := scanTop; state != nil; state = state(l) {
		if l.done {
			break
		}
	}
}

// scan is a generator function that scans the input string and yields tokens.
func scan(s string) iter.Seq[token] {
	return func(yield func(token) bool) {
		l := &scanner{
			input: s,
			yield: yield,
		}
		runPattern(l)
	}
}

// Decode decodes a string into a byte slice using the specified encoding rules.
func Decode(s string) ([]byte, error) {
	var out []byte
	for t := range scan(s) {
		switch t.typ {
		case tokenEOF:
			break
		case tokenError:
			return nil, fmt.Errorf("error at position %d: %s", t.pos, t.val)
		case tokenByte:
			out = append(out, t.val...)
		}
	}

	return out, nil
}
