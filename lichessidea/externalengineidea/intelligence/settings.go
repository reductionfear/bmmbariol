package intelligence

// IntelligenceSettings contains all intelligence feature settings
// Ported from Python models.py IntelligenceSettings
type IntelligenceSettings struct {
	// Master toggle
	IntelligenceEnabled bool `json:"intelligence_enabled"`

	// Avoidance control
	AvoidLowIntelligence     bool    `json:"avoid_low_intelligence"`
	LowIntelligenceThreshold float64 `json:"low_intelligence_threshold"` // -3.0 to -1.0

	// Move multipliers
	AggressivenessContempt float64 `json:"aggressiveness_contempt"`
	PassivenessContempt    float64 `json:"passiveness_contempt"`
	TradingPreference      float64 `json:"trading_preference"`
	CapturePreference      float64 `json:"capture_preference"`
	CastlePreference       float64 `json:"castle_preference"`
	EnPassantPreference    float64 `json:"en_passant_preference"`
	PromotionPreference    float64 `json:"promotion_preference"`

	// Boolean preferences
	PreferEarlyCastling bool    `json:"prefer_early_castling"`
	PreferPins          bool    `json:"prefer_pins"`
	PreferSideCastle    bool    `json:"prefer_side_castle"`
	CastleSide          *string `json:"castle_side"` // "kingside" or "queenside"

	// Piece movement preferences (multipliers)
	PawnPreference   float64 `json:"pawn_preference"`
	KnightPreference float64 `json:"knight_preference"`
	BishopPreference float64 `json:"bishop_preference"`
	RookPreference   float64 `json:"rook_preference"`
	QueenPreference  float64 `json:"queen_preference"`
	KingPreference   float64 `json:"king_preference"`

	// Special behaviors
	StayEqual            bool    `json:"stay_equal"`
	StalemateProbability float64 `json:"stalemate_probability"`
	AlwaysPromoteQueen   bool    `json:"always_promote_queen"`
	CheckmateImmediately bool    `json:"checkmate_immediately"`
}

// NewDefaultIntelligenceSettings returns default intelligence settings
func NewDefaultIntelligenceSettings() *IntelligenceSettings {
	return &IntelligenceSettings{
		IntelligenceEnabled:      false,
		AvoidLowIntelligence:     false,
		LowIntelligenceThreshold: -1.5,
		AggressivenessContempt:   1.0,
		PassivenessContempt:      1.0,
		TradingPreference:        0.0,
		CapturePreference:        1.0,
		CastlePreference:         1.0,
		EnPassantPreference:      1.0,
		PromotionPreference:      1.0,
		PreferEarlyCastling:      false,
		PreferPins:               false,
		PreferSideCastle:         false,
		CastleSide:               nil,
		PawnPreference:           1.0,
		KnightPreference:         1.0,
		BishopPreference:         1.0,
		RookPreference:           1.0,
		QueenPreference:          1.0,
		KingPreference:           1.0,
		StayEqual:                false,
		StalemateProbability:     0.0,
		AlwaysPromoteQueen:       false,
		CheckmateImmediately:     false,
	}
}

// IsFullyDisabled checks if intelligence is completely disabled
func (s *IntelligenceSettings) IsFullyDisabled() bool {
	return !s.IntelligenceEnabled
}

// ShouldAvoidLowIntelligence checks if low intelligence avoidance is enabled
func (s *IntelligenceSettings) ShouldAvoidLowIntelligence() bool {
	return s.IntelligenceEnabled && s.AvoidLowIntelligence
}

// GetClampedThreshold returns the threshold clamped to valid range
func (s *IntelligenceSettings) GetClampedThreshold() float64 {
	if s.LowIntelligenceThreshold < -3.0 {
		return -3.0
	}
	if s.LowIntelligenceThreshold > -1.0 {
		return -1.0
	}
	return s.LowIntelligenceThreshold
}
