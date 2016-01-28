package main

import (
	"flag"
//	"log"
)

var (
    isTest bool
)

func init() {
	flag.BoolVar(&isTest, "image", false, "Run as Test")
}

func main() {
	flag.Parse()
    if isTest {
        log.Println("Runnig image test")
        runImageTest()
    } else {
        port := ":8080"
        log.Println("Runnig webserver on port " + port)
        runWs(port)
    }
}