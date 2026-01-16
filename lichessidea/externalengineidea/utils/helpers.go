package utils

import (
	"fmt"
	"math"
)

// ClampFloat64 clamps a float64 value between min and max
func ClampFloat64(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ParseUCIMove parses a UCI move string (e.g., "e2e4" or "e7e8q")
func ParseUCIMove(uci string) (from, to string, promotion rune, err error) {
	if len(uci) < 4 || len(uci) > 5 {
		return "", "", 0, fmt.Errorf("invalid UCI move length: %s", uci)
	}
	
	from = uci[0:2]
	to = uci[2:4]
	
	if len(uci) == 5 {
		promotion = rune(uci[4])
	}
	
	return from, to, promotion, nil
}

// IsCheckmate checks if a score indicates checkmate
func IsCheckmate(scorePawns float64) bool {
	return math.Abs(scorePawns) >= 100.0
}

// IsForcedMate checks if a score indicates a forced mate
func IsForcedMate(mateIn *int) bool {
	return mateIn != nil && *mateIn != 0
}
