package main

import (
	"reflect"
	"testing"
)

func TestEnergyFromKcal(t *testing.T) {
	type args struct {
		kcal float64
	}
	tests := []struct {
		name string
		args args
		want EnergyKJ
	}{{name: "Zero",
		want: EnergyKJ(0),
	}, {name: "Standard",
		args: struct{ kcal float64 }{kcal: 1},
		want: EnergyKJ(4.184),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EnergyFromKcal(tt.args.kcal); got != tt.want {
				t.Errorf("EnergyFromKcal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnergyKJ_GetPoints(t *testing.T) {
	type args struct {
		st ScoreType
	}
	tests := []struct {
		name string
		e    EnergyKJ
		args args
		want int
	}{
		{name: "Zero", e: EnergyFromKcal(0), args: struct{ st ScoreType }{st: Water}, want: 0},
		{name: "Beverage 10", e: EnergyFromKcal(100), args: struct{ st ScoreType }{st: Beverage}, want: 10},
		{name: "Non Beverage 10", e: EnergyFromKcal(100), args: struct{ st ScoreType }{st: Food}, want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.GetPoints(tt.args.st); got != tt.want {
				t.Errorf("GetPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFiberGram_GetPoints(t *testing.T) {
	tests := []struct {
		name string
		fg   FiberGram
		want int
	}{
		{name: "5 Point", fg: FiberGram(1000), want: 5},
		{name: "0 Point", fg: FiberGram(0.9), want: 0},
		{name: "1 Point", fg: FiberGram(1.9), want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fg.GetPoints(); got != tt.want {
				t.Errorf("GetPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFruitsPercent_GetPoints(t *testing.T) {
	type args struct {
		st ScoreType
	}
	tests := []struct {
		name string
		f    FruitsPercent
		args args
		want int
	}{
		{name: "Nothing", f: 0, args: struct{ st ScoreType }{st: Food}, want: 0},
		{name: "Beverage 10", f: 10, args: struct{ st ScoreType }{st: Beverage}, want: 0},
		{name: "Beverage 20", f: 41, args: struct{ st ScoreType }{st: Beverage}, want: 2},
		{name: "Beverage 61", f: 61, args: struct{ st ScoreType }{st: Beverage}, want: 4},
		{name: "Beverage 90", f: 90, args: struct{ st ScoreType }{st: Beverage}, want: 10},
		{name: "Other", f: 10, args: struct{ st ScoreType }{st: Food}, want: 0},
		{name: "Other", f: 41, args: struct{ st ScoreType }{st: Food}, want: 1},
		{name: "Other", f: 61, args: struct{ st ScoreType }{st: Food}, want: 2},
		{name: "Other", f: 81, args: struct{ st ScoreType }{st: Food}, want: 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.GetPoints(tt.args.st); got != tt.want {
				t.Errorf("GetPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNutritionalScore(t *testing.T) {
	type args struct {
		n  NutritionalData
		st ScoreType
	}
	tests := []struct {
		name string
		args args
		want NutritionalScore
	}{
		{name: "Nutella", args: struct {
			n  NutritionalData
			st ScoreType
		}{n: struct {
			Energy              EnergyKJ
			Sugars              SugarGram
			SaturatedFattyAcids SaturatedFattyAcids
			Sodium              SodiumMilligram
			Fruits              FruitsPercent
			Fiber               FiberGram
			Protein             ProteinsGram
			IsWater             bool
		}{Energy: 200, Sugars: 21, SaturatedFattyAcids: 4, Sodium: 15, Fruits: 0, Fiber: 0, Protein: 2, IsWater: false}, st: Food}, want: struct {
			Value     int
			Positive  int
			Negative  int
			ScoreType ScoreType
		}{Value: 6, Positive: 1, Negative: 7, ScoreType: Food}},
		{name: "Water", args: struct {
			n  NutritionalData
			st ScoreType
		}{n: struct {
			Energy              EnergyKJ
			Sugars              SugarGram
			SaturatedFattyAcids SaturatedFattyAcids
			Sodium              SodiumMilligram
			Fruits              FruitsPercent
			Fiber               FiberGram
			Protein             ProteinsGram
			IsWater             bool
		}{Energy: 0, Sugars: 0, SaturatedFattyAcids: 0, Sodium: 0, Fruits: 0, Fiber: 0, Protein: 0, IsWater: true}, st: Water}, want: struct {
			Value     int
			Positive  int
			Negative  int
			ScoreType ScoreType
		}{Value: 0, Positive: 0, Negative: 0, ScoreType: Water}},

		{name: "Cheese", args: struct {
			n  NutritionalData
			st ScoreType
		}{n: struct {
			Energy              EnergyKJ
			Sugars              SugarGram
			SaturatedFattyAcids SaturatedFattyAcids
			Sodium              SodiumMilligram
			Fruits              FruitsPercent
			Fiber               FiberGram
			Protein             ProteinsGram
			IsWater             bool
		}{Energy: 356, Sugars: 2.22, SaturatedFattyAcids: 27.4, Sodium: 819, Fruits: 0, Fiber: 0, Protein: 24.9, IsWater: false}, st: Cheese}, want: struct {
			Value     int
			Positive  int
			Negative  int
			ScoreType ScoreType
		}{Value: 15, Positive: 5, Negative: 20, ScoreType: Cheese}},

		{name: "Steak", args: struct {
			n  NutritionalData
			st ScoreType
		}{n: struct {
			Energy              EnergyKJ
			Sugars              SugarGram
			SaturatedFattyAcids SaturatedFattyAcids
			Sodium              SodiumMilligram
			Fruits              FruitsPercent
			Fiber               FiberGram
			Protein             ProteinsGram
			IsWater             bool
		}{Energy: 179, Sugars: 0, SaturatedFattyAcids: 7.6, Sodium: 60, Fruits: 0, Fiber: 0, Protein: 26.0, IsWater: false}, st: Food}, want: struct {
			Value     int
			Positive  int
			Negative  int
			ScoreType ScoreType
		}{Value: 2, Positive: 5, Negative: 7, ScoreType: Food}},
		{name: "Apple", args: struct {
			n  NutritionalData
			st ScoreType
		}{n: struct {
			Energy              EnergyKJ
			Sugars              SugarGram
			SaturatedFattyAcids SaturatedFattyAcids
			Sodium              SodiumMilligram
			Fruits              FruitsPercent
			Fiber               FiberGram
			Protein             ProteinsGram
			IsWater             bool
		}{Energy: 104, Sugars: 20.8, SaturatedFattyAcids: 0.2, Sodium: 2, Fruits: 100, Fiber: 4.8, Protein: 0.5, IsWater: false}, st: Food}, want: struct {
			Value     int
			Positive  int
			Negative  int
			ScoreType ScoreType
		}{Value: -6, Positive: 10, Negative: 4, ScoreType: Food}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNutritionalScore(tt.args.n, tt.args.st); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNutritionalScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNutritionalScore_GetNutriScore(t *testing.T) {
	type fields struct {
		Value     int
		Positive  int
		Negative  int
		ScoreType ScoreType
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "Nutella", fields: struct {
			Value     int
			Positive  int
			Negative  int
			ScoreType ScoreType
		}{Value: 6, Positive: 1, Negative: 7, ScoreType: Food}, want: "C"},
		{name: "Water", fields: struct {
			Value     int
			Positive  int
			Negative  int
			ScoreType ScoreType
		}{Value: 0, Positive: 0, Negative: 0, ScoreType: Water}, want: "A"},

		{name: "Beverage", fields: struct {
			Value     int
			Positive  int
			Negative  int
			ScoreType ScoreType
		}{Value: 2, Positive: 0, Negative: 0, ScoreType: Beverage}, want: "C"},

		{name: "Cheese", fields: struct {
			Value     int
			Positive  int
			Negative  int
			ScoreType ScoreType
		}{Value: 15, Positive: 5, Negative: 20, ScoreType: Cheese}, want: "E"},

		{name: "Steak", fields: struct {
			Value     int
			Positive  int
			Negative  int
			ScoreType ScoreType
		}{Value: 2, Positive: 5, Negative: 7, ScoreType: Food}, want: "B"},
		{name: "Apple", fields: struct {
			Value     int
			Positive  int
			Negative  int
			ScoreType ScoreType
		}{Value: -6, Positive: 10, Negative: 4, ScoreType: Food}, want: "A"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns := NutritionalScore{
				Value:     tt.fields.Value,
				Positive:  tt.fields.Positive,
				Negative:  tt.fields.Negative,
				ScoreType: tt.fields.ScoreType,
			}
			if got := ns.GetNutriScore(); got != tt.want {
				t.Errorf("GetNutriScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProteinGram_GetPoints(t *testing.T) {
	tests := []struct {
		name string
		pg   ProteinsGram
		want int
	}{
		{name: "Default", pg: 10, want: 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pg.GetPoints(); got != tt.want {
				t.Errorf("GetPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSaturatedFattyAcids_GetPoints(t *testing.T) {
	tests := []struct {
		name string
		sfa  SaturatedFattyAcids
		want int
	}{
		{name: "10", sfa: 10, want: 9},
		{name: "Default", sfa: 0, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sfa.GetPoints(); got != tt.want {
				t.Errorf("GetPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSodiumFromSalt(t *testing.T) {
	type args struct {
		saltMg float64
	}
	tests := []struct {
		name string
		args args
		want SodiumMilligram
	}{
		{name: "Default", args: struct{ saltMg float64 }{saltMg: 100.0}, want: 40},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SodiumFromSalt(tt.args.saltMg); got != tt.want {
				t.Errorf("SodiumFromSalt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSodiumMilligram_GetPoints(t *testing.T) {
	tests := []struct {
		name string
		s    SodiumMilligram
		want int
	}{
		{name: "Zero", s: 0, want: 0},
		{name: "10", s: 2120, want: 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetPoints(); got != tt.want {
				t.Errorf("GetPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSugarGram_GetPoints(t *testing.T) {
	type args struct {
		st ScoreType
	}
	tests := []struct {
		name string
		sg   SugarGram
		args args
		want int
	}{
		{name: "Beverage", sg: 10, args: struct{ st ScoreType }{st: Beverage}, want: 7},
		{name: "Other", sg: 10, args: struct{ st ScoreType }{st: Food}, want: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sg.GetPoints(tt.args.st); got != tt.want {
				t.Errorf("GetPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPointsFromRange(t *testing.T) {
	type args struct {
		v     float64
		steps []float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "Simple", args: struct {
			v     float64
			steps []float64
		}{v: 12000, steps: energyLevels}, want: 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPointsFromRange(tt.args.v, tt.args.steps); got != tt.want {
				t.Errorf("getPointsFromRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
