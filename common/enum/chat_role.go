package enum

import "github.com/Riven-Spell/enum/v2"

type eChatRole struct {
	enum.EnumImpl[ChatRole, eChatRole]
}

var EChatRole eChatRole

type ChatRole string

func (r ChatRole) String() string {
	return string(r)
}

func (eChatRole) System() ChatRole    { return "system" }
func (eChatRole) User() ChatRole      { return "user" }
func (eChatRole) Assistant() ChatRole { return "assistant" }
