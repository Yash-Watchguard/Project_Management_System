package repository

import (
	"encoding/json"
	
	"fmt"
	"io/ioutil"
	"os"
	"sync"
    "github.com/fatih/color"
	"github.com/Yash-Watchguard/Tasknest/model"
	"golang.org/x/crypto/bcrypt"
	"errors"
)

var (
	userFile = "C:/Users/ygoyal/Desktop/PMS_Project/Pms/data/user.json"
	mu       sync.Mutex
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (repo *UserRepo) SaveUser(user *model.User) error {
	mu.Lock()
	defer mu.Unlock()

	var users []model.User

	// Read existing users
	data, err := os.ReadFile(userFile)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, &users)
		if err != nil {
			fmt.Println("⚠️ Error parsing JSON:", err)
		}
	}

	// Append new user
	users = append(users, *user)

	// Marshal and write
	out, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(userFile, out, 0644)
	if err != nil {
		return err
	}

	fmt.Println("✅ User written to user.json")
	return nil
}

func (repo *UserRepo) IsUserPresent(name, email, password string) (*model.User, error) {
    data, err := ioutil.ReadFile(userFile)
    if err != nil {
        return nil, err
    }

    var users []model.User
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
func(repo *UserRepo)ViewProfile(user *model.User){
color.Cyan("----------- %s Profile -----------",user.Role)
color.Yellow("ID     : %d", user.Id)
color.Yellow("Name   : %s", user.Name)
color.Yellow("Email  : %s", user.Email)
color.Yellow("Role   : %s", user.Role)
color.Cyan("-------------------------------------")
fmt.Println("Press ENTER to return to dashboard...")
fmt.Scanln()
}

func (ur *UserRepo) GetAllUsers() []model.User {
	data, err := ioutil.ReadFile(userFile)
	if err != nil {
		return []model.User{}
	}

	var users []model.User
	_ = json.Unmarshal(data, &users)
	return users
}

func(ur *UserRepo)DeleteUserById(userId string)error{
	data,err := ioutil.ReadFile(userFile)

	if err!=nil{
		return err
	}

	var users []model.User
	err=json.Unmarshal(data,&users)
	if err != nil{
		return err
	}

	flag:=false
    newUsers := []model.User{}
	for _,user := range users {
		if user.Id!=userId{
          newUsers=append(newUsers, user)
		}else{
			flag=true
		}
	}

	if !flag {
		return errors.New("user not found")
	}

	updatedData, err := json.MarshalIndent(newUsers, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(userFile, updatedData, 0644)
	if err != nil {
		return err
	}

	return nil
}
func (ur *UserRepo) GetAllManager()error {
	Data,err:=ioutil.ReadFile(userFile)

	if err!=nil{
		return err
	}

	var users []model.User
	_=json.Unmarshal(Data,&users)
    counter:=1
	 flag:=false
	for _,user:=range users{
		if user.Role=="Manager"{
			flag=true
			fmt.Printf("%d. Name :%s , UserId : %s \n",counter,user.Name,user.Id)
			counter++
		}
	}
	if !flag{
		return errors.New("There is No Manager For assign the Task")
	}
	return nil
}




