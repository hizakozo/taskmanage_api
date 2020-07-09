package data

import ()

type Ticket struct {
	ID          int    `gorm:"column:ticket_id;PRIMARY_KEY"`
	ProjectId   int    `gorm:"column:project_id"`
	Title       string `gorm:"column:title"`
	Explanation string `gorm:"column:explanation"`
	Reporter    int    `gorm:"column:reporter"`
	Worker      int    `gorm:"column:worker"`
}

type TicketImg struct {
	ID            int    `gorm:"column:ticket_img_id;PRIMARY_KEY"`
	TicketId      int    `gorm:"column:ticket_id"`
	TicketImgPath string `gorm:"column:ticket_img_path"`
}

func InsertTicket(ticket Ticket) int {
	Db.Create(&ticket)
	return ticket.ID
}

func TicketByProjectIdStatusId(projectId int, statusId int) []Ticket {
	var tickets []Ticket
	Db.Table("ticket").Select("ticket.ticket_id, project_id, title, explanation, reporter, worker").
		Joins("join ticket_status on ticket.ticket_id = ticket_status.ticket_id").
		Where("ticket.project_id = ? AND ticket_status.status_id = ?", projectId, statusId).
		Scan(&tickets)
	return tickets
}

func TicketById(ticketId int) (Ticket, error) {
	ticket := Ticket{ID: ticketId}
	err := Db.Select("ticket_id, project_id, title, explanation, reporter, worker").
		Find(&ticket).Error
	return ticket, err
}

func UpdateTicket(ticket Ticket) error {
	err := Db.Save(&ticket).Error
	return err
}

func TicketImgById(ticketId int) ([]TicketImg, error) {
	var ticketImgs []TicketImg
	err := Db.Table("ticket_img").
		Select("ticket_img_id, ticket_id, ticket_img_path").
		Where("ticket_id = ?", ticketId).
		Scan(&ticketImgs).Error
	return ticketImgs, err
}

func DeleteTicket(ticketId int) error {
	ticket := Ticket{ID: ticketId}
	err := Db.Delete(&ticket).Error
	return err
}

func TicketByStatusId(statusId int) []Ticket {
	var tickets []Ticket
	Db.Select("t.ticket_id, project_id, title, explanation, reporter, worker").
		Table("ticket t").
		Joins("join ticket_status ts on t.ticket_id = ts.ticket_id").
		Where("ts.status_id = ?", statusId).
		Scan(&tickets)
	return tickets
}