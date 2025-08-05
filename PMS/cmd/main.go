package main

import (
	"fmt"
	
	"log"

	"github.com/Yash-Watchguard/Tasknest/handler"
	"github.com/Yash-Watchguard/Tasknest/internal/constants"
	"github.com/fatih/color"
)

func main() {
	err := RunApp()
	if err != nil {
		log.Fatal(err)
	}
}
var(
	
)

func RunApp() error {
	for {
		color.Red(constants.WelcomeMSG)
		color.Blue(constants.SignupChoice)

		color.Blue(constants.LoginChoice)
		color.Blue(constants.ExitChoice)

		var Choice int
		fmt.Scanln(&Choice)

		switch Choice {
		case 1:
			err := handler.Signup()
			if err!=nil{
				fmt.Println(err)
			}

		case 2:
			 err:= handler.Login()
			if err!=nil{
				fmt.Println(err)
			}

		case 3:
			fmt.Println(constants.GoodByeMsg)
			return nil
		default:
			color.Red(constants.InvalidChoice)
		}
	}

}
