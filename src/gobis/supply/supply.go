package supply

import (
	"fmt"
	"github.com/cloudfoundry/libbuildpack/checksum"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cloudfoundry/libbuildpack"
)

type Stager interface {
	// TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/stager.go
	BuildDir() string
	DepDir() string
	DepsIdx() string
	DepsDir() string
	LinkDirectoryInDepDir(string, string) error
	WriteProfileD(scriptName, scriptContents string) error
}

type Manifest interface {
	// TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/manifest.go
	AllDependencyVersions(string) []string
	DefaultVersion(string) (libbuildpack.Dependency, error)
}

type Installer interface {
	// TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/installer.go
	InstallDependency(libbuildpack.Dependency, string) error
	InstallOnlyVersion(string, string) error
}

type Command interface {
	// TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/command.go
	Execute(string, io.Writer, io.Writer, string, ...string) error
	Output(dir string, program string, args ...string) (string, error)
}

type Supplier struct {
	Manifest  Manifest
	Installer Installer
	Stager    Stager
	Command   Command
	Log       *libbuildpack.Logger
}

func (s *Supplier) Run() error {
	return checksum.Do(s.Stager.BuildDir(), s.Log.Debug, func() error {
		s.Log.BeginStep("Installing gobis-server")
		err := s.InstallGobisServer("/tmp/gobis-server")
		if err != nil {
			return err
		}
		return s.Stager.WriteProfileD("gobis.sh",
			fmt.Sprintf(`
export PATH=$PATH:"$HOME/bin":%s
export GOBIS_PORT=8081
export PORT=8081
export VCAP_APP_PORT=8081
`,
				filepath.Join("$DEPS_DIR", s.Stager.DepsIdx(), "bin")))
	})
}

func (s *Supplier) InstallGobisServer(tempDir string) error {
	dep, err := s.Manifest.DefaultVersion("gobis-server")
	if err != nil {
		return err
	}
	gServerName := "gobis-server"
	if runtime.GOOS == "windows" {
		gServerName += ".exe"
	}
	installDir := filepath.Join(filepath.Join(s.Stager.DepDir(), "bin", gServerName))

	if err := s.Installer.InstallDependency(dep, tempDir); err != nil {
		return err
	}

	binName := "gobis-server_linux_amd64"
	if runtime.GOOS == "windows" {
		binName = "gobis-server_windows_amd64.exe"
	}

	if err := os.Rename(filepath.Join(tempDir, binName), installDir); err != nil {
		return err
	}

	return os.Setenv("PATH", fmt.Sprintf("%s:%s", os.Getenv("PATH"), filepath.Join(s.Stager.DepDir(), "bin")))
}
