package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed assets/translation_files.zip
var translationFilesZip []byte

type PckExplorerService struct {
	ctx   context.Context
	state *InstallerState
}

func NewPckExplorerService(state *InstallerState) *PckExplorerService {
	return &PckExplorerService{state: state}
}

func (s *PckExplorerService) startup(ctx context.Context) {
	s.ctx = ctx
}

func (s *PckExplorerService) RunInstallation() {
	go s.startInstallProcess()
}

func (s *PckExplorerService) startInstallProcess() {
	targetPckPath, isDemo, makeBackup := s.state.GetState()

	gameDir := filepath.Dir(targetPckPath)
	modifiedPckPath := filepath.Join(gameDir, "ModifiedPCK.pck")
	backupPckPath := filepath.Join(gameDir, "UntilThen.pck.bak")

	tempDir, err := os.MkdirTemp("", "untilthen_patcher_*")
	if err != nil {
		s.failAndLog(modifiedPckPath, fmt.Errorf("falha ao criar pasta temporária: %w", err))
		return
	}
	defer os.RemoveAll(tempDir)

	wails.EventsEmit(s.ctx, "install_step", "Extraindo ferramentas de patch...")
	if err := s.unzipFromMemory(pckExplorerBinZip, tempDir, "unzip_bin_progress"); err != nil {
		s.failAndLog(modifiedPckPath, err)
		return
	}

	binPath := filepath.Join(tempDir, pckBinName)

	if runtime.GOOS != "windows" {
		os.Chmod(binPath, 0755)
	}

	wails.EventsEmit(s.ctx, "install_step", "Preparando arquivos de tradução...")
	if err := s.unzipFromMemory(translationFilesZip, tempDir, "unzip_trans_progress"); err != nil {
		s.failAndLog(modifiedPckPath, err)
		return
	}

	translationFolder := "full"
	if isDemo {
		translationFolder = "demo"
	}
	translationFilesPath := filepath.Join(tempDir, translationFolder)

	wails.EventsEmit(s.ctx, "install_step", "Aplicando tradução (isso pode levar alguns instantes)...")

	cmd := exec.CommandContext(s.ctx, binPath, "-pc", targetPckPath, translationFilesPath, modifiedPckPath, "2.2.4.1")

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		s.failAndLog(modifiedPckPath, fmt.Errorf("falha ao iniciar processo: %w", err))
		return
	}

	go s.streamLogs(stdout, "install_log")
	go s.streamLogs(stderr, "install_error")

	if err := cmd.Wait(); err != nil {
		s.failAndLog(modifiedPckPath, fmt.Errorf("erro durante o patch: %w", err))
		return
	}

	wails.EventsEmit(s.ctx, "install_step", "Finalizando instalação...")

	if makeBackup {
		os.Remove(backupPckPath)
		if err := os.Rename(targetPckPath, backupPckPath); err != nil {
			s.failAndLog(modifiedPckPath, fmt.Errorf("erro no backup: %w", err))
			return
		}
	} else {
		os.Remove(targetPckPath)
	}

	if err := os.Rename(modifiedPckPath, targetPckPath); err != nil {
		if makeBackup {
			os.Rename(backupPckPath, targetPckPath)
		}
		s.failAndLog(modifiedPckPath, fmt.Errorf("erro ao renomear arquivo final: %w", err))
		return
	}

	wails.EventsEmit(s.ctx, "install_success", "Sucesso!")
}

func (s *PckExplorerService) unzipFromMemory(data []byte, dest, eventName string) error {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	total := len(reader.File)
	for i, f := range reader.File {
		progress := int(float64(i+1) / float64(total) * 100)
		wails.EventsEmit(s.ctx, eventName, progress)

		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func (s *PckExplorerService) failAndLog(modifiedPck string, err error) {
	wails.EventsEmit(s.ctx, "install_error", err.Error())
	if modifiedPck != "" {
		os.Remove(modifiedPck)
	}
}

func (s *PckExplorerService) streamLogs(pipe io.ReadCloser, eventName string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		wails.EventsEmit(s.ctx, eventName, scanner.Text())
	}
}
