package main

import (
	"flag"
	"log"
	"math/rand"
	"os"

	"github.com/corentings/chess/v2"
	"github.com/google/uuid"
)

func main() {
	batchId := uuid.New().String()
	batchSize := flag.Int("batch-size", 1_000, "size of batch")
	var outputFile = flag.String("output", "", "Output PGN file")
	flag.Parse()

	if *outputFile == "" {
		log.Fatal("Missing -output for pgn output path")
	}

	f, err := os.Create(*outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for range *batchSize {
		gameId := uuid.New().String()

		game := chess.NewGame()
		game.AddTagPair("Event", "Random Games")
		game.AddTagPair("Source", "random-games/main.go")
		game.AddTagPair("BatchID", batchId)
		game.AddTagPair("GameID", gameId)

		// generate moves until game is over
		for game.Outcome() == chess.NoOutcome {
			draws := game.EligibleDraws()
			for _, draw := range draws {
				if draw != chess.DrawOffer {
					game.Draw(draw)
				}
			}
			// select a random move
			moves := game.ValidMoves()
			move := moves[rand.Intn(len(moves))]
			game.Move(&move, nil)
		}

		game.AddTagPair("Termination", game.Method().String())

		if _, err := f.WriteString(game.String() + "\n"); err != nil {
			log.Fatalf("Error writing game: %v", err)
		}
	}
}
