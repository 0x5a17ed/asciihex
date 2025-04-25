package asciihex

import (
	"fmt"
	"strings"
)

// Encode converts binary data into a readable string representation using the following rules:
//  - Printable ASCII characters (0x20-0x7E) except '~' and '^' are kept as-is
//  - Control characters (0x00-0x1F) are encoded as '^' followed by their corresponding control picture character
//  - The DEL character (0x7F) is encoded as "^?"
//  - The '~' character is encoded as "~~"
//  - The '^' character is encoded as "~^"
//  - All other bytes (>= 0x80) are encoded as "~" followed by their two-digit hexadecimal representation
//
// This encoding ensures that binary data can be safely represented as printable ASCII text,
// similar to quoted-printable or ASCII armor encodings.
func Encode(data []byte) string {
	var builder strings.Builder
	for _, b := range data {
		switch {
		case b == '~':
			builder.WriteString("~~")
		case b == '^':
			builder.WriteString("~^")
		case b >= 0x20 && b <= 0x7E:
			builder.WriteByte(b)
		case b < 0x20:
			builder.WriteString("^" + string(b+'@'))
		case b == 0x7F:
			builder.WriteString("^?")
		default:
			builder.WriteString(fmt.Sprintf("~%02X", b))
		}
	}
	return builder.String()
}
