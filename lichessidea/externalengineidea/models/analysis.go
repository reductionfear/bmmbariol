package models

import (
	"time"
)

// AnalysisSession tracks a single analysis session with depth progression
// Ported from Python models.py AnalysisSession
type AnalysisSession struct {
	TargetDepth  int              `json:"target_depth"`
	CurrentDepth int              `json:"current_depth"`
	StartTime    time.Time        `json:"start_time"`
	Candidates   []*MoveCandidate `json:"candidates"`
	NodesSearched int64           `json:"nodes_searched"`
	IsComplete   bool             `json:"is_complete"`
	MateFound    bool             `json:"mate_found"`
	BestMove     string           `json:"best_move"`

	// Intelligence tracking
	IntelligenceApplied       bool   `json:"intelligence_applied"`
	IntelligenceAvoided       bool   `json:"intelligence_avoided"`
	OriginalBestMove          string `json:"original_best_move"`
	IntelligenceDecisionReason string `json:"intelligence_decision_reason"`
}

// NewAnalysisSession creates a new analysis session
func NewAnalysisSession(targetDepth int) *AnalysisSession {
	return &AnalysisSession{
		TargetDepth:                targetDepth,
		CurrentDepth:               0,
		StartTime:                  time.Now(),
		Candidates:                 make([]*MoveCandidate, 0),
		NodesSearched:              0,
		IsComplete:                 false,
		MateFound:                  false,
		BestMove:                   "",
		IntelligenceApplied:        false,
		IntelligenceAvoided:        false,
		OriginalBestMove:           "",
		IntelligenceDecisionReason: "",
	}
}

// UpdateProgress updates analysis progress
func (s *AnalysisSession) UpdateProgress(depth int, candidates []*MoveCandidate) {
	if depth > s.CurrentDepth {
		s.CurrentDepth = depth
	}
	if candidates != nil {
		s.Candidates = candidates
	}

	// Check for completion
	if s.CurrentDepth >= s.TargetDepth {
		s.IsComplete = true
	}
}

// SetIntelligenceDecision tracks intelligence decision for this analysis
func (s *AnalysisSession) SetIntelligenceDecision(applied, avoided bool, reason, originalMove string) {
	s.IntelligenceApplied = applied
	s.IntelligenceAvoided = avoided
	s.IntelligenceDecisionReason = reason
	s.OriginalBestMove = originalMove
}

// GetProgressPercentage returns progress as percentage
func (s *AnalysisSession) GetProgressPercentage() float64 {
	if s.TargetDepth == 0 {
		return 0.0
	}
	progress := (float64(s.CurrentDepth) / float64(s.TargetDepth)) * 100.0
	if progress > 100.0 {
		return 100.0
	}
	return progress
}

// GetElapsedTime returns elapsed time in seconds
func (s *AnalysisSession) GetElapsedTime() float64 {
	return time.Since(s.StartTime).Seconds()
}
