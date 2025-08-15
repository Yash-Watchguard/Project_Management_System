package repository

import (
	"encoding/json"

	"errors"
	"fmt"

	"os"
	"sync"

	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"golang.org/x/crypto/bcrypt"
)

var (
	userFile = "C:/Users/ygoyal/Desktop/PMS_Project/Pms/internal/data/user.json"
	mu       sync.Mutex
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (repo *UserRepo) SaveUser(newUser *user.User) error {
	mu.Lock()
	defer mu.Unlock()

	var users []user.User

	data, err := os.ReadFile(userFile)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, &users)
		if err != nil {
			fmt.Println(" Error parsing JSON:", err)
		}
	}

	users = append(users, *newUser)

	out, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(userFile, out, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (repo *UserRepo) IsUserPresent(name, email, password string) (*user.User, error) {
	data, err := os.ReadFile(userFile)
	if err != nil {
		return nil, err
	}

	var users []user.User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if u.Name == name && u.Email == email {
			// compare stored hashed password with plain password
			err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
			if err == nil {
				return &u, nil
			}
		}
	}

	return nil, errors.New("invalid details")
}
func (repo *UserRepo) ViewProfile(userId string) ([]user.User, error) {
	data, err := os.ReadFile(userFile)
	if err != nil {
		return nil, errors.New("error in readfile")
	}

	var users []user.User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, errors.New("error in unmarshal")
	}

	for _, u := range users {
		if u.Id == userId {
			return []user.User{u}, nil
		}
	}
	return nil, errors.New("user not found")
}

func (ur *UserRepo) GetAllUsers() ([]user.User,error) {
	data, err := os.ReadFile(userFile)
	if err != nil {
		return []user.User{},err
	}
    if len(data)==0{
        return []user.User{},errors.New("no users present")
	}
	var users []user.User
	err= json.Unmarshal(data, &users)
	if err!=nil{
		return users,err
	}

	return users,err
}

func (ur *UserRepo) DeleteUserById(userId string) error {
	data, err := os.ReadFile(userFile)

	if err != nil {
		return err
	}

	var users []user.User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return err
	}

	flag := false
	newUsers := []user.User{}
	for _, user := range users {
		if user.Id != userId {
			newUsers = append(newUsers, user)
		} else {
			flag = true
		}
	}

	if !flag {
		return errors.New("user not found")
	}

	updatedData, err := json.MarshalIndent(newUsers, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(userFile, updatedData, 0644)
	if err != nil {
		return err
	}

	return nil
}
func (ur *UserRepo) GetAllManager() error {
	Data, err := os.ReadFile(userFile)

	if err != nil {
		return err
	}

	var users []user.User
	_ = json.Unmarshal(Data, &users)
	counter := 1
	flag := false
	for _, user := range users {
		if user.Role == roles.Manager { // TODO: Use enum
			flag = true
			fmt.Printf("%d. Name :%s , UserId : %s \n", counter, user.Name, user.Id)
			counter++
		}
	}
	if !flag {
		return errors.New("no manager found")
	}
	return nil
}

func (ur *UserRepo) UpdateProfile(userId, name, email, password, number string) error {
	data, err := os.ReadFile(userFile)
	if err != nil {
		return err
	}

	var users []user.User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return err
	}

	found := false
	for i, user := range users {
		if user.Id == userId {
			users[i].Name = name
			users[i].Email = email
			users[i].Password = password
			users[i].PhoneNumber = number
			found = true
			break
		}
	}

	if !found {
		return errors.New("user not found")
	}

	updatedData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(userFile, updatedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (ad *UserRepo) PromoteEmployee(employeeId string) error {

	data, err := os.ReadFile(userFile)
	if err != nil {
		return errors.New("error in readfile")
	}

	var users []user.User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return errors.New("error in unmarshal")
	}

	found := false
	for i, u := range users {
		if u.Id == employeeId {
			users[i].Role = roles.Manager
			found = true
			break
		}
	}

	if !found {
		return errors.New("employee not found")
	}

	
	updatedData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return errors.New("error in marshal")
	}

	
	err = os.WriteFile(userFile, updatedData, 0644)
	if err != nil {
		return errors.New("error in writing file")
	}

	return nil
}
func (ur *UserRepo) ViewAllEmployee() ([]user.User, error) {
	var users []user.User
	var employees []user.User

	data, err := os.ReadFile(userFile)
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
