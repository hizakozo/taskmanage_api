package domain

type TicketStatus struct {
	ID       int `gorm:"column:ticket_status_id;PRIMARY_KEY"`
	TicketId int `gorm:"column:ticket_id"`
	StatusId int `gorm:"column:status_id"`
}

type TicketStatusRepository interface {
	InsertTicketStatus(ticketStatus TicketStatus)
	UpdateTicketStatus(ticketStatusId int, statusId int)
	TicketStatusByTicketId(ticketId int) (*TicketStatus, error)
}
