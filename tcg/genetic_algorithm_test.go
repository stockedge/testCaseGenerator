package tcg

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestFactorRange(t *testing.T) {
	factorInfo := FactorInfo{
		min: 0,
		max: 4,
	}
	if factorRange(factorInfo) != 4 {
		t.Fail()
	}
}

func TestFactorTotal(t *testing.T) {
	factorInfos := []FactorInfo{
		{
			min: 0,
			max: 4,
		},
		{
			min: 4,
			max: 8,
		},
	}
	if factorTotal(factorInfos) != 8 {
		t.Fail()
	}
}

func TestGenerateTestCase(t *testing.T) {
	factorInfos := []FactorInfo{
		{
			min: 0,
			max: 4,
		},
	}

	var gaSetting = GASetting{
		populationSize:       50,
		eliteSize:            15,
		crossoverProbability: 0.1,
		mutationProbability:  0.01,
		factorInfos:          factorInfos,
		testCaseSize:         150,
		maxGeneration:        20000,
		printStep:            100,
		GASelector:           &RouletteSelector{},
		GAMutation:           &IdiotMutation{},
		GACrossover:          &SingleCrossover{},
	}

	generatedValues := [...]int{0, 0, 0, 0}
	for i := 0; i < 1000; i++ {
		testCase := generateTestCase(gaSetting)
		generatedValues[testCase[0]] += 1
	}

	for i := 0; i < 4; i++ {
		if generatedValues[i] == 0 {
			t.Fail()
		}
	}
}

func TestEvaluate(t *testing.T) {
	var factorInfos = []FactorInfo{
		{min: 0, max: 4},
		{min: 4, max: 8},
		{min: 8, max: 12},
		{min: 12, max: 20},
		{min: 20, max: 30},
	}
	rand.Seed(time.Now().Unix())
	var gaSetting = GASetting{
		populationSize:       50,
		eliteSize:            15,
		crossoverProbability: 0.1,
		mutationProbability:  0.01,
		factorInfos:          factorInfos,
		testCaseSize:         150,
		maxGeneration:        20000,
		printStep:            100,
		GASelector:           &RouletteSelector{},
		GAMutation:           &IdiotMutation{},
		GACrossover:          &SingleCrossover{},
	}

	var population = generatePopulation(gaSetting)
	evaluate(&population, factorInfos)

	for i := 0; i < len(population)-1; i++ {
		if population[i].fitness < population[i+1].fitness {
			t.Fail()
		}
	}
}

func TestCrossover(t *testing.T) {
	var factorInfos = []FactorInfo{
		{min: 0, max: 4},
		{min: 4, max: 8},
		{min: 8, max: 12},
		{min: 12, max: 20},
		{min: 20, max: 30},
	}
	var gaSetting = GASetting{
		populationSize:       50,
		eliteSize:            15,
		crossoverProbability: 0.1,
		mutationProbability:  0.01,
		factorInfos:          factorInfos,
		testCaseSize:         150,
		maxGeneration:        20000,
		printStep:            100,
		GASelector:           &RouletteSelector{},
		GAMutation:           &IdiotMutation{},
		GACrossover:          &SingleCrossover{},
	}

	rand.Seed(time.Now().Unix())
	parent1 := generateChrom(gaSetting)
	parent2 := generateChrom(gaSetting)
	child := gaSetting.GACrossover.crossover(parent1, parent2)
	for test := range child.testCases {
		for i := range child.testCases[test] {
			if parent1.testCases[test][i] != child.testCases[test][i] && parent2.testCases[test][i] != child.testCases[test][i] {
				t.Fail()
			}
		}
	}
}

func TestMakeSelector(t *testing.T) {
	var factorInfos = []FactorInfo{
		{min: 0, max: 4},
		{min: 4, max: 8},
		{min: 8, max: 12},
		{min: 12, max: 20},
		{min: 20, max: 30},
	}

	var gaSetting = GASetting{
		populationSize:       50,
		eliteSize:            15,
		crossoverProbability: 0.1,
		mutationProbability:  0.01,
		factorInfos:          factorInfos,
		testCaseSize:         150,
		maxGeneration:        20000,
		printStep:            100,
		GASelector:           &RouletteSelector{},
		GAMutation:           &IdiotMutation{},
		GACrossover:          &SingleCrossover{},
	}

	rand.Seed(time.Now().Unix())
	var population = generatePopulation(gaSetting)
	evaluate(&population, factorInfos)
	PrintStatistics(population)

	selector := gaSetting.GASelector.makeSelector(&population)

	var totalFitness = 0
	for i := 0; i < 100; i++ {
		totalFitness += selector().fitness
	}
	fmt.Print(totalFitness / 100)

}

func TestMutation(t *testing.T) {
	//TODO
}

func TestEvolution(t *testing.T) {
	Evolution(nil)
}
