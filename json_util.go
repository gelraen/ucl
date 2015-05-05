package ucl

import (
	"fmt"
	"strings"
)

/* TODO(imax): code below does not take into account escaping of
 * non-BMP characters as specified by RFC 4627.
 *
 *   To escape an extended character that is not in the Basic Multilingual
 *   Plane, the character is represented as a twelve-character sequence,
 *   encoding the UTF-16 surrogate pair.  So, for example, a string
 *   containing only the G clef character (U+1D11E) may be represented as
 *   "\uD834\uDD1E".
 *
 * Such surrogate pairs will be unescaped as 2 adjacent UTF-8 sequences.
 */

func json_escape(s string) string {
	r := ""
	for _, c := range s {
		switch c {
		case '"':
			r += `\"`
		case '\\':
			r += `\\`
		case '/':
			r += `\/`
		case '\b':
			r += `\b`
		case '\f':
			r += `\f`
		case '\n':
			r += `\n`
		case '\r':
			r += `\r`
		case '\t':
			r += `\t`
		default:
			r += string(c)
		}
	}
	return strings.Replace(s, "\"", "\\\"", -1)
}

func isHexDigit(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func hexValue(c rune) rune {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	}
	return '\uFFFD'
}

func json_unescape(s string) (string, error) {
	r := ""
	start := 0
	const (
		RecordStart    = iota // Beginning of the string or right after the escape sequence.
		Regular               // After a regular unescaped character.
		AfterBackslash        // After a backslash.
		AfterU                // After \u, a separate counter is used to eat exactly 4 next characters.
	)
	state := RecordStart
	var ucount uint
	var u rune
	// Iteration over a string interpretes it as UTF-8 and produces Unicode
	// runes. i is the index of the first byte of the rune, c is the rune.
	for i, c := range s {
		switch state {
		case RecordStart:
			switch c {
			case '\\':
				state = AfterBackslash
			default:
				start = i
				state = Regular
			}
		case Regular:
			switch c {
			case '\\':
				r += s[start:i]
				state = AfterBackslash
			}
		case AfterBackslash:
			switch c {
			case '"':
				r += "\""
				state = RecordStart
			case '\\':
				r += "\\"
				state = RecordStart
			case '/':
				r += "/"
				state = RecordStart
			case 'b':
				r += "\b"
				state = RecordStart
			case 'f':
				r += "\f"
				state = RecordStart
			case 'n':
				r += "\n"
				state = RecordStart
			case 'r':
				r += "\r"
				state = RecordStart
			case 't':
				r += "\t"
				state = RecordStart
			case 'u':
				ucount = 0
				u = 0
				state = AfterU
			default:
				return "", fmt.Errorf("invalid escape sequence %q at %d", c, i)
			}
		case AfterU:
			if !isHexDigit(c) {
				return "", fmt.Errorf("invalid hex digit %q at %d", c, i)
			}
			v := hexValue(c)
			u |= v << (4 * (3 - ucount))
			ucount++
			if ucount == 4 {
				r += string(u)
				state = RecordStart
			}
		}
	}
	if state == Regular {
		r += s[start:len(s)]
	}
	if state != Regular && state != RecordStart {
		return "", fmt.Errorf("incomplete escape sequence at the end of the string")
	}
	return r, nil
}
