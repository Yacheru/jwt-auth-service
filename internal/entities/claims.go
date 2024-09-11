package entities

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	ExpiresAt int64  `json:"expires_at"`
	IpAddress string `json:"ip_address"`
	UserID    string `json:"user_id"`
}
