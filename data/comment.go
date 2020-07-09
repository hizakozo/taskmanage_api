package data

type Comment struct {
	ID       int    `gorm:"column:comment_id;PRIMARY_KEY"`
	UserId   int    `gorm:"column:user_id"`
	TicketId int    `gorm:"column:ticket_id"`
	Comment  string `gorm:"column:comment"`
}

type CommentImg struct {
	ID             int    `gorm:"column:comment_img_id;PRIMARY_KEY"`
	CommentId      int    `gorm:"column:comment_id"`
	CommentImgPath string `gorm:"column:comment_img_path"`
}

func CommentByTicketId(ticketId int) ([]Comment, error) {
	var comments []Comment
	err := Db.Table("comment").
		Select("comment_id, user_id, ticket_id, comment").
		Where("ticket_id = ?", ticketId).
		Scan(&comments).Error
	return comments, err
}

func CommentImgByCommentId(commentId int) ([]CommentImg, error) {
	var commentImgs []CommentImg
	err := Db.Table("comment_img").
		Select("comment_img_id, comment_id, comment_img_path").
		Where("comment_id = ?", commentId).
		Scan(&commentImgs).Error
	return commentImgs, err
}