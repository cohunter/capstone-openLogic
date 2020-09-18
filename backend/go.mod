module backend

go 1.14

replace (
	local/datastore => ./datastore
	local/google-token-auth => ./google-token-auth
)

require (
	local/datastore v0.0.0
	local/google-token-auth v0.0.0
)
