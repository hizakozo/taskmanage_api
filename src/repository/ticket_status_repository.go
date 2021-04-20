package repository

import (
	"github.com/jinzhu/gorm"
	"taskmanage_api/src/domain"
)

type ticketStatusRepository struct {
	Db *gorm.DB
}

func NewTicketStatusRepository(Db *gorm.DB) domain.TicketStatusRepository {
	return &ticketStatusRepository{
		Db: Db,
	}
}

func (tsr *ticketStatusRepository) InsertTicketStatus(ticketStatus domain.TicketStatus) {
	tsr.Db.Create(&ticketStatus)
}

func (tsr *ticketStatusRepository) UpdateTicketStatus(ticketStatusId int, statusId int) {
	ticketStatus := domain.TicketStatus{ID: ticketStatusId}
	tsr.Db.Model(&ticketStatus).
		Update("status_id", statusId)
}

func (tsr *ticketStatusRepository) TicketStatusByTicketId(ticketId int) (*domain.TicketStatus, error) {
	ticketStatus := domain.TicketStatus{}
	err := tsr.Db.Where("ticket_id = ?", ticketId).Find(&ticketStatus).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &ticketStatus, err
}
