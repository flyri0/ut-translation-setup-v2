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
	"time"

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
	wails.LogInfo(s.ctx, "PckExplorerService Iniciado")
}

func (s *PckExplorerService) RunInstallation() {
	wails.LogInfo(s.ctx, "Solicitação de instalação recebida. Iniciando processo em background...")
	go s.startInstallProcess()
}

func (s *PckExplorerService) startInstallProcess() {
	targetPckPath, isDemo, makeBackup := s.state.GetState()
	wails.LogInfo(s.ctx, fmt.Sprintf("Parâmetros de instalação - Path: %s, Demo: %t, Backup: %t", targetPckPath, isDemo, makeBackup))

	gameDir := filepath.Dir(targetPckPath)
	modifiedPckPath := filepath.Join(gameDir, "ModifiedPCK.pck")
	backupPckPath := filepath.Join(gameDir, "UntilThen.pck.bak")

	tempDir, err := os.MkdirTemp("", "untilthen_patcher_*")
	if err != nil {
		s.failAndLog(modifiedPckPath, fmt.Errorf("falha ao criar pasta temporária: %w", err))
		return
	}
	// Garante a limpeza do diretório temporário no final do processo
	defer func() {
		wails.LogInfo(s.ctx, fmt.Sprintf("Limpando diretório temporário: %s", tempDir))
		os.RemoveAll(tempDir)
	}()

	wails.LogInfo(s.ctx, fmt.Sprintf("Diretório temporário criado com sucesso em: %s", tempDir))

	wails.EventsEmit(s.ctx, "install_step", "Extraindo ferramentas de patch...")
	wails.LogInfo(s.ctx, "Iniciando extração do binário pckExplorerBinZip...")
	if err := s.unzipFromMemory(pckExplorerBinZip, tempDir, "unzip_bin_progress"); err != nil {
		s.failAndLog(modifiedPckPath, fmt.Errorf("erro ao extrair ferramentas: %w", err))
		return
	}

	binPath := filepath.Join(tempDir, pckBinName)
	wails.LogInfo(s.ctx, fmt.Sprintf("Caminho do binário definido: %s", binPath))

	if runtime.GOOS != "windows" {
		wails.LogInfo(s.ctx, "Sistema não-Windows detectado, aplicando permissões de execução (0755) ao binário.")
		os.Chmod(binPath, 0755)
	}

	wails.EventsEmit(s.ctx, "install_step", "Preparando arquivos de tradução...")
	wails.LogInfo(s.ctx, "Iniciando extração dos arquivos de tradução (translationFilesZip)...")
	if err := s.unzipFromMemory(translationFilesZip, tempDir, "unzip_trans_progress"); err != nil {
		s.failAndLog(modifiedPckPath, fmt.Errorf("erro ao extrair arquivos de tradução: %w", err))
		return
	}

	translationFolder := "full"
	if isDemo {
		translationFolder = "demo"
	}
	translationFilesPath := filepath.Join(tempDir, translationFolder)
	wails.LogInfo(s.ctx, fmt.Sprintf("Modo de tradução selecionado: %s (Caminho: %s)", translationFolder, translationFilesPath))

	wails.EventsEmit(s.ctx, "install_step", "Aplicando tradução (isso pode levar alguns instantes)...")

	wails.LogInfo(s.ctx, fmt.Sprintf("Iniciando comando: %s -pc %s %s %s 2.2.4.1", binPath, targetPckPath, translationFilesPath, modifiedPckPath))
	cmd := exec.CommandContext(s.ctx, binPath, "-pc", targetPckPath, translationFilesPath, modifiedPckPath, "2.2.4.1")

	// This prevents the console window from appearing on Windows
	cmd.SysProcAttr = getSysProcAttr()

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		s.failAndLog(modifiedPckPath, fmt.Errorf("falha ao iniciar processo de patch: %w", err))
		return
	}

	wails.LogInfo(s.ctx, "Processo de patch rodando. Aguardando conclusão...")
	go s.streamLogs(stdout, "install_log")
	go s.streamLogs(stderr, "install_error")

	if err := cmd.Wait(); err != nil {
		s.failAndLog(modifiedPckPath, fmt.Errorf("erro executando o patcher: %w", err))
		return
	}

	wails.LogInfo(s.ctx, "Processo de patch concluído com sucesso.")
	wails.EventsEmit(s.ctx, "install_step", "Finalizando instalação...")

	if makeBackup {
		wails.LogInfo(s.ctx, fmt.Sprintf("Realizando backup do PCK original para: %s", backupPckPath))
		os.Remove(backupPckPath) // Remove existing if any
		if err := os.Rename(targetPckPath, backupPckPath); err != nil {
			s.failAndLog(modifiedPckPath, fmt.Errorf("erro no backup: %w", err))
			return
		}
	} else {
		wails.LogInfo(s.ctx, "Backup desativado. Removendo arquivo PCK original...")
		os.Remove(targetPckPath)
	}

	wails.LogInfo(s.ctx, fmt.Sprintf("Renomeando PCK modificado de %s para %s", modifiedPckPath, targetPckPath))
	if err := os.Rename(modifiedPckPath, targetPckPath); err != nil {
		wails.LogError(s.ctx, fmt.Sprintf("Falha ao renomear PCK final: %v", err))
		if makeBackup {
			wails.LogInfo(s.ctx, "Tentando restaurar backup devido à falha de renomeação...")
			os.Rename(backupPckPath, targetPckPath)
		}
		s.failAndLog(modifiedPckPath, fmt.Errorf("erro ao renomear arquivo final: %w", err))
		return
	}

	wails.LogInfo(s.ctx, "Instalação finalizada com sucesso!")
	wails.EventsEmit(s.ctx, "install_success", "Sucesso!")
}

func (s *PckExplorerService) unzipFromMemory(data []byte, dest, eventName string) error {
	wails.LogInfo(s.ctx, fmt.Sprintf("Descompactando %d bytes para %s...", len(data), dest))
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
			wails.LogError(s.ctx, fmt.Sprintf("Erro ao criar estrutura de diretório para %s: %v", fpath, err))
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			wails.LogError(s.ctx, fmt.Sprintf("Erro ao preparar arquivo %s para escrita: %v", fpath, err))
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			wails.LogError(s.ctx, fmt.Sprintf("Erro ao ler arquivo do zip %s: %v", f.Name, err))
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			wails.LogError(s.ctx, fmt.Sprintf("Erro ao copiar os dados extraídos para %s: %v", fpath, err))
			return err
		}
	}

	wails.LogInfo(s.ctx, fmt.Sprintf("Descompactação de %d arquivos concluída.", total))
	return nil
}

func (s *PckExplorerService) failAndLog(modifiedPck string, err error) {
	wails.LogError(s.ctx, fmt.Sprintf("Falha Crítica na instalação: %v", err))
	wails.EventsEmit(s.ctx, "install_error", err.Error())

	if modifiedPck != "" {
		wails.LogInfo(s.ctx, fmt.Sprintf("Excluindo ModifiedPCK parcial gerado por conta da falha: %s", modifiedPck))
		os.Remove(modifiedPck)
	}

	wails.MessageDialog(s.ctx, wails.MessageDialogOptions{
		Type:    wails.ErrorDialog,
		Title:   "Erro durante a Instalação",
		Message: fmt.Sprintf("Um erro inesperado aconteceu durante a instalação\nUm arquivo de log foi criado em: %s", GetLogFilePath()),
	})
}

func (s *PckExplorerService) streamLogs(pipe io.ReadCloser, eventName string) {
	defer pipe.Close()
	scanner := bufio.NewScanner(pipe)

	var lastLine string
	hasNewContent := false

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	go func() {
		for scanner.Scan() {
			lastLine = scanner.Text()
			hasNewContent = true
		}
	}()

	for {
		select {
		case <-ticker.C:
			if hasNewContent {
				wails.EventsEmit(s.ctx, eventName, lastLine)
				hasNewContent = false
			}
		case <-s.ctx.Done():
			return
		}
	}
}
