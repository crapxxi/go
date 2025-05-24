// calculator
package main

import (
	"fmt"
	"time"
)

func plus(num1 int, num2 int) int {
	return num1 + num2
}
func minus(num1 int, num2 int) int {
	return num1 - num2
}
func multiply(num1 int, num2 int) int {
	return num1 * num2
}
func divide(num1 int, num2 int) int {
	if num2 == 0 {
		printerror("division by zero!")
		return 0
	}
	return num1 / num2
}
func printerror(msg string) {
	fmt.Println("Error: ", msg)
}
func calc() {
	for {
		clearterminal()
		fmt.Println("--------------\nCalculator1.0\n--------------")
		var fnum int
		var snum int
		var operation string
		fmt.Println("Write the numbers: ")
		_, err := fmt.Scan(&fnum, &operation, &snum)
		if err != nil {
			printerror("scan error")
			time.Sleep(5 * time.Second)
			return
		}
		switch operation {
		case "+":
			fmt.Println(plus(fnum, snum))
			time.Sleep(5 * time.Second)

		case "-":
			fmt.Println(minus(fnum, snum))
			time.Sleep(5 * time.Second)

		case "*":
			fmt.Println(multiply(fnum, snum))
			time.Sleep(5 * time.Second)

		case "/":
			fmt.Println(divide(fnum, snum))
			time.Sleep(5 * time.Second)

		default:
			printerror("exiting the program")
			time.Sleep(5 * time.Second)
		}
	}
}
