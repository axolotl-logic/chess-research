package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/corentings/chess/v2"
)

func main() {
	var inputFile = flag.String("input", "", "Input PGN file")
	var outputFile = flag.String("output", "", "Output CSV file")
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
	scanner := chess.NewScanner(f)

	// Read all games
	for scanner.HasNext() {
		game, err := scanner.ParseNext()
		if err != nil {
			log.Fatalf("Failed to parse game: %v", err)
		}

		for _, pos := range game.Positions() {
			writePosition(w, pos)
		}
	}
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

func writePosition(writer *csv.Writer, pos *chess.Position) {
	counts := newPieceCounts(pos)
	record := []string{
		// white_to_play
		strconv.FormatBool(pos.Turn() == chess.White),
		// white_castle_kingside_rights
		strconv.FormatBool(pos.CastleRights().CanCastle(chess.White, chess.Side(chess.KingSide))),
		// white_castle_queenside_rights
		strconv.FormatBool(pos.CastleRights().CanCastle(chess.White, chess.Side(chess.QueenSide))),
		// black_castle_kingside_rights
		strconv.FormatBool(pos.CastleRights().CanCastle(chess.Black, chess.Side(chess.KingSide))),
		// black_castle_queenside_rights
		strconv.FormatBool(pos.CastleRights().CanCastle(chess.Black, chess.Side(chess.QueenSide))),
		// half_move_clock
		strconv.Itoa(pos.HalfMoveClock()),
		// ply
		strconv.Itoa(pos.Ply()),
		// fen
		pos.String(),
		// white_queen_count
		strconv.Itoa(counts.whiteQueenCount),
		// white_knight_count
		strconv.Itoa(counts.whiteKnightCount),
		// white_pawn_count
		strconv.Itoa(counts.whitePawnCount),
		// white_rook_count
		strconv.Itoa(counts.whiteRookCount),
		// white_dark_bishop_count
		strconv.Itoa(counts.whiteDarkBishopCount),
		//white_light_bishop_count
		strconv.Itoa(counts.whiteLightBishopCount),
		// black_queen_count
		strconv.Itoa(counts.blackQueenCount),
		// black_knight_count
		strconv.Itoa(counts.blackKnightCount),
		// black_pawn_count
		strconv.Itoa(counts.blackPawnCount),
		// black_rook_count
		strconv.Itoa(counts.blackRookCount),
		// black_dark_bishop_count
		strconv.Itoa(counts.blackDarkBishopCount),
		// black_light_bishop_count
		strconv.Itoa(counts.blackLightBishopCount),
	}

	if err := writer.Write(record); err != nil {
		log.Fatalf("Error writing record: %v", err)
	}
}

func isDark(sq chess.Square) bool {
	if ((sq / 8) % 2) == (sq % 2) {
		return true
	}
	return false
}
