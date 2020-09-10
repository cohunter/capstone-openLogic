package datastore

import (
	"log"
	"os"
	"reflect"
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

// Make Proof implement the UserWithEmail interface for testing
func (p Proof) GetEmail() string {
	return p.UserSubmitted
}

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

// Check if two Proof structs are the same, excluding the Id and TimeSubmitted fields
func proofIsEqual(a, b Proof) bool {
	valueA := reflect.ValueOf(a)
	valueB := reflect.ValueOf(b)

	for i := 0; i < valueA.NumField(); i++ {
		log.Printf("%+v", valueA.Field(i).Type())
		switch field := reflect.TypeOf(a).Field(i).Name; field {
		case "Id":
			// pass
		case "TimeSubmitted":
			// pass
		default:
			if valueA.FieldByName(field).String() != valueB.FieldByName(field).String() {
				return false
			}
		}
	}
	return true
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

			err, proofs := proofStore.GetUserProofs(tc.testData)
			if err != nil {
				t.Error(err)
			}
			if len(proofs) == 0 {
				t.Error("No proof retrieved after Store")
			}
			for _, proof := range proofs {
				if proofIsEqual(tc.testData, proof) {
					return
				}
			}
			t.Error("No matching proof retrieved after Store")
		})
	}
}