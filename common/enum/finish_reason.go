package enum

import "github.com/Riven-Spell/enum/v2"

type eFinishReason struct {
	enum.EnumImpl[FinishReason, eFinishReason]
}

var EFinishReason eFinishReason

type FinishReason string

func (r FinishReason) String() string {
	return string(r)
}

func (eFinishReason) Stop() FinishReason  { return "stop" }
func (eFinishReason) Length() FinishReason { return "length" }
