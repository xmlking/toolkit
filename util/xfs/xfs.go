package xfs

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"golang.org/x/tools/go/packages"
)

type hybridFS struct {
	ofs fs.FS
	efs fs.FS
}

func (f hybridFS) Open(name string) (file fs.File, err error) {
	if filepath.IsAbs(name) {
		file, err = os.DirFS("").Open(name[1:]) // FIXME: what for windows?
		log.Debug().Str("file", name).Str("FileSystem", "OS").Err(err).Msg("loading from")
		return
	}

	file, err = f.ofs.Open(name)
	log.Debug().Str("file", name).Str("FileSystem", "OS").Err(err).Msg("loading from")
	if err == nil {
		return
	}

	file, err = f.efs.Open(name)
	log.Debug().Str("file", name).Str("FileSystem", "Embed").Err(err).Msg("loading from")
	return
}

func FS(efs fs.FS) fs.FS {
	root, err := getGoModuleDir()
	if err != nil || root == "" {
		root = "."
		log.Debug().Err(err).Msgf("got no module path. using FileSystem root as: '%s/'", root)

	} else {
		log.Debug().Msgf("got module path. using FileSystem root as: '%s'", root)
	}
	ofs := os.DirFS(root)
	return &hybridFS{ofs, efs}
}

// TODO: remove this method.
// this only works for non-workspace projects (go 1.18)
func getGoModuleDirForNonWorkspace() (path string, err error) {
	cmd := exec.Command("go", "list", "-json", "-m")
	//cmd.Env = append(os.Environ(), "GO111MODULE=on")
	var out []byte
	if out, err = cmd.Output(); err != nil {
		// error: xec: \"go\": executable file not found in $PATH"
		// means running in in docker/pod
		return
	}

	var mod struct {
		Dir string
	}
	if err = json.Unmarshal(out, &mod); err != nil {
		log.Error().Err(err).Msg("error Unmarshal 'go list -json -m' output")
		return
	}
	path = mod.Dir
	return
}

func getGoModuleDir() (path string, err error) {
	cfg := &packages.Config{
		Mode: packages.NeedModule,
	}
	root, err := packages.Load(cfg, "")
	if err != nil {
		return "", fmt.Errorf("load packages error: %v", err)
	}
	if len(root) != 1 {
		return "", fmt.Errorf("unsupported packages number: %d", len(root))
	}
	packages.PrintErrors(root)
	log.Debug().Msgf("%v", root[0].Module.Dir)
	return root[0].Module.Dir, nil
}
