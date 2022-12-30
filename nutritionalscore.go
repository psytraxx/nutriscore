package main

// https://github.com/PacktPublishing/Learn-GO-by-Building-Three-Simple-Golang-Projects/blob/main/Project%201%20-%20Nutri%20Score/nutriscore.go

type ScoreType int

const (
	Food ScoreType = iota
	Beverage
	Water
	Cheese
)

type NutritionalScore struct {
	Value     int
	Positive  int
	Negative  int
	ScoreType ScoreType
}

var scoreToLetter = []string{"A", "B", "C", "D", "E"}

type EnergyKJ float64
type SugarGram float64
type SaturatedFattyAcids float64
type FruitsPercent float64
type SodiumMilligram float64
type FiberGram float64
type ProteinsGram float64

type NutritionalData struct {
	Energy              EnergyKJ
	Sugars              SugarGram
	SaturatedFattyAcids SaturatedFattyAcids
	Sodium              SodiumMilligram
	Fruits              FruitsPercent
	Fiber               FiberGram
	Protein             ProteinsGram
	IsWater             bool
}

var energyLevels = []float64{3350, 3015, 2680, 2345, 2010, 1675, 1340, 1005, 670, 335}
var sugarsLevels = []float64{45, 40, 36, 31, 27, 22.5, 18, 13.5, 9, 4.5}
var saturatedFattyAcidsLevels = []float64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
var sodiumLevels = []float64{900, 810, 720, 630, 540, 450, 360, 270, 180, 90}
var fibreLevels = []float64{4.7, 3.7, 2.8, 1.9, 0.9}
var proteinLevels = []float64{8, 6.4, 4.8, 3.2, 1.6}

var energyLevelsBeverage = []float64{270, 240, 210, 180, 150, 120, 90, 60, 30, 0}
var sugarsLevelsBeverage = []float64{13.5, 12, 10.5, 9, 7.5, 6, 4.5, 3, 1.5, 0}

func (e EnergyKJ) GetPoints(st ScoreType) int {
	if st == Beverage {
		return getPointsFromRange(float64(e), energyLevelsBeverage)
	} else {
		return getPointsFromRange(float64(e), energyLevels)
	}
}

func (sg SugarGram) GetPoints(st ScoreType) int {
	if st == Beverage {
		return getPointsFromRange(float64(sg), sugarsLevelsBeverage)
	} else {
		return getPointsFromRange(float64(sg), sugarsLevels)
	}
}

func (sfa SaturatedFattyAcids) GetPoints() int {
	return getPointsFromRange(float64(sfa), saturatedFattyAcidsLevels)
}

func (s SodiumMilligram) GetPoints() int {
	return getPointsFromRange(float64(s), sodiumLevels)
}

func (f FruitsPercent) GetPoints(st ScoreType) int {
	if st == Beverage {
		if f > 80 {
			return 10
		} else if f > 60 {
			return 4
		} else if f > 40 {
			return 2
		}
		return 0
	}
	if f > 80 {
		return 5
	} else if f > 60 {
		return 2
	} else if f > 40 {
		return 1
	}
	return 0
}

func (fg FiberGram) GetPoints() int {
	return getPointsFromRange(float64(fg), fibreLevels)
}

func (pg ProteinsGram) GetPoints() int {
	return getPointsFromRange(float64(pg), proteinLevels)
}

func EnergyFromKcal(kcal float64) EnergyKJ {
	return EnergyKJ(kcal * 4.184)
}

func SodiumFromSalt(saltMg float64) SodiumMilligram {
	return SodiumMilligram(saltMg / 2.5)
}

func GetNutritionalScore(n NutritionalData, st ScoreType) NutritionalScore {
	value := 0
	positive := 0
	negative := 0

	if st != Water {
		fruitPoints := n.Fruits.GetPoints(st)
		fibrePoints := n.Fiber.GetPoints()
		negative = n.Energy.GetPoints(st) + n.Sugars.GetPoints(st) + n.SaturatedFattyAcids.GetPoints() + n.Sodium.GetPoints()
		positive = fibrePoints + fruitPoints + n.Protein.GetPoints()
		if st == Cheese {
			value = negative - positive
		} else {
			if negative >= 11 && fruitPoints < 5 {
				value = negative - positive - fruitPoints
			} else {
				value = negative - positive
			}
		}
	}

	return NutritionalScore{
		Value:     value,
		Positive:  positive,
		Negative:  negative,
		ScoreType: st,
	}
}

func (ns NutritionalScore) GetNutriScore() string {
	if ns.ScoreType == Food {
		return scoreToLetter[getPointsFromRange(float64(ns.Value), []float64{18, 10, 2, -1})]
	}
	if ns.ScoreType == Water {
		return scoreToLetter[0]
	}
	return scoreToLetter[getPointsFromRange(float64(ns.Value), []float64{9, 5, 1, -1})]
}

func getPointsFromRange(v float64, steps []float64) int {
	lenSteps := len(steps)
	for i, l := range steps {
		if v > l {
			return lenSteps - i
		}
	}
	return 0
}