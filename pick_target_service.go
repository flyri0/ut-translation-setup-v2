package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/andygrunwald/vdf"
)

// Tenta encontrar o local de instalalção do jogo, recebe o local de instalação da steam e um id de jogo da steam como argumento
// e retorna uma string vazia caso não encontre.
func findGamePath(steamPath string, gameID int) string {
	if steamPath == "" {
		return ""
	}

	libraryFoldersPath := filepath.Join(steamPath, "steamapps", "libraryfolders.vdf")
	file, err := os.Open(libraryFoldersPath)
	if err != nil {
		return ""
	}

	data, err := vdf.NewParser(file).Parse()
	file.Close()
	if err != nil {
		return ""
	}

	folders, ok := data["libraryfolders"].(map[string]any)
	if !ok {
		return ""
	}

	var libraries []string
	for _, entry := range folders {
		entryMap, ok := entry.(map[string]any)
		if !ok {
			continue
		}

		if path, ok := entryMap["path"].(string); ok && path != "" {
			libraries = append(libraries, path)
		}
	}
	if len(libraries) == 0 {
		libraries = []string{steamPath}
	}

	for _, lib := range libraries {
		manifestPath := filepath.Join(lib, "steamapps", fmt.Sprintf("appmanifest_%d.acf", gameID))

		manifestFile, err := os.Open(manifestPath)
		if err != nil {
			continue
		}

		manifest, err := vdf.NewParser(manifestFile).Parse()
		manifestFile.Close()
		if err != nil {
			continue
		}

		appState, ok := manifest["AppState"].(map[string]any)
		if !ok {
			continue
		}

		installDir, ok := appState["installdir"].(string)
		if !ok || installDir == "" {
			continue
		}

		gamePath := filepath.Join(lib, "steamapps", "common", installDir)
		if info, err := os.Stat(gamePath); err == nil && info.IsDir() {
			absolutePath, _ := filepath.Abs(gamePath)
			return absolutePath
		}
	}

	return ""
}

// Tenta encontrar o local de instalação da steam em diferentes plataformas.
// Retorna uma string vazia caso não encontre.
func findSteamPath() string {
	var candidates []string

	switch runtime.GOOS {
	case "linux":
		home, _ := os.UserHomeDir()
		candidates = []string{
			filepath.Join(home, ".steam", "steam"),
			filepath.Join(home, ".local", "share", "Steam"),
			filepath.Join(home, "snap", "steam", "common", ".local", "share", "Steam"),
			filepath.Join(home, ".var", "app", "com.valvesoftware.Steam", "data", "Steam"),
		}

	case "windows":
		if path := steamPathFromRegistry(); path != "" {
			candidates = append(candidates, path)
		}

		candidates = append(candidates,
			filepath.Join(os.Getenv("ProgramFiles(x86)"), "Steam"),
			filepath.Join(os.Getenv("ProgramFiles"), "Steam"),
		)

	default:
		return ""
	}

	return ""
}
