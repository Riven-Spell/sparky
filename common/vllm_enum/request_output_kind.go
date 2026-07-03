package vllm_enum

import "github.com/Riven-Spell/enum/v2"

type eRequestOutputKind struct {
	enum.EnumImpl[RequestOutputKind, eRequestOutputKind]
}

var ERequestOutputKind eRequestOutputKind

type RequestOutputKind int

func (k RequestOutputKind) String() string {
	return ERequestOutputKind.String(k)
}

func (eRequestOutputKind) Text() RequestOutputKind   { return 0 }
func (eRequestOutputKind) Tokens() RequestOutputKind  { return 1 }
func (eRequestOutputKind) Full() RequestOutputKind    { return 2 }
