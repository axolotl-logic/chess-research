package main

import (
	"strconv"

	"github.com/corentings/chess/v2"
)

type PawnFileCounts struct {
	whitePawnFileA int
	whitePawnFileB int
	whitePawnFileC int
	whitePawnFileD int
	whitePawnFileE int
	whitePawnFileF int
	whitePawnFileG int
	whitePawnFileH int
	blackPawnFileA int
	blackPawnFileB int
	blackPawnFileC int
	blackPawnFileD int
	blackPawnFileE int
	blackPawnFileF int
	blackPawnFileG int
	blackPawnFileH int
}

type PawnFileCountDiffs struct {
	whitePawnFileADiff int
	whitePawnFileBDiff int
	whitePawnFileCDiff int
	whitePawnFileDDiff int
	whitePawnFileEDiff int
	whitePawnFileFDiff int
	whitePawnFileGDiff int
	whitePawnFileHDiff int
	blackPawnFileADiff int
	blackPawnFileBDiff int
	blackPawnFileCDiff int
	blackPawnFileDDiff int
	blackPawnFileEDiff int
	blackPawnFileFDiff int
	blackPawnFileGDiff int
	blackPawnFileHDiff int
}

func getPawnFileDiffsHeaders() []string {
	return []string{
		"whitePawnFileADiff",
		"whitePawnFileBDiff",
		"whitePawnFileCDiff",
		"whitePawnFileDDiff",
		"whitePawnFileEDiff",
		"whitePawnFileFDiff",
		"whitePawnFileGDiff",
		"whitePawnFileHDiff",
		"blackPawnFileADiff",
		"blackPawnFileBDiff",
		"blackPawnFileCDiff",
		"blackPawnFileDDiff",
		"blackPawnFileEDiff",
		"blackPawnFileFDiff",
		"blackPawnFileGDiff",
		"blackPawnFileHDiff",
	}
}

func getPawnFileDiffsValues(diffs PawnFileCountDiffs) []string {
	return []string{
		strconv.Itoa(diffs.whitePawnFileADiff),
		strconv.Itoa(diffs.whitePawnFileBDiff),
		strconv.Itoa(diffs.whitePawnFileCDiff),
		strconv.Itoa(diffs.whitePawnFileDDiff),
		strconv.Itoa(diffs.whitePawnFileEDiff),
		strconv.Itoa(diffs.whitePawnFileFDiff),
		strconv.Itoa(diffs.whitePawnFileGDiff),
		strconv.Itoa(diffs.whitePawnFileHDiff),
		strconv.Itoa(diffs.blackPawnFileADiff),
		strconv.Itoa(diffs.blackPawnFileBDiff),
		strconv.Itoa(diffs.blackPawnFileCDiff),
		strconv.Itoa(diffs.blackPawnFileDDiff),
		strconv.Itoa(diffs.blackPawnFileEDiff),
		strconv.Itoa(diffs.blackPawnFileFDiff),
		strconv.Itoa(diffs.blackPawnFileGDiff),
		strconv.Itoa(diffs.blackPawnFileHDiff),
	}
}

func diffPawnFileCounts(a *PawnFileCounts, b *PawnFileCounts) PawnFileCountDiffs {
	return PawnFileCountDiffs{
		whitePawnFileADiff: a.whitePawnFileA - b.whitePawnFileA,
		whitePawnFileBDiff: a.whitePawnFileB - b.whitePawnFileB,
		whitePawnFileCDiff: a.whitePawnFileC - b.whitePawnFileC,
		whitePawnFileDDiff: a.whitePawnFileD - b.whitePawnFileD,
		whitePawnFileEDiff: a.whitePawnFileE - b.whitePawnFileE,
		whitePawnFileFDiff: a.whitePawnFileF - b.whitePawnFileF,
		whitePawnFileGDiff: a.whitePawnFileG - b.whitePawnFileG,
		whitePawnFileHDiff: a.whitePawnFileH - b.whitePawnFileH,
		blackPawnFileADiff: a.blackPawnFileA - b.blackPawnFileA,
		blackPawnFileBDiff: a.blackPawnFileB - b.blackPawnFileB,
		blackPawnFileCDiff: a.blackPawnFileC - b.blackPawnFileC,
		blackPawnFileDDiff: a.blackPawnFileD - b.blackPawnFileD,
		blackPawnFileEDiff: a.blackPawnFileE - b.blackPawnFileE,
		blackPawnFileFDiff: a.blackPawnFileF - b.blackPawnFileF,
		blackPawnFileGDiff: a.blackPawnFileG - b.blackPawnFileG,
		blackPawnFileHDiff: a.blackPawnFileH - b.blackPawnFileH,
	}
}

func newPawnFileCounts(pos *chess.Position) PawnFileCounts {
	counts := PawnFileCounts{}
	board := pos.Board()

	for sq, piece := range board.SquareMap() {
		file := sq.File()
		switch piece {
		case chess.WhitePawn:
			switch file {
			case chess.FileA:
				counts.whitePawnFileA++
			case chess.FileB:
				counts.whitePawnFileB++
			case chess.FileC:
				counts.whitePawnFileC++
			case chess.FileD:
				counts.whitePawnFileD++
			case chess.FileE:
				counts.whitePawnFileE++
			case chess.FileF:
				counts.whitePawnFileF++
			case chess.FileG:
				counts.whitePawnFileG++
			case chess.FileH:
				counts.whitePawnFileH++
			}
		case chess.BlackPawn:
			switch file {
			case chess.FileA:
				counts.blackPawnFileA++
			case chess.FileB:
				counts.blackPawnFileB++
			case chess.FileC:
				counts.blackPawnFileC++
			case chess.FileD:
				counts.blackPawnFileD++
			case chess.FileE:
				counts.blackPawnFileE++
			case chess.FileF:
				counts.blackPawnFileF++
			case chess.FileG:
				counts.blackPawnFileG++
			case chess.FileH:
				counts.blackPawnFileH++
			}
		}
	}

	return counts
}
