package response

type TicketList struct {
	Project  IdName   `json:"project"`
	Statuses []Status `json:"statuses"`
}

type Status struct {
	Id       int      `json:"id"`
	Progress int      `json:"progress"`
	Name     string   `json:"name"`
	Tickets  []Ticket `json:"tickets,omitempty"`
}

type Ticket struct {
	Id     int    `json:"id"`
	Title  string `json:"title,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

type TicketDetail struct {
	TicketId    int         `json:"ticket_id"`
	Title       string      `json:"title"`
	Explanation string      `json:"explanation"`
	Status      IdName      `json:"status"`
	Worker      IdName      `json:"worker"`
	Reporter    IdName      `json:"reporter"`
	TicketImgs  []TicketImg `json:"ticket_imgs"`
	Comments    []Comment   `json:"comments"`
}

type TicketImg struct {
	Id   int    `json:"id"`
	Path string `json:"path"`
}

type Comment struct {
	Id          int          `json:"id"`
	UserName    string       `json:"user_name"`
	Comment     string       `json:"comment"`
	CommentImgs []CommentImg `json:"comment_imgs"`
}

type CommentImg struct {
	Id   int    `json:"id"`
	Path string `json:"path"`
}

type StatusList struct {
	Statuses []Status `json:"statuses"`
}

type IdName struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type UserProject struct {
	UserId    int `json:"user_id"`
	ProjectId int `json:"project_id"`
}

type ProjectList struct {
	User     User      `json:"user"`
	Projects []Project `json:"projects"`
}

type User struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type Project struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
}
