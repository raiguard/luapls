// Utility functions for interpreting ASCII text. Lua does not support
// unicode, so we use byte instead of rune.

package lua

// Represents the end of a file. ASCII code 28 is the file separator.
const eof byte = 28

func isAlnum(c byte) bool {
    return isDigit(c) || isLetter(c)
}

func isExponentLiteral(c byte, isHexNum bool) bool {
    if isHexNum {
        return toLower(c) == 'p'
    } else {
        return toLower(c) == 'e'
    }
}

func isDigit(c byte) bool {
    return '0' <= c && c <= '9'
}

func isEOF(c byte) bool {
    return c == eof
}

func isHex(c byte) bool {
    return isDigit(c) || ('a' <= toLower(c) && toLower(c) <= 'f')
}

func isLetter(c byte) bool {
    return 'a' <= toLower(c)  && toLower(c) <= 'z'
}

func isSpace(c byte) bool {
    // TODO: Is this all of them?
    return c == '\n' || c == ' ' || c == '\t' || c == '\r'
}

func toLower(c byte) byte {
    return ('a' - 'A') | c
}
