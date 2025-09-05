package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"math"
)

// Takes a slice of ints and returns sum
func SumInts(ints []int) int {
	var result int
	for _, val := range ints {
		result += val
	}
	return result
}

// Compares two ints and returns the lowest
func MinInt(i, j int) int {
	if i < j {
		return i
	} else {
		return j
	}
}

// Compares two ints and returns the highest
func MaxInt(i, j int) int {
	if i > j {
		return i
	} else {
		return j
	}
}

// Takes a slice of ints and returns the lowest value and the corresponding key
func MinInts(ints []int) (int, int) {
	result := ints[0]
	rKey := 0
	for key, val := range ints {
		if val < result {
			result = val
			rKey = key
		}
	}
	return result, rKey
}

// Takes a slice of ints and returns the highest value and the corresponding key
func MaxInts(ints []int) (int, int) {
	result := ints[0]
	rKey := 0
	for key, val := range ints {
		if val > result {
			result = val
			rKey = key
		}
	}
	return result, rKey
}

// Check if error, if error print it
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// Takes the current day "dayX" as a string, and if test it checks the file "test_input.txt", returns a slice of string with the content of the file
func ReadFileSlice(day string, test bool) []string {
	fmt.Println(os.Getwd())
	path := ""
	if test {
		path = fmt.Sprintf("../%v/test_input.txt", day)
	} else {
		path = fmt.Sprintf("../%v/input.txt", day)
	}

	file, err := os.ReadFile(path)
	CheckErr(err)

	s := strings.Split(string(file), "\n")

	return s[:len(s)-1]
}

// Checks slice of string if it contains string i, return True of False
func ContainsString(slice []string, i string) bool {
	res := false
	for _, num := range slice {
		if num == i {
			res = true
			break
		}
	}
	return res

}

// Checks slice of int if it contains int i, return True of False
func ContainsInt(slice []int, i int) bool {
	res := false
	for _, num := range slice {
		if num == i {
			res = true
			break
		}
	}
	return res

}

// Converts string s to int, returns int
// Removes trailing whitespace
func StrToInt(s string) int {
	s = strings.Replace(s, " ", "", -1)
	i, err := strconv.Atoi(s)
	CheckErr(err)
	return i
}

// Converts slice of string to slice of int
// If only white space, ignore
func StringsToInt(s []string) []int {
	res := []int{}
	for _, v := range s {
		if v == " " || v == "" {
			continue
		}
		res = append(res, StrToInt(v))
	}
	return res
}

// Makes absolute value of int
func IntAbs(i int) int {
	if i < 0 {
		i = -i
	}
	return i
}

// Convert string to int64
func StrToInt64(s string) int64 {
	i, err := strconv.Atoi(s)
	CheckErr(err)
	return int64(i)
}

// Converts a slice of strings to a slice of ints
func ConvertSliceStringToInt(slice []string) []int {
	si := []int{}
	fmt.Printf("slice = %+v\n", slice)
	for _, num := range slice {
		nn, err := strconv.Atoi(num)
		CheckErr(err)
		si = append(si, nn)
	}
	return si
}

// Converts a slice of strings to a slive of int64
func ConvertSliceStringToInt64(slice []string) []int64 {
	si := []int64{}
	fmt.Printf("slice = %+v\n", slice)
	for _, num := range slice {
		nn, err := strconv.Atoi(num)
		CheckErr(err)
		si = append(si, int64(nn))
	}
	return si
}

// A map with numbers matched to english words
// Returns map with key as the string and value int
func mapNumtextAndNums() map[string]int {
	numbers := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
		"six":   6,
		"seven": 7,
		"eight": 8,
		"nine":  9,
		"zero":  0,
	}
	return numbers
}

// A map where the alphabet is mapped to incremental numbers
// Returns map with key as the string and value int
func mapAlphaWithNum() map[string]int {
	var A = map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7, "h": 8, "i": 9, "j": 10, "k": 11, "l": 12, "m": 13, "n": 14, "o": 15, "p": 16, "q": 17, "r": 18, "s": 19, "t": 20, "u": 21, "v": 22, "w": 23, "x": 24, "y": 25, "z": 26}
	return A
}

// Check if string can be converted to int
func CheckIfInt(s string) bool {
	_, err := strconv.Atoi(s)
	if err != nil {
		return false
	} else {
		return true
	}
}

// Check all the keys of the map if value exists
func MapContainsKey(indata map[int]int, value int) bool {
	if indata[value] == 0 {
		return true
	} else {
		return false
	}
}

// Remove key i from slice
func RemoveFromSliceInt(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// Remove key i from slice but order is importnant
func RemoveFromSliceOrder(s []int, index int) []int {
	ret := make([]int, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

// Take a byte of length 8 and extract [8]bool
func GetBoolsFromByte(b byte) [8]bool {
	var res [8]bool
	for i := range 8{
		bf := math.Pow(2.0, float64(i))
		bb := byte(bf)
		if bb & b >> i == 1 {
			res[len(res)-i-1] =  true
		} else {
			res[len(res)-i-1] = false
		}
	}
	return res
}
// Takes a slice of bools, returns the byte of the 8 first bits
func GetByteFromBools(b []bool) byte {
	res := 0
	for k := range 8 {
	//for k, v := range b {
		if b[k] {
			res += int(math.Pow(2.0, float64(7-k)))
		}
	}
	return byte(res)
}
