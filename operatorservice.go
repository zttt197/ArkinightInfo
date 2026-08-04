package main

import (
	"sync"

	"arkinightinfo/data"
)

// OperatorService exposes operator data and update functions to the frontend.
type OperatorService struct {
	mu       sync.Mutex
	dataRoot string
	ops      []data.Operator
	version  string
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

// CheckUpdate checks if a newer data version is available on GitHub.
func (s *OperatorService) CheckUpdate() *data.CheckResult {
	result, err := data.CheckUpdate(s.dataRoot)
	if err != nil {
		return &data.CheckResult{HasUpdate: false, Message: "检查更新失败: " + err.Error()}
	}
	return result
}

// DoUpdate downloads the latest data and reloads operators.
func (s *OperatorService) DoUpdate() error {
	err := data.Update(s.dataRoot, func(msg string) {
		// Progress will be logged; frontend can poll GetDataVersion
	})
	if err != nil {
		return err
	}
	return nil
}
