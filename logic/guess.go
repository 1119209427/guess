package logic

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

//实现猜数游戏的逻辑

// CheckGuess 游戏规则:电脑随机生成四个数字，用户输入四个数字。电脑反馈XAYB，x代表数字完全猜对的数量，y代表数字猜对但位置不对的数量。
//根据用户输入的难度:简单生成两位数 中等生成三位数
//如电脑为1234，用户输入1243，返回2A2B.
func CheckGuess(guess string, secretNumber int) string {
	return reply(guess, secretNumber)
}

func strToByte(str string) []string {
	var result []string
	bytes := []byte(str)
	for i := 0; i < len(bytes); i++ {
		result = append(result, string(bytes[i]))
	}
	return result
}

func check(strGuess []string, stcNumber string) string {
	a, b := 0, 0
	for i := 0; i < len(strGuess); i++ {
		if strings.Count(stcNumber, strGuess[i]) == 1 {
			if strings.Contains(stcNumber, strGuess[i]) && strings.Index(stcNumber, strGuess[i]) == i {
				a++
			}
			if strings.Contains(stcNumber, strGuess[i]) && strings.Index(stcNumber, strGuess[i]) != i {
				b++
			}
		}
		if strings.Count(stcNumber, strGuess[i]) > 1 {
			if strings.Contains(stcNumber, strGuess[i]) && strings.Index(stcNumber, strGuess[i]) == i {
				a++
			}
			if strings.Contains(stcNumber, strGuess[i]) && strings.Index(stcNumber, strGuess[i]) != i {
				b++
			}
			stcNumber = strings.TrimPrefix(stcNumber, strGuess[i])
			strGuess = strGuess[i+1:]

			i = i - 1
		}

	}
	return fmt.Sprintf("A%dB%d", a, b)
}

func reply(guess string, secretNumber int) string {

	stcNumber := strconv.Itoa(secretNumber)
	iGuess, err := strconv.Atoi(guess)
	if err != nil {
		fmt.Println(err.Error())
	}

	strGuess := strToByte(guess)
	result := check(strGuess, stcNumber)
	var str string

	fmt.Println(secretNumber)
	if iGuess > secretNumber {
		str = fmt.Sprintf("Your guess is bigger than the secret number. Please try again %s", result)
	} else if iGuess < secretNumber {
		str = fmt.Sprintf("Your guess is smaller than the secret number. Please try again %s", result)
	} else {
		str = fmt.Sprintf("Correct, you Legend! %s", result)
	}
	return str
}

func Easy() int {
	maxNum := 99
	rand.Seed(time.Now().UnixNano())
	secretNumber := rand.Intn(maxNum)
	return secretNumber
}

func Medium() int {
	maxNum := 999
	rand.Seed(time.Now().UnixNano())
	secretNumber := rand.Intn(maxNum)
	return secretNumber
}

func Hard() int {
	maxNum := 9999
	rand.Seed(time.Now().UnixNano())
	secretNumber := rand.Intn(maxNum)
	return secretNumber
}
