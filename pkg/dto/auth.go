package dto

type Credentials struct {
	Address  string `json:"address"`
	Password string `json:"password"`
}

type Token struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}
