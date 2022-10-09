package logic

import (
	"fmt"
	"testing"
)

func TestCheckGuess(t *testing.T) {
	str := CheckGuess("50", 50)
	fmt.Println(str)
	if str != "Correct, you Legend! A2B0" {
		t.Error("猜数逻辑错误")
	}
	str1 := CheckGuess("91", 19)
	fmt.Println(str1)
	if str1 != "Your guess is bigger than the secret number. Please try again A0B2" {
		t.Error("猜数逻辑错误")
	}
	str2 := CheckGuess("99", 99)
	fmt.Println(str2)
	if str2 != "Correct, you Legend! A2B0" {
		t.Error("猜数逻辑错误")
	}
	str3 := CheckGuess("505", 505)
	fmt.Println(str3)
	if str3 != "Correct, you Legend! A3B0" {
		t.Error("猜数逻辑错误")
	}
	str4 := CheckGuess("919", 991)
	fmt.Println(str4)
	if str4 != "Your guess is smaller than the secret number. Please try again A1B2" {
		t.Error("猜数逻辑错误")
	}
	str5 := CheckGuess("999", 999)
	fmt.Println(str5)
	if str5 != "Correct, you Legend! A3B0" {
		t.Error("猜数逻辑错误")
	}
}
func TestStr(t *testing.T) {
	s1 := "hello"
	b := []byte(s1)

	// []byte to string
	s2 := string(b[:])
	fmt.Println(s2)

}
