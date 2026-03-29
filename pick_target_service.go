package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/andygrunwald/vdf"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

type FileValidationResult struct {
	Valid  bool   `json:"valid"`
	IsDemo bool   `json:"isDemo"`
	Path   string `json:"path"`
	Error  string `json:"error,omitempty"` // "not_found" | "not_selected" | "invalid_file"
}

type PickTargetService struct {
	ctx context.Context
}

func NewPickTargetService() *PickTargetService {
	return &PickTargetService{}
}

func (s *PickTargetService) startup(ctx context.Context) {
	s.ctx = ctx
}

func (s *PickTargetService) QuickFind() FileValidationResult {
	steamPath := s.findSteamPath()

	for _, gameID := range []int{1574820, 2296400} {
		installDir := s.findGamePath(steamPath, gameID)
		if installDir == "" {
			continue
		}

		candidate := filepath.Join(installDir, "UntilThen.pck")
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return s.validateFile(candidate)
		}
	}

	return FileValidationResult{Valid: false, Error: "not_found"}
}

func (s *PickTargetService) OpenFilePicker() FileValidationResult {
	path, err := wails.OpenFileDialog(s.ctx, wails.OpenDialogOptions{
		Title: "Selecione UntilThen.pck",
		Filters: []wails.FileFilter{
			{DisplayName: "UntilThen.pck (*.pck)", Pattern: "*.pck"},
		},
	})

	if err != nil || path == "" {
		return FileValidationResult{Valid: false, Error: "not_selected"}
	}

	return s.validateFile(path)
}

func (s *PickTargetService) validateFile(path string) FileValidationResult {
	if path == "" {
		return FileValidationResult{Valid: false, Error: "invalid_file"}
	}

	info, err := os.Stat(path)
	if err != nil || info.IsDir() || strings.ToLower(filepath.Ext(path)) != ".pck" {
		return FileValidationResult{Valid: false, Error: "invalid_file"}
	}

	isDemo := strings.Contains(filepath.Dir(path), "Until Then Demo")

	return FileValidationResult{
		Valid:  true,
		IsDemo: isDemo,
		Path:   path,
	}
}

func (s *PickTargetService) findGamePath(steamPath string, gameID int) string {
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

func (s *PickTargetService) findSteamPath() string {
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

	for _, p := range candidates {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return p
		}
	}

	return ""
}
