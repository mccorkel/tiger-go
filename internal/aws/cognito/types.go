package cognito

// AuthResponse represents the authentication response
type AuthResponse struct {
	AccessToken  string
	IdToken      string
	TokenType    string
	ExpiresIn    int32
	RefreshToken string
}
