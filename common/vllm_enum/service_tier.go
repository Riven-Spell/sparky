package vllm_enum

import "github.com/Riven-Spell/enum/v2"

type eServiceTier struct {
	enum.EnumImpl[ServiceTier, eServiceTier]
}

var EServiceTier eServiceTier

type ServiceTier string

func (t ServiceTier) String() string {
	return string(t)
}

func (eServiceTier) Auto() ServiceTier     { return "auto" }
func (eServiceTier) Default() ServiceTier  { return "default" }
func (eServiceTier) Flex() ServiceTier     { return "flex" }
func (eServiceTier) Scale() ServiceTier    { return "scale" }
func (eServiceTier) Priority() ServiceTier { return "priority" }
