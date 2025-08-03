package main

import (
	"fmt"
	"log"

	
	"github.com/fatih/color"
)

func main(){
    err:= RunApp()
	if err!=nil{
		log.Fatal(err)
	}
}
func RunApp()error{
for{
	color.Red("----------------------------------Welcome in Tasknest🫶----------------------")
	color.Blue("For Signup press 1\n")
	
    color.Blue("For Login  press 2\n")
    color.Blue("For Exit Press 3\n")

    var Choice int
    fmt.Scanln(&Choice)
	
	switch Choice{
	case 1:
		_=SignupC()

	case 2:
		_=LoginCli()
		
	case 3:
		fmt.Println("Bye Bye 👋")
		return nil
	default:
		color.Red(" Invalid choice. Try again.")
	}
	}

}