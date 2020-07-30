package response

type LoginResponse struct {
	UserToken string `json:"user_token"`
	UserId    int    `json:"user_id"`
}
type TicketList struct {
	Project  IdName   `json:"project"`
	Statuses []Status `json:"statuses"`
}

type Status struct {
	Id       int      `json:"id"`
	Progress int      `json:"progress"`
	Name     string   `json:"name"`
	Tickets  []Ticket `json:"tickets"`
}

type Ticket struct {
	Id     int    `json:"id"`
	Title  string `json:"title,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

type TicketDetail struct {
	Project     IdName      `json:"project"`
	TicketId    int         `json:"ticket_id"`
	Title       string      `json:"title"`
	Explanation string      `json:"explanation"`
	Status      IdName      `json:"status"`
	Worker      User        `json:"worker,omitempty"`
	Reporter    IdName      `json:"reporter,omitempty"`
	TicketImgs  []TicketImg `json:"ticket_imgs"`
}

type TicketImg struct {
	Id   int    `json:"id"`
	Path string `json:"path"`
}

type Comment struct {
	Id          int          `json:"id"`
	User        IdName       `json:"user"`
	Comment     string       `json:"comment"`
	CommentImgs []CommentImg `json:"comment_imgs,omitempty"`
}

type CommentImg struct {
	Id   int    `json:"id"`
	Path string `json:"path"`
}

type StatusList struct {
	Statuses []Status `json:"statuses"`
}

type IdName struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
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
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	MailAddress string `json:"mail_address,omitempty"`
}

type Project struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
}

type UserList struct {
	Users []User `json:"users"`
}

type CommentList struct {
	Comments []Comment `json:"comments"`
}

type CommentCreate struct {
	TicketId int     `json:"ticket_id"`
	Comment  Comment `json:"comment"`
}

type UserProfile struct {
	User     User             `json:"user"`
	Projects []ProjectTickets `json:"projects"`
}

type ProjectTickets struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Tickets []Ticket `json:"tickets,omitempty"`
}
