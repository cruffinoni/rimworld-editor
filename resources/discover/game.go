package discover

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

var ErrPathNotFound = errors.New("game path not found")

func GetGamePath(opeSystem string) (string, error) {
	var gameDir string

	if opeSystem == "" {
		opeSystem = runtime.GOOS
	}
	switch opeSystem {
	case Windows:
		// Check all possible locations on Windows
		possibleDirs := []string{
			os.Getenv("ProgramFiles(x86)"),
			os.Getenv("ProgramFiles"),
			"D:\\Program Files (x86)",
			"D:\\Program Files",
			"E:\\Program Files (x86)",
			"E:\\Program Files",
		}
		for _, dir := range possibleDirs {
			path := filepath.Join(dir, "Steam", "steamapps", "common", "RimWorld")
			if _, err := os.Stat(path); err == nil {
				gameDir = path
				break
			}
		}
	case Darwin:
		gameDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Steam", "steamapps", "common", "RimWorld")
	case Linux:
		gameDir = filepath.Join(os.Getenv("HOME"), ".steam", "steam", "steamapps", "common", "RimWorld")
	default:
		return "", ErrOSNotSupported
	}

	if gameDir == "" {
		return "", ErrPathNotFound
	}

	return gameDir, nil
}
