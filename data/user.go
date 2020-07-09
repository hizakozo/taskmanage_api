package data

import(
)

type User struct {
    ID   int `gorm:"column:user_id;PRIMARY_KEY"`
    Name string `gorm:"column:user_name"`
    Avatar string
	Isdelete int `gorm:"column:is_delete;default:'galeone'"`
}

func UserById(userId int) (User, error){
	user := User{ID: userId}
	err := Db.Find(&user).Error
	return user, err
}

func InsertUser(user User) int{
	Db.Create(&user)
	return user.ID
}