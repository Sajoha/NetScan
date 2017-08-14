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
		message := ui.NewLabel(fmt.Sprintf("Current IP: %s", getIP()))
		start := ui.NewEntry()
		end := ui.NewEntry()
		button := ui.NewButton("Start")
		box := ui.NewVerticalBox()
		box.Append(message, false)
		box.Append(start, false)
		box.Append(end, false)
		box.Append(button, false)
		box.SetPadded(true)
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

/*******************************************************************************
* Retrieve the current IP of the host device on the network. (Not yet tested
* where multiple networks might be in use).
*******************************************************************************/
func getIP() string {
	addrs, err := net.InterfaceAddrs()
	checkErr(err)

	for _, ip := range addrs {
		if ipnet, ok := ip.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return "N/A"
}

/*******************************************************************************
* Takes the two IP addresses, and then scans all IPs between the two to see if
* the host responds. This dials port 1 of the host to look for a response.
*
* Future improvements: • Switch from dialing port 1 to using ICMP
*											 • Multithread the scan for a faster result
*******************************************************************************/
func scan(startIP, endIP string) {
	fmt.Println("Starting scan")
	timeout := 50 * time.Millisecond

	startSplit := strings.SplitN(startIP, ".", 4)
	endSplit := strings.SplitN(endIP, ".", 4)

	startRej := fmt.Sprintf("%s.%s.%s", startSplit[0], startSplit[1], startSplit[2])

	startVal, err := strconv.Atoi(startSplit[3])
	checkErr(err)

	endVal, err := strconv.Atoi(endSplit[3])
	checkErr(err)

	for i := startVal; i <= endVal; i++ {
		ip := fmt.Sprintf("%s.%v", startRej, i)
		ipPort := fmt.Sprintf("%s:1", ip)
		_, err := net.DialTimeout("tcp", ipPort, timeout)
		checkNetRes(err, ip)
	}
}

/*******************************************************************************
* Check the response from the dial request for specific types of response.
*******************************************************************************/
func checkNetRes(err error, ip string) {
	if err != nil {
		if strings.HasSuffix(err.Error(), "connection refused") {
			fmt.Printf("Response from %s\n", ip)
		} else if strings.HasSuffix(err.Error(), "permission denied") {
			fmt.Printf("Denied by %s\n", ip)
		} else if strings.HasSuffix(err.Error(), "i/o timeout") {
			// Do Nothing
		} else {
			fmt.Printf("A different error: %s\n", err)
		}
	}
}

/*******************************************************************************
* Takes the two IP addresses, and validate whether they're in a valid format,
* in a similar subnet and that the IP is legitimate. Runs before the scan to
* warn the user of issues before the scan starts.
*
* Future Improvements: • Allow a more flexible IP input, currently only the last
*												 chunk can differ
*******************************************************************************/
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

/*******************************************************************************
* Lazy error checker, just to save lines elsewhere.
*******************************************************************************/
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
