package sparkrun

import "github.com/Riven-Spell/enum/v2"

type eTransferMode struct {
	enum.EnumImpl[TransferMode, eTransferMode]
}

var ETransferMode eTransferMode

// TransferMode is the resource transfer mode passed to
// `sparkrun cluster create|update --transfer-mode`.
type TransferMode string

func (m TransferMode) String() string {
	return string(m)
}

func (eTransferMode) Auto() TransferMode      { return "auto" }
func (eTransferMode) Local() TransferMode     { return "local" }
func (eTransferMode) Push() TransferMode      { return "push" }
func (eTransferMode) Delegated() TransferMode { return "delegated" }
