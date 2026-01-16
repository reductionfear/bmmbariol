package chess

import (
	"github.com/notnil/chess"
)

// ThreatDetector provides threat and pin detection functionality
type ThreatDetector struct{}

// NewThreatDetector creates a new threat detector
func NewThreatDetector() *ThreatDetector {
	return &ThreatDetector{}
}

// Threat represents a threat on the board
type Threat struct {
	FromSquare   chess.Square
	ToSquare     chess.Square
	AttackingPiece chess.Piece
	TargetPiece    chess.Piece
	ThreatType     string // "capture", "check", "pin"
	Value          int    // Threat value in centipawns
}

// DetectThreats detects all threats in the current position
func (d *ThreatDetector) DetectThreats(game *chess.Game) []Threat {
	threats := make([]Threat, 0)
	
	// Get all valid moves for current side
	for _, move := range game.ValidMoves() {
		board := game.Position().Board()
		fromPiece := board.Piece(move.S1())
		toPiece := board.Piece(move.S2())
		
		// Check for capture threats
		if toPiece != chess.NoPiece {
			threat := Threat{
				FromSquare:     move.S1(),
				ToSquare:       move.S2(),
				AttackingPiece: fromPiece,
				TargetPiece:    toPiece,
				ThreatType:     "capture",
				Value:          d.getPieceValue(toPiece.Type()),
			}
			threats = append(threats, threat)
		}
		
		// Check for check threats
		gameCopy := game.Copy()
		if err := gameCopy.Move(move); err == nil {
			if gameCopy.Position().InCheck() {
				threat := Threat{
					FromSquare:     move.S1(),
					ToSquare:       move.S2(),
					AttackingPiece: fromPiece,
					ThreatType:     "check",
					Value:          1000, // Check is highly valuable tactically
				}
				threats = append(threats, threat)
			}
		}
	}
	
	return threats
}

// DetectPins detects all pins in the current position
func (d *ThreatDetector) DetectPins(game *chess.Game) []Pin {
	pins := make([]Pin, 0)
	
	// This is a simplified pin detection
	// A full implementation would check for pieces that are pinned to the king
	// or other valuable pieces along ranks, files, and diagonals
	
	board := game.Position().Board()
	currentTurn := game.Position().Turn()
	
	// Find the king
	var kingSquare chess.Square
	for sq := chess.A1; sq <= chess.H8; sq++ {
		piece := board.Piece(sq)
		if piece.Type() == chess.King && piece.Color() == currentTurn {
			kingSquare = sq
			break
		}
	}
	
	// Check for pieces pinned to the king
	// This is a placeholder - a full implementation would check all directions
	_ = kingSquare // Use the king square to detect pins
	
	return pins
}

// Pin represents a pinned piece
type Pin struct {
	PinnedSquare   chess.Square
	PinnedPiece    chess.Piece
	PinningSquare  chess.Square
	PinningPiece   chess.Piece
	TargetSquare   chess.Square
	TargetPiece    chess.Piece
	Direction      string // "vertical", "horizontal", "diagonal"
}

// IsSquarePinned checks if a piece on a square is pinned
func (d *ThreatDetector) IsSquarePinned(game *chess.Game, square chess.Square) bool {
	pins := d.DetectPins(game)
	for _, pin := range pins {
		if pin.PinnedSquare == square {
			return true
		}
	}
	return false
}

// CreatesPin checks if a move creates a pin
func (d *ThreatDetector) CreatesPin(game *chess.Game, move *chess.Move) bool {
	// Make the move on a copy
	gameCopy := game.Copy()
	if err := gameCopy.Move(move); err != nil {
		return false
	}
	
	// Check if any new pins were created
	// This is simplified - a full implementation would compare pins before and after
	pins := d.DetectPins(gameCopy)
	return len(pins) > 0
}

// getPieceValue returns the centipawn value of a piece
func (d *ThreatDetector) getPieceValue(pieceType chess.PieceType) int {
	values := map[chess.PieceType]int{
		chess.Pawn:   100,
		chess.Knight: 320,
		chess.Bishop: 330,
		chess.Rook:   500,
		chess.Queen:  900,
		chess.King:   0,
	}
	return values[pieceType]
}

// GetHighestValueThreat returns the threat with the highest value
func (d *ThreatDetector) GetHighestValueThreat(threats []Threat) *Threat {
	if len(threats) == 0 {
		return nil
	}
	
	highest := &threats[0]
	for i := range threats {
		if threats[i].Value > highest.Value {
			highest = &threats[i]
		}
	}
	
	return highest
}

// FilterThreatsByType filters threats by type
func (d *ThreatDetector) FilterThreatsByType(threats []Threat, threatType string) []Threat {
	filtered := make([]Threat, 0)
	for _, threat := range threats {
		if threat.ThreatType == threatType {
			filtered = append(filtered, threat)
		}
	}
	return filtered
}

// CountThreats counts the number of threats
func (d *ThreatDetector) CountThreats(threats []Threat) int {
	return len(threats)
}
