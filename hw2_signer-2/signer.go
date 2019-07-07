package main

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"sync"
)
var mu = sync.Mutex{}
// сюда писать код
func getCrc32(getChan chan<- string, number string) {
	getChan <- DataSignerCrc32(number)
}

func getMd5(getChan chan<- string, number string) {
	mu.Lock()
	defer mu.Unlock()
	getChan <- DataSignerCrc32(DataSignerMd5(number))
}

func getSingleString(number string, out chan interface{}, wg *sync.WaitGroup){
	defer wg.Done()
	getChanCrc32 := make(chan string)
	getChanMd5 := make(chan string)
	go getCrc32(getChanCrc32, number)
	go getMd5(getChanMd5, number)

	crc32Result := <-getChanCrc32
	md5Result := <-getChanMd5

	out <- crc32Result + "~" + md5Result
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	fmt.Printf("start single\n")
	for input := range in {
		fmt.Printf("inSingle\n")
		number := input.(int)
		stringNumber := strconv.Itoa(number)

		wg.Add(1)
		go getSingleString(stringNumber, out, wg)
		fmt.Printf("outSingle\n")
	}
	wg.Wait()
	close(out)
	fmt.Printf("end single\n")
}

func getSignerCrc32(input string, output chan string) {
	result := DataSignerCrc32(input)
	output <- result
}

func multiAtomic(input interface{}, out chan<- interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("inMulti\n")
	inputLine := input.(string)
	result := ""
	threads := []string{"0", "1", "2", "3", "4", "5"}
	var resultChannels = []chan string {
		make(chan string),
		make(chan string),
		make(chan string),
		make(chan string),
		make(chan string),
		make(chan string),
	}
	for i, th := range threads{
		buf := bytes.Buffer{}
		buf.WriteString(th)
		buf.WriteString(inputLine)
		go getSignerCrc32(buf.String(), resultChannels[i])
	}
	for i := 0; i < 6; i++{
		result += <-resultChannels[i]
	}
	out <- result
	fmt.Printf("outMulti\n")
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for input := range in {
		wg.Add(1)
		multiAtomic(input, out, wg)
	}
	wg.Wait()
	close(out)
}

func CombineResults(in, out chan interface{}) {
	results := make([]string, 0, 8)
	for input := range in{
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
	close(out)
}


func ExecutePipeline(inputJobs ...job) {
	channels := []chan interface{}{}
	for i := 0; i < len(inputJobs) + 1; i++{
		channels = append(channels, make(chan interface{}, 1000))
	}
	for i, job := range inputJobs {
		if i == len(inputJobs) - 1{
			job(channels[i], channels[i + 1])
		}else{
			go job(channels[i], channels[i + 1])
		}
	}
}
