package enum

import (
	"github.com/Riven-Spell/enum/v2"
	"os"
	"path/filepath"
	"runtime"
)

type eEnvironmentVariable struct {
	enum.EnumImpl[EnvironmentVariable, eEnvironmentVariable]
}

var EEnvironmentVariable eEnvironmentVariable

type EnvironmentVariable struct {
	Name        string
	Description string
	Default     string
	Secret      bool
	Hidden      bool
}

func (e EnvironmentVariable) String() string {
	return e.Name
}

func (e EnvironmentVariable) Lookup() (string, bool) {
	val, ok := os.LookupEnv(e.Name)
	return val, ok
}

func (e EnvironmentVariable) Get() string {
	val, ok := e.Lookup()
	if !ok {
		return e.Default
	}
	return val
}

func (e EnvironmentVariable) Clear() error {
	return os.Unsetenv(e.Name)
}

func (eEnvironmentVariable) AppDir() EnvironmentVariable {
	return EnvironmentVariable{
		Name:        "SPARKY_APP_DIR",
		Description: "Application directory for Sparky persistent state",
		Default:     filepath.Join(resolveHome(), ".sparky"),
		Secret:      false,
		Hidden:      false,
	}
}

func (eEnvironmentVariable) SparkyConfig() EnvironmentVariable {
	return EnvironmentVariable{
		Name:        "SPARKY_CONFIG",
		Description: "Manager configuration file, relative to AppDir",
		Default:     filepath.Join(EEnvironmentVariable.AppDir().Get(), "manager.yaml"),
		Secret:      false,
		Hidden:      false,
	}
}

func resolveHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		if runtime.GOOS == "windows" {
			return "C:\\"
		}
		return "/"
	}
	return home
}