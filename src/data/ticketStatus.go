package data

import "github.com/jinzhu/gorm"

type TicketStatus struct {
	ID       int `gorm:"column:ticket_status_id;PRIMARY_KEY"`
	TicketId int `gorm:"column:ticket_id"`
	StatusId int `gorm:"column:status_id"`
}

func InsertTicketStatus(ticketStatus TicketStatus) {
	Db.Create(&ticketStatus)
}

func UpdateTicketStatus(ticketStatusId int, statusId int) {
	ticketStatus := TicketStatus{ID: ticketStatusId}
	Db.Model(&ticketStatus).
		Update("status_id", statusId)
}

func TicketStatusByTicketId(ticketId int) (*TicketStatus, error) {
	ticketStatus := TicketStatus{}
	err := Db.Where("ticket_id = ?", ticketId).Find(&ticketStatus).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &ticketStatus, err
}
