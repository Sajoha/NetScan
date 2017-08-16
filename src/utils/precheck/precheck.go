package precheck

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"utils/errorcheck"
)

/*******************************************************************************
* Takes the two IP addresses, and validate whether they're in a valid format,
* in a similar subnet and that the IP is legitimate. Runs before the scan to
* warn the user of issues before the scan starts.
*
* Future Improvements: • Allow a more flexible IP input, currently only the last
*												 chunk can differ
*******************************************************************************/
func Validate(startIP, endIP string) error {
	// Validate whether or not the two IPs are in the correct IPv4 format
	ipExpression := "^([0-9]{1,3})[.]([0-9]{1,3})[.]([0-9]{1,3})[.]([0-9]{1,3})$"

	startMatch, err := regexp.MatchString(ipExpression, startIP)
	errorcheck.CheckErr(err)

	endMatch, err := regexp.MatchString(ipExpression, endIP)
	errorcheck.CheckErr(err)

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
		errorcheck.CheckErr(err)

		endVal, err := strconv.Atoi(endSplit[i])
		errorcheck.CheckErr(err)

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
