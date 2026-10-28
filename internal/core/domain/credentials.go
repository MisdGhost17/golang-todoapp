package domain

type Credentials struct {
	UserID       int
	PasswordHash string
}

func NewCredentials(
	userID int,
	passwordHash string,
) Credentials {
	return Credentials{
		UserID:       userID,
		PasswordHash: passwordHash,
	}
}
