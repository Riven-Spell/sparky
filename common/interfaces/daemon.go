package interfaces

type Daemon interface {
	Start(onExit ...func()) error
	Stop() error
	Running() bool
	Error() error
}
