package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	utilErrorHeader string = "UtilError : "
)

func today() string {
	return time.Now().Format("2006-01-02")
}

func aday(num int) string {
	return time.Now().Add(time.Duration(num) * time.Hour * 24).Format("2006-01-02")
}

func checkDayFormat(str string) bool {
	var layout string = "2006-01-02"
	_, err := time.Parse(layout, str)
	if err != nil {
		return false
	}
	return true
}

func debugPrint(str ...interface{}) {
	fmt.Println(str)
	return
}

func getNumSlice(str string, sep string) ([]int, error) {
	strNums := strings.Split(str, sep)
	var nums []int
	for _, strNum := range strNums {
		if num, err := strconv.Atoi(strNum); err == nil {
			nums = append(nums, num)
		} else {
			return nil, getUtilError("Any strings that cant be translated to number is included.")
		}
	}
	return nums, nil
}

func getUtilError(eBody string) error {
	return errors.New(utilErrorHeader + eBody)
}

func exists(path string) bool {
	_, err := os.Stat(transEnvPath(path))
	return err == nil
}

func transEnvPath(path string) string {
	ss := strings.Split(path, "/")
	if len(ss[0]) != 0 && ss[0][:1] == "$" {
		env := os.Getenv(ss[0][1:])
		ss[0] = env
	}
	return strings.Join(ss, "/")
}
