package models

import (
	"time"
)

// GameState tracks the complete server-side game state
// Ported from Python models.py GameState
type GameState struct {
	// Board state
	CurrentFEN      string   `json:"current_fen"`
	PositionHistory []string `json:"position_history"`

	// Game information
	MoveNumber int    `json:"move_number"`
	Turn       string `json:"turn"`       // 'w' or 'b'
	GamePhase  string `json:"game_phase"` // "opening", "middlegame", "endgame"
	PieceCount int    `json:"piece_count"`

	// Game status
	IsCheck     bool   `json:"is_check"`
	IsCheckmate bool   `json:"is_checkmate"`
	IsStalemate bool   `json:"is_stalemate"`
	IsDraw      bool   `json:"is_draw"`
	GameResult  string `json:"game_result"` // "*", "1-0", "0-1", "1/2-1/2"

	// Analysis state
	CurrentEvaluation float64   `json:"current_evaluation"` // In pawns
	EvaluationHistory []float64 `json:"evaluation_history"`
	LastMove          string    `json:"last_move"`
	LastMoveQuality   string    `json:"last_move_quality"`

	// Intelligence tracking
	IntelligenceActive       bool   `json:"intelligence_active"`
	LastIntelligenceDecision string `json:"last_intelligence_decision"` // "used", "avoided", "disabled"
	IntelligenceAvoidedCount int    `json:"intelligence_avoided_count"`
	TotalMoveCount           int    `json:"total_move_count"`

	// Timing state
	MoveTimes      []float64 `json:"move_times"` // Move times in seconds
	TotalGameTime  float64   `json:"total_game_time"`
	LastMoveTime   float64   `json:"last_move_time"`

	// Opening state
	OpeningName   string `json:"opening_name"`
	MovesFromBook int    `json:"moves_from_book"`
	StillInBook   bool   `json:"still_in_book"`

	// Server state
	LastUpdateTime time.Time `json:"last_update_time"`
	AnalysisActive bool      `json:"analysis_active"`
	WaitingForMove bool      `json:"waiting_for_move"`
}

// NewGameState creates a new game state with default starting position
func NewGameState() *GameState {
	return &GameState{
		CurrentFEN:               "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		PositionHistory:          make([]string, 0),
		MoveNumber:               1,
		Turn:                     "w",
		GamePhase:                "opening",
		PieceCount:               32,
		IsCheck:                  false,
		IsCheckmate:              false,
		IsStalemate:              false,
		IsDraw:                   false,
		GameResult:               "*",
		CurrentEvaluation:        0.0,
		EvaluationHistory:        make([]float64, 0),
		LastMove:                 "",
		LastMoveQuality:          "Unknown",
		IntelligenceActive:       false,
		LastIntelligenceDecision: "",
		IntelligenceAvoidedCount: 0,
		TotalMoveCount:           0,
		MoveTimes:                make([]float64, 0),
		TotalGameTime:            0.0,
		LastMoveTime:             0.0,
		OpeningName:              "Unknown",
		MovesFromBook:            0,
		StillInBook:              true,
		LastUpdateTime:           time.Now(),
		AnalysisActive:           false,
		WaitingForMove:           false,
	}
}

// TrackIntelligenceDecision tracks intelligence decision for statistics
func (g *GameState) TrackIntelligenceDecision(decision string) {
	g.LastIntelligenceDecision = decision
	g.TotalMoveCount++

	if decision == "avoided" {
		g.IntelligenceAvoidedCount++
	}
}

// GetIntelligenceUsageRate returns percentage of moves where intelligence was used
func (g *GameState) GetIntelligenceUsageRate() float64 {
	if g.TotalMoveCount == 0 {
		return 0.0
	}

	intelligenceUsed := g.TotalMoveCount - g.IntelligenceAvoidedCount
	return (float64(intelligenceUsed) / float64(g.TotalMoveCount)) * 100.0
}

// AddEvaluationPoint adds an evaluation point to history
func (g *GameState) AddEvaluationPoint(evaluation float64) {
	g.CurrentEvaluation = evaluation
	g.EvaluationHistory = append(g.EvaluationHistory, evaluation)

	// Limit history size
	if len(g.EvaluationHistory) > 200 {
		g.EvaluationHistory = g.EvaluationHistory[len(g.EvaluationHistory)-100:]
	}
}

// AddMoveTime adds move time to history
func (g *GameState) AddMoveTime(moveTime float64) {
	g.LastMoveTime = moveTime
	g.MoveTimes = append(g.MoveTimes, moveTime)
	g.TotalGameTime += moveTime

	// Limit history size
	if len(g.MoveTimes) > 100 {
		g.MoveTimes = g.MoveTimes[len(g.MoveTimes)-50:]
	}
}
