package main

import (
	"strconv"

	"github.com/corentings/chess/v2"
)

type PieceSightCount struct {
	whiteQueenSightCount       int
	whiteKnightSightCount      int
	whitePawnSightCount        int
	whiteRookSightCount        int
	whiteDarkBishopSightCount  int
	whiteLightBishopSightCount int
	blackQueenSightCount       int
	blackKnightSightCount      int
	blackPawnSightCount        int
	blackRookSightCount        int
	blackDarkBishopSightCount  int
	blackLightBishopSightCount int
}

type PieceSightCountDiffs struct {
	whiteQueenSightCountDiff       int
	whiteKnightSightCountDiff      int
	whitePawnSightCountDiff        int
	whiteRookSightCountDiff        int
	whiteDarkBishopSightCountDiff  int
	whiteLightBishopSightCountDiff int
	blackQueenSightCountDiff       int
	blackKnightSightCountDiff      int
	blackPawnSightCountDiff        int
	blackRookSightCountDiff        int
	blackDarkBishopSightCountDiff  int
	blackLightBishopSightCountDiff int
}

func getPieceSightDiffHeaders() []string {
	return []string{
		"whiteQueenSightCountDiff",
		"whiteKnightSightCountDiff",
		"whitePawnSightCountDiff",
		"whiteRookSightCountDiff",
		"whiteDarkBishopSightCountDiff",
		"whiteLightBishopSightCountDiff",
		"blackQueenSightCountDiff",
		"blackKnightSightCountDiff",
		"blackPawnSightCountDiff",
		"blackRookSightCountDiff",
		"blackDarkBishopSightCountDiff",
		"blackLightBishopSightCountDiff",
	}
}

func getPieceSightDiffValues(diffs PieceSightCountDiffs) []string {
	return []string{
		strconv.Itoa(diffs.whiteQueenSightCountDiff),
		strconv.Itoa(diffs.whiteKnightSightCountDiff),
		strconv.Itoa(diffs.whitePawnSightCountDiff),
		strconv.Itoa(diffs.whiteRookSightCountDiff),
		strconv.Itoa(diffs.whiteDarkBishopSightCountDiff),
		strconv.Itoa(diffs.whiteLightBishopSightCountDiff),
		strconv.Itoa(diffs.blackQueenSightCountDiff),
		strconv.Itoa(diffs.blackKnightSightCountDiff),
		strconv.Itoa(diffs.blackPawnSightCountDiff),
		strconv.Itoa(diffs.blackRookSightCountDiff),
		strconv.Itoa(diffs.blackDarkBishopSightCountDiff),
		strconv.Itoa(diffs.blackLightBishopSightCountDiff),
	}
}

func diffPieceSightCounts(a *PieceSightCount, b *PieceSightCount) PieceSightCountDiffs {
	return PieceSightCountDiffs{
		whiteQueenSightCountDiff:       a.whiteQueenSightCount - b.whiteQueenSightCount,
		whiteKnightSightCountDiff:      a.whiteKnightSightCount - b.whiteKnightSightCount,
		whitePawnSightCountDiff:        a.whitePawnSightCount - b.whitePawnSightCount,
		whiteRookSightCountDiff:        a.whiteRookSightCount - b.whiteRookSightCount,
		whiteDarkBishopSightCountDiff:  a.whiteDarkBishopSightCount - b.whiteDarkBishopSightCount,
		whiteLightBishopSightCountDiff: a.whiteLightBishopSightCount - b.whiteLightBishopSightCount,
		blackQueenSightCountDiff:       a.blackQueenSightCount - b.blackQueenSightCount,
		blackKnightSightCountDiff:      a.blackKnightSightCount - b.blackKnightSightCount,
		blackPawnSightCountDiff:        a.blackPawnSightCount - b.blackPawnSightCount,
		blackRookSightCountDiff:        a.blackRookSightCount - b.blackRookSightCount,
		blackDarkBishopSightCountDiff:  a.blackDarkBishopSightCount - b.blackDarkBishopSightCount,
		blackLightBishopSightCountDiff: a.blackLightBishopSightCount - b.blackLightBishopSightCount,
	}
}

func newPieceSightCounts(pos *chess.Position) PieceSightCount {
	counts := PieceSightCount{}
	board := pos.Board()

	// Get all valid moves to determine piece sight/attack patterns
	validMoves := pos.ValidMoves()
	pos.ChangeTurn()
	validMoves = append(validMoves, pos.ValidMoves()...)
	pos.ChangeTurn()

	for _, move := range validMoves {
		piece := board.Piece(move.S1())
		targetSq := move.S2()

		switch piece {
		case chess.WhitePawn:
			counts.whitePawnSightCount++
		case chess.WhiteQueen:
			counts.whiteQueenSightCount++
		case chess.WhiteRook:
			counts.whiteRookSightCount++
		case chess.WhiteKnight:
			counts.whiteKnightSightCount++
		case chess.WhiteBishop:
			if isDark(targetSq) {
				counts.whiteDarkBishopSightCount++
			} else {
				counts.whiteLightBishopSightCount++
			}
		case chess.BlackPawn:
			counts.blackPawnSightCount++
		case chess.BlackQueen:
			counts.blackQueenSightCount++
		case chess.BlackRook:
			counts.blackRookSightCount++
		case chess.BlackKnight:
			counts.blackKnightSightCount++
		case chess.BlackBishop:
			if isDark(targetSq) {
				counts.blackDarkBishopSightCount++
			} else {
				counts.blackLightBishopSightCount++
			}
		}
	}
	return counts
}
