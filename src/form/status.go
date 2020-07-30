package form

type CreateStatusForm struct {
	ProjectId  int    `json:"project_id"  validate:"required"`
	StatusName string `json:"status_name"  validate:"required"`
}

type UpdateStatusForm struct {
	StatusId   int    `json:"status_id"`
	ProjectId  int    `json:"project_id"`
	Progress   int    `json:"progress"`
	StatusName string `json:"status_name"`
}
