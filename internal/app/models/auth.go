package models

type User struct {
	Login    string
	Password string
}

type Token struct {
	AccessToken  string
	RefreshToken string
}
