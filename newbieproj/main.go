package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func clearterminal() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	for {
		clearterminal()
		fmt.Println("Choose 1 of the projects:\n[1]Calculator\n[2]ToDoList\n[0]Exit")
		var operation int
		_, err := fmt.Scan(&operation)
		if err != nil {
			printerror(err.Error())
			time.Sleep(5 * time.Second)
			return
		}
		switch operation {
		case 0:
			fmt.Println("BYE BYE!!!")
			return
		case 1:
			calc()
		case 2:
			todolist()
		default:
			printerror("Invalid operation")
			time.Sleep(5 * time.Second)
		}
	}
}
