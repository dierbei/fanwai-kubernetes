package main

import (
	"log"

	"githup.com/dierbei/fanwai-kubernetes/kubectl-plugin/lib"
)

func main() {
	if err := lib.RunCmd(); err != nil {
		log.Fatal(err)
	}
}
