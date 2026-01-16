package chess

import (
	"fmt"

	"github.com/notnil/chess"
)

// MoveUtils provides utility functions for chess move operations
type MoveUtils struct{}

// NewMoveUtils creates a new move utilities instance
func NewMoveUtils() *MoveUtils {
	return &MoveUtils{}
}

// ParseUCIMove parses a UCI move string and returns a move
func (m *MoveUtils) ParseUCIMove(game *chess.Game, uci string) (*chess.Move, error) {
	move, err := chess.UCINotation{}.Decode(game.Position(), uci)
	if err != nil {
		return nil, fmt.Errorf("failed to parse UCI move %s: %w", uci, err)
	}
	return move, nil
}

// MoveToUCI converts a move to UCI notation
func (m *MoveUtils) MoveToUCI(move *chess.Move) string {
	return chess.UCINotation{}.Encode(chess.Position{}, move)
}

// IsCapture checks if a move is a capture
func (m *MoveUtils) IsCapture(game *chess.Game, move *chess.Move) bool {
	return game.Position().Board().Piece(move.S2()) != chess.NoPiece || move.HasTag(chess.EnPassant)
}

// IsCastling checks if a move is castling
func (m *MoveUtils) IsCastling(move *chess.Move) bool {
	return move.HasTag(chess.KingSideCastle) || move.HasTag(chess.QueenSideCastle)
}

// IsKingsideCastling checks if a move is kingside castling
func (m *MoveUtils) IsKingsideCastling(move *chess.Move) bool {
	return move.HasTag(chess.KingSideCastle)
}

// IsQueensideCastling checks if a move is queenside castling
func (m *MoveUtils) IsQueensideCastling(move *chess.Move) bool {
	return move.HasTag(chess.QueenSideCastle)
}

// IsEnPassant checks if a move is en passant
func (m *MoveUtils) IsEnPassant(move *chess.Move) bool {
	return move.HasTag(chess.EnPassant)
}

// IsPromotion checks if a move is a pawn promotion
func (m *MoveUtils) IsPromotion(move *chess.Move) bool {
	return move.Promo() != chess.NoPiece
}

// GetPromotionPiece returns the piece type of a promotion
func (m *MoveUtils) GetPromotionPiece(move *chess.Move) chess.PieceType {
	return move.Promo()
}

// GivesCheck checks if a move gives check
func (m *MoveUtils) GivesCheck(game *chess.Game, move *chess.Move) bool {
	// Make move on a copy and check if it results in check
	gameCopy := game.Copy()
	if err := gameCopy.Move(move); err != nil {
		return false
	}
	return gameCopy.Position().InCheck()
}

// IsCheckmate checks if a move results in checkmate
func (m *MoveUtils) IsCheckmate(game *chess.Game, move *chess.Move) bool {
	gameCopy := game.Copy()
	if err := gameCopy.Move(move); err != nil {
		return false
	}
	return gameCopy.Outcome() == chess.WhiteWon || gameCopy.Outcome() == chess.BlackWon
}

// GetMovedPiece returns the piece being moved
func (m *MoveUtils) GetMovedPiece(game *chess.Game, move *chess.Move) chess.Piece {
	return game.Position().Board().Piece(move.S1())
}

// GetCapturedPiece returns the piece being captured (if any)
func (m *MoveUtils) GetCapturedPiece(game *chess.Game, move *chess.Move) chess.Piece {
	return game.Position().Board().Piece(move.S2())
}

// GetFromSquare returns the source square of a move
func (m *MoveUtils) GetFromSquare(move *chess.Move) chess.Square {
	return move.S1()
}

// GetToSquare returns the destination square of a move
func (m *MoveUtils) GetToSquare(move *chess.Move) chess.Square {
	return move.S2()
}

// GetFromSquareName returns the source square name (e.g., "e2")
func (m *MoveUtils) GetFromSquareName(move *chess.Move) string {
	sq := move.S1()
	file := sq % 8
	rank := sq / 8
	return string(rune('a'+file)) + string(rune('1'+rank))
}

// GetToSquareName returns the destination square name (e.g., "e4")
func (m *MoveUtils) GetToSquareName(move *chess.Move) string {
	sq := move.S2()
	file := sq % 8
	rank := sq / 8
	return string(rune('a'+file)) + string(rune('1'+rank))
}

// IsLegalMove checks if a move is legal in the current position
func (m *MoveUtils) IsLegalMove(game *chess.Game, move *chess.Move) bool {
	for _, validMove := range game.ValidMoves() {
		if validMove.S1() == move.S1() && validMove.S2() == move.S2() && validMove.Promo() == move.Promo() {
			return true
		}
	}
	return false
}

// GetMoveString returns a human-readable move string
func (m *MoveUtils) GetMoveString(move *chess.Move) string {
	return move.String()
}

// IsPawnMove checks if a move is a pawn move
func (m *MoveUtils) IsPawnMove(game *chess.Game, move *chess.Move) bool {
	piece := game.Position().Board().Piece(move.S1())
	return piece.Type() == chess.Pawn
}

// IsKnightMove checks if a move is a knight move
func (m *MoveUtils) IsKnightMove(game *chess.Game, move *chess.Move) bool {
	piece := game.Position().Board().Piece(move.S1())
	return piece.Type() == chess.Knight
}

// IsBishopMove checks if a move is a bishop move
func (m *MoveUtils) IsBishopMove(game *chess.Game, move *chess.Move) bool {
	piece := game.Position().Board().Piece(move.S1())
	return piece.Type() == chess.Bishop
}

// IsRookMove checks if a move is a rook move
func (m *MoveUtils) IsRookMove(game *chess.Game, move *chess.Move) bool {
	piece := game.Position().Board().Piece(move.S1())
	return piece.Type() == chess.Rook
}

// IsQueenMove checks if a move is a queen move
func (m *MoveUtils) IsQueenMove(game *chess.Game, move *chess.Move) bool {
	piece := game.Position().Board().Piece(move.S1())
	return piece.Type() == chess.Queen
}

// IsKingMove checks if a move is a king move
func (m *MoveUtils) IsKingMove(game *chess.Game, move *chess.Move) bool {
	piece := game.Position().Board().Piece(move.S1())
	return piece.Type() == chess.King
}
