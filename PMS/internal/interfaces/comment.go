package interfaces
import(
"github.com/Yash-Watchguard/Tasknest/internal/model/comment"
)

type CommentRepo interface {
	ViewAllComments(projectId string) ([]comment.Comment,error)
	UpdateComment(userId string,commentId string,updatedComment string)error
	AddComment(newComment comment.Comment)error
	DeleteComment(userId string, commentId string)error
}