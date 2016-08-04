package emenv

import (
	"io"
	"strconv"
	"strings"
	"unicode"
)

func NewTokenizer(input []byte) Tokenizer {
	r := strings.NewReader(string(input))
	return Tokenizer{r: r}
}

func (tk *Tokenizer) SkipComments(top rune) (rune, error) {
	newline := false
	if top != ';' {
		return top, nil
	}
	for {
		r, _, err := tk.r.ReadRune()
		if err != nil {
			return 0, err
		}
		if r == ';' {
			newline = false
		}
		if newline && !unicode.IsSpace(r) {
			return r, nil
		}
		if r == '\n' {
			newline = true
		}
	}
	return 0, UnreachableError
}

func (tk *Tokenizer) NextRune() (rune, error) {
	for {
		r, _, err := tk.r.ReadRune()
		if err != nil {
			return 0, err
		}
		if r, err = tk.SkipComments(r); err != nil {
			return 0, err
		}
		if !unicode.IsSpace(r) {
			return r, nil
		}
	}
	return 0, UnreachableError
}

func SymbolTokenFromString(s string) (Token, error) {
	if s == "nil" {
		return Token{Type: NilToken}, nil
	}

	i, err := strconv.Atoi(s)
	if err == nil {
		return Token{Type: NumberToken, String: s, Number: i}, nil
	}

	if strings.HasPrefix(s, ":") {
		s = strings.TrimPrefix(s, ":")
		return Token{Type: KeywordToken, String: s}, nil
	}

	return Token{Type: SymbolToken, String: s}, nil
}

func (tk *Tokenizer) LastRune(r rune) (bool, error) {
	if r == '[' || r == ']' ||
		r == '(' || r == ')' ||
		r == '.' || unicode.IsSpace(r) {

		err := tk.r.UnreadRune()
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (tk *Tokenizer) NextToken() (Token, error) {
	buf := make([]rune, 0)

	r, err := tk.NextRune()
	if err != nil {
		if err == io.EOF {
			return Token{Type: EOFToken}, nil
		}
		return Token{}, err
	}

	switch {
	case r == '[':
		return Token{Type: OpenVectorToken}, nil
	case r == ']':
		return Token{Type: CloseVectorToken}, nil
	case r == '(':
		return Token{Type: OpenParToken}, nil
	case r == ')':
		return Token{Type: CloseParToken}, nil
	case r == '\'':
		return Token{Type: QuoteToken}, nil
	case r == '.':
		return Token{Type: DotToken}, nil
	case r == '"':
		for {
			escaped := false
			subr, _, err := tk.r.ReadRune()
			if err != nil {
				return Token{}, err
			}
			if subr == '\\' {
				subr, _, err = tk.r.ReadRune()
				if err != nil {
					return Token{}, err
				}
				escaped = true
			}
			if subr == '"' && escaped == false {
				return Token{Type: StringToken, String: string(buf)}, nil
			}
			buf = append(buf, subr)
		}
		break
	case r == ':':
		for {
			subr, _, err := tk.r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return Token{Type: KeywordToken, String: string(buf)}, nil
				}
				return Token{}, err
			}

			isLast, err := tk.LastRune(subr)
			if err != nil {
				return Token{}, err
			}

			if isLast {
				return Token{Type: KeywordToken, String: string(buf)}, nil
			}
			buf = append(buf, subr)
		}
		break
	default:
		buf = append(buf, r)
		for {
			subr, _, err := tk.r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return SymbolTokenFromString(string(buf))
				}
				return Token{}, err
			}

			isLast, err := tk.LastRune(subr)
			if err != nil {
				return Token{}, err
			}
			if isLast {
				return SymbolTokenFromString(string(buf))
			}
			buf = append(buf, subr)
		}
		break
	}
	return Token{}, UnreachableError
}

func (tk *Tokenizer) Tokenize() ([]Token, error) {
	tokens := make([]Token, 0)

	for {
		token, err := tk.NextToken()
		if err != nil {
			return nil, err
		}
		if token.Type == EOFToken {
			return tokens, nil
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func ParseTokens(input []byte) ([]Token, error) {
	tk := NewTokenizer(input)
	return tk.Tokenize()
}
