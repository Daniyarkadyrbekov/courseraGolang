package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"./easyJson"
)

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	r := regexp.MustCompile("@")
	seenBrowsers := []string{}
	uniqueBrowsers := 0
	//foundUsers := ""

	scanner := bufio.NewScanner(file)
	i := 0
	firstLine := true
	for scanner.Scan() {
		//user := make(map[string]interface{})
		jsonStruct := easyJson.EasyJsonStruct{}
		err := jsonStruct.UnmarshalJSON([]byte(scanner.Text()))
		if err != nil {
			panic(err)
		}
		//fmt.Printf("jsonStruct = %v\n", jsonStruct)

		isAndroid := false
		isMSIE := false

		//browsers, ok := user["browsers"].([]interface{})
		//if !ok {
		//	// log.Println("cant cast browsers")
		//	continue
		//}

		for _, browser := range jsonStruct.Browsers {
			//browser, ok := browserRaw.(string)
			//if !ok {
			//	// log.Println("cant cast browser to string")
			//	continue
			//}
			if strings.Contains(browser, "Android") {
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		for _, browser := range jsonStruct.Browsers {
			//browser, ok := browserRaw.(string)
			//if !ok {
			//	// log.Println("cant cast browser to string")
			//	continue
			//}
			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			i++
			continue
		}

		email := r.ReplaceAllString(user["email"].(string), " [at] ")
		if firstLine{
			fmt.Fprintln(out, "found users:")
			firstLine = false
		}
		//foundUsers := fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email)
		fmt.Fprintln(out,fmt.Sprintf("[%d] %s <%s>", i, user["name"], email))
		i++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	//fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}
