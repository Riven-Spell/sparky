package vllm_enum

import "github.com/Riven-Spell/enum/v2"

type eComputerEnvironment struct {
	enum.EnumImpl[ComputerEnvironment, eComputerEnvironment]
}

var EComputerEnvironment eComputerEnvironment

type ComputerEnvironment string

func (e ComputerEnvironment) String() string {
	return string(e)
}

func (eComputerEnvironment) Windows() ComputerEnvironment  { return "windows" }
func (eComputerEnvironment) Mac() ComputerEnvironment      { return "mac" }
func (eComputerEnvironment) Linux() ComputerEnvironment    { return "linux" }
func (eComputerEnvironment) Ubuntu() ComputerEnvironment   { return "ubuntu" }
func (eComputerEnvironment) Browser() ComputerEnvironment  { return "browser" }
