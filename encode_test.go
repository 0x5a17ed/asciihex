package asciihex

import (
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		inp  []byte
		want string
	}{
		// plain ASCII
		{[]byte("Hello"), "Hello"},
		{[]byte("Hi there!"), "Hi there!"},

		// control characters
		{[]byte{0x00}, "^@"},
		{[]byte{0x01}, "^A"},
		{[]byte{0x1E}, "^^"},
		{[]byte{0x1F}, "^_"},
		{[]byte{0x7F}, "^?"},

		// non-ASCII bytes
		{[]byte{0x80}, "~80"},
		{[]byte{0xFF}, "~FF"},
		{[]byte{0x8F, 0x00}, "~8F^@"},

		// escaped caret and tilde
		{[]byte{'~'}, "~~"},
		{[]byte{'^'}, "~^"},
		{[]byte{'~', '^'}, "~~~^"},

		// mixed content
		{[]byte("Hi\n^~"), "Hi^J~^~~"},
		{[]byte{0x01, '^', 0xFF}, "^A~^~FF"},
		{[]byte("~^@~"), "~~~^@~~"},
	}
	for _, tt := range tests {
		got := Encode(tt.inp)
		if got != tt.want {
			t.Errorf("Encode(%v) = %q, want %q", tt.inp, got, tt.want)
		}
	}
}
