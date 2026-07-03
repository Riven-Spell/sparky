package enum

import "github.com/Riven-Spell/enum/v2"

type eObjectType struct {
	enum.EnumImpl[ObjectType, eObjectType]
}

var EObjectType eObjectType

type ObjectType string

func (t ObjectType) String() string {
	return string(t)
}

func (eObjectType) ChatCompletion() ObjectType { return "chat.completion" }
func (eObjectType) TextCompletion() ObjectType { return "text_completion" }
func (eObjectType) Model() ObjectType          { return "model" }
func (eObjectType) List() ObjectType           { return "list" }
