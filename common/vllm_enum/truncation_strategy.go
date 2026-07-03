package vllm_enum

import "github.com/Riven-Spell/enum/v2"

type eTruncationStrategy struct {
	enum.EnumImpl[TruncationStrategy, eTruncationStrategy]
}

var ETruncationStrategy eTruncationStrategy

type TruncationStrategy string

func (s TruncationStrategy) String() string {
	return string(s)
}

func (eTruncationStrategy) Auto() TruncationStrategy     { return "auto" }
func (eTruncationStrategy) Disabled() TruncationStrategy { return "disabled" }
