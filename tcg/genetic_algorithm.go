package tcg

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type GASelector interface {
	makeSelector(population *[]Chrom) func() Chrom
}

type GACrossover interface {
	crossover(parent1 Chrom, parent2 Chrom) Chrom
}

type GAMutation interface {
	mutation(gaSetting GASetting, chrom Chrom) Chrom
}

type GASetting struct {
	populationSize       int
	eliteSize            int
	crossoverProbability float64
	mutationProbability  float64
	factorInfos          []FactorInfo
	testCaseSize         int
	maxGeneration        int
	printStep            int
	GASelector
	GAMutation
	GACrossover
}

type Factor = int

type FactorInfo struct {
	min int
	max int
}

type TestCase = []Factor

type Chrom struct {
	testCases []TestCase
	fitness   int
}

type RouletteSelector struct{}
type TournamentSelector struct{}
type SingleCrossover struct{}
type IdiotMutation struct{}

func generateTestCase(gaSetting GASetting) TestCase {
	var gene = []int{}
	for _, factorInfo := range gaSetting.factorInfos {
		gene = append(gene, generateFactor(factorInfo))
	}
	return gene
}

func generateFactor(factorInfo FactorInfo) int {
	return rand.Int()%factorRange(factorInfo) + factorInfo.min
}

func generateChrom(gaSetting GASetting) Chrom {
	var testCases = []TestCase{}
	for i := 0; i < gaSetting.testCaseSize; i++ {
		testCases = append(testCases, generateTestCase(gaSetting))
	}
	return Chrom{
		testCases: testCases,
		fitness:   0,
	}
}

func generatePopulation(gaSetting GASetting) []Chrom {
	var population = make([]Chrom, gaSetting.populationSize)
	for i := range population {
		population[i] = generateChrom(gaSetting)
	}
	return population
}

func factorRange(factorInfo FactorInfo) int {
	return factorInfo.max - factorInfo.min
}

func factorTotal(factorInfos []FactorInfo) int {
	var sum = 0
	for _, factorInfo := range factorInfos {
		sum += factorRange(factorInfo)
	}
	return sum
}

func calculateFitnessSum(chrom Chrom, factorInfos []FactorInfo) int {
	var result = 0
	for t := range chrom.testCases {
		for i := range chrom.testCases[t] {
			result += chrom.testCases[t][i]
		}
	}
	return result
}

func calculateFitness(chrom Chrom, factorInfos []FactorInfo) int {
	var factorTotal = factorTotal(factorInfos)
	var pairCount = make([][]int, factorTotal)
	for i := 0; i < factorTotal; i++ {
		pairCount[i] = make([]int, factorTotal)
		for j := 0; j < factorTotal; j++ {
			pairCount[i][j] = 0
		}
	}

	for _, t := range chrom.testCases {
		for i := 0; i < len(t); i++ {
			for j := 0; j < i; j++ {
				pairCount[t[i]][t[j]] += 1
			}
		}
	}

	var result = 0
	var repeated = 0
	for i := 0; i < factorTotal; i++ {
		for j := 0; j < factorTotal; j++ {
			if pairCount[i][j] > 0 {
				result += 1
			}
			if pairCount[i][j] > 1 {
				repeated += 1
			}
		}
	}

	return result - repeated
}

func evaluate(population *[]Chrom, factorInfos []FactorInfo) {
	var wg sync.WaitGroup
	wg.Add(len(*population))
	for i := range *population {
		go func(i int) {
			(*population)[i].fitness = calculateFitness((*population)[i], factorInfos)
			wg.Done()
		}(i)
	}
	wg.Wait()

	sort.Slice(*population, func(i, j int) bool {
		return (*population)[i].fitness > (*population)[j].fitness
	})
}

func (idiotCrossover *IdiotMutation) mutation(gaSetting GASetting, chrom Chrom) Chrom {
	var child = copyChrom(chrom)
	for t := range chrom.testCases {
		for i := range chrom.testCases[t] {
			if rand.Float64() < gaSetting.mutationProbability {
				child.testCases[t][i] = generateFactor(gaSetting.factorInfos[i])
			}
		}
	}
	return child
}

func copyChrom(chrom Chrom) Chrom {
	child := chrom
	child.testCases = make([]TestCase, len(chrom.testCases))

	for t := range chrom.testCases {
		child.testCases[t] = make([]int, len(chrom.testCases[t]))
		copy(child.testCases[t], chrom.testCases[t])
	}

	return child
}

func (singleCrossover *SingleCrossover) crossover(parent1 Chrom, parent2 Chrom) Chrom {
	child := copyChrom(parent1)
	crossoverPoint := rand.Intn(len(parent1.testCases))

	for t := range parent2.testCases {
		for i := range parent2.testCases[t] {
			if i < crossoverPoint {
				child.testCases[t][i] = parent2.testCases[t][i]
			}
		}
	}
	return child
}

func (rouletteSelector *RouletteSelector) makeSelector(population *[]Chrom) func() Chrom {
	totalFitness := 0
	for i := range *population {
		totalFitness += (*population)[i].fitness
	}

	return func() Chrom {
		r := rand.Intn(totalFitness) + 1
		index := 0

		for i := range *population {
			if r <= (*population)[i].fitness {
				index = i
				break
			}
			r -= (*population)[i].fitness
		}

		return (*population)[index]
	}
}

func (tournamentSelector *TournamentSelector) makeSelector(population *[]Chrom) func() Chrom {
	return func() Chrom {
		result := []Chrom{}
		for i := 0; i < 10; i++ {
			r := rand.Intn(len(*population))
			result = append(result, (*population)[r])
		}

		sort.Slice(result, func(i, j int) bool {
			return result[i].fitness > result[j].fitness
		})

		return result[0]
	}
}

func alternate(population *[]Chrom, gaSetting GASetting) []Chrom {
	selector := gaSetting.GASelector.makeSelector(population)
	result := make([]Chrom, len(*population))
	var wg sync.WaitGroup
	wg.Add(len(*population))
	for i := range *population {
		go func(i int) {
			var child Chrom
			if len(*population)-1 <= i {
				child = generateChrom(gaSetting)
			} else {
				if i < gaSetting.eliteSize {
					child = (*population)[i]
				}
				if gaSetting.eliteSize <= i {
					parent1 := selector()
					parent2 := selector()
					child = gaSetting.GACrossover.crossover(parent1, parent2)
					child = gaSetting.GAMutation.mutation(gaSetting, child)
				}
			}
			result[i] = child
			wg.Done()
		}(i)
	}
	wg.Wait()
	evaluate(&result, gaSetting.factorInfos)
	return result
}

func PrintStatistics(population []Chrom) {
	fitnessAvg := 0.0
	for i := range population {
		fitnessAvg += float64(population[i].fitness)
	}
	fitnessAvg = fitnessAvg / float64(len(population))
	fmt.Printf("fitnessAvg: %#v,", fitnessAvg)
	fmt.Printf("fitnessMax: %#v,", (population)[0].fitness)
	fmt.Printf("fitnessMin: %#v\n", (population)[len(population)-1].fitness)
}

func buildFactorInfos(params int, levels int) []FactorInfo {
	max := 0
	var factorInfos = []FactorInfo{}
	for i := 0; i < params; i++ {
		factorInfos = append(factorInfos, FactorInfo{min: max, max: max + levels})
		max += levels
	}
	return factorInfos
}

func Evolution(factorInfos []FactorInfo) []TestCase {
	rand.Seed(time.Now().Unix())
	var gaSetting = GASetting{
		populationSize:       50,
		eliteSize:            15,
		crossoverProbability: 0.1,
		mutationProbability:  0.01,
		factorInfos:          factorInfos,
		testCaseSize:         15,
		maxGeneration:        20000,
		printStep:            100,
		GASelector:           &RouletteSelector{},
		GAMutation:           &IdiotMutation{},
		GACrossover:          &SingleCrossover{},
	}
	var population = generatePopulation(gaSetting)
	evaluate(&population, gaSetting.factorInfos)
	for i := 0; i < gaSetting.maxGeneration; i++ {
		population = alternate(&population, gaSetting)
		if i%gaSetting.printStep == 0 {
			PrintStatistics(population)
		}
	}
	return population[0].testCases
}
