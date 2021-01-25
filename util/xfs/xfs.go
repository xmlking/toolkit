package xfs

import (
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

type hybridFS struct {
	ofs fs.FS
	efs fs.FS
}

func (f hybridFS) Open(name string) (fs.File, error) {
	file, err1 := f.ofs.Open(name)
	if err1 == nil {
		log.Debug().Msgf("loading from OS FileSystem: (%s)", name)
		return file, nil
	}

	file, err := f.efs.Open(name)
	log.Debug().Err(err1).Msgf("loading form Embed FileSystem: (%s)", name)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func FS(efs fs.FS) fs.FS {
	root, err := getGoModuleDir()
	if err != nil || root == "" {
		root = "."
		log.Info().Err(err).Msgf("got no module path. using FileSystem root as: '%s/'", root)

	} else {
		log.Info().Msgf("got module path. using FileSystem root as: '%s'", root)
	}
	ofs := os.DirFS(root)
	return &hybridFS{ofs, efs}
}

func getGoModuleDir() (path string, err error) {
	cmd := exec.Command("go", "list", "-json", "-m")
	//cmd.Env = append(os.Environ(), "GO111MODULE=on")
	var out []byte
	if out, err = cmd.Output(); err != nil {
		log.Error().Err(err).Msg("error running 'go list -json -m'")
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
