package dto

// TokenResponseDto used to serialize response with JWT tokens
type TokenResponseDto struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// NewTokenResponseDto generate new TokenResponseDto
func NewTokenResponseDto(accessToken, refreshToken string) *TokenResponseDto {
	return &TokenResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
