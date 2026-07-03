package vllm_enum

import "github.com/Riven-Spell/enum/v2"

type eVerbosity struct {
	enum.EnumImpl[Verbosity, eVerbosity]
}

var EVerbosity eVerbosity

type Verbosity string

func (v Verbosity) String() string {
	return string(v)
}

func (eVerbosity) Low() Verbosity   { return "low" }
func (eVerbosity) Medium() Verbosity { return "medium" }
func (eVerbosity) High() Verbosity  { return "high" }
