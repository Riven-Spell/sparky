package vllm_enum

import "github.com/Riven-Spell/enum/v2"

type eMessagePhase struct {
	enum.EnumImpl[MessagePhase, eMessagePhase]
}

var EMessagePhase eMessagePhase

type MessagePhase string

func (p MessagePhase) String() string {
	return string(p)
}

func (eMessagePhase) Commentary() MessagePhase  { return "commentary" }
func (eMessagePhase) FinalAnswer() MessagePhase { return "final_answer" }
