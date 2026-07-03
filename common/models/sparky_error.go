package models

import "github.com/Riven-Spell/sparky/common/enum"

type SparkyError struct {
	ErrorCode    enum.SparkyErrorCode `json:"ErrorCode"`
	ErrorDetail  string               `json:"ErrorDetail"`
}
