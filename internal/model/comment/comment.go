package comment

type Comment struct{
	CommentId string `json:"commentid" db:"commentid"`
	Content string  `json:"content" db:"content"`
	CreatedBy string `json:"createdby" db:"createdby"`
	TaskId    string  `json:"taskid" db:"taskid"`
}