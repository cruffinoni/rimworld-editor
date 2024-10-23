package discover

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sort"
)

var ErrOSNotSupported = errors.New("operating system not supported")

const (
	Windows = "windows"
	Darwin  = "macos"
	Linux   = "linux"
)

func GetSavegamePath(opeSystem string) (string, error) {
	var (
		path string
		err  error
	)

	if opeSystem == "" {
		opeSystem = runtime.GOOS
	}
	switch opeSystem {
	case Windows:
		localLowAppData := os.Getenv("LOCALLOWAPPDATA")
		if localLowAppData != "" {
			path = filepath.Join(localLowAppData, "Ludeon Studios", "RimWorld by Ludeon Studios", "Saves")
		} else {
			appData := os.Getenv("LOCALAPPDATA")
			path = filepath.Join(appData, "..", "LocalLow", "Ludeon Studios", "RimWorld by Ludeon Studios", "Saves")
		}
	case Darwin:
		path = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "RimWorld", "Saves")
	case Linux:
		path = filepath.Join(os.Getenv("HOME"), ".config", "unity3d", "Ludeon Studios", "RimWorld", "Saves")
	default:
		return "", ErrOSNotSupported
	}

	return path, err
}

func GetLatestSavegameFiles(maxFile int, savePath string) ([]os.FileInfo, error) {
	dir, err := os.Open(savePath)
	if err != nil {
		return nil, err
	}

	defer dir.Close()
	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	for _, entry := range files {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".rws" {
			files = append(files, entry)
		}
	}

	// Sort files by modification time (newest first)
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})

	if len(files) > maxFile {
		files = files[:maxFile]
	}

	return files, nil
}
