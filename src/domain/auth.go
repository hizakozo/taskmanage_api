package domain

type Auth struct {
	ID          int    `gorm:"column:auth_id;PRIMARY_KEY"`
	UserId      int    `gorm:"column:user_id"`
	LoginId     string `gorm:"column:login_id"`
	Password    string
	MailAddress string
}

type AuthRepository interface {
	InsertAuth(auth Auth)
	AuthByLoginId(loginId string) (Auth, error)
	AuthByUserId(userId int) (Auth, error)
	AuthByMailAddress(mailAddress string) (Auth, error)
}