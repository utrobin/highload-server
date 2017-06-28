package main

import (
	"flag"
	"log"
)

var CPU_NUM = flag.Int("c", 0, "Num of CPU's used in server")
var ROOT_PATH = flag.String("r", "", "Path to root directory")


func main() {
	flag.Parse()

	if *CPU_NUM <= 0 {
		log.Println("Invalid CPU number! Use -c for choose NCPU")
		return
	}
	if *ROOT_PATH == "" {
		log.Println("Invalid root path! Use -r to select root path. Choosing default value...")
		*ROOT_PATH = "DOCUMENT_ROOT";
	}

	startServer()
}



