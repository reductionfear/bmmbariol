package intelligence

import (
	"github.com/notnil/chess"
	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/models"
)

// Analyzer analyzes chess positions
type Analyzer struct {
	settings *IntelligenceSettings
}

// NewAnalyzer creates a new position analyzer
func NewAnalyzer(settings *IntelligenceSettings) *Analyzer {
	return &Analyzer{
		settings: settings,
	}
}

// AnalyzeMoveCharacteristics analyzes move characteristics
func (a *Analyzer) AnalyzeMoveCharacteristics(candidate *models.MoveCandidate, game *chess.Game) {
	move, err := chess.UCINotation{}.Decode(game.Position(), candidate.Move)
	if err != nil {
		return
	}

	// Check if capture
	candidate.IsCapture = game.Position().Board().Piece(move.S2()) != chess.NoPiece

	// Check if castling
	candidate.IsCastling = move.HasTag(chess.KingSideCastle) || move.HasTag(chess.QueenSideCastle)

	// Check if en passant
	candidate.IsEnPassant = move.HasTag(chess.EnPassant)

	// Check if promotion
	candidate.IsPromotion = move.Promo() != chess.NoPieceType

	// Make move on copy to check for check
	gameCopy := game.Clone()
	if err := gameCopy.Move(move); err == nil {
		// Check if opponent is in check after the move
		candidate.IsCheck = len(gameCopy.ValidMoves()) > 0 && gameCopy.Method() == chess.NoMethod && gameCopy.Outcome() == chess.NoOutcome

		// Determine if aggressive or passive
		if candidate.IsCapture || candidate.IsCheck {
			candidate.IsAggressive = true
			candidate.IsTactical = true
		} else {
			candidate.IsPassive = true
			candidate.IsPositional = true
		}
	}
}

// IsCheckmate checks if a position is checkmate
func (a *Analyzer) IsCheckmate(game *chess.Game) bool {
	return game.Outcome() == chess.WhiteWon || game.Outcome() == chess.BlackWon
}

// IsCriticalPosition determines if a position requires careful analysis
func (a *Analyzer) IsCriticalPosition(game *chess.Game) bool {
	// Very few legal moves (forced positions)
	if len(game.ValidMoves()) < 3 {
		return true
	}

	return false
}

// GetGamePhase determines the current game phase
func (a *Analyzer) GetGamePhase(game *chess.Game) string {
	// Count pieces
	pieceCount := 0
	boardMap := game.Position().Board()
	for sq := chess.A1; sq <= chess.H8; sq++ {
		if boardMap.Piece(sq) != chess.NoPiece {
			pieceCount++
		}
	}

	// Determine phase based on moves and piece count
	moveHistory := len(game.Moves())
	if moveHistory <= 20 {
		return "opening"
	} else if pieceCount <= 12 {
		return "endgame"
	}
	return "middlegame"
}

// CountPieces counts the number of pieces on the board
func (a *Analyzer) CountPieces(game *chess.Game) int {
	count := 0
	boardMap := game.Position().Board()
	for sq := chess.A1; sq <= chess.H8; sq++ {
		if boardMap.Piece(sq) != chess.NoPiece {
			count++
		}
	}
	return count
}
