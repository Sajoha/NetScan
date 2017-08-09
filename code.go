package main

import (
	"errors"
	"fmt"
	"github.com/andlabs/ui"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	err := ui.Main(func() {
		start := ui.NewEntry()
		end := ui.NewEntry()
		button := ui.NewButton("Start")
		box := ui.NewVerticalBox()
		box.Append(start, false)
		box.Append(end, false)
		box.Append(button, false)
		window := ui.NewWindow("Port Scan", 200, 100, false)
		window.SetChild(box)
		button.OnClicked(func(*ui.Button) {
			ipValid := validate(start.Text(), end.Text())
			if ipValid == nil {
				scan(start.Text(), end.Text())
			}
		})
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}

func checkNetRes(err error, ip string) {
	if err != nil {
		if strings.HasSuffix(err.Error(), "connection refused") {
			fmt.Printf("Response from %s\n", ip)
		} else if strings.HasSuffix(err.Error(), "permission denied") {
			fmt.Printf("Denied by %s\n", ip)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func scan(startIP, endIP string) {
	fmt.Println("Starting scan")
	timeout := 100 * time.Millisecond
	startSplit := strings.SplitN(startIP, ".", 4)
	startRej := fmt.Sprintf("%s.%s.%s", startSplit[0], startSplit[1], startSplit[2])
	for i := 1; i <= 255; i++ {
		ip := fmt.Sprintf("%s.%v", startRej, i)
		ipPort := fmt.Sprintf("%s:1", ip)
		_, err := net.DialTimeout("tcp", ipPort, timeout)
		checkNetRes(err, ip)
	}
}

func validate(startIP, endIP string) error {
	// Validate whether or not the two IPs are in the correct IPv4 format
	ipExpression := "^([0-9]{1,3})[.]([0-9]{1,3})[.]([0-9]{1,3})[.]([0-9]{1,3})$"

	startMatch, err := regexp.MatchString(ipExpression, startIP)
	checkErr(err)

	endMatch, err := regexp.MatchString(ipExpression, endIP)
	checkErr(err)

	if startMatch == false {
		return errors.New("Starting IP is not a valid IP")
	} else if endMatch == false {
		return errors.New("End IP is not a valid IP")
	}

	// Check to see whether the two IPs are within a similar subnet
	startSplit := strings.SplitN(startIP, ".", 4)
	endSplit := strings.SplitN(endIP, ".", 4)

	startRej := startSplit[0] + startSplit[1] + startSplit[2]
	endRej := endSplit[0] + endSplit[1] + startSplit[2]

	compareRes := strings.Compare(startRej, endRej)

	if compareRes != 0 {
		return errors.New("IP subnets do not match!")
	}

	// Check that the IP chunks are in the valid range (Assuming we have an IP format from the regular expression)
	for i := 0; i <= 3; i++ {
		startVal, err := strconv.Atoi(startSplit[i])
		checkErr(err)

		endVal, err := strconv.Atoi(endSplit[i])
		checkErr(err)

		if startVal < 0 || startVal > 255 {
			return errors.New("Starting IP is not a valid IP")
		}
		if endVal < 0 || endVal > 255 {
			return errors.New("End IP is not a valid IP")
		}

		if startVal >= endVal && i == 3 {
			return errors.New("Starting IP is lower than end IP")
		}
	}

	// Everything went well
	return nil
}
