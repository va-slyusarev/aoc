package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"time"
)

var input = flag.String("input", "../input.txt", "Input data file path.")

type log struct {
	time time.Time
	do   string
}

func main() {
	flag.Parse()

	f, err := ioutil.ReadFile(*input)
	if err != nil {
		fmt.Printf("error read input data file: %v\n", err)
		return
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(f))

	var re = regexp.MustCompile("\\[(\\d{4}-\\d{2}-\\d{2}\\s\\d{2}:\\d{2})\\]\\s(.+)") // [1518-11-01 00:55] wakes up
	var logs []log

	for scanner.Scan() {
		value := scanner.Text()
		values := re.FindStringSubmatch(value)
		if len(values) != 3 {
			fmt.Printf("broken data %q from input data file: %v\n", value, err)
			return
		}

		t, err := time.Parse("2006-01-02 15:04", values[1])
		if err != nil {
			fmt.Printf("broken data time %q from input data file: %v\n", values[1], err)
			return
		}

		logs = append(logs, log{t, values[2]})
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error scan input data file: %v\n", err)
		return
	}

	sort.SliceStable(logs, func(i, j int) bool { return logs[i].time.Before(logs[j].time) })

	// key: guard id, value: sleep count between 00:00 - 00:59
	guards := make(map[int][60]int)

	preGuard := -1
	preMinute := 0
	for _, log := range logs {

		hour := log.time.Hour()
		minute := log.time.Minute()

		// guards can leave until 00:00
		if hour == 23 {
			hour = 0
			minute = 0
		}

		var guardId int
		if _, err := fmt.Sscanf(log.do, "Guard #%d begins shift", &guardId); err == nil {
			if _, ok := guards[guardId]; !ok {
				guards[guardId] = [60]int{}
			}
			preGuard = guardId
			preMinute = 0
			continue
		}

		if log.do == "falls asleep" {
			table := guards[preGuard]
			table[minute] += 1
			guards[preGuard] = table
			preMinute = minute
			continue
		}
		if log.do == "wakes up" {
			table := guards[preGuard]
			for i := preMinute + 1; i < minute; i++ {
				table[i] += 1
			}
			guards[preGuard] = table
			continue
		}
	}

	max := 0
	guardId := -1

	for k, v := range guards {

		sum := 0
		for i := range v {
			sum += v[i]
		}

		if sum >= max {
			max = sum
			guardId = k
		}
	}

	maxIndex := 0
	max = 0
	for i := range guards[guardId] {

		if guards[guardId][i] >= max {
			max = guards[guardId][i]
			maxIndex = i
		}
	}

	fmt.Printf("Target value is: %d\n", maxIndex*guardId)
}
