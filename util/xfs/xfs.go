package xfs

import (
	"fmt"
	"io/fs"
	"os"
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
		//file, err = os.DirFS("").Open(name[1:]) // FIXME: what for windows?
		root := "/"
		if vol := filepath.VolumeName(name); vol != "" {
			root = vol + "\\"
		}
		var rel string
		if rel, err = filepath.Rel(root, name); err != nil {
			return
		}
		rel = filepath.ToSlash(rel)
		file, err = os.DirFS(root).Open(rel)
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
