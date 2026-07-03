package vllm_enum

import "github.com/Riven-Spell/enum/v2"

type eImageDetail struct {
	enum.EnumImpl[ImageDetail, eImageDetail]
}

var EImageDetail eImageDetail

type ImageDetail string

func (d ImageDetail) String() string {
	return string(d)
}

func (eImageDetail) Auto() ImageDetail     { return "auto" }
func (eImageDetail) Low() ImageDetail      { return "low" }
func (eImageDetail) High() ImageDetail     { return "high" }
func (eImageDetail) Original() ImageDetail { return "original" }
