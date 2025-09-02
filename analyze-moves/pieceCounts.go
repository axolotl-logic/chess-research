package main

import (
	"strconv"

	"github.com/corentings/chess/v2"
)

type PieceCounts struct {
	whiteQueenCount       int
	whiteKnightCount      int
	whitePawnCount        int
	whiteRookCount        int
	whiteDarkBishopCount  int
	whiteLightBishopCount int
	blackQueenCount       int
	blackKnightCount      int
	blackPawnCount        int
	blackRookCount        int
	blackDarkBishopCount  int
	blackLightBishopCount int
}

type PieceCountDiffs struct {
	whiteQueenCountDiff       int
	whiteKnightCountDiff      int
	whitePawnCountDiff        int
	whiteRookCountDiff        int
	whiteDarkBishopCountDiff  int
	whiteLightBishopCountDiff int
	blackQueenCountDiff       int
	blackKnightCountDiff      int
	blackPawnCountDiff        int
	blackRookCountDiff        int
	blackDarkBishopCountDiff  int
	blackLightBishopCountDiff int
}

func getPieceCountDiffHeaders() []string {
	return []string{
		"whiteQueenCountDiff",
		"whiteKnightCountDiff",
		"whitePawnCountDiff",
		"whiteRookCountDiff",
		"whiteDarkBishopCountDiff",
		"whiteLightBishopCountDiff",
		"blackQueenCountDiff",
		"blackKnightCountDiff",
		"blackPawnCountDiff",
		"blackRookCountDiff",
		"blackDarkBishopCountDiff",
		"blackLightBishopCountDiff",
	}
}

func getPieceCountDiffValues(diffs PieceCountDiffs) []string {
	return []string{
		strconv.Itoa(diffs.whiteQueenCountDiff),
		strconv.Itoa(diffs.whiteKnightCountDiff),
		strconv.Itoa(diffs.whitePawnCountDiff),
		strconv.Itoa(diffs.whiteRookCountDiff),
		strconv.Itoa(diffs.whiteDarkBishopCountDiff),
		strconv.Itoa(diffs.whiteLightBishopCountDiff),
		strconv.Itoa(diffs.blackQueenCountDiff),
		strconv.Itoa(diffs.blackKnightCountDiff),
		strconv.Itoa(diffs.blackPawnCountDiff),
		strconv.Itoa(diffs.blackRookCountDiff),
		strconv.Itoa(diffs.blackDarkBishopCountDiff),
		strconv.Itoa(diffs.blackLightBishopCountDiff),
	}
}

func diffPieceCounts(a *PieceCounts, b *PieceCounts) PieceCountDiffs {
	return PieceCountDiffs{
		whiteQueenCountDiff:       a.whiteQueenCount - b.whiteQueenCount,
		whiteKnightCountDiff:      a.whiteKnightCount - b.whiteKnightCount,
		whitePawnCountDiff:        a.whitePawnCount - b.whitePawnCount,
		whiteRookCountDiff:        a.whiteRookCount - b.whiteRookCount,
		whiteDarkBishopCountDiff:  a.whiteDarkBishopCount - b.whiteDarkBishopCount,
		whiteLightBishopCountDiff: a.whiteLightBishopCount - b.whiteLightBishopCount,
		blackQueenCountDiff:       a.blackQueenCount - b.blackQueenCount,
		blackKnightCountDiff:      a.blackKnightCount - b.blackKnightCount,
		blackPawnCountDiff:        a.blackPawnCount - b.blackPawnCount,
		blackRookCountDiff:        a.blackRookCount - b.blackRookCount,
		blackDarkBishopCountDiff:  a.blackDarkBishopCount - b.blackDarkBishopCount,
		blackLightBishopCountDiff: a.blackLightBishopCount - b.blackLightBishopCount,
	}
}

func newPieceCounts(pos *chess.Position) PieceCounts {
	counts := PieceCounts{}
	board := pos.Board()

	for sq, piece := range board.SquareMap() {
		switch piece {
		case chess.WhitePawn:
			counts.whitePawnCount++
		case chess.WhiteQueen:
			counts.whiteQueenCount++
		case chess.WhiteRook:
			counts.whiteRookCount++
		case chess.WhiteKnight:
			counts.whiteKnightCount++
		case chess.WhiteBishop:
			if isDark(sq) {
				counts.whiteDarkBishopCount++
			} else {
				counts.whiteLightBishopCount++
			}
		case chess.BlackPawn:
			counts.blackPawnCount++
		case chess.BlackQueen:
			counts.blackQueenCount++
		case chess.BlackRook:
			counts.blackRookCount++
		case chess.BlackKnight:
			counts.blackKnightCount++
		case chess.BlackBishop:
			if isDark(sq) {
				counts.blackDarkBishopCount++
			} else {
				counts.blackLightBishopCount++
			}
		}

	}

	return counts
}
