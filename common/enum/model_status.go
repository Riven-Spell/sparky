package enum

import "github.com/Riven-Spell/enum/v2"

type eModelStatus struct {
	enum.EnumImpl[ModelStatus, eModelStatus]
}

var EModelStatus eModelStatus

type ModelStatus string

func (s ModelStatus) String() string {
	return string(s)
}

func (eModelStatus) Loaded() ModelStatus   { return "loaded" }
func (eModelStatus) Unloaded() ModelStatus { return "unloaded" }
func (eModelStatus) Loading() ModelStatus  { return "loading" }
func (eModelStatus) Evicting() ModelStatus { return "evicting" }
func (eModelStatus) Crashing() ModelStatus { return "crashing" }
