package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/Yash-Watchguard/Tasknest/internal/constants"
	"github.com/Yash-Watchguard/Tasknest/internal/model/comment"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"github.com/Yash-Watchguard/Tasknest/internal/repository"
	"github.com/Yash-Watchguard/Tasknest/internal/service"
	"github.com/fatih/color"
)

func employeeDashboard(ctx context.Context,user *user.User){

	userRepo:=repository.NewUserRepo()
	taskRepo:=repository.NewTaskRepo()
	commentRepo:=repository.NewCommentRepo()
	empRepo:=repository.NewEmployeeRepo()

	EmployeeService :=service.NewEmpService(userRepo,taskRepo,commentRepo,empRepo)

	for{
        color.Blue(constants.EmployDashbEntry)
		color.Blue("1. View Profile")
		color.Blue("2. View All Assigned Task")
		color.Blue("3  Update Task Status")
		color.Blue("4. Logout")

		var choice int
		fmt.Print(color.CyanString("Enter your choice: "))
		fmt.Scanln(&choice)

		switch choice{
		case 1:
			exist, err := viewEmpProfile(EmployeeService, ctx, user)
			if err != nil {
				color.Red("%v", err)
			}
			if exist {
				return
			}
		case 2:
			err:=viewAssignedTask(EmployeeService,ctx,user.Id)
			if err != nil {
				color.Red("%v", err)
			}
		case 3:
			taskId,err:=GetInput("Enter Task Id : ")
			if err!=nil{
				color.Red("error in getting input")
			}
			err=updateTaskStatus(EmployeeService,ctx,taskId)
			if err != nil {
				color.Red("%v", err)
			}
		case 4:
			color.Green("Logging Out......")
			return 
		}
	}
}

func viewAssignedTask(emp *service.EmployeeService, ctx context.Context, empId string)error{

   tasks,err:=emp.GetAssigenedTask(ctx,empId)
   if err != nil {
		color.Red(" Failed to fetch tasks: %v", err)
		return err
	}

	if len(tasks) == 0 {
		color.Yellow("No tasks Assigned.")
		return nil
	}

	for i, task := range tasks {
		color.Cyan("------------ Task %d ------------", i+1)
		color.Cyan("Task ID        : %v", task.TaskId)
		color.Cyan("Titel          : %v", task.Tile)
		color.Cyan("Description    : %v", task.Description)
		color.Cyan("Task Priority  : %v",task.TaskPriority)
		color.Cyan("Assigned To    : %v", task.AssignedTo)
		color.Cyan("Status         : %v", task.TaskStatus)
		color.Cyan("Deadline       : %v", task.Deadline)

		fmt.Println()
	}
	color.Blue("1.Enter task Id for More operation")
	color.Blue("2.Enter to go back")
	taskId,err:=GetInput("")

	if err!=nil{
		return nil
	}
	if taskId == "" {
		return nil 
	}
    for{
	color.Cyan("1.Update Task Status")
	color.Cyan("2.View All Comments")
	color.Cyan("3.Add Comment")
	color.Cyan("4.Update Comment")
	color.Cyan("5.Delete Comment")
	color.Cyan("6.Return Back")
	var choice int
	fmt.Println(color.BlueString("Enter choice:"))
	fmt.Scan(&choice)

	switch choice{
	case 1:
		err:=updateTaskStatus(emp,ctx,taskId)
        if err!=nil{
			color.Red("%v",err)
		}
	case 2:
		err:=viewAllComments(emp,taskId)
		if err!=nil{
			color.Red("%v",err)
		}
	case 3:
		err:=addNewComment(emp,ctx)
		if err!=nil{
			color.Red("%v",err)
		}
	case 4:
		err:=editComment(emp,ctx)
		if err!=nil{
			color.Red("%v",err)
		}
	case 5:
		err:=removeComment(emp,ctx)
		if err!=nil{
			color.Red("%v",err)
		}
	case 6:
		return nil
	default:
		color.Red("Enter Valid choice")
	}
}
}
func updateTaskStatus(emp *service.EmployeeService,ctx context.Context,taskId string)error{
   
   updatedStatus,err:=GetInput("Enter Updated status (pending/in progress/done) : ")
   if err!=nil{
	color.Red("error in getting input")
   }
   newStatus:=status.GetStatusFromString(updatedStatus)
   err=emp.UpdateTaskStatus(ctx,taskId,newStatus)
   if err!=nil{
	return err
   }
   color.Green("Task Status Updated Successfully")
   color.Blue("press enter to going back")
   fmt.Scanln()
   return nil
}
func viewEmpProfile(emp *service.EmployeeService, ctx context.Context, employee *user.User) (bool, error) {
	userProfiles, err := emp.ViewProfile(ctx, employee.Id)
	if err != nil {
		return false, err
	}

	if len(userProfiles) == 0 {
		return false, errors.New("no user profile")
	}

	user := userProfiles[0]
	color.Cyan("----------- %s Profile -----------", user.Name)
	color.Yellow("ID           : %s", user.Id)
	color.Yellow("Name         : %s", user.Name)
	color.Yellow("Email        : %s", user.Email)
	color.Yellow("Phone Number : %s", user.PhoneNumber)
	color.Yellow("Role         : %v", user.Role)
	color.Cyan("----------------------------------")

	color.Blue("1. Update Profile")
	color.Blue("2. Delete Profile")
	color.Blue("3. Go Back")

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		err := updateEmpProfile(emp, ctx, &user)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		return false, nil

	case 2:
		err := emp.DeleteEmp(ctx, user.Id)
		if err != nil {
			fmt.Printf("%v\n", err)
			return false, err
		}
		color.Red("Profile Deleted Successfully")
		return true, nil

	case 3:
		return false, nil

	default:
		color.Red("Invalid choice.")
		return false, nil
	}
}

func updateEmpProfile(emp*service.EmployeeService, ctx context.Context, user *user.User) error {
	for {
		color.Blue("1. Update Name")
		color.Blue("2. Update Email")
		color.Blue("3. Update Password")
		color.Blue("4. Update Contact")
		color.Blue("5. Go Back")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			name, err := GetInput("Enter Name:")
			if err != nil {
				color.Red("Error reading name: %v", err)
				continue
			}
			user.Name = name

		case 2:
			email, err := GetValidEmail()
			if err != nil {
				color.Red("Invalid email: %v", err)
				continue
			}
			user.Email = email

		case 3:
			password, err := GetValidPassword()
			if err != nil {
				color.Red("Invalid password: %v", err)
				continue
			}
			user.Password = password

		case 4:
			contact, err := GetValidPhoneNumber()
			if err != nil {
				color.Red("Invalid phone number: %v", err)
				continue
			}
			user.PhoneNumber = contact

		case 5:
			return nil

		default:
			color.Red("Invalid choice.")
			continue
		}

		err := emp.UpdateProfile(user.Id, ctx, user.Name, user.Email, user.Password, user.PhoneNumber)
		if err != nil {
			color.Red("Update failed: %v", err)
		} else {
			color.Green("User updated successfully!")
		}
		
	}
}

func removeComment(emp *service.EmployeeService, ctx context.Context)error{
	commentId,err:=GetInput("Enter Comment id for delete the comment:")
	if err!=nil{
		return nil
	}
	err= emp.DeleteComment(ctx,commentId)
    
	if err!=nil{
		return err
		
	}
	color.Green("comment deleted succesfully")
	color.Blue("press enter to going back")
	fmt.Scanln()
    return nil
}

func addNewComment(emp *service.EmployeeService, ctx context.Context)error{
	
	var taskId string
	for{
	taskId, err := GetInput("Enter Task ID to comment on:")
	if err != nil || taskId == "" {
		color.Red("Invalid task ID")
	}else{
		break
	}
    }

	content, err := GetInput("Enter your comment:")
	if err != nil {
		fmt.Print(err)
	}

	createdBy:=ctx.Value(ContextKey.UserId).(string)
    
	commentId := GenerateUUID()

	newComment := comment.Comment{
		CommentId: commentId,
		Content:   content,
		CreatedBy: createdBy,
		TaskId:    taskId,
	}

	err = emp.AddComment(newComment)
	if err != nil {
		return err
	}

	color.Green("✅ Comment added successfully!")
	color.Blue("press enter to going back")
	fmt.Scanln()
	return nil
}

func viewAllComments(emp *service.EmployeeService,taskId string)error{
    
    comments,err:=emp.ViewAllComment(taskId)
    if err!=nil{
		return err
	}

	if len(comments)==0{
		return errors.New("no comments found for this Task")
	}

	for _,comment:=range comments{
		color.Blue("Comment for %v (Created by-%v)",comment.TaskId,comment.CreatedBy)
		color.Cyan("%v",comment.Content)
		color.Blue("----------------------------------------------------------------")
	}
    color.Green("Press Enter to return to the previous menu...")
	fmt.Scanln()
	return nil
}

func editComment(emp *service.EmployeeService,ctx context.Context)error{
	commentId,err:=GetInput("Enter Comment id for update the comment:")
	if err!=nil{
		return nil
	}
	updatedComment,err:=GetInput("Enter Updated Comment:")
    if err!=nil{
		return nil
	}
	err=emp.UpdateComment(ctx,commentId,updatedComment)
    if err!=nil{
		return err
	}
	color.Green("Comment Updated Successfully......")
	color.Cyan("Press enter for going back")
	fmt.Scanln()
	return nil
}