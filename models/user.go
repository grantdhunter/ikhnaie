package models

type User struct {
}

var users = make(map[string]*User)

func GetUser(token string) *User {
	return users[token]
}

func SetUser(token string, user *User) {
	users[token] = user
}
