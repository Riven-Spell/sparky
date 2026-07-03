package vllm_enum

import "github.com/Riven-Spell/enum/v2"

type eExecutionType struct {
	enum.EnumImpl[ExecutionType, eExecutionType]
}

var EExecutionType eExecutionType

type ExecutionType string

func (t ExecutionType) String() string {
	return string(t)
}

func (eExecutionType) Server() ExecutionType { return "server" }
func (eExecutionType) Client() ExecutionType { return "client" }
