package main

import (
	"bytes"
	"sort"
	"strconv"
	"sync"
)

// сюда писать код

func getCrc32(getChan chan<- string, number string) {
	getChan <- DataSignerCrc32(number)
}

var mu = sync.Mutex{}

func getMd5(getChan chan<- string, number string) {
	mu.Lock()
	md5Result := DataSignerMd5(number)
	mu.Unlock()
	getChan <- DataSignerCrc32(md5Result)
}

func getSingleString(number string, out chan interface{}, wg *sync.WaitGroup) {
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
	for input := range in {
		number := input.(int)
		stringNumber := strconv.Itoa(number)

		wg.Add(1)
		go getSingleString(stringNumber, out, wg)
	}
	wg.Wait()
}

func getSignerCrc32(input string, output chan string) {
	result := DataSignerCrc32(input)
	output <- result
}

const threadsNumber = 6

func multiAtomic(input interface{}, out chan<- interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	inputLine := input.(string)
	result := ""

	threads := make([]string, threadsNumber)
	resultChannels := make([]chan string, threadsNumber)

	for i := 0; i < threadsNumber; i++{
		threads[i] = strconv.Itoa(i)
		resultChannels[i] = make(chan string)
	}
	for i, th := range threads {
		buf := bytes.Buffer{}
		buf.WriteString(th)
		buf.WriteString(inputLine)
		go getSignerCrc32(buf.String(), resultChannels[i])
	}
	for i := 0; i < 6; i++ {
		result += <-resultChannels[i]
	}
	out <- result
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for input := range in {
		wg.Add(1)
		go multiAtomic(input, out, wg)
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	results := make([]string, 0, 8)
	for input := range in {
		line := input.(string)
		results = append(results, line)
	}
	sort.Strings(results)
	result := ""
	for i, line := range results {
		if i == 0 {
			result = line
		} else {
			result += "_" + line
		}
	}
	out <- result
}

func ExecutePipeline(inputJobs ...job) {
	channels := []chan interface{}{}
	for i := 0; i < len(inputJobs)+1; i++ {
		channels = append(channels, make(chan interface{}))
	}
	wg := &sync.WaitGroup{}
	for i, job := range inputJobs {
		wg.Add(1)
		go worker(wg, channels[i], channels[i+1], job)
	}
	wg.Wait()
}

func worker(wg *sync.WaitGroup, in, out chan interface{}, jobFunc job) {
	defer close(out)
	defer wg.Done()
	jobFunc(in, out)
}
