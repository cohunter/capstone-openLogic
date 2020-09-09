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

type UserWithEmail interface {
	GetEmail() string
}

// IProofStore is the interface through which the main program interacts with the
// datastore. To add functionality, define a method in the interface and on the struct
// which implements the interface.
type IProofStore interface {
	// The main program should call the Close() function when it is done using the datastore.
	Close() error

	// Delete all stored records from the datastore
	Empty() error

	// Get all user attempts of "repo problems" (those specified as repo problems by admin users)
	// The returned array of Proof objects contains only those that meet the conditions:
	// 1) The problem was started by a user selecting it via the 'repo problems' menu
	// 2) The Premise and Conclusion have not changed from the repo problem's Premise and Conclusion
	// 3) The repo problem was submitted by a user who is currently an admin user
	GetAllAttemptedRepoProofs() (error, []Proof)

	// Get public repo problems submitted by admin users, for display to users to choose from
	GetRepoProofs() (error, []Proof)

	// Get all non-completed proofs submitted by a user
	GetUserProofs(user UserWithEmail) (error, []Proof)

	// Get all completed proofs submitted by a user
	GetUserCompletedProofs(user UserWithEmail) (error, []Proof)

	// Save a Proof to the database. The ID will be ignored on input, IDs are created by auto-increment.
	// Updates are done by a UNIQUE index on ProofName and UserSubmitted. New entries with the same
	// ProofName and UserSubmitted overwrite prior entries.
	Store(Proof) error

	// Updates the admin users in the database — used for ensuring that public repo problems
	// were submitted by users who are administrators
	UpdateAdmins(adminUsers map[string]bool)
}
