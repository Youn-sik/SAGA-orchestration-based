package main

import "saga/router"

func main() {
	hsvr := router.SetUp()
	router.RunHTTP(hsvr)
}

/*



 */
