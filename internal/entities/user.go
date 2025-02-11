package entities

const (
	RoleUser  string = "user"
	RoleAdmin string = "admin"
)

type User struct {
	ID             int
	Username       string
	HashedPassword string
	Role           string
}

type UserData struct {
	Username       string
	HashedPassword string
	Role           string
}
