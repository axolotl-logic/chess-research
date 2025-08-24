package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/corentings/chess/v2"
)

func main() {
	var inputFile = flag.String("input", "", "Input TSV file")
	var outputFile = flag.String("output", "", "Output TSV file")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		flag.Usage()
		os.Exit(-1)
	}

	f, err := os.Open(*inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	out_f, err := os.Create(*outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer out_f.Close()

	w := csv.NewWriter(out_f)
	w.Comma = '\t' // Set delimiter to tab for TSV
	writeHeaders(w)

	// Read TSV file
	reader := csv.NewReader(f)
	reader.Comma = '\t' // Set delimiter to tab for TSV

	// Read header row to find column indices
	headers, err := reader.Read()
	if err != nil {
		log.Fatal("Failed to read headers:", err)
	}

	columnIndices := make(map[string]int)
	for i, header := range headers {
		columnIndices[header] = i
	}

	// Verify required columns exist
	requiredColumns := []string{"moves", "black_clock", "white_clock", "white", "white_elo", "black", "black_elo", "result", "termination", "date"}
	for _, col := range requiredColumns {
		if _, exists := columnIndices[col]; !exists {
			log.Fatalf("Required column '%s' not found in TSV", col)
		}
	}

	// Process each row
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatal("Failed to read record:", err)
		}

		gameInfo := parseGameRecord(record, columnIndices)
		processGame(w, &gameInfo)
	}

	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

type GameInfo struct {
	White       string
	WhiteElo    string
	Black       string
	BlackElo    string
	Result      string
	Termination string
	Date        string
	Moves       []string
	WhiteClocks []string
	BlackClocks []string
}

func processGame(w *csv.Writer, gameInfo *GameInfo) {
	if len(gameInfo.Moves) == 0 {
		return
	}

	// Initialize chess game
	game := chess.NewGame()

	// Process each move
	for moveIndex, moveStr := range gameInfo.Moves {
		// Get current position before move
		currentPos := game.Position()

		err := game.PushNotationMove(moveStr, chess.AlgebraicNotation{}, nil)
		if err != nil {
			log.Fatalf("Failed to make move '%s': %v", moveStr, err)
		}

		// Get new position after move
		newPos := game.Position()

		// Determine clock time for this move
		clockTime := getClockTimeForMove(gameInfo, moveIndex, currentPos.Turn())

		// Write position data
		writePosition(w, gameInfo, currentPos, newPos, clockTime)
	}
}

func getClockTimeForMove(gameInfo *GameInfo, moveIndex int, turn chess.Color) string {
	// Calculate which move number this is (considering both white and black moves)
	moveNumber := moveIndex / 2

	if turn == chess.White {
		if moveNumber < len(gameInfo.WhiteClocks) {
			return gameInfo.WhiteClocks[moveNumber]
		}
	} else {
		if moveNumber < len(gameInfo.BlackClocks) {
			return gameInfo.BlackClocks[moveNumber]
		}
	}

	return ""
}

func parseGameRecord(record []string, columnIndices map[string]int) GameInfo {
	gameInfo := GameInfo{
		White:       getFieldValue(record, columnIndices, "white"),
		WhiteElo:    getFieldValue(record, columnIndices, "whiteElo"),
		Black:       getFieldValue(record, columnIndices, "black"),
		BlackElo:    getFieldValue(record, columnIndices, "blackElo"),
		Result:      getFieldValue(record, columnIndices, "result"),
		Termination: getFieldValue(record, columnIndices, "termination"),
		Date:        getFieldValue(record, columnIndices, "date"),
	}

	// Parse moves
	movesStr := getFieldValue(record, columnIndices, "moves")
	if movesStr != "" {
		gameInfo.Moves = strings.Fields(movesStr)
	}

	// Parse clock times
	whiteClockStr := getFieldValue(record, columnIndices, "whiteClock")
	blackClockStr := getFieldValue(record, columnIndices, "blackClock")

	if whiteClockStr != "" {
		gameInfo.WhiteClocks = strings.Fields(whiteClockStr)
	}
	if blackClockStr != "" {
		gameInfo.BlackClocks = strings.Fields(blackClockStr)
	}

	return gameInfo
}

func getFieldValue(record []string, columnIndices map[string]int, fieldName string) string {
	if index, exists := columnIndices[fieldName]; exists && index < len(record) {
		return record[index]
	}
	return ""
}

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

func newPieceCounts(pos *chess.Position) PieceCounts {
	counts := PieceCounts{}
	board := pos.Board()

	for sq, piece := range board.SquareMap() {
		switch piece {
		case chess.WhitePawn:
			counts.whitePawnCount++
		case chess.WhiteQueen:
			counts.whiteQueenCount++
			break
		case chess.WhiteRook:
			counts.whiteRookCount++
			break
		case chess.WhiteKnight:
			counts.whiteKnightCount++
			break
		case chess.WhiteBishop:
			if isDark(sq) {
				counts.whiteDarkBishopCount++
			} else {
				counts.whiteLightBishopCount++
			}
			break
		case chess.BlackPawn:
			counts.blackPawnCount++
		case chess.BlackQueen:
			counts.blackQueenCount++
			break
		case chess.BlackRook:
			counts.blackRookCount++
			break
		case chess.BlackKnight:
			counts.blackKnightCount++
			break
		case chess.BlackBishop:
			if isDark(sq) {
				counts.blackDarkBishopCount++
			} else {
				counts.blackLightBishopCount++
			}
			break
		}

	}

	return counts
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

func writePosition(writer *csv.Writer, gameInfo *GameInfo, lastPos *chess.Position, pos *chess.Position, clock string) {
	// Get counts from old position
	pieceCounts := newPieceCounts(lastPos)
	pieceSightCounts := newPieceSightCounts(lastPos)
	pawnFileCounts := newPawnFileCounts(lastPos)
	castlingRights := newCastlingRights(lastPos)

	// Get counts from new position to calculate diffs
	newPieceCounts := newPieceCounts(pos)
	newPieceSightCounts := newPieceSightCounts(pos)
	newPawnFileCounts := newPawnFileCounts(pos)
	newCastlingRights := newCastlingRights(pos)

	// Calculate diffs
	pieceCountDiffs := diffPieceCounts(&newPieceCounts, &pieceCounts)
	pieceSightDiffs := diffPieceSightCounts(&newPieceSightCounts, &pieceSightCounts)
	pawnFileCountDiffs := diffPawnFileCounts(&newPawnFileCounts, &pawnFileCounts)
	castlingRightsDiffs := diffCastlingRights(&newCastlingRights, &castlingRights)

	// Build CSV row
	record := []string{
		// white_to_play
		strconv.FormatBool(pos.Turn() == chess.White),
		// ply
		strconv.Itoa(pos.Ply()),
		// fen
		pos.String(),
		// clock for current player,
		clock,
		// Basic game info
		gameInfo.White,
		gameInfo.WhiteElo,
		gameInfo.Black,
		gameInfo.BlackElo,
		gameInfo.Result,
		gameInfo.Termination,
		gameInfo.Date,
		// PieceCounts (old position)
		fmt.Sprintf("%d", pieceCounts.whiteQueenCount),
		fmt.Sprintf("%d", pieceCounts.whiteKnightCount),
		fmt.Sprintf("%d", pieceCounts.whitePawnCount),
		fmt.Sprintf("%d", pieceCounts.whiteRookCount),
		fmt.Sprintf("%d", pieceCounts.whiteDarkBishopCount),
		fmt.Sprintf("%d", pieceCounts.whiteLightBishopCount),
		fmt.Sprintf("%d", pieceCounts.blackQueenCount),
		fmt.Sprintf("%d", pieceCounts.blackKnightCount),
		fmt.Sprintf("%d", pieceCounts.blackPawnCount),
		fmt.Sprintf("%d", pieceCounts.blackRookCount),
		fmt.Sprintf("%d", pieceCounts.blackDarkBishopCount),
		fmt.Sprintf("%d", pieceCounts.blackLightBishopCount),

		// PieceCountDiffs (new - old)
		fmt.Sprintf("%d", pieceCountDiffs.whiteQueenCountDiff),
		fmt.Sprintf("%d", pieceCountDiffs.whiteKnightCountDiff),
		fmt.Sprintf("%d", pieceCountDiffs.whitePawnCountDiff),
		fmt.Sprintf("%d", pieceCountDiffs.whiteRookCountDiff),
		fmt.Sprintf("%d", pieceCountDiffs.whiteDarkBishopCountDiff),
		fmt.Sprintf("%d", pieceCountDiffs.whiteLightBishopCountDiff),
		fmt.Sprintf("%d", pieceCountDiffs.blackQueenCountDiff),
		fmt.Sprintf("%d", pieceCountDiffs.blackKnightCountDiff),
		fmt.Sprintf("%d", pieceCountDiffs.blackPawnCountDiff),
		fmt.Sprintf("%d", pieceCountDiffs.blackRookCountDiff),
		fmt.Sprintf("%d", pieceCountDiffs.blackDarkBishopCountDiff),
		fmt.Sprintf("%d", pieceCountDiffs.blackLightBishopCountDiff),

		// PieceSightCounts (old position)
		fmt.Sprintf("%d", pieceSightCounts.whiteQueenSightCount),
		fmt.Sprintf("%d", pieceSightCounts.whiteKnightSightCount),
		fmt.Sprintf("%d", pieceSightCounts.whitePawnSightCount),
		fmt.Sprintf("%d", pieceSightCounts.whiteRookSightCount),
		fmt.Sprintf("%d", pieceSightCounts.whiteDarkBishopSightCount),
		fmt.Sprintf("%d", pieceSightCounts.whiteLightBishopSightCount),
		fmt.Sprintf("%d", pieceSightCounts.blackQueenSightCount),
		fmt.Sprintf("%d", pieceSightCounts.blackKnightSightCount),
		fmt.Sprintf("%d", pieceSightCounts.blackPawnSightCount),
		fmt.Sprintf("%d", pieceSightCounts.blackRookSightCount),
		fmt.Sprintf("%d", pieceSightCounts.blackDarkBishopSightCount),
		fmt.Sprintf("%d", pieceSightCounts.blackLightBishopSightCount),

		// PieceSightSeeDiffs (new - old)
		fmt.Sprintf("%d", pieceSightDiffs.whiteQueenSightCountDiff),
		fmt.Sprintf("%d", pieceSightDiffs.whiteKnightSightCountDiff),
		fmt.Sprintf("%d", pieceSightDiffs.whitePawnSightCountDiff),
		fmt.Sprintf("%d", pieceSightDiffs.whiteRookSightCountDiff),
		fmt.Sprintf("%d", pieceSightDiffs.whiteDarkBishopSightCountDiff),
		fmt.Sprintf("%d", pieceSightDiffs.whiteLightBishopSightCountDiff),
		fmt.Sprintf("%d", pieceSightDiffs.blackQueenSightCountDiff),
		fmt.Sprintf("%d", pieceSightDiffs.blackKnightSightCountDiff),
		fmt.Sprintf("%d", pieceSightDiffs.blackPawnSightCountDiff),
		fmt.Sprintf("%d", pieceSightDiffs.blackRookSightCountDiff),
		fmt.Sprintf("%d", pieceSightDiffs.blackDarkBishopSightCountDiff),
		fmt.Sprintf("%d", pieceSightDiffs.blackLightBishopSightCountDiff),

		//  Pawns count
		fmt.Sprintf("%d", pawnFileCounts.whitePawnFileA),
		fmt.Sprintf("%d", pawnFileCounts.whitePawnFileB),
		fmt.Sprintf("%d", pawnFileCounts.whitePawnFileC),
		fmt.Sprintf("%d", pawnFileCounts.whitePawnFileD),
		fmt.Sprintf("%d", pawnFileCounts.whitePawnFileE),
		fmt.Sprintf("%d", pawnFileCounts.whitePawnFileF),
		fmt.Sprintf("%d", pawnFileCounts.whitePawnFileG),
		fmt.Sprintf("%d", pawnFileCounts.whitePawnFileH),
		fmt.Sprintf("%d", pawnFileCounts.blackPawnFileA),
		fmt.Sprintf("%d", pawnFileCounts.blackPawnFileB),
		fmt.Sprintf("%d", pawnFileCounts.blackPawnFileC),
		fmt.Sprintf("%d", pawnFileCounts.blackPawnFileD),
		fmt.Sprintf("%d", pawnFileCounts.blackPawnFileE),
		fmt.Sprintf("%d", pawnFileCounts.blackPawnFileF),
		fmt.Sprintf("%d", pawnFileCounts.blackPawnFileG),
		fmt.Sprintf("%d", pawnFileCounts.blackPawnFileH),

		// PawnFileCountDiffs (new - old)
		fmt.Sprintf("%d", pawnFileCountDiffs.whitePawnFileADiff),
		fmt.Sprintf("%d", pawnFileCountDiffs.whitePawnFileBDiff),
		fmt.Sprintf("%d", pawnFileCountDiffs.whitePawnFileCDiff),
		fmt.Sprintf("%d", pawnFileCountDiffs.whitePawnFileDDiff),
		fmt.Sprintf("%d", pawnFileCountDiffs.whitePawnFileEDiff),
		fmt.Sprintf("%d", pawnFileCountDiffs.whitePawnFileFDiff),
		fmt.Sprintf("%d", pawnFileCountDiffs.whitePawnFileGDiff),
		fmt.Sprintf("%d", pawnFileCountDiffs.whitePawnFileHDiff),
		fmt.Sprintf("%d", pawnFileCountDiffs.blackPawnFileADiff),
		fmt.Sprintf("%d", pawnFileCountDiffs.blackPawnFileBDiff),
		fmt.Sprintf("%d", pawnFileCountDiffs.blackPawnFileCDiff),
		fmt.Sprintf("%d", pawnFileCountDiffs.blackPawnFileDDiff),
		fmt.Sprintf("%d", pawnFileCountDiffs.blackPawnFileEDiff),
		fmt.Sprintf("%d", pawnFileCountDiffs.blackPawnFileFDiff),
		fmt.Sprintf("%d", pawnFileCountDiffs.blackPawnFileGDiff),
		fmt.Sprintf("%d", pawnFileCountDiffs.blackPawnFileHDiff),

		// CastlingRights (old position)
		fmt.Sprintf("%t", castlingRights.whiteKingsideCastling),
		fmt.Sprintf("%t", castlingRights.whiteQueensideCastling),
		fmt.Sprintf("%t", castlingRights.blackKingsideCastling),
		fmt.Sprintf("%t", castlingRights.blackQueensideCastling),

		// CastlingRightsDiffs (changed or not)
		fmt.Sprintf("%t", castlingRightsDiffs.whiteKingsideCastlingDiff),
		fmt.Sprintf("%t", castlingRightsDiffs.whiteQueensideCastlingDiff),
		fmt.Sprintf("%t", castlingRightsDiffs.blackKingsideCastlingDiff),
		fmt.Sprintf("%t", castlingRightsDiffs.blackQueensideCastlingDiff),
	}

	writer.Write(record)
}
func writeHeaders(writer *csv.Writer) {
	headers := []string{
		"white_to_play",
		"ply",
		"fen",
		"clock",
		// Basic game info
		"white",
		"white_elo",
		"black",
		"black_elo",
		"result",
		"termination",
		"date",
		// PieceCounts (old position)
		"white_queen_count",
		"white_knight_count",
		"white_pawn_count",
		"white_rook_count",
		"white_dark_bishop_count",
		"white_light_bishop_count",
		"black_queen_count",
		"black_knight_count",
		"black_pawn_count",
		"black_rook_count",
		"black_dark_bishop_count",
		"black_light_bishop_count",

		// PieceCountDiffs (new - old)
		"white_queen_count_diff",
		"white_knight_count_diff",
		"white_pawn_count_diff",
		"white_rook_count_diff",
		"white_dark_bishop_count_diff",
		"white_light_bishop_count_diff",
		"black_queen_count_diff",
		"black_knight_count_diff",
		"black_pawn_count_diff",
		"black_rook_count_diff",
		"black_dark_bishop_count_diff",
		"black_light_bishop_count_diff",

		// PieceSightCounts (old position)
		"white_queen_sight_count",
		"white_knight_sight_count",
		"white_pawn_sight_count",
		"white_rook_sight_count",
		"white_dark_bishop_sight_count",
		"white_light_bishop_sight_count",
		"black_queen_sight_count",
		"black_knight_sight_count",
		"black_pawn_sight_count",
		"black_rook_sight_count",
		"black_dark_bishop_sight_count",
		"black_light_bishop_sight_count",

		// PieceSightCountDiffs (new - old)
		"white_queen_sight_count_diff",
		"white_knight_sight_count_diff",
		"white_pawn_sight_count_diff",
		"white_rook_sight_count_diff",
		"white_dark_bishop_sight_count_diff",
		"white_light_bishop_sight_count_diff",
		"black_queen_sight_count_diff",
		"black_knight_sight_count_diff",
		"black_pawn_sight_count_diff",
		"black_rook_sight_count_diff",
		"black_dark_bishop_sight_count_diff",
		"black_light_bishop_sight_count_diff",

		// PawnFileCounts (old position)
		"white_pawn_file_a",
		"white_pawn_file_b",
		"white_pawn_file_c",
		"white_pawn_file_d",
		"white_pawn_file_e",
		"white_pawn_file_f",
		"white_pawn_file_g",
		"white_pawn_file_h",
		"black_pawn_file_a",
		"black_pawn_file_b",
		"black_pawn_file_c",
		"black_pawn_file_d",
		"black_pawn_file_e",
		"black_pawn_file_f",
		"black_pawn_file_g",
		"black_pawn_file_h",

		// PawnFileCountDiffs (new - old)
		"white_pawn_file_a_diff",
		"white_pawn_file_b_diff",
		"white_pawn_file_c_diff",
		"white_pawn_file_d_diff",
		"white_pawn_file_e_diff",
		"white_pawn_file_f_diff",
		"white_pawn_file_g_diff",
		"white_pawn_file_h_diff",
		"black_pawn_file_a_diff",
		"black_pawn_file_b_diff",
		"black_pawn_file_c_diff",
		"black_pawn_file_d_diff",
		"black_pawn_file_e_diff",
		"black_pawn_file_f_diff",
		"black_pawn_file_g_diff",
		"black_pawn_file_h_diff",

		// CastlingRights (old position)
		"white_kingside_castling",
		"white_queenside_castling",
		"black_kingside_castling",
		"black_queenside_castling",

		// CastlingRightsDiffs (changed or not)
		"white_kingside_castling_diff",
		"white_queenside_castling_diff",
		"black_kingside_castling_diff",
		"black_queenside_castling_diff",
	}

	writer.Write(headers)
}

func isDark(sq chess.Square) bool {
	if ((sq / 8) % 2) == (sq % 2) {
		return true
	}
	return false
}
