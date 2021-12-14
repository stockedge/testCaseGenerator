package main

import (
	"flag"
	"fmt"
)

import "testCaseGenerator/tcg"

func main() {

	flag.Parse()
	factorInfoMap := tcg.ParseFactorInfo(tcg.ReadFile(flag.Args()[0]))
	factorInfo := tcg.ConvertFactorInfoMapToStruct(factorInfoMap)

	fmt.Println(tcg.Evolution(factorInfo))
}
