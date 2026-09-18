//go:build cgo
// +build cgo

package sqlite3

import (
	"database/sql"
	"encoding/binary"
	"math"
	"testing"
)

func TestVec0IsAvailableOnNewConnections(t *testing.T) {
	database, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer database.Close()

	if _, err := database.Exec("CREATE VIRTUAL TABLE vectors USING vec0(embedding float[2])"); err != nil {
		t.Fatalf("create vec0 table: %v", err)
	}
	for rowID, embedding := range [][]float32{{0, 0}, {1, 1}} {
		if _, err := database.Exec("INSERT INTO vectors(rowid, embedding) VALUES (?, ?)", rowID+1, serializeFloat32(embedding)); err != nil {
			t.Fatalf("insert vector: %v", err)
		}
	}

	var nearestRowID int
	if err := database.QueryRow("SELECT rowid FROM vectors WHERE embedding MATCH ? ORDER BY distance LIMIT 1", serializeFloat32([]float32{0.9, 0.9})).Scan(&nearestRowID); err != nil {
		t.Fatalf("query nearest vector: %v", err)
	}
	if nearestRowID != 2 {
		t.Errorf("nearest vector row ID = %d, want 2", nearestRowID)
	}
}

func TestFTS5IsAvailableOnNewConnections(t *testing.T) {
	database, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer database.Close()

	if _, err := database.Exec("CREATE VIRTUAL TABLE documents USING fts5(content)"); err != nil {
		t.Fatalf("create FTS5 table: %v", err)
	}
}

func serializeFloat32(vector []float32) []byte {
	serialized := make([]byte, len(vector)*4)
	for index, value := range vector {
		binary.LittleEndian.PutUint32(serialized[index*4:], math.Float32bits(value))
	}
	return serialized
}
