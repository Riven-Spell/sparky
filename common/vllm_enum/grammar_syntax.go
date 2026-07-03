package vllm_enum

import "github.com/Riven-Spell/enum/v2"

type eGrammarSyntax struct {
	enum.EnumImpl[GrammarSyntax, eGrammarSyntax]
}

var EGrammarSyntax eGrammarSyntax

type GrammarSyntax string

func (s GrammarSyntax) String() string {
	return string(s)
}

func (eGrammarSyntax) Lark() GrammarSyntax  { return "lark" }
func (eGrammarSyntax) Regex() GrammarSyntax { return "regex" }
