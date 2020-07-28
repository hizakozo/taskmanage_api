package form

type SignUpForm struct {
	Name        string `json:"name"  validate:"required"`
	LoginId     string `json:"login_id" form:"login_id"  validate:"required"`
	Password    string `json:"password"  validate:"required"`
	MailAddress string `json:"mail_address" form:"mail_address"  validate:"required,email"`
	Avatar      string `json:"avatar"`
}

type LoginForm struct {
	LoginId  string `json:"login_id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LogOutForm struct {
	UserId int `json:"user_id"`
}
