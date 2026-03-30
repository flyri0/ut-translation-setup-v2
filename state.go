package main

import "sync"

type InstallerState struct {
	mu         sync.Mutex
	TargetPath string
	isDemo     bool
	MakeBackup bool
}

func NewInstallerState() *InstallerState {
	return &InstallerState{}
}

func (s *InstallerState) SetState(path string, isDemo bool, makeBackup bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.TargetPath = path
	s.isDemo = isDemo
	s.MakeBackup = makeBackup
}

func (s *InstallerState) GetState() (string, bool, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.TargetPath, s.isDemo, s.MakeBackup
}
