package http

type signUpRequest struct {
	Username string  `json:"username" binding:"required,min=3,max=50"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=8"`
	Age      int     `json:"age" binding:"required,gte=13,lte=100"`
	Language string  `json:"language" binding:"required,len=2"` //типо eu
	Region   string  `json:"region" binding:"required"`
	Discord  *string `json:"discord" binding:"omitempty,max=100"`
	Telegram *string `json:"telegram" binding:"omitempty,max=100"`
}

type signInRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Region   string `json:"region"`
}
