package entities

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type Session struct {
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}
