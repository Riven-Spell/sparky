package vllm_enum

import "github.com/Riven-Spell/enum/v2"

type eTruncationSide struct {
	enum.EnumImpl[TruncationSide, eTruncationSide]
}

var ETruncationSide eTruncationSide

type TruncationSide string

func (s TruncationSide) String() string {
	return string(s)
}

func (eTruncationSide) Left() TruncationSide  { return "left" }
func (eTruncationSide) Right() TruncationSide { return "right" }
