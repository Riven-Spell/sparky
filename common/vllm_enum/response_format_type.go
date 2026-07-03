package vllm_enum

import "github.com/Riven-Spell/enum/v2"

type eResponseFormatType struct {
	enum.EnumImpl[ResponseFormatType, eResponseFormatType]
}

var EResponseFormatType eResponseFormatType

type ResponseFormatType string

func (t ResponseFormatType) String() string {
	return string(t)
}

func (eResponseFormatType) Text() ResponseFormatType       { return "text" }
func (eResponseFormatType) JSONObject() ResponseFormatType { return "json_object" }
func (eResponseFormatType) JSONSchema() ResponseFormatType { return "json_schema" }
