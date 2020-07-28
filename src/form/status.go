package form

type CreateStatusForm struct {
	ProjectId  int    `json:"project_id"`
	StatusName string `json:"status_name"`
}

type UpdateStatusForm struct {
	StatusId   int    `json:"status_id"`
	ProjectId  int    `json:"project_id"`
	Progress   int    `json:"progress"`
	StatusName string `json:"status_name"`
}
