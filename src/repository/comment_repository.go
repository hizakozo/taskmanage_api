package repository

import (
	"github.com/jinzhu/gorm"
	"taskmanage_api/src/domain"
)

type commentRepository struct {
	Db *gorm.DB
}

func NewCommentRepository(Db *gorm.DB) domain.CommentRepository {
	return &commentRepository{
		Db: Db,
	}
}

func (cr *commentRepository)CommentByTicketId(ticketId int) []domain.Comment {
	var comments []domain.Comment
	cr.Db.Table("comment").
		Select("comment_id, user_id, ticket_id, comment").
		Where("ticket_id = ?", ticketId).
		Scan(&comments)
	return comments
}

func (cr *commentRepository)CommentImgByCommentId(commentId int) []domain.CommentImg {
	var commentImgs []domain.CommentImg
	cr.Db.Table("comment_img").
		Select("comment_img_id, comment_id, comment_img_path").
		Where("comment_id = ?", commentId).
		Scan(&commentImgs)
	return commentImgs
}

func (cr *commentRepository)UpdateComment(commentId int, comment string) {
	updateComment := domain.Comment{ID: commentId}
	cr.Db.Model(&updateComment).Update("comment", comment)
}

func (cr *commentRepository)InsertComment(comment domain.Comment) domain.Comment {
	cr.Db.Create(&comment)
	return comment
}