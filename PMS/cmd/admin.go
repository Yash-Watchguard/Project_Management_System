package main

import (
	"fmt"
     "errors"
	"github.com/Yash-Watchguard/Tasknest/model"
	"github.com/Yash-Watchguard/Tasknest/repository"
	"github.com/Yash-Watchguard/Tasknest/service"
	"github.com/fatih/color"
)

func AdminDashboard(admin *model.User) {
	// Repository Initialization
	userRepo := repository.NewUserRepo()
	projectRepo := repository.NewProjectRepo()

	// Inject dependencies into AdminService
	adminService := service.NewAdminServices(userRepo,projectRepo)

	for {
		color.Blue("\n------------------ Admin Dashboard ------------------")
		color.Blue("1. View Profile")
		color.Blue("2. View All Users")
		color.Blue("3. Delete User")
		color.Blue("4. Add New Project")
		color.Blue("5. View All Projects")
		color.Blue("6. Delete Project")
		color.Blue("7. Logout")

		var choice int
		fmt.Print("Enter your choice: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			adminService.ViewProfile(admin)

		case 2:
			adminService.ViewAllUsers()

		case 3:
			err := deleteUser(adminService)
			if err != nil {
				color.Red("Error: %v", err)
			}

		case 4:
			err := addNewProject(adminService,*admin)
			if err != nil {
				color.Red("%v", err)
			}

		case 5:
			err := viewAllProjects(adminService,admin)
			if err != nil {
				color.Red("Error: %v", err)
			}
			 fmt.Println("Press ENTER to return to dashboard...")
	         fmt.Scanln()

		case 6:
			err := deleteProject(adminService,admin)
			if err != nil {
				color.Red("Error: %v", err)
			}
			 fmt.Println("Press ENTER to return to dashboard...")
	         fmt.Scanln()


		case 7:
			color.Green("Logging out...")
		default:
			color.Red("Invalid choice. Please try again.")
		}
	}
}
func deleteUser(ad *service.AdminService)error{
	ad.ViewAllUsers()
	fmt.Println("Enter User Id of user:")
	var userId string
	_,err:=fmt.Scanln(&userId)

	if err!=nil{
       return err
	}
	err=ad.DeleteUser(userId)
	if err != nil {
		return err
	}

	color.Green("User deleted successfully!")
	fmt.Println("Press ENTER to return to dashboard...")
	fmt.Scanln()

	return nil
}

func addNewProject(ad *service.AdminService,admin model.User)error{
	createdBy:=admin.Id
	projectId:=GenerateUUID()
	projectName,err:=GetInput("Enter Project Name:")
	if err!=nil{
		return err
	}

	projectDescription,err:=GetInput("Enter Project Description")
	if err!=nil{
		return err
	}
    deadline,_:=GetInput("Enter deadline (YYYY-MM-DD):")
	actualdeadline,err:=TimeParser(deadline)
	if err!=nil{
		return err
	}

	color.Blue("select Manager from Given list 👇")
	err=ad.GetAllManager()
	if err!=nil{
		color.Red("%s",err)
		return errors.New("Project Add Faild")
	}
	var managerId string
	color.Blue("Enter Manager Id:")
	fmt.Scanln(&managerId)

	project:=model.Project{
        ProjectId: projectId,
		ProjectName:projectName ,
		ProjectDescription: projectDescription,
		Deadline: actualdeadline,
		CreatedBy:createdBy,
		AssignedManager: managerId,
	}
    err=ad.AddProject(project)
	if err != nil {
		color.Red("Error adding project: %v", err)
	} else {
		color.Green("✅ Project added successfully!")
	}
	fmt.Println("Press ENTER to return to dashboard...")
	fmt.Scanln()
	return nil
}

func viewAllProjects(ad *service.AdminService,admin *model.User)error{
	 var projects []model.Project
     projects,err:=ad.ViewAllProjects()
     
	 if err!=nil{
		return err
	 }
	 counter:=1
	 for _,project:=range projects{
		if project.CreatedBy==admin.Id{
        color.Cyan("----------------%d----------------",counter)
		color.Cyan("Project Name: %v\n",project.ProjectName)
		color.Cyan("Project Id: %v\n",project.ProjectId)
		color.Cyan("Project Description: %v\n",project.ProjectDescription)
		color.Cyan("Project Deadline: %v\n",project.Deadline)
		color.Cyan("Assigned To: %v\n",project.AssignedManager)
		counter++
		}
	 }
	 return nil
}
func deleteProject(ad *service.AdminService,admin *model.User)error{
	color.Green("---------Project List👇----------")
	_=viewAllProjects(ad,admin)
	var proId string
	color.Blue("Enter the Project id which you want to delete")
    fmt.Scanln(&proId)
	err:= ad.DeleteProject(proId)
	color.Green("Project Deleted Succesfully ✅")
	return err

}



