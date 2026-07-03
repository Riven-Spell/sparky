package vllm_enum

import "github.com/Riven-Spell/enum/v2"

type eToolChoiceType struct {
	enum.EnumImpl[ToolChoiceType, eToolChoiceType]
}

var EToolChoiceType eToolChoiceType

type ToolChoiceType string

func (t ToolChoiceType) String() string {
	return string(t)
}

func (eToolChoiceType) None() ToolChoiceType     { return "none" }
func (eToolChoiceType) Auto() ToolChoiceType     { return "auto" }
func (eToolChoiceType) Required() ToolChoiceType { return "required" }
