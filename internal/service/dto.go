package service

type SignUpInput struct {
	Username string
	Email    string
	Password string
	Age      int
	Language string
	Discord  string
	Telegram string
	Region   string
}

type SignInInput struct {
	Email    string
	Password string
}
