package auth

// 刷新 Access Token
func RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	return GenerateAccessToken(claims.Username)
}
