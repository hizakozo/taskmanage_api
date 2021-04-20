package domain

type UserProject struct {
	ID        int `gorm:"column:user_project_id;PRIMARY_KEY"`
	UserId    int `gorm:"column:user_id"`
	ProjectId int `gorm:"column:project_id"`
}

type UserProjectRepository interface {
	InsertUserProject(userProject UserProject)
	UserProjectByUserIdProjectId(userId int, projectId int) []UserProject
}
