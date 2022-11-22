package lua

// Scanner lexes Lua source code into a series of Tokens.
// TODO: Create multiple lexers for different Lua flavors.
type Scanner struct {
	src []byte

	ch   byte
	line int
	pos  int
}

// Initializes the scanner.
func (s *Scanner) Init(src []byte) {
	s.src = src
}

// Scans and returns the next token in the source.
func (s *Scanner) Scan() (pos int, tok Token, lit string) {
	s.advanceSpace()

	pos = s.pos

	switch ch := s.ch; {
	case isLetter(ch):
		lit = s.scanIdentifier()
		// TODO: Keywords are all > 1 letter
		if realTok, isKeyword := tokens[lit]; isKeyword {
			tok = realTok
		} else {
			tok = IDENTIFIER
		}
	case isDigit(ch):
		tok = NUMBER
		lit = s.scanNumber()
	case isEOF(ch):
    	tok = EOF
    case ch == '\'' || ch == '"':
        tok = STRING
        lit = s.scanString()
	default:
    	// Always advance
    	s.next()
	}

	return
}

// Advances to the next byte in the source.
func (s *Scanner) next() {
	s.pos++
	if s.pos == len(s.src) {
    	// TODO: Handle this
		s.ch = eof
	} else {
		s.ch = s.src[s.pos]
	}
}

// Advances to the next non-numeric character.
func (s *Scanner) advanceNumber(isHexNum bool) {
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

// Advances to the next non-whitespace character.
// TODO: Store leading and trailing space in each Node?
func (s *Scanner) advanceSpace() {
	for {
		if isSpace(s.ch) {
			if s.ch == '\n' {
				s.line++
			}
			s.next()
		} else {
			break
		}
	}
}

func (s *Scanner) scanIdentifier() string {
	start := s.pos

	// First character was already read
	s.next()

	for {
		ch := s.ch
		if !isAlnum(ch) && ch != '_' {
			break
		}
		s.next()
	}

	return string(s.src[start:s.pos])
}

func (s *Scanner) scanNumber() string {
	isHexNum := false
	start := s.pos

	// First digit is always a number
	s.next()

	// Hexadecimal numbers
	if s.ch == 'x' || s.ch == 'X' {
		isHexNum = true
		s.next()
	}

	s.advanceNumber(isHexNum)

	// Decimal part
	if s.ch == '.' {
		s.next()
		s.advanceNumber(isHexNum)
	}

	// Exponent part
	if isExponentLiteral(s.ch, isHexNum) {
		s.next()
		// Exponent may have a sign
		if s.ch == '+' || s.ch == '-' {
			s.next()
		}
		s.advanceNumber(isHexNum)
	}

	return string(s.src[start:s.pos])
}

func (s *Scanner) scanString() string {
    start := s.pos
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

	return string(s.src[start:s.pos])
}
