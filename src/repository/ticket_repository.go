package repository

import (
	"github.com/jinzhu/gorm"
	"taskmanage_api/src/domain"
)

type ticketRepository struct {
	Db *gorm.DB
}

func NewTicketRepository(Db *gorm.DB) domain.TicketRepository {
	return &ticketRepository{
		Db: Db,
	}
}

func (tr *ticketRepository) InsertTicket(ticket domain.Ticket) int {
	tr.Db.Create(&ticket)
	return ticket.ID
}

func (tr *ticketRepository) TicketByProjectIdStatusId(projectId int, statusId int) []domain.Ticket {
	var tickets []domain.Ticket
	tr.Db.Table("ticket t").Select("t.ticket_id, project_id, title, explanation, reporter, worker").
		Joins("join ticket_status ts on t.ticket_id = ts.ticket_id").
		Where("t.project_id = ? AND ts.status_id = ?", projectId, statusId).
		Order("t.update_at desc").
		Scan(&tickets)
	return tickets
}

func (tr *ticketRepository) TicketById(ticketId int) (*domain.Ticket, error) {
	ticket := domain.Ticket{}
	err := tr.Db.Select("ticket_id, project_id, title, explanation, reporter, worker").
		Table("ticket").
		Where("ticket_id = ?", ticketId).
		Find(&ticket).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &ticket, err
}

func (tr *ticketRepository) TicketByProjectIdWorker(projectId int, worker int) []domain.Ticket {
	var tickets []domain.Ticket
	tr.Db.Select("ticket_id, project_id, title, explanation, reporter, worker").
		Table("ticket").
		Where("project_id = ? AND worker = ?", projectId, worker).
		Scan(&tickets)
	return tickets
}

func (tr *ticketRepository) UpdateTicket(ticket *domain.Ticket, ticketStatusId int, statusId int) {
	tx := tr.Db.Begin()
	if err := tx.Save(&ticket).Error; err != nil {
		tx.Rollback()
	}
	ticketStatus := domain.TicketStatus{ID: ticketStatusId}
	if err := tx.Model(&ticketStatus).Update("status_id", statusId).Error; err != nil {
		tx.Rollback()
	}
	tx.Commit()
}

func (tr *ticketRepository) TicketImgById(ticketId int) []domain.TicketImg {
	var ticketImgs []domain.TicketImg
	tr.Db.Table("ticket_img").
		Select("ticket_img_id, ticket_id, ticket_img_path").
		Where("ticket_id = ?", ticketId).
		Scan(&ticketImgs)
	return ticketImgs
}

func (tr *ticketRepository) DeleteTicket(ticketId int) {
	ticket := domain.Ticket{ID: ticketId}
	tr.Db.Delete(&ticket)
}

func (tr *ticketRepository) TicketByStatusId(statusId int) []domain.Ticket {
	var tickets []domain.Ticket
	tr.Db.Select("t.ticket_id, project_id, title, explanation, reporter, worker").
		Table("ticket t").
		Joins("join ticket_status ts on t.ticket_id = ts.ticket_id").
		Where("ts.status_id = ?", statusId).
		Scan(&tickets)
	return tickets
}

func (tr *ticketRepository) TicketByIdUserId(ticketId int, userId int) error {
	ticket := domain.Ticket{}
	err := tr.Db.Select("ticket_id, ticket_id, p.project_id, title, explanation, reporter, worker").
		Table("ticket t").
		Joins("join project p on ticket.project_id = p.project_id").
		Joins("join user_project up on p.project_id = up.project_id").
		Where("ticket_id = ? AND user_id = ?", ticketId, userId).
		Find(&ticket).Error
	if gorm.IsRecordNotFoundError(err) {
		return err
	}
	return nil
}