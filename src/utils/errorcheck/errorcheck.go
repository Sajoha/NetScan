package errorcheck

import (
	"fmt"
	"strings"
	"utils/output"
)

/*******************************************************************************
* Check the return from the dial request for specific types of response.
*******************************************************************************/
func CheckNetRes(err error, ip string) {
	if err != nil {
		if strings.HasSuffix(err.Error(), "connection refused") {
			output.Info(fmt.Sprintf("Response from %s", ip))
		} else if strings.HasSuffix(err.Error(), "permission denied") {
			output.Warn(fmt.Sprintf("Denied by %s", ip))
		} else if strings.HasSuffix(err.Error(), "i/o timeout") {
			// Didn't recieve a response, do Nothing
		} else {
			CheckErr(err)
		}
	}
}

/*******************************************************************************
* Lazy error checker, just to save lines elsewhere.
*******************************************************************************/
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
