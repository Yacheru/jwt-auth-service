package entities

type User struct {
	UserID       string `json:"user_id" db:"uuid" swaggerignore:"true"`
	Nickname     string `json:"nickname" binding:"required" db:"nickname"`
	Email        string `json:"email" binding:"required,email" db:"email"`
	Password     string `json:"password" binding:"required" db:"password"`
	IpAddr       string `json:"ip_addr" db:"ip" swaggerignore:"true"`
	RefreshToken string `json:"refresh_token,omitempty" db:"refresh_token" swaggerignore:"true"`
	ExpiresIn    int64  `json:"expires_at,omitempty" db:"expires_in" swaggerignore:"true"`
}

type UserLogin struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	IpAddress string `json:"ip_address" swaggerignore:"true"`
}
