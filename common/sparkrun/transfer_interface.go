package sparkrun

import "github.com/Riven-Spell/enum/v2"

type eTransferInterface struct {
	enum.EnumImpl[TransferInterface, eTransferInterface]
}

var ETransferInterface eTransferInterface

// TransferInterface is the network interface used for transfers,
// passed to `sparkrun cluster create|update --transfer-interface`.
type TransferInterface string

func (i TransferInterface) String() string {
	return string(i)
}

func (eTransferInterface) Auto() TransferInterface { return "auto" }
func (eTransferInterface) Cx7() TransferInterface  { return "cx7" }
func (eTransferInterface) Mgmt() TransferInterface { return "mgmt" }
