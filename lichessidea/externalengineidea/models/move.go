package models

// MoveCandidate represents a detailed move candidate with analysis data
// Ported from Python models.py MoveCandidate
type MoveCandidate struct {
	Move       string  `json:"move"`        // UCI format (e.g., "e2e4")
	FromSquare string  `json:"from_square"` // e.g., "e2"
	ToSquare   string  `json:"to_square"`   // e.g., "e4"
	ScoreCP    int     `json:"score_cp"`    // Centipawns
	ScorePawns float64 `json:"score_pawns"` // Pawns (cp/100)
	Depth      int     `json:"depth"`
	MateIn     *int    `json:"mate_in,omitempty"`
	PVLine     []string `json:"pv_line"`
	Nodes      int64   `json:"nodes"`

	// Move characteristics
	IsAggressive bool    `json:"is_aggressive"`
	IsPassive    bool    `json:"is_passive"`
	IsTactical   bool    `json:"is_tactical"`
	IsPositional bool    `json:"is_positional"`
	IsCapture    bool    `json:"is_capture"`
	IsCheck      bool    `json:"is_check"`
	IsBookMove   bool    `json:"is_book_move"`
	IsCastling   bool    `json:"is_castling"`
	IsEnPassant  bool    `json:"is_en_passant"`
	IsPromotion  bool    `json:"is_promotion"`
	CreatesPin   bool    `json:"creates_pin"`
	IsTrade      bool    `json:"is_trade"`
	TradeValue   float64 `json:"trade_value"`

	// Intelligence tracking
	OriginalScorePawns     float64 `json:"original_score_pawns"`
	IntelligenceModified   bool    `json:"intelligence_modified"`
	IntelligenceMultiplier float64 `json:"intelligence_multiplier"`

	// Quality assessment
	QualityScore float64 `json:"quality_score"` // 0-100
	QualityLabel string  `json:"quality_label"` // "Excellent", "Good", etc.
}

// NewMoveCandidate creates a new move candidate with default values
func NewMoveCandidate(move, fromSquare, toSquare string) *MoveCandidate {
	return &MoveCandidate{
		Move:                   move,
		FromSquare:             fromSquare,
		ToSquare:               toSquare,
		ScoreCP:                0,
		ScorePawns:             0.0,
		Depth:                  0,
		MateIn:                 nil,
		PVLine:                 make([]string, 0),
		Nodes:                  0,
		IsAggressive:           false,
		IsPassive:              false,
		IsTactical:             false,
		IsPositional:           false,
		IsCapture:              false,
		IsCheck:                false,
		IsBookMove:             false,
		IsCastling:             false,
		IsEnPassant:            false,
		IsPromotion:            false,
		CreatesPin:             false,
		IsTrade:                false,
		TradeValue:             0.0,
		OriginalScorePawns:     0.0,
		IntelligenceModified:   false,
		IntelligenceMultiplier: 1.0,
		QualityScore:           0.0,
		QualityLabel:           "Unknown",
	}
}

// ApplyIntelligenceModification applies intelligence modification to this candidate
func (m *MoveCandidate) ApplyIntelligenceModification(newScorePawns, multiplier float64) {
	// Store original if not already stored
	if !m.IntelligenceModified {
		m.OriginalScorePawns = m.ScorePawns
	}
	m.ScorePawns = newScorePawns
	m.ScoreCP = int(newScorePawns * 100)
	m.IntelligenceModified = true
	m.IntelligenceMultiplier = multiplier
}
