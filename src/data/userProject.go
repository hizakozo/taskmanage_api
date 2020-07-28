package data

type UserProject struct {
	ID        int `gorm:"column:user_project_id;PRIMARY_KEY"`
	UserId    int `gorm:"column:user_id"`
	ProjectId int `gorm:"column:project_id"`
}

func InsertUserProject(userProject UserProject) {
	Db.Create(&userProject)
}

func UserProjectByUserIdProjectId(userId int, projectId int) []UserProject {
	var userProject []UserProject
	Db.Select("user_project_id, user_id, project_id").
		Table("user_project").
		Where("user_id = ? AND project_id = ?", userId, projectId).
		Scan(&userProject)
	return userProject
}
