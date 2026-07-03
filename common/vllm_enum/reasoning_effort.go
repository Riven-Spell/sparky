package vllm_enum

import "github.com/Riven-Spell/enum/v2"

type eReasoningEffort struct {
	enum.EnumImpl[ReasoningEffort, eReasoningEffort]
}

var EReasoningEffort eReasoningEffort

type ReasoningEffort string

func (e ReasoningEffort) String() string {
	return string(e)
}

func (eReasoningEffort) None() ReasoningEffort    { return "none" }
func (eReasoningEffort) Minimal() ReasoningEffort { return "minimal" }
func (eReasoningEffort) Low() ReasoningEffort     { return "low" }
func (eReasoningEffort) Medium() ReasoningEffort  { return "medium" }
func (eReasoningEffort) High() ReasoningEffort    { return "high" }
func (eReasoningEffort) XHigh() ReasoningEffort   { return "xhigh" }
func (eReasoningEffort) Max() ReasoningEffort     { return "max" }
