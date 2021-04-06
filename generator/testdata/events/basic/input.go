package basic

type UserRegisteredEvent struct {
	Username string
	Email    string
}

type UserLoggedOutEvent struct {
	UserID string
}
