package requests

type SignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RefreshToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
