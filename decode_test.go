package asciihex

import (
	"bytes"
	"testing"
)

func TestDecode(t *testing.T) {
	tt := []struct {
		input    string
		expected []byte
		wantErr  bool
	}{
		{"Hello", []byte("Hello"), false},
		{"Hi there!", []byte("Hi there!"), false},

		// Literal escapes
		{"~~", []byte{'~'}, false},
		{"~^", []byte{'^'}, false},

		// Control characters
		{"^@", []byte{0x00}, false},
		{"^A", []byte{0x01}, false},
		{"^Z", []byte{0x1A}, false},
		{"^?", []byte{0x7F}, false},

		// Hex escapes
		{"~8F", []byte{0x8F}, false},
		{"~00", []byte{0x00}, false},
		{"~FF", []byte{0xFF}, false},

		// Mixed content
		{"Hey^Jyou~^~~", []byte{'H', 'e', 'y', 0x0A, 'y', 'o', 'u', '^', '~'}, false},
		{"A^B~7E~^C~~Z", []byte{'A', 0x02, 0x7E, '^', 'C', '~', 'Z'}, false},

		// Invalid sequences
		{"~", nil, true},
		{"~F", nil, true},  // not enough hex digits / early end of input.
		{"~XY", nil, true}, // not valid hex.
		{"^", nil, true},
		{"^!", nil, true}, // not a valid control sequence
		{"~@", nil, true}, // not ~^ or ~XX or ~~
	}

	for _, tc := range tt {
		t.Run(tc.input, func(t *testing.T) {
			switch got, err := Decode(tc.input); {
			case tc.wantErr && err == nil:
				t.Errorf("expected error, got nil")
			case !tc.wantErr && err != nil:
				t.Errorf("unexpected error: %v", err)
			case !bytes.Equal(got, tc.expected):
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
