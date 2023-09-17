package test

import (
	"fmt"
	"testing"
)

/*
*
质数
*/
func TestA(t *testing.T) {
	for number := 1; number <= 20; number++ {
		if findprimes(number) {
			fmt.Printf("%v ", number)
		}
	}
}
func findprimes(number int) bool {
	for i := 2; i < number; i++ {
		if number%i == 0 {
			return false
		}
	}
	if number > 1 {
		return true
	} else {
		return false
	}
}

/*
*
3、5整除
*/
func TestB(t *testing.T) {
	FizzBuzz()
}
func FizzBuzz() {
	for i := 1; i <= 100; i++ {
		switch 0 {
		case i % 15:
			fmt.Println("FIzzBuzz")
		case i % 3:
			fmt.Println("Fizz")
		case i % 5:
			fmt.Println("Buzz")
		default:
			fmt.Println(i)
		}
	}
}

func TestC(t *testing.T) {
	getNum()
}

func getNum() {
	val := 0
	for {
		fmt.Print("Enter number: ")
		fmt.Scanf("%d", &val)
		switch {
		case val < 0:
			panic("You entered a negative number!")
		case val == 0:
			fmt.Println("0 is neither negative nor positive")
		default:
			fmt.Println("You entered:", val)
		}
	}
}
