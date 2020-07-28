package data

type User struct {
	ID       int    `gorm:"column:user_id;PRIMARY_KEY"`
	Name     string `gorm:"column:user_name"`
	Avatar   string
	Isdelete int `gorm:"column:is_delete;default:'galeone'"`
}

func UserById(userId int) (User, error) {
	user := User{}
	err := Db.Select("user_id, user_name, avatar").
		Table("user").
		Where("user_id  = ?", userId).
		Find(&user).Error
	return user, err
}

func InsertUser(user User) int {
	Db.Create(&user)
	return user.ID
}

func UserByProjectId(projectId int) []User {
	var users []User
	Db.Select("u.user_id, user_name, avatar").
		Table("user u").
		Joins("join user_project up on u.user_id = up.user_id").
		Where("up.project_id = ?", projectId).
		Scan(&users)
	return users
}
