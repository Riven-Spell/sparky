package models

import "github.com/Riven-Spell/sparky/common/enum"

type ChatMessage struct {
	Role    enum.ChatRole `json:"role"`
	Content string        `json:"content"`
}
