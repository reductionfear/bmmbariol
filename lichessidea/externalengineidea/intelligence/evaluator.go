package intelligence

import (
	"github.com/notnil/chess"
	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/models"
	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/utils"
)

// Evaluator applies intelligence modifications to move candidates
type Evaluator struct {
	settings *IntelligenceSettings
}

// NewEvaluator creates a new intelligence evaluator
func NewEvaluator(settings *IntelligenceSettings) *Evaluator {
	return &Evaluator{
		settings: settings,
	}
}

// ApplyIntelligence applies intelligence modifications to move candidates
func (e *Evaluator) ApplyIntelligence(candidates []*models.MoveCandidate, game *chess.Game) []*models.MoveCandidate {
	// Check if intelligence is disabled
	if e.settings.IsFullyDisabled() {
		return candidates
	}

	// Store original candidates for threshold checking
	originalCandidates := e.cloneCandidates(candidates)

	// Apply multipliers to each candidate
	for i, candidate := range candidates {
		// Skip checkmate moves - they get absolute evaluation
		if utils.IsForcedMate(candidate.MateIn) {
			candidate.ScorePawns = 1000.0
			if candidate.MateIn != nil && *candidate.MateIn > 0 {
				candidate.ScorePawns = 1000.0
			} else {
				candidate.ScorePawns = -1000.0
			}
			candidate.ScoreCP = int(candidate.ScorePawns * 100)
			continue
		}

		// Apply multipliers for non-mate moves
		candidates[i] = e.applyMultipliers(candidate, game)
	}

	// Check avoidance threshold
	if e.settings.ShouldAvoidLowIntelligence() && len(candidates) > 0 {
		bestScore := candidates[0].ScorePawns
		threshold := e.settings.GetClampedThreshold()

		if bestScore <= threshold {
			utils.Logger.Infof("Intelligence avoided: best score %.2f <= threshold %.2f", bestScore, threshold)
			return originalCandidates
		}
	}

	return candidates
}

// applyMultipliers applies multipliers to a single candidate
func (e *Evaluator) applyMultipliers(candidate *models.MoveCandidate, game *chess.Game) *models.MoveCandidate {
	totalMultiplier := 1.0

	// Determine if position is critical (forced moves with few options)
	isCritical := len(game.ValidMoves()) < 3

	// Apply piece preference multipliers
	totalMultiplier *= e.getPieceMultiplier(candidate)

	// Apply move characteristic multipliers
	totalMultiplier *= e.getMoveCharacteristicMultiplier(candidate)

	// Apply the multiplier to the score
	newScore := applyChessMultiplier(candidate.ScorePawns, totalMultiplier, isCritical)

	// Update candidate
	candidate.ApplyIntelligenceModification(newScore, totalMultiplier)

	return candidate
}

// getPieceMultiplier returns the multiplier based on the piece being moved
func (e *Evaluator) getPieceMultiplier(candidate *models.MoveCandidate) float64 {
	// Parse move to determine piece type (this is simplified)
	// In a full implementation, we'd need board context
	
	// For now, return 1.0 as we need board context to determine piece type
	// This would be enhanced with proper board analysis
	return 1.0
}

// getMoveCharacteristicMultiplier returns multiplier based on move characteristics
func (e *Evaluator) getMoveCharacteristicMultiplier(candidate *models.MoveCandidate) float64 {
	multiplier := 1.0

	// Capture preference
	if candidate.IsCapture {
		multiplier *= e.settings.CapturePreference
	}

	// Castling preference
	if candidate.IsCastling {
		multiplier *= e.settings.CastlePreference
	}

	// En passant preference
	if candidate.IsEnPassant {
		multiplier *= e.settings.EnPassantPreference
	}

	// Promotion preference
	if candidate.IsPromotion {
		multiplier *= e.settings.PromotionPreference
	}

	// Check preference (tactical)
	if candidate.IsCheck {
		multiplier *= e.settings.AggressivenessContempt
	}

	return multiplier
}

// applyChessMultiplier applies a multiplier to a chess evaluation
func applyChessMultiplier(eval, multiplier float64, isCritical bool) float64 {
	if multiplier == 1.0 || eval == 0.0 {
		return eval
	}

	// Reduce effect in critical positions (check, forced moves)
	if isCritical {
		if multiplier > 1.0 {
			multiplier = 1.0 + (multiplier-1.0)*0.3
		} else {
			multiplier = 1.0 - (1.0-multiplier)*0.3
		}
	}

	// Handle sign correctly
	if eval > 0 {
		return eval * multiplier
	}
	return eval / multiplier
}

// cloneCandidates creates a deep copy of candidates
func (e *Evaluator) cloneCandidates(candidates []*models.MoveCandidate) []*models.MoveCandidate {
	cloned := make([]*models.MoveCandidate, len(candidates))
	for i, c := range candidates {
		copy := *c
		cloned[i] = &copy
	}
	return cloned
}
