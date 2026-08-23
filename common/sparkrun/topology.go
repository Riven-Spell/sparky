package sparkrun

import "github.com/Riven-Spell/enum/v2"

type eTopology struct {
	enum.EnumImpl[Topology, eTopology]
}

var ETopology eTopology

// Topology is the CX7 topology passed to
// `sparkrun cluster update --topology`.
type Topology string

func (t Topology) String() string {
	return string(t)
}

func (eTopology) None() Topology   { return "none" }
func (eTopology) Direct() Topology { return "direct" }
func (eTopology) Switch() Topology { return "switch" }
func (eTopology) Ring() Topology   { return "ring" }
