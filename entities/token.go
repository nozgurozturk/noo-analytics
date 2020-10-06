package entities

type Token struct {
	Token   string
	Uuid    string
	Expires int64
	UserId  string
	Role    string
	Type    string
}

type TokenResponse struct {
	AccessToken  *Token `json:"at"`
	RefreshToken *Token `json:"rt"`
}
