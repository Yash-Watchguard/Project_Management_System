package comment

type Comment struct {
	CommentId string `json:"commentid" db:"commentid"`
	Content   string `json:"content" db:"content"`
	CreatedBy string `json:"createdby" db:"createdby"`
	TaskId    string `json:"taskid" db:"taskid"`
}

type DynamoComment struct {
	PK        string `json:"PK" dynamodbav:"PK"`
	SK        string `json:"SK" dynamodbav:"SK"`
	CommentId string `json:"CommentId" dynamodbav:"CommentId"`
	Content   string `json:"Content" dynamodbav:"Content"`
	CreatedBy string `json:"CreatedBy" dynamodbav:"CreatedBy"`
	TaskId    string `json:"TaskId" dynamodbav:"TaskId"`
}
