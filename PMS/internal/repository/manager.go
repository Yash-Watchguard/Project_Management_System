package repository

import (
	"encoding/json"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"os"
)

type ManagerRepo struct {
	filepath string
}

// PromoteEmployee implements interfaces.ManagerRepository.
func (managerRepo *ManagerRepo) PromoteEmployee(employeeId string) error {
	panic("unimplemented")
}

func NewManagerRepo() *ManagerRepo {
	return &ManagerRepo{filepath: "C:/Users/ygoyal/Desktop/PMS_Project/Pms/internal/data/user.json"}
}

func (managerRepo *ManagerRepo) ViewAllEmployee() ([]user.User, error) {
	var users []user.User
	var employees []user.User

	data, err := os.ReadFile(managerRepo.filepath)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []user.User{}, nil
	}

	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Role == 2 {
			employees = append(employees, user)
		}
	}

	return employees, nil
}
