package main

import "fmt"

func main() {
	ns := GetNutritionalScore(NutritionalData{
		Energy:              EnergyFromKcal(100),
		Sugars:              SugarGram(10),
		SaturatedFattyAcids: SaturatedFattyAcids(2),
		Sodium:              SodiumMilligram(500),
		Fruits:              FruitsPercent(60),
		Fiber:               FiberGram(4),
		Protein:             ProteinsGram(2),
	}, Food)

	fmt.Printf("Nutritional Score: %d\n", ns.Value)
	fmt.Printf("NutritiScore: %s\n", ns.GetNutriScore())
}
