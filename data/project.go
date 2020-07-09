package data

import ()

type Project struct {
	ID            int    `gorm:"column:project_id;PRIMARY_KEY"`
	ProjectName   string `gorm:"column:project_name"`
	Description   string `gorm:"column:description"`
	ProjectAvatar string `gorm:"column:project_avatar"`
}

func ProjectsByUserId(userId int) []Project {
	var projects []Project
	Db.Table("project").Select("project.project_id, project_name, description, project_avatar").
		Joins("join user_project on project.project_id = user_project.project_id").
		Where("user_project.user_id = ?", userId).
		Scan(&projects)
	return projects
}

func ProjectById(projectId int) (Project, error) {
	project := Project{ID: projectId}
	err := Db.Find(&project).Error
	return project, err
}

func InsertProject(project Project) int {
	Db.Create(&project)
	return project.ID
}

func UpdateProject(project Project) error {
	err := Db.Save(&project).Error
	return err
}

func DeleteProject(projectId int) error {
	project := Project{ID: projectId}
	err := Db.Delete(&project).Error
	return err
}
