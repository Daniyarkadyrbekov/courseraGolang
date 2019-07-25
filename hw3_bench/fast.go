package main

import (
	"./easyJson"
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// вам надо написать более быструю оптимальную этой функции

var dataPool = sync.Pool{
	New: func() interface{} {
		return new(easyJson.EasyJsonStruct)
	},
}

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	seenBrowsers := []string{}
	uniqueBrowsers := 0

	scanner := bufio.NewScanner(file)
	i := 0
	firstLine := true
	for scanner.Scan() {
		var jsonStruct  = easyJson.EasyJsonStruct{}
		err := jsonStruct.UnmarshalJSON([]byte(scanner.Text()))
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		for _, browser := range jsonStruct.Browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		for _, browser := range jsonStruct.Browsers {
			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			i++
			continue
		}
		
		email := strings.Replace(jsonStruct.Email, "@", " [at] ", -1)
		if firstLine{
			fmt.Fprintln(out, "found users:")
			firstLine = false
		}
		fmt.Fprintln(out,fmt.Sprintf("[%d] %s <%s>", i, jsonStruct.Name, email))
		i++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}
