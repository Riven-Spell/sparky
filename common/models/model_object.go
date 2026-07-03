package models

import "github.com/Riven-Spell/sparky/common/enum"

type ModelObject struct {
	Id       string          `json:"id"`
	Object   enum.ObjectType `json:"object"`
	Created  int64           `json:"created"`
	OwnedBy  string          `json:"owned_by"`
}
