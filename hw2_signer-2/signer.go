package main

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
)

// сюда писать код
func SingleHash(in, out chan interface{}) {

}

func MultiHash(in, out chan interface{}) {

}

func CombineResults(in, out chan interface{}) {

}

func ExecutePipeline(inputJob ...job) {

}

func single(n interface{})string {
	//crc32(data)+"~"+crc32(md5(data))
	number, ok := n.(int)
	if !ok {
		fmt.Printf("can't convert to int\n")
		return ""
	}
	stringNumber := strconv.Itoa(number)
	return DataSignerCrc32(stringNumber) + "~" + DataSignerCrc32(DataSignerMd5(stringNumber))
}

func multi(n interface{}) string{
	inputLine, ok := n.(string)
	if !ok {
		fmt.Printf("Can't convert to string multi input")
	}
	result := ""
	threads := []string{"0", "1", "2", "3", "4", "5",}
	for _, th := range threads{
		buf := bytes.Buffer{}
		buf.WriteString(th)
		buf.WriteString(inputLine)
		result += DataSignerCrc32(buf.String())
	}
	return result
}

func combine(results []string)string {
	sort.Strings(results)
	result := ""
	for i, line := range results{
		if i == 0{
			result = line
		}else{
			result += "_" + line
		}
	}

	return result
}

func MyHashFunction(fibNumbers []int) string {
	results := []string{}
	for _, number := range fibNumbers {
		singleLine := single(number)
		multiLine := multi(singleLine)
		results = append(results, multiLine)
	}
	result := combine(results)

	return result
}

func main() {
	fmt.Printf("hello world!\n")
}