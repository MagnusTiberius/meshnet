package lex

import (
	"strings"
	"unicode/utf8"
)

//PrvLexer lexer
type PrvLexer struct {
	s     string
	pos   int
	width int
	start int
}

func (l *PrvLexer) Error(s string) {
	//fmt.Printf("%s\nLine: %d\n", s, lineCounter)
	panic("Error:")
}

func (l *PrvLexer) peek() int32 {
	rune := l.next()
	l.backup()
	return rune
}

func (l *PrvLexer) backup() {
	l.pos -= l.width
}

func (l *PrvLexer) ignore() {
	l.start = l.pos
}

func (l *PrvLexer) next() (rune rune) {
	if l.pos >= len(l.s) {
		l.width = 0
		return 0
	}
	rune, l.width = utf8.DecodeRuneInString(l.s[l.pos:])
	l.pos += l.width
	return rune
}

func (l *PrvLexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *PrvLexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

//Lex lexer
func (l *PrvLexer) Lex(lval []byte) int {
	return 0
}
