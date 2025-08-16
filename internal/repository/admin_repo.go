package repository
// import(
// 	"encoding/json"
// 	"os"
// 	"errors"
// 	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
// 	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
// )
// type AdminRepo struct{
// 	filepath string
// }
// func NewAdminRepo()*AdminRepo{
// 	return &AdminRepo{filepath:"C:/Users/ygoyal/Desktop/PMS_Project/Pms/internal/data/user.json"}
// }
// func (ad *AdminRepo) PromoteEmployee(employeeId string) error {

// 	data, err := os.ReadFile(ad.filepath)
// 	if err != nil {
// 		return errors.New("error in readfile")
// 	}

// 	var users []user.User
// 	err = json.Unmarshal(data, &users)
// 	if err != nil {
// 		return errors.New("error in unmarshal")
// 	}

// 	found := false
// 	for i, u := range users {
// 		if u.Id == employeeId {
// 			users[i].Role = roles.Manager
// 			found = true
// 			break
// 		}
// 	}

// 	if !found {
// 		return errors.New("employee not found")
// 	}

	
// 	updatedData, err := json.MarshalIndent(users, "", "  ")
// 	if err != nil {
// 		return errors.New("error in marshal")
// 	}

	
// 	err = os.WriteFile(ad.filepath, updatedData, 0644)
// 	if err != nil {
// 		return errors.New("error in writing file")
// 	}

// 	return nil
// }
