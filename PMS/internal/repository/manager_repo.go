package repository

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
)

type manngerRepo struct {
	dbProject string
	dbTask    string
}

func NewManagerRepo() *manngerRepo {
	return &manngerRepo{dbProject: "C:/Users/ygoyal/Desktop/PMS_Project/Pms/internal/data/project.json",
		dbTask: "C:/Users/ygoyal/Desktop/PMS_Project/Pms/internal/data/task.json"}

}

func(manager *manngerRepo)ViewAssignedProject(ctx context.Context)([]project.Project,error){
	managerId:=ctx.Value(ContextKey.UserId).(string)
	var projects []project.Project
	data,err:=os.ReadFile(manager.dbProject)
	if err!=nil{
		return projects,errors.New("error in getting projects")
	}

	err=json.Unmarshal(data,&projects)
	if err!=nil{
		return projects,errors.New("error in getting projects")
	}
    var ManagerProjects []project.Project
	for _,value:=range projects{
		if value.AssignedManager==managerId{
          ManagerProjects=append(ManagerProjects, value)
		}
	}

    return ManagerProjects,nil
}


