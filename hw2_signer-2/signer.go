package main

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"time"
)

// сюда писать код
func SingleHash(in, out chan interface{}) {
	fmt.Printf("start single\n")
	for input := range in {
		fmt.Printf("inSingle\n")
		number := input.(int)
		stringNumber := strconv.Itoa(number)
		 out <- DataSignerCrc32(stringNumber) + "~" + DataSignerCrc32(DataSignerMd5(stringNumber))
		fmt.Printf("outSingle\n")
	}
	close(out)
	fmt.Printf("end single\n")
}

func MultiHash(in, out chan interface{}) {
	fmt.Printf("start Multi\n")
	for input := range in {
		fmt.Printf("inMulti\n")
		inputLine := input.(string)
		result := ""
		threads := []string{"0", "1", "2", "3", "4", "5",}
		for _, th := range threads{
			buf := bytes.Buffer{}
			buf.WriteString(th)
			buf.WriteString(inputLine)
			result += DataSignerCrc32(buf.String())
		}
		out <- result
		fmt.Printf("outMulti\n")
	}
	close(out)
	fmt.Printf("end Multi\n")
}

func CombineResults(in, out chan interface{}) {
	fmt.Printf("start Combine\n")
	results := make([]string, 0, 8)
	for input := range in{
		fmt.Printf("inCombine\n")
		line := input.(string)
		results = append(results, line)
	}
	sort.Strings(results)
	result := ""
	for i, line := range results{
		if i == 0{
			result = line
		}else{
			result += "_" + line
		}
	}
	out <- result
	fmt.Printf("outCombine\n")
	close(out)
	fmt.Printf("end Combine\n")
}


func ExecutePipeline(inputJobs ...job) {
	channels := []chan interface{}{}
	for i := 0; i < len(inputJobs) + 1; i++{
		channels = append(channels, make(chan interface{}, 1000))
	}
	fmt.Printf("channels = %v\n", channels)
	for i, job := range inputJobs {
		go job(channels[i], channels[i + 1])
	}
	time.Sleep(30 * time.Second)
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
