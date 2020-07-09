package data

import (
	"github.com/jinzhu/gorm"
	"fmt"
)

type UserProject struct {
	ID int `gorm:"column:user_project_id;PRIMARY_KEY"`
	UserId int `gorm:"column:user_id"`
	ProjectId int `gorm:"column:project_id"`
}

func InsertUserProject(userProject UserProject) {
	Db.Create(&userProject)
}

func UserProjectByUserIdProjectId(userId int, projectId int) error {
	userProject := UserProject{}
	err := Db.Select("user_project_id, user_id, project_id").
	Where("user_id = ? AND project_id = ?", userId, projectId).
	Find(&userProject).Error
	if gorm.IsRecordNotFoundError(err) {
		fmt.Println(err)
	}
	return err
}