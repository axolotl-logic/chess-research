package main

import (
	"strconv"

	"github.com/corentings/chess/v2"
)

type CastlingRights struct {
	whiteKingsideCastling  bool
	whiteQueensideCastling bool
	blackKingsideCastling  bool
	blackQueensideCastling bool
}

type CastlingRightsDiffs struct {
	whiteKingsideCastlingDiff  bool
	whiteQueensideCastlingDiff bool
	blackKingsideCastlingDiff  bool
	blackQueensideCastlingDiff bool
}

func getCastlingDiffHeaders() []string {
	return []string{
		"whiteKingsideCastlingDiff",
		"whiteQueensideCastlingDiff",
		"blackKingsideCastlingDiff",
		"blackQueensideCastlingDiff",
	}
}

func getCastlingDiffValues(diffs CastlingRightsDiffs) []string {
	return []string{
		strconv.FormatBool(diffs.whiteKingsideCastlingDiff),
		strconv.FormatBool(diffs.whiteQueensideCastlingDiff),
		strconv.FormatBool(diffs.blackKingsideCastlingDiff),
		strconv.FormatBool(diffs.blackQueensideCastlingDiff),
	}
}

func diffCastlingRights(a *CastlingRights, b *CastlingRights) CastlingRightsDiffs {
	return CastlingRightsDiffs{
		whiteKingsideCastlingDiff:  a.whiteKingsideCastling != b.whiteKingsideCastling,
		whiteQueensideCastlingDiff: a.whiteQueensideCastling != b.whiteQueensideCastling,
		blackKingsideCastlingDiff:  a.blackKingsideCastling != b.blackKingsideCastling,
		blackQueensideCastlingDiff: a.blackQueensideCastling != b.blackQueensideCastling,
	}
}

func newCastlingRights(pos *chess.Position) CastlingRights {
	rights := CastlingRights{}
	castlingRights := pos.CastleRights()

	rights.whiteKingsideCastling = castlingRights.CanCastle(chess.White, chess.KingSide)
	rights.whiteQueensideCastling = castlingRights.CanCastle(chess.White, chess.QueenSide)
	rights.blackKingsideCastling = castlingRights.CanCastle(chess.Black, chess.KingSide)
	rights.blackQueensideCastling = castlingRights.CanCastle(chess.Black, chess.QueenSide)

	return rights
}
