package xfs

import (
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type hybridFS struct {
	ofs fs.FS
	efs fs.FS
}

func (f hybridFS) Open(name string) (fs.File, error) {
	if filepath.IsAbs(name) {
		log.Debug().Str("file", name).Str("FileSystem", "OS").Msg("trying to open")
		return os.DirFS("").Open(name[1:]) // FIXME: what for windows?
	}

	log.Debug().Str("file", name).Str("FileSystem", "OS").Msg("trying to open")
	if file, err := f.ofs.Open(name); err == nil {
		return file, nil
	} else {
		log.Debug().Str("error", err.Error()).Msgf("Got error. will try Embed FS next")
	}

	log.Debug().Str("file", name).Str("FileSystem", "Embed").Msg("trying to open")
	file, err := f.efs.Open(name)
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
