package form

type SignUpForm struct {
	Name        string `json:"name"`
	LoginId     string `json:"login_id"`
	Password    string `json:"password"`
	MailAddress string `json:"mail_address"`
}

type LoginForm struct {
	LoginId  string `json:"login_id"`
	Password string `json:"password"`
}

type LogOutForm struct {
	UserId int `json:"user_id"`
}