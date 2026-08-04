package main

import (
	"sync"

	"arkinightinfo/data"
)

// OperatorService exposes operator data and update functions to the frontend.
type OperatorService struct {
	mu         sync.Mutex
	dataRoot   string
	ops        []data.Operator
	version    string
	progress   string
	onProgress func(string)
}

// SetProgressCallback sets the function called during download to push progress to the frontend.
func (s *OperatorService) SetProgressCallback(fn func(string)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onProgress = fn
}

func (s *OperatorService) emitProgress(msg string) {
	s.mu.Lock()
	s.progress = msg
	cb := s.onProgress
	s.mu.Unlock()
	if cb != nil {
		cb(msg)
	}
}

// LoadOperators loads or reloads operator data from disk.
func (s *OperatorService) LoadOperators() ([]data.Operator, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	game, err := data.Load(s.dataRoot)
	if err != nil {
		return nil, err
	}
	s.ops = game.Ops
	s.version = game.Version
	return s.ops, nil
}

// GetDataVersion returns the current data version string.
func (s *OperatorService) GetDataVersion() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.version
}

// GetProgress returns the current download progress message.
func (s *OperatorService) GetProgress() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.progress
}

// CheckUpdate checks if a newer data version is available on GitHub.
func (s *OperatorService) CheckUpdate() *data.CheckResult {
	result, err := data.CheckUpdate(s.dataRoot)
	if err != nil {
		return &data.CheckResult{HasUpdate: false, Message: "检查更新失败: " + err.Error()}
	}
	return result
}

// DoUpdate downloads the latest data in a goroutine and pushes progress via callback.
func (s *OperatorService) DoUpdate() {
	s.emitProgress("正在下载数据…")
	go func() {
		err := data.Update(s.dataRoot, func(msg string) {
			s.emitProgress(msg)
		})
		if err != nil {
			s.emitProgress("下载失败: " + err.Error())
		} else {
			s.emitProgress("下载完成，正在加载…")
		}
	}()
}
