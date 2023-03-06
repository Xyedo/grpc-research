package auth

type Tokenizer interface {
	GenerateAccessToken(id string) (string, error)
	GenerateRefreshToken(id string) (string, error)
	ValidateRefreshToken(token string) (string, error)
	ValidateAccessToken(token string) (string, error)
}
