package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/corentings/chess/v2"
)

func main() {
	inputFile := flag.String("input", "", "Input CSV file")
	outputFile := flag.String("output", "", "Output CSV file")
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
	writeHeaders(w)

	// Read TSV file
	reader := csv.NewReader(f)

	// Read header row to find column indices
	headers, err := reader.Read()
	if err != nil {
		log.Fatal("Failed to read headers:", err)
	}

	columnIndices := make(map[string]int)
	for i, header := range headers {
		columnIndices[header] = i
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
	white       string
	whiteElo    string
	black       string
	blackElo    string
	result      string
	termination string
	date        string
	moves       []string
	event       string
	timeControl string
}

func getGameInfoHeaders() []string {
	return []string{
		"white",
		"whiteElo",
		"black",
		"blackElo",
		"result",
		"termination",
		"date",
		"event",
		"timeControl",
	}
}

type OutputRow struct {
	playedByWhite  bool
	ply            int
	fromFen        string
	toFen          string
	game           GameInfo
	castlingDiff   CastlingRightsDiffs
	pawnFileDiff   PawnFileCountDiffs
	pieceSightDiff PieceSightCountDiffs
	pieceCountDiff PieceCountDiffs
}

func writeHeaders(w *csv.Writer) {
	headers := []string{
		"played_by_white",
		"ply",
		"from_fen",
		"to_fen",
	}

	headers = append(headers, getGameInfoHeaders()...)
	headers = append(headers, getCastlingDiffHeaders()...)
	headers = append(headers, getPawnFileDiffsHeaders()...)
	headers = append(headers, getPieceSightDiffHeaders()...)
	headers = append(headers, getPieceCountDiffHeaders()...)

	w.Write(headers)
}

func writeOutputRow(w *csv.Writer, out *OutputRow) {
	values := []string{
		strconv.FormatBool(out.playedByWhite),
		strconv.Itoa(out.ply),
		out.fromFen,
		out.toFen,
	}

	values = append(values, getCastlingDiffValues(out.castlingDiff)...)
	values = append(values, getPawnFileDiffsValues(out.pawnFileDiff)...)
	values = append(values, getPieceSightDiffValues(out.pieceSightDiff)...)
	values = append(values, getPieceCountDiffValues(out.pieceCountDiff)...)

	w.Write(values)
}

func processGame(w *csv.Writer, gameInfo *GameInfo) {
	if len(gameInfo.moves) == 0 {
		return
	}

	// Initialize chess game
	game := chess.NewGame()

	// Process each move
	for _, moveStr := range gameInfo.moves {
		// Get position before move
		lastPos := game.Position()

		err := game.PushNotationMove(moveStr, chess.AlgebraicNotation{}, nil)
		if err != nil {
			log.Fatalf("Failed to make move '%s': %v", moveStr, err)
		}

		// Get new position after move
		newPos := game.Position()

		// Write position data
		writePosition(w, gameInfo, lastPos, newPos)
	}
}

func parseGameRecord(record []string, columnIndices map[string]int) GameInfo {
	gameInfo := GameInfo{
		white:       getFieldValue(record, columnIndices, "White"),
		whiteElo:    getFieldValue(record, columnIndices, "WhiteElo"),
		black:       getFieldValue(record, columnIndices, "Black"),
		blackElo:    getFieldValue(record, columnIndices, "BlackElo"),
		result:      getFieldValue(record, columnIndices, "Result"),
		termination: getFieldValue(record, columnIndices, "Termination"),
		date:        getFieldValue(record, columnIndices, "Date"),
		event:       getFieldValue(record, columnIndices, "Event"),
		timeControl: getFieldValue(record, columnIndices, "TimeControl"),
	}

	// Parse moves
	movesStr := getFieldValue(record, columnIndices, "Moves")
	if movesStr != "" {
		gameInfo.moves = strings.Fields(movesStr)
	}

	return gameInfo
}

func writePosition(writer *csv.Writer, gameInfo *GameInfo, lastPos *chess.Position, pos *chess.Position) {
	// Get counts from old position
	pieceCounts := newPieceCounts(lastPos)
	pieceSightCounts := newPieceSightCounts(lastPos)
	pawnFileCounts := newPawnFileCounts(lastPos)
	castlingRights := newCastlingRights(lastPos)

	// Get counts from new position
	newPieceCounts := newPieceCounts(pos)
	newPieceSightCounts := newPieceSightCounts(pos)
	newPawnFileCounts := newPawnFileCounts(pos)
	newCastlingRights := newCastlingRights(pos)

	// Calculate diffs
	pieceCountDiffs := diffPieceCounts(&newPieceCounts, &pieceCounts)
	pieceSightDiffs := diffPieceSightCounts(&newPieceSightCounts, &pieceSightCounts)
	pawnFileCountDiffs := diffPawnFileCounts(&newPawnFileCounts, &pawnFileCounts)
	castlingRightsDiffs := diffCastlingRights(&newCastlingRights, &castlingRights)

	out := OutputRow{
		ply:            lastPos.Ply(),
		playedByWhite:  lastPos.Turn() == chess.White,
		fromFen:        lastPos.String(),
		toFen:          pos.String(),
		game:           *gameInfo,
		castlingDiff:   castlingRightsDiffs,
		pawnFileDiff:   pawnFileCountDiffs,
		pieceSightDiff: pieceSightDiffs,
		pieceCountDiff: pieceCountDiffs,
	}

	writeOutputRow(writer, &out)
}
