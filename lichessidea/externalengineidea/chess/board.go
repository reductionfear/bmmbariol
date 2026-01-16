package chess

import (
	"github.com/notnil/chess"
)

// BoardUtils provides utility functions for chess board operations
type BoardUtils struct{}

// NewBoardUtils creates a new board utilities instance
func NewBoardUtils() *BoardUtils {
	return &BoardUtils{}
}

// ParseFEN creates a board from FEN notation
func (b *BoardUtils) ParseFEN(fen string) (*chess.Game, error) {
	fenFunc, err := chess.FEN(fen)
	if err != nil {
		return nil, err
	}
	return chess.NewGame(fenFunc), nil
}

// GetPieceAt returns the piece at a square
func (b *BoardUtils) GetPieceAt(game *chess.Game, square chess.Square) chess.Piece {
	return game.Position().Board().Piece(square)
}

// IsSquareAttacked checks if a square is under attack
func (b *BoardUtils) IsSquareAttacked(game *chess.Game, square chess.Square, byColor chess.Color) bool {
	// Get all moves for the attacking color
	// This is a simplified check - a full implementation would be more thorough
	pos := game.Position()
	
	// Create a temporary board with the opposite turn to check attacks
	// In the notnil/chess library, we can check if a square is attacked by examining moves
	for _, move := range game.ValidMoves() {
		if pos.Turn() == byColor && move.S2() == square {
			// This square is a valid destination, meaning it's "attacked"
			return true
		}
	}
	
	return false
}

// GetSquareFromName converts square name (e.g., "e4") to Square
func (b *BoardUtils) GetSquareFromName(name string) (chess.Square, error) {
	// Parse square name
	if len(name) != 2 {
		return chess.NoSquare, chess.ErrInvalidNotation
	}
	
	file := int(name[0] - 'a')
	rank := int(name[1] - '1')
	
	if file < 0 || file > 7 || rank < 0 || rank > 7 {
		return chess.NoSquare, chess.ErrInvalidNotation
	}
	
	return chess.Square(rank*8 + file), nil
}

// GetSquareName converts Square to name (e.g., "e4")
func (b *BoardUtils) GetSquareName(square chess.Square) string {
	if square == chess.NoSquare || square < 0 || square > 63 {
		return ""
	}
	
	file := square % 8
	rank := square / 8
	
	return string(rune('a'+file)) + string(rune('1'+rank))
}

// IsPieceType checks if a piece is of a specific type
func (b *BoardUtils) IsPieceType(piece chess.Piece, pieceType chess.PieceType) bool {
	return piece.Type() == pieceType
}

// GetPieceColor returns the color of a piece
func (b *BoardUtils) GetPieceColor(piece chess.Piece) chess.Color {
	return piece.Color()
}

// CountPiecesOfType counts pieces of a specific type and color
func (b *BoardUtils) CountPiecesOfType(game *chess.Game, pieceType chess.PieceType, color chess.Color) int {
	count := 0
	board := game.Position().Board()
	
	for sq := chess.A1; sq <= chess.H8; sq++ {
		piece := board.Piece(sq)
		if piece != chess.NoPiece && piece.Type() == pieceType && piece.Color() == color {
			count++
		}
	}
	
	return count
}

// GetAllPiecesOfType returns all squares containing pieces of a specific type and color
func (b *BoardUtils) GetAllPiecesOfType(game *chess.Game, pieceType chess.PieceType, color chess.Color) []chess.Square {
	squares := make([]chess.Square, 0)
	board := game.Position().Board()
	
	for sq := chess.A1; sq <= chess.H8; sq++ {
		piece := board.Piece(sq)
		if piece != chess.NoPiece && piece.Type() == pieceType && piece.Color() == color {
			squares = append(squares, sq)
		}
	}
	
	return squares
}

// IsEndgame determines if the position is an endgame
func (b *BoardUtils) IsEndgame(game *chess.Game) bool {
	// Count total pieces (excluding kings)
	board := game.Position().Board()
	pieceCount := 0
	
	for sq := chess.A1; sq <= chess.H8; sq++ {
		piece := board.Piece(sq)
		if piece != chess.NoPiece && piece.Type() != chess.King {
			pieceCount++
		}
	}
	
	// Endgame threshold: 6 or fewer pieces (excluding kings)
	return pieceCount <= 6
}

// IsOpening determines if the position is in the opening phase
func (b *BoardUtils) IsOpening(game *chess.Game) bool {
	// Simple heuristic: opening if move number <= 10
	return game.Position().MoveCount() <= 20 // Full moves, so 10 moves per side
}

// GetMaterialBalance calculates material balance (positive = white ahead)
func (b *BoardUtils) GetMaterialBalance(game *chess.Game) int {
	pieceValues := map[chess.PieceType]int{
		chess.Pawn:   100,
		chess.Knight: 320,
		chess.Bishop: 330,
		chess.Rook:   500,
		chess.Queen:  900,
		chess.King:   0,
	}
	
	balance := 0
	board := game.Position().Board()
	
	for sq := chess.A1; sq <= chess.H8; sq++ {
		piece := board.Piece(sq)
		if piece != chess.NoPiece {
			value := pieceValues[piece.Type()]
			if piece.Color() == chess.White {
				balance += value
			} else {
				balance -= value
			}
		}
	}
	
	return balance
}
