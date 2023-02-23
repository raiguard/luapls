package lua

// scanner lexes Lua source code into a series of tokens.
// TODO: Create multiple lexers for different Lua flavors.
type scanner struct {
	src []byte

	ch   byte
	col  int
	line int
	pos  int
}

// Initializes the scanner.
func (s *scanner) init(src []byte) {
	s.src = src
	s.ch = src[0]
}

// Scans and returns the next token in the source.
func (s *scanner) scan() (Token, string) {
	s.nextWhitespace()

	col, line, startPos := s.col, s.line, s.pos
	kind := INVALID

	switch ch := s.ch; {
	case isLetter(ch):
		kind = IDENT
		s.nextIdentifier()
	case isDigit(ch):
		kind = NUMBER
		s.nextNumber()
	case isEOF(ch):
		kind = EOF
	case ch == '\'' || ch == '"':
		kind = STRING
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
				kind = EQL
				s.next()
			default:
				kind = ASSIGN
			}
		case '<':
			switch ch2 {
			case '=':
				kind = LEQ
				s.next()
			default:
				kind = LSS
			}
		case '>':
			switch ch2 {
			case '=':
				kind = GEQ
				s.next()
			default:
				kind = GTR
			}
		case '~':
			if ch2 == '=' {
				kind = NEQ
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
						kind = LABEL
						s.next()
						break
					}
				}
				// TODO: Error
			} else {
				kind = COLON
			}
		case '.':
			if ch2 == '.' {
				s.next()
				if s.ch == '.' {
					kind = VARARG
				} else {
					kind = CONCAT
				}
			} else {
				kind = DOT
			}
		case '-':
			if ch2 == '-' {
				kind = COMMENT
				// TODO: Raw string comments
				// Advance to the next newline
				for {
					s.next()
					if s.ch == '\n' || s.ch == eof {
						break
					}
				}
			} else {
				kind = SUB
			}
		default:
			realTok, reserved := tokenKinds[string(ch)]
			if reserved {
				kind = realTok
			} else {
				// TODO: Error?
			}
		}
	}

	raw := string(s.src[startPos:s.pos])
	// Convert identifier to reserved token if it is one
	if kind == IDENT {
		if realTok, reserved := tokenKinds[raw]; reserved {
			kind = realTok
		}
	}

	return Token{
		Kind: kind,
		Pos:  TokenPos{Line: line, Col: col},
	}, raw
}

// HELPER FUNCTIONS

// Advances to the next byte in the source.
func (s *scanner) next() {
	s.pos++
	s.col++
	if s.pos == len(s.src) {
		// TODO: Handle this
		s.ch = eof
	} else {
		s.ch = s.src[s.pos]
	}
}

func (s *scanner) nextIdentifier() {
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

func (s *scanner) nextNumber() {
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

func (s *scanner) nextNumeric(isHexNum bool) {
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

func (s *scanner) nextString() {
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

func (s *scanner) nextWhitespace() {
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
