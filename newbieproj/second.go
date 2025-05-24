package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

func todolist() {
	tasks := make([]string, 0, 5)
	reader := bufio.NewReader(os.Stdin)
	for {
		clearterminal()
		fmt.Println("---TODOLIST---")
		var operation int
		if len(tasks) != 0 {
			for i := 0; i < len(tasks); i++ {
				fmt.Printf("[%d]%s\n", i, tasks[i])
			}
		} else {
			fmt.Println("List is empty!")
		}
		fmt.Println("--------------")
		fmt.Println("[1]New task\n[2]Delete task\n[0]Exit")
		_, err := fmt.Scan(&operation)
		if err != nil {
			printerror("invalid operation")
			return
		}
		switch operation {
		case 1:
			clearterminal()
			fmt.Println("Please, enter the task: ")
			text, err := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			if err != nil {
				printerror(err.Error())
				time.Sleep(time.Second * 2)
			} else {
				tasks = append(tasks, text)
				fmt.Println("Task added succesfully!")
				time.Sleep(time.Second * 3)
			}
		case 2:
			if len(tasks) == 0 {
				printerror("You should have at least 1 task!")
				time.Sleep(time.Second * 3)
				continue
			}
			fmt.Println("Please enter the index of the task")
			var dindex int
			_, err := fmt.Scan(&dindex)
			if err != nil {
				printerror(err.Error())
				time.Sleep(time.Second * 2)
			} else {
				if dindex >= len(tasks) || dindex < 0 {
					printerror("Number is out of range")
					time.Sleep(time.Second * 3)
					continue
				}
				tasks = slices.Delete(tasks, dindex, dindex+1)
			}
		case 0:
			fmt.Println("BYEBYE!")
			time.Sleep(time.Second * 3)
			return
		}
	}
}
