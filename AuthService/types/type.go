package types

type LoginResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
}

