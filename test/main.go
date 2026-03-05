package main

import ("fmt")

func passed(msg string, a ...any) {
	fmt.Printf("[\x1b[32mGOOD\x1b[0m]:  " + msg + "\n", a...)
}
func task(msg string, a ...any) {
	fmt.Printf("\n[\x1b[36mtask\x1b[0m]:  " + msg + "\n", a...)
}
func failed(msg string, a ...any) {
	fmt.Printf("[\x1b[31;1mFAILED\x1b[0m]:  " + msg + "\n", a...)
}
