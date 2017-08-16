package main

import (
	"fmt"
	"github.com/andlabs/ui"
	"net"
	"strconv"
	"strings"
	"time"
	"utils/errorcheck"
	"utils/output"
	"utils/precheck"
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
			ipValid := precheck.Validate(start.Text(), end.Text())
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
	errorcheck.CheckErr(err)

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
	output.Banner(fmt.Sprintf("Starting scan between %s and %s", startIP, endIP))
	timeout := 50 * time.Millisecond

	startSplit := strings.SplitN(startIP, ".", 4)
	endSplit := strings.SplitN(endIP, ".", 4)

	startRej := fmt.Sprintf("%s.%s.%s", startSplit[0], startSplit[1], startSplit[2])

	startVal, err := strconv.Atoi(startSplit[3])
	errorcheck.CheckErr(err)

	endVal, err := strconv.Atoi(endSplit[3])
	errorcheck.CheckErr(err)

	for i := startVal; i <= endVal; i++ {
		ip := fmt.Sprintf("%s.%v", startRej, i)
		ipPort := fmt.Sprintf("%s:1", ip)
		_, err := net.DialTimeout("tcp", ipPort, timeout)
		errorcheck.CheckNetRes(err, ip)
	}

	fmt.Println("")
	output.Info("Done\n")
}
