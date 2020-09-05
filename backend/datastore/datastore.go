package datastore

// The Proof struct is used for input/output of proofs to the datastore.
// The Id field is ignored on input; proofs will be stored with an
// auto-increment ID, or updated if there is an existing proof with the
// same values for UserSubmitted and ProofName.
type Proof struct {
	Id             string   // SQL ID
	EntryType      string   // 'proof'
	UserSubmitted  string	// Used for results, ignored on user input
	ProofName      string   // user-chosen name (repo problems start with 'Repository - ')
	ProofType      string   // 'prop' (propositional/tfl) or 'fol' (first order logic)
	Premise        []string // premises of the proof; an array of WFFs
	Logic          []string // body of the proof; a JSON-encoded string
	Rules          []string // deprecated; now always an empty string
	ProofCompleted string   // 'true', 'false', or 'error'
	Conclusion     string   // conclusion of the proof
	RepoProblem    string   // 'true' if problem started from a repo problem, else 'false'
	TimeSubmitted  string
}

//type ProofStore interface {
//	GetByUser(string) Proof
//}

type UserWithEmail interface {
	GetEmail() string
}

// IProofStore is the interface through which the main program interacts with the
// datastore. To add functionality, define a method in the interface and on the struct
// which implements the interface.
type IProofStore interface {
	// The main program should call the Close() function when it is done using the datastore.
	Close() error
	Empty() error
	GetAllAttemptedRepoProofs() (error, []Proof)
	GetRepoProofs() (error, []Proof)
	GetUserProofs(user UserWithEmail) (error, []Proof)
	GetUserCompletedProofs(user UserWithEmail) (error, []Proof)
	Store(Proof) error
	UpdateAdmins(adminUsers map[string]bool)
}
