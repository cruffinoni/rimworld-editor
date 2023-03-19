package discover

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sort"
)

var ErrOSNotSupported = errors.New("operating system not supported")

func GetSavegamePath() (string, error) {
	var path string
	var err error

	switch runtime.GOOS {
	case "windows":
		localLowAppData := os.Getenv("LOCALLOWAPPDATA")
		if localLowAppData != "" {
			path = filepath.Join(localLowAppData, "Ludeon Studios", "RimWorld by Ludeon Studios", "Saves")
		} else {
			appData := os.Getenv("LOCALAPPDATA")
			path = filepath.Join(appData, "..", "LocalLow", "Ludeon Studios", "RimWorld by Ludeon Studios", "Saves")
		}
	case "darwin":
		path = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "RimWorld", "Saves")
	case "linux":
		path = filepath.Join(os.Getenv("HOME"), ".config", "unity3d", "Ludeon Studios", "RimWorld", "Saves")
	default:
		return "", ErrOSNotSupported
	}

	return path, err
}

func GetLatestSavegameFiles(maxFile int) ([]os.FileInfo, error) {
	savegamePath, err := GetSavegamePath()
	if err != nil {
		return nil, err
	}

	dir, err := os.Open(savegamePath)
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
