package lua

import protocol "github.com/tliron/glsp/protocol_3_16"

// Scanner lexes Lua source code into a series of Tokens.
// TODO: Create multiple lexers for different Lua flavors.
type Scanner struct {
	src []byte

	ch   byte
	col  int
	line int
	pos  int
}

// Initializes the scanner.
func (s *Scanner) Init(src []byte) {
	s.src = src
	s.ch = src[0]
}

// Scans and returns the next token in the source.
func (s *Scanner) Scan() (protocol.Range, Token, string) {
	s.nextWhitespace()

	startCol, startLine, startPos := s.col, s.line, s.pos
	tok := INVALID

	switch ch := s.ch; {
	case isLetter(ch):
		tok = IDENTIFIER
		s.nextIdentifier()
	case isDigit(ch):
		tok = NUMBER
		s.nextNumber()
	case isEOF(ch):
		tok = EOF
	case ch == '\'' || ch == '"':
		tok = STRING
		s.nextString()
	default:
		// Always advance
		s.next()

		// Compound symbols
		ch2 := s.ch
		switch ch {
		case '=':
			switch ch2 {
			case '=':
				tok = EQL
				s.next()
			default:
				tok = ASSIGN
			}
		case '<':
			switch ch2 {
			case '=':
				tok = LEQ
				s.next()
			default:
				tok = LSS
			}
		case '>':
			switch ch2 {
			case '=':
				tok = GEQ
				s.next()
			default:
				tok = GTR
			}
		case '~':
			if ch2 == '=' {
				tok = NEQ
				s.next()
			} else {
				// TODO: Error
			}
		case ':':
			if ch2 == ':' {
				s.next()
				s.nextIdentifier()
				if s.ch == ':' {
					s.next()
					if s.ch == ':' {
						tok = LABEL
						s.next()
						break
					}
				}
				// TODO: Error
			} else {
				tok = COLON
			}
		case '.':
			if ch2 == '.' {
				s.next()
				if s.ch == '.' {
					tok = VARARG
				} else {
					tok = CONCAT
				}
			} else {
				tok = DOT
			}
		case '-':
			if ch2 == '-' {
				tok = COMMENT
				// TODO: Raw string comments
				// Advance to the next newline
				for {
					s.next()
					if s.ch == '\n' || s.ch == eof {
						break
					}
				}
			} else {
				tok = SUB
			}
		default:
			realTok, reserved := tokens[string(ch)]
			if reserved {
				tok = realTok
			} else {
				// TODO: Error?
			}
		}
	}

	lit := string(s.src[startPos:s.pos])
	// Convert identifier to reserved token if it is one
	if tok == IDENTIFIER {
		if realTok, reserved := tokens[lit]; reserved {
			tok = realTok
		}
	}

	return protocol.Range{
		Start: protocol.Position{Line: uint32(startLine), Character: uint32(startCol)},
		End:   protocol.Position{Line: uint32(s.line), Character: uint32(s.col)},
	}, tok, lit
}

// Advances to the next byte in the source.
func (s *Scanner) next() {
	s.pos++
	s.col++
	if s.pos == len(s.src) {
		// TODO: Handle this
		s.ch = eof
	} else {
		s.ch = s.src[s.pos]
	}
}

func (s *Scanner) nextIdentifier() {
	// First character was already read
	s.next()

	for {
		ch := s.ch
		if !isAlnum(ch) && ch != '_' {
			break
		}
		s.next()
	}
}

func (s *Scanner) nextNumber() {
	isHexNum := false

	// First digit is always a number
	s.next()

	// Hexadecimal numbers
	if s.ch == 'x' || s.ch == 'X' {
		isHexNum = true
		s.next()
	}

	s.nextNumeric(isHexNum)

	// Decimal part
	if s.ch == '.' {
		s.next()
		s.nextNumeric(isHexNum)
	}

	// Exponent part
	if isExponentLiteral(s.ch, isHexNum) {
		s.next()
		// Exponent may have a sign
		if s.ch == '+' || s.ch == '-' {
			s.next()
		}
		s.nextNumeric(isHexNum)
	}
}

func (s *Scanner) nextNumeric(isHexNum bool) {
	for {
		switch ch := s.ch; {
		case isHexNum && !isHex(ch):
			return
		case !isHexNum && !isDigit(ch):
			return
		}
		s.next()
	}
}

func (s *Scanner) nextString() {
	delim := s.ch

	// Collect all bytes until the next occurrence of the delimiter
	for {
		s.next()
		if s.ch == delim {
			// Advance once more to capture the closing delimiter
			s.next()
			break
		} else if s.ch == '\n' {
			// TODO: Error on unterminated string
			break
		}
	}
}

func (s *Scanner) nextWhitespace() {
	for {
		if isSpace(s.ch) {
			if s.ch == '\n' {
				s.col = -1
				s.line++
			}
			s.next()
		} else {
			break
		}
	}
}
