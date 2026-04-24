// Package parse turns a (potentially multiline) string into a set of tokens.
//
// Parse rules:
//
//  1. Spaces, tabs, or carriage returns can be used to separate tokens
//  2. Single or double quotes can be used to desgniate strings
//  3. A backslash can be used to cancel the meaning of proceeded character.
//     e.g.  'It\'s good' (or just say "It's good")
//
// 4 A # can be used for comments, which last until the end of the line
//
// Implementation is via a finite state machine
package parse

import (
	"errors"
	"fmt"
	"mattwach/rpngo/common/elog"
)

var (
	ErrUnterminatedDouble      = errors.New("unterminated double quote")
	ErrUnterminatedSingleQuote = errors.New("unterminated single quote")
	ErrUnterminatedBrace       = errors.New("unterminatd brace")
)

type State uint8

const (
	WHITESPACE State = iota
	TOKEN
	STRING_DOUBLE
	STRING_SINGLE
	STRING_BRACES
	COMMENT
)

type parseData struct {
	s             State
	t             []rune
	nextIsLiteral bool
	braceDepth    int
}

const defaultStaticRunes = 64

// Use a static instance to avoid heap allocations on every Fields call
var parse parseData = parseData{
	t: make([]rune, defaultStaticRunes),
}

func (p *parseData) init() {
	p.s = WHITESPACE
	if cap(p.t) > defaultStaticRunes {
		elog.Heap("alloc: /parse/parse.go:51: p.t = make([]rune, defaultStaticRunes)")
		p.t = make([]rune, defaultStaticRunes) // object allocated on the heap: escapes at line 51
	}
	p.t = p.t[:0]
}

func (p *parseData) whitespace(c rune) {
	if isWhitespace(c) {
		return
	}
	if c == '\\' {
		p.nextIsLiteral = true
		p.s = TOKEN
		return
	}
	if c == '#' {
		p.s = COMMENT
		return
	}
	if len(p.t) > 0 {
		p.t = p.t[:1]
		p.t[0] = c
	} else {
		p.t = append(p.t, c)
	}
	switch c {
	case '\'':
		p.s = STRING_SINGLE
	case '"':
		p.s = STRING_DOUBLE
	case '{':
		p.s = STRING_BRACES
	default:
		p.s = TOKEN
	}
}

func (p *parseData) token(c rune, fn func(string) error) error {
	if p.nextIsLiteral {
		p.t = append(p.t, c)
		p.nextIsLiteral = false
		return nil
	}
	if c == '\\' {
		p.nextIsLiteral = true
		return nil
	}
	if isWhitespace(c) {
		token := string(p.t)
		p.t = p.t[:0]
		p.s = WHITESPACE

		return fn(token)
	}
	p.t = append(p.t, c)
	return nil
}

func (p *parseData) str(c rune, fn func(string) error) error {
	if p.nextIsLiteral {
		switch c {
		case 'n':
			c = '\n'
		case 't':
			c = '\t'
		}
		p.t = append(p.t, c)
		p.nextIsLiteral = false
		return nil
	}
	if c == '\\' {
		p.nextIsLiteral = true
		return nil
	}
	p.t = append(p.t, c)
	var callFn bool
	switch p.s {
	case STRING_DOUBLE:
		callFn = c == '"'
	case STRING_SINGLE:
		callFn = c == '\''
	case STRING_BRACES:
		if c == '{' {
			p.braceDepth++
		} else if c == '}' {
			if p.braceDepth == 0 {
				callFn = true
			} else {
				p.braceDepth--
			}
		}
	}
	if callFn {
		token := string(p.t)
		p.t = p.t[:0]
		p.s = WHITESPACE
		return fn(token)
	}
	return nil
}

func (p *parseData) comment(c rune) {
	if c == '\n' {
		p.s = WHITESPACE
	}
}

// Fields breaks a string into fields, calling fn for each parsed field
// (to avoid memory allocation)
func Fields(m string, fn func(string) error) error {
	parse.init()
	var err error
	starti := 0
	for i, c := range m {
		switch parse.s {
		case WHITESPACE:
			parse.whitespace(c)
			if parse.s != WHITESPACE {
				starti = i
			}
		case TOKEN:
			err = parse.token(c, fn)
		case STRING_SINGLE, STRING_DOUBLE, STRING_BRACES:
			err = parse.str(c, fn)
		case COMMENT:
			parse.comment(c)
		}
		if err != nil {
			elog.Heap("alloc: /parse/parse.go:163: buildContextString(m, starti, i),")
			return fmt.Errorf(
				"%s: %w",
				buildContextString(m, starti, i), // object allocated on the heap: escapes at line 163
				err)
		}
	}

	switch parse.s {
	case TOKEN:
		err = parse.token('\n', fn)
	case STRING_SINGLE:
		err = ErrUnterminatedSingleQuote
	case STRING_DOUBLE:
		err = ErrUnterminatedDouble
	case STRING_BRACES:
		err = ErrUnterminatedBrace
	}

	if err != nil {
		elog.Heap("alloc: /parse/parse.go:180: buildContextString(m, starti, len(m)),")
		return fmt.Errorf(
			"%s: %w",
			buildContextString(m, starti, len(m)), // object allocated on the heap: escapes at line 180
			err)
	}
	return nil
}

// set a max context length so that the error location is not buried
const maxContextLength = 80

func buildContextString(m string, starti, endi int) string {
	if len(m) > maxContextLength {
		m, starti, endi = truncateString(m, starti, endi, maxContextLength)
	}
	if starti >= len(m) {
		return m + "<-"
	}
	return m[:starti] + "->" + m[starti:endi] + "<-" + m[endi:]
}

// If the error message is too long, this function tries to intelligently,
// reduce it's size so that the problem argument is still present. examples:
//
// Try to center ->the<- span (m[starti:endi])
// But sometimes  it's  close  to  ->the<- end
// Or ->the<- start of  the  truncated  string
//
// Use maxlen characters for all cases.  The function assumes
// that the caller checked for len(m) > maxlen.
func truncateString(m string, starti, endi, maxlen int) (string, int, int) {
	if endi-starti >= maxlen {
		// The span exceeds maxlen, just show part of it.
		return m[starti : starti+maxlen], 0, maxlen
	}

	// The target center of the truncated string
	center := (starti + endi) / 2

	// start and end use the center point
	starts := center - maxlen/2
	ends := starts + maxlen

	if starts < 0 {
		// the center to too close to the start of the string
		// so clamp the start to zero.
		starts = 0
		ends = starts + maxlen
	} else if starts > 0 {
		// the start is beyond zero so the front of the
		// string will be trimmed and the indexes need to
		// be shifted to be relative to the trimmed string
		starti -= starts
		endi -= starts
	}

	if ends > len(m) {
		// The span was close to the end of the string
		// so we clamp to the end of the string.
		delta := ends - len(m)
		starti += delta
		endi += delta
		starts -= delta
		ends -= delta
	}

	return m[starts:ends], starti, endi
}

func isWhitespace(c rune) bool {
	return (c == ' ') || (c == '\t') || (c == '\n')
}
