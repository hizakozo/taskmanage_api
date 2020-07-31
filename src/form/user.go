package form

type SignUpForm struct {
	Name        string `form:"name" validate:"required"`
	LoginId     string `form:"login_id"  validate:"required"`
	Password    string `form:"password" validate:"required"`
	MailAddress string `form:"mail_address" validate:"required,email"`
	Avatar      string `form:"avatar"`
}

type LoginForm struct {
	LoginId  string `json:"login_id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LogOutForm struct {
	UserId int `json:"user_id"`
}
