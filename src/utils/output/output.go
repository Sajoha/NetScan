package output

import "fmt"

// ANSI colour representations
var Green = "\033[0;32m"
var Yellow = "\033[0;33m"
var Red = "\033[0;31m"
var Reset = "\033[0m"

func Banner(message string) {
	fmt.Println("-----------------------------------------------------------------")
	Info(message)
	fmt.Println("-----------------------------------------------------------------\n")
}

/*******************************************************************************
* Print a given message to the terminal with a green info tag
*******************************************************************************/
func Info(message string) {
	fmt.Printf("%s[INFO]%s %s\n", Green, Reset, message)
}

/*******************************************************************************
* Print a given message to the terminal with a yellow warning tag
*******************************************************************************/
func Warn(message string) {
	fmt.Printf("%s[WARN]%s %s\n", Yellow, Reset, message)
}

/*******************************************************************************
* Print a given message to the terminal with a red error tag
*******************************************************************************/
func Err(message string) {
	fmt.Printf("%s[ERROR]%s %s\n", Red, Reset, message)
}
