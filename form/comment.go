package form

type CreateCommentForm struct {
	TicketId int    `json:"ticket_id"`
	Comment  string `json:"comment"`
}
