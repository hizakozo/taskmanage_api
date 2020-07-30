package form

type CreateProjectForm struct {
	ProjectName string `json:"project_name"  validate:"required"`
	Description string `json:"description"`
}

type UpdateProjectForm struct {
	ProjectId     int    `json:"project_id"  validate:"required"`
	ProjectName   string `json:"project_name"  validate:"required"`
	Description   string `json:"description"`
	ProjectAvatar string `json:"project_avatar"`
}

type InviteProjectForm struct {
	ProjectId   int    `json:"project_id"`
	MailAddress string `json:"mail_address"`
}

type JoinProjectForm struct {
	Token string `json:"token"`
}
