package intelligence

import (
	"github.com/notnil/chess"
)

// MultiplierCalculator calculates multipliers for pieces and moves
type MultiplierCalculator struct {
	settings *IntelligenceSettings
}

// NewMultiplierCalculator creates a new multiplier calculator
func NewMultiplierCalculator(settings *IntelligenceSettings) *MultiplierCalculator {
	return &MultiplierCalculator{
		settings: settings,
	}
}

// GetPieceMultiplier returns the multiplier for a specific piece type
func (m *MultiplierCalculator) GetPieceMultiplier(pieceType chess.PieceType) float64 {
	switch pieceType {
	case chess.Pawn:
		return m.settings.PawnPreference
	case chess.Knight:
		return m.settings.KnightPreference
	case chess.Bishop:
		return m.settings.BishopPreference
	case chess.Rook:
		return m.settings.RookPreference
	case chess.Queen:
		return m.settings.QueenPreference
	case chess.King:
		return m.settings.KingPreference
	default:
		return 1.0
	}
}

// GetCaptureMultiplier returns the multiplier for captures
func (m *MultiplierCalculator) GetCaptureMultiplier() float64 {
	return m.settings.CapturePreference
}

// GetCastlingMultiplier returns the multiplier for castling moves
func (m *MultiplierCalculator) GetCastlingMultiplier(isKingside bool) float64 {
	baseMultiplier := m.settings.CastlePreference

	// Apply side preference if enabled
	if m.settings.PreferSideCastle && m.settings.CastleSide != nil {
		if *m.settings.CastleSide == "kingside" && isKingside {
			baseMultiplier *= 1.2
		} else if *m.settings.CastleSide == "queenside" && !isKingside {
			baseMultiplier *= 1.2
		}
	}

	return baseMultiplier
}

// GetPromotionMultiplier returns the multiplier for pawn promotions
func (m *MultiplierCalculator) GetPromotionMultiplier(promoPiece chess.PieceType) float64 {
	baseMultiplier := m.settings.PromotionPreference

	// If always promote to queen, boost queen promotions
	if m.settings.AlwaysPromoteQueen && promoPiece == chess.Queen {
		baseMultiplier *= 1.5
	}

	return baseMultiplier
}

// GetEnPassantMultiplier returns the multiplier for en passant captures
func (m *MultiplierCalculator) GetEnPassantMultiplier() float64 {
	return m.settings.EnPassantPreference
}

// GetAggressivenessMultiplier returns the multiplier for aggressive moves
func (m *MultiplierCalculator) GetAggressivenessMultiplier() float64 {
	return m.settings.AggressivenessContempt
}

// GetPassivenessMultiplier returns the multiplier for passive moves
func (m *MultiplierCalculator) GetPassivenessMultiplier() float64 {
	return m.settings.PassivenessContempt
}

// GetTradingMultiplier returns the multiplier for trading pieces
func (m *MultiplierCalculator) GetTradingMultiplier() float64 {
	return m.settings.TradingPreference
}

// CalculateTotalMultiplier calculates the total multiplier for a move
func (m *MultiplierCalculator) CalculateTotalMultiplier(
	pieceType chess.PieceType,
	isCapture, isCastling, isEnPassant, isPromotion, isCheck bool,
	isKingside bool,
	promoPiece chess.PieceType,
) float64 {
	multiplier := 1.0

	// Apply piece multiplier
	multiplier *= m.GetPieceMultiplier(pieceType)

	// Apply move characteristic multipliers
	if isCapture {
		multiplier *= m.GetCaptureMultiplier()
	}

	if isCastling {
		multiplier *= m.GetCastlingMultiplier(isKingside)
	}

	if isEnPassant {
		multiplier *= m.GetEnPassantMultiplier()
	}

	if isPromotion {
		multiplier *= m.GetPromotionMultiplier(promoPiece)
	}

	if isCheck {
		multiplier *= m.GetAggressivenessMultiplier()
	}

	return multiplier
}

// ApplyMultiplierReduction reduces multiplier effect in critical positions
func (m *MultiplierCalculator) ApplyMultiplierReduction(multiplier float64, isCritical bool) float64 {
	if !isCritical {
		return multiplier
	}

	// Reduce multiplier effect by 70% in critical positions
	if multiplier > 1.0 {
		return 1.0 + (multiplier-1.0)*0.3
	}
	return 1.0 - (1.0-multiplier)*0.3
}
