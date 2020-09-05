package datastore

import (
	"log"
	"os"
	"testing"
)

const (
	// Use an in-memory database for running tests
	// Must specify cache=shared to prevent multiple connections from getting different DBs
	test_dsn = "file::memory:?cache=shared"
)

var (
	proofStore *SQLProofStore
)

func getNewDS(t *testing.T) *SQLProofStore {
	ds, err := NewSQLite(test_dsn)
	if err != nil {
		t.Fatal(err)
	}
	return ds
}

func TestMain(m *testing.M) {
	log.Print("main")
	var err error
	proofStore, err = NewSQLite(test_dsn)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

func TestEmpty(t *testing.T) {
	ds := getNewDS(t)
	err := ds.Empty()
	if err != nil {
		t.Error(err)
	}
	return
}

func TestStore(t *testing.T) {
	testCases := []struct{
		testName string
		testData Proof
	}{
		{
		"Test Proof 1", Proof{
			EntryType: "proof",
			UserSubmitted: "cohunter@csumb.edu",
			ProofName: "Test Proof 1",
			ProofType: "prop",
			Premise: []string{"A", "A -> B"},
			Logic: []string{},
			Rules: []string{},
			ProofCompleted: "false",
			Conclusion: "B",
			TimeSubmitted: "2020-07-29T04:27:32.592+0000",
			RepoProblem: "false",
			},
		},
		{
			"Test Empty Proof", Proof{},
		},
		{
			"Test Proof with ID #7", Proof{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			err := proofStore.Store(tc.testData)
			if err != nil {
				t.Error(err)
			}
		})
	}
}