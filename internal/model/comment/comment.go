package comment

type Comment struct{
	CommentId string `json:"commentid"`
	Content string  `json:"content"`
	CreatedBy string `json:"createdby"`
	TaskId    string  `json:"taskid"`
	CreatorName string `json:"creatorname"`
}