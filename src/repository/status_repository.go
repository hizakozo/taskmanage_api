package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"taskmanage_api/src/domain"
)

type statusRepository struct {
	Db *gorm.DB
}

func NewStatusRepository(Db *gorm.DB) domain.StatusRepository {
	return &statusRepository{
		Db: Db,
	}
}

func (sr *statusRepository) InsertStatus(status domain.Status) domain.Status {
	sr.Db.Create(&status)
	return status
}

func (sr *statusRepository) StatusByProjectId(projectId int) []domain.Status {
	var statuses []domain.Status
	sr.Db.Table("status").
		Select("status_id, project_id, progress, status_name").
		Where("project_id = ?", projectId).
		Order("progress").
		Scan(&statuses)
	return statuses
}

func (sr *statusRepository) StatusById(statusId int) (*domain.Status, error) {
	status := domain.Status{ID: statusId}
	err := sr.Db.Select("status_id, project_id, progress, status_name").
		Find(&status).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &status, err
}

func (sr *statusRepository) StatusByTicketId(ticketId int) domain.Status {
	status := domain.Status{}
	sr.Db.Table("status s").Select("s.status_id, project_id, progress, status_name").
		Joins("join ticket_status ts on s.status_id = ts.status_id").
		Where("ts.ticket_id = ?", ticketId).
		Find(&status)
	return status
}

func (sr *statusRepository) MaxProgressByProjectId(projectId int) int {
	status := domain.Status{}
	sr.Db.Select("status_id, project_id, progress, status_name").
		Table("status").
		Where("project_id = ?", projectId).
		Order("progress desc").Limit(1).
		Find(&status)
	return status.Progress
}

func (sr *statusRepository) UpdateStatus(status domain.Status) domain.Status {
	sr.Db.Model(&status).Updates(status)
	return status
}

func (sr *statusRepository) StatusByIdProjectId(statusId int, projectId int) (domain.Status, error) {
	status := domain.Status{}
	err := sr.Db.Select("status_id, project_id, progress, status_name").
		Table("status").
		Where("status_id = ? AND project_id = ?", statusId, projectId).
		Find(&status).Error
	return status, err
}

func (sr *statusRepository) UpdateProgress(statusId int, progress int) {
	status := domain.Status{ID: statusId}
	sr.Db.Model(&status).Update("progress", progress)
}

func (sr *statusRepository) DeleteStatusTransaction(statusId int, statuses []domain.Status, trgProgress int) {
	deleteStatus := domain.Status{ID: statusId}
	tx := sr.Db.Begin()
	if err := tx.Delete(&deleteStatus).Error; err != nil {
		tx.Rollback()
	}
	var err error
	for _, status := range statuses {
		if trgProgress < status.Progress {
			updateStatus := domain.Status{ID: status.ID}
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
