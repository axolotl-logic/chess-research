package main

import (
	"log"

	"github.com/corentings/chess/v2"
)

func isDark(sq chess.Square) bool {
	if ((sq / 8) % 2) == (sq % 2) {
		return true
	}
	return false
}

func getFieldValue(record []string, columnIndices map[string]int, fieldName string) string {
	idx, ok := columnIndices[fieldName]
	if !ok || idx >= len(record) {
		log.Fatalf("Unable to find column %s", fieldName)
	}

	return record[idx]
}
