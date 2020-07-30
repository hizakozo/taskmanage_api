package form

type CreateTicketForm struct {
	ProjectId   int    `json:"project_id"`
	StatusId    int    `json:"status_id"`
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	Worker      int    `json:"worker"`
}

type ChangeStatusForm struct {
	ProjectId int `json:"project_id"`
	TicketId  int `json:"ticket_id"`
	StatusId  int `json:"status_id"`
}

type UpdateTicketForm struct {
	TicketId    int    `json:"ticket_id"`
	Title       string `json:"title"`
	StatusId    int    `json:"status_id"`
	Explanation string `json:"explanation"`
	Worker      int    `json:"worker"`
}
