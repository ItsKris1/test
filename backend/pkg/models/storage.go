package models

// Repositories contains all the repo structs
type Repositories struct {
	UserRepo    UserRepository
	SessionRepo SessionRepository
	GroupRepo   GroupRepository
	PostRepo    PostRepository
}
