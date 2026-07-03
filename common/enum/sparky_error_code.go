package enum

import "github.com/Riven-Spell/enum/v2"

type eSparkyErrorCode struct {
	enum.EnumImpl[SparkyErrorCode, eSparkyErrorCode]
}

var ESparkyErrorCode eSparkyErrorCode

type SparkyErrorCode string

func (c SparkyErrorCode) String() string {
	return string(c)
}

func (eSparkyErrorCode) ModelNotRegistered() SparkyErrorCode { return "ModelNotRegistered" }
func (eSparkyErrorCode) ModelTooLarge() SparkyErrorCode      { return "ModelTooLarge" }
func (eSparkyErrorCode) CannotEvict() SparkyErrorCode        { return "CannotEvict" }
func (eSparkyErrorCode) ModelLoading() SparkyErrorCode       { return "ModelLoading" }
