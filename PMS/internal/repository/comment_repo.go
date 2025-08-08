package repository

import ("errors"
"os"
"encoding/json"
"github.com/Yash-Watchguard/Tasknest/internal/model/comment"
)

type CommentRepo struct {
	filepPath string
}

func NewCommentRepo() *CommentRepo {
	return &CommentRepo{filepPath: "C:/Users/ygoyal/Desktop/PMS_Project/PMS/internal/data/comment.json"}
}

func (cr *CommentRepo) ViewAllComments(taskId string) ([]comment.Comment, error) {


	data, err := os.ReadFile(cr.filepPath)
	if err != nil {
		return nil, errors.New("failed to read comments file")
	}
    
	var allComments []comment.Comment
	if err := json.Unmarshal(data, &allComments); err != nil {
		return nil, errors.New("failed to parse comments")
	}

	var comments []comment.Comment
	for _, comment := range allComments {
		if comment.TaskId == taskId {
			comments= append(comments, comment)
		}
	}

	return comments, nil
}
func(cr *CommentRepo)UpdateComment(userId string,commentId string,updatedComm string)error{
	commentData,err:=os.ReadFile(cr.filepPath)
	if err!=nil{
        return errors.New("error in reading data")
	}
	var comments []comment.Comment
    if len(commentData) > 0 {
		err = json.Unmarshal(commentData, &comments)
		if err != nil {
			return err
		}
	}
	
	var updatedComment []comment.Comment
	flound:=false
	for _,comment:=range comments{
		if comment.CommentId==commentId {
			if userId!=comment.CreatedBy{
				return errors.New("you are not authorized to update")
			}
			comment.Content=updatedComm
			flound=true
		}
		updatedComment=append(updatedComment, comment)
	}
	if !flound{
		return errors.New("enter valid comment id")
	}

	updatedData, err := json.MarshalIndent(updatedComment, "", "  ")
	if err != nil {
		return err
	}
    err = os.WriteFile(cr.filepPath, updatedData, 0644)
	if err != nil {
		return err
	}
    return nil
}

func (cr *CommentRepo) AddComment(newComment comment.Comment) error {
    commentData, err := os.ReadFile(cr.filepPath)
    if err != nil {
       return errors.New("eror in reading filr")
    }

    var comments []comment.Comment
    if len(commentData)!=0{
    json.Unmarshal(commentData, &comments)
	}
    comments = append(comments, newComment)

    updatedData, err := json.MarshalIndent(comments, "", "  ")
    if err != nil {
        return errors.New("error marshaling")
    }

    if err := os.WriteFile(cr.filepPath, updatedData, 0644); err != nil {
        return errors.New("error in writing")
    }

    return nil
}


func(cr *CommentRepo)DeleteComment(userId string,commentId string)error{
	commentData,err:=os.ReadFile(cr.filepPath)
	if err!=nil{
        return errors.New("error in reading data")
	}
    if len(commentData)==0{
		return errors.New("no comment for delete")
	}
	var comments []comment.Comment

	if err :=json.Unmarshal(commentData,&comments);err!=nil{
		return errors.New("error in unmarshal")
	}

	var updatedComment []comment.Comment
	flound:=false
	for _,comment:=range comments{
		if comment.CommentId==commentId {
			if userId!=comment.CreatedBy{
				return errors.New("you are not authorized to Delete")
			}
			flound=true
		}else{
		updatedComment=append(updatedComment, comment)
		}
	}
	if !flound{
		return errors.New("enter valid comment id")
	}

	updatedData, err := json.MarshalIndent(updatedComment, "", "  ")
	if err != nil {
		return err
	}
    err = os.WriteFile(cr.filepPath, updatedData, 0644)
	if err != nil {
		return err
	}
    return nil
}
