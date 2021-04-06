package basic

// UserRegisteredEvent topic:user
type UserRegisteredEvent struct {
	Username string
	Email    string
}

// UserLoggedOutEvent topic:user
type UserLoggedOutEvent struct {
	UserID string
}
