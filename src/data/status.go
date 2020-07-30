package data

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Status struct {
	ID         int    `gorm:"column:status_id;PRIMARY_KEY"`
	ProjectId  int    `gorm:"column:project_id"`
	Progress   int    `gorm:"column:progress"`
	StatusName string `gorm:"column:status_name"`
}

func InsertStatus(status Status) Status {
	Db.Create(&status)
	return status
}

func StatusByProjectId(projectId int) []Status {
	var statuses []Status
	Db.Table("status").
		Select("status_id, project_id, progress, status_name").
		Where("project_id = ?", projectId).
		Order("progress").
		Scan(&statuses)
	return statuses
}

func StatusById(statusId int) (*Status, error) {
	status := Status{ID: statusId}
	err := Db.Select("status_id, project_id, progress, status_name").
		Find(&status).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &status, err
}

func StatusByTicketId(ticketId int) Status {
	status := Status{}
	Db.Table("status s").Select("s.status_id, project_id, progress, status_name").
		Joins("join ticket_status ts on s.status_id = ts.status_id").
		Where("ts.ticket_id = ?", ticketId).
		Find(&status)
	return status
}

func MaxProgressByProjectId(projectId int) int {
	status := Status{}
	Db.Select("status_id, project_id, progress, status_name").
		Table("status").
		Where("project_id = ?", projectId).
		Order("progress desc").Limit(1).
		Find(&status)
	return status.Progress
}

func UpdateStatus(status Status) Status {
	Db.Model(&status).Updates(status)
	return status
}

func StatusByIdProjectId(statusId int, projectId int) (Status, error) {
	status := Status{}
	err := Db.Select("status_id, project_id, progress, status_name").
		Table("status").
		Where("status_id = ? AND project_id = ?", statusId, projectId).
		Find(&status).Error
	return status, err
}

func UpdateProgress(statusId int, progress int) {
	status := Status{ID: statusId}
	Db.Model(&status).Update("progress", progress)
}

func DeleteStatusTransaction(statusId int, statuses []Status, trgProgress int) {
	deleteStatus := Status{ID: statusId}
	tx := Db.Begin()
	if err := tx.Delete(&deleteStatus).Error; err != nil {
		tx.Rollback()
	}
	var err error
	for _, status := range statuses {
		if trgProgress < status.Progress {
			updateStatus := Status{ID: status.ID}
			afterProgress := status.Progress - 1
			err = tx.Model(&updateStatus).Update("progress", afterProgress).Error
		}
	}
	if err != nil {
		fmt.Print(123)
		tx.Rollback()
	}
	tx.Commit()
}
