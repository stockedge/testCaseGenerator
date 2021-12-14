package tcg

import (
	"io/ioutil"
	"strings"
)

type FactorInfoMap = map[string][]string

func ParseFactorInfo(input string) FactorInfoMap {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	factorInfoMap := map[string][]string{}

	for _, line := range lines {
		splitted := strings.SplitN(line, ":", 2)
		factor := strings.TrimSpace(splitted[0])
		var levels []string
		for _, level := range strings.Split(splitted[1], ",") {
			levels = append(levels, strings.TrimSpace(level))
		}
		factorInfoMap[factor] = levels
	}

	return factorInfoMap
}

func ReadFile(fileName string) string {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func ConvertFactorInfoMapToStruct(factorInfoMap FactorInfoMap) []FactorInfo {
	var result []FactorInfo
	var size = 0
	var prevSize = 0
	for i := range factorInfoMap {
		size += len(factorInfoMap[i])
		result = append(result,
			FactorInfo{
				min: prevSize,
				max: size,
			})
		prevSize = size
	}
	return result
}
