package main

import (
	"fmt"
	"math"
	"sort"

	"sync"
)

var wg sync.WaitGroup

//Starting point of segment
var a float64

//Ending point of segment
var b float64

//Precision
var eps float64

//Reliability parameter
var r float64

var maxIndex_ int = -1

var funcs []func(x float64) float64 = make([]func(x float64) float64, 0, 3)

var paramsZ_ []float64 = make([]float64, cap(funcs))

var maxValuesDifference_ []float64 = make([]float64, cap(funcs))

var trials_ []pointTrial

var fixedIndex_ [][]pointTrial = make([][]pointTrial, cap(funcs))

func f1(x float64) float64 {
	return math.Exp(-0.5*x) * math.Sin(6*x-1.5)
}

// Function 2
func f2(x float64) float64 {
	return math.Abs(x) * math.Sin(2*math.Pi*x-0.5)
}

// Function 3
func fi(x float64) float64 {
	return math.Cos(18*x-3)*math.Sin(10*x-7) + 1.5
}

type pointTrial struct {
	x     float64
	value float64
	index int
}

func InputCheck(a, b, eps, r float64) string {
	if a > b {
		return "Left bound is more than right bound"
	}
	if eps <= 0 {
		return "Epsilon is less then or equal to 0"
	}
	if r <= 1 {
		return "Parameter of reliability is less then or equal to 1"
	}
	return "OK"
}

func NewTrial(x float64, c chan pointTrial) {

	var index int = -1
	var value float64 = 0
	for i, vall := range funcs {
		index = i
		value = vall(x)
		if value >= 0 {

			break
		}
	}
	c <- pointTrial{x, value, index} // >??????????????
}

func NewTrialInt(vector []pointTrial, c chan pointTrial) {
	var NewX float64
	if vector[0].index != vector[1].index {
		NewX = float64(vector[0].x+vector[1].x) / 2
		NewTrial(NewX, c)

	} else {
		var index_ int = vector[0].index
		var newX float64 = float64(vector[0].x+vector[1].x)/2 - float64(vector[0].value-vector[1].value)/(2*r*maxValuesDifference_[index_])
		NewTrial(newX, c)
	}
}

func calculateFixedIndex() {
	defer wg.Done()
	for i := 0; i < len(fixedIndex_); i++ {
		fixedIndex_[i] = fixedIndex_[i][:0]
	}
	sort.Slice(trials_, func(i, j int) (less bool) { return trials_[i].x < trials_[j].x })
	for i := 1; i < len(trials_)-1; i++ {
		fixedIndex_[trials_[i].index] = append(fixedIndex_[trials_[i].index], trials_[i])
	}

}

func calculateMaxValuesDifference() {
	defer wg.Done()
	maxM := float64(1)
	for i := 0; i < len(fixedIndex_); i++ {
		if len(fixedIndex_[i]) < 2 {
			maxValuesDifference_[i] = 1
		} else {
			maxM = -10000000
			for j := 0; j < len(fixedIndex_[i])-1; j++ {
				itPrev := fixedIndex_[i][j]
				itCurr := fixedIndex_[i][j+1]
				tempM := math.Abs(itCurr.value-itPrev.value) / (itCurr.x - itPrev.x)
				if tempM > maxM {
					maxM = tempM
				}
			}
			if maxM > 0 {
				maxValuesDifference_ = append(maxValuesDifference_, maxM)
			} else {
				maxValuesDifference_ = append(maxValuesDifference_, 1)
			}
		}
	}
}

func calculateZ(bestTrial pointTrial) {
	defer wg.Done()
	for i := 0; i < len(paramsZ_); i++ {
		if i == maxIndex_ {
			paramsZ_[i] = bestTrial.value
		} else {
			paramsZ_[i] = 0
		}
	}

}

func calculateMaxR(ch chan []pointTrial) {

	bestPrev := pointTrial{0, 0, -1}
	bestCurr := pointTrial{0, 0, -1}
	maxR := float64(-100000000)
	index := 1
	res := make([]pointTrial, 0)

	for i := 0; i < len(trials_)-1; i++ {

		pointPrev := trials_[i]
		pointCurr := trials_[i+1]
		currR := float64(0)
		delta := float64(pointCurr.x - pointPrev.x)

		switch {
		case pointPrev.index == pointCurr.index:
			index = pointPrev.index
			currR = delta + math.Pow(pointCurr.value-pointPrev.value, 2)/(delta*math.Pow(r, 2)*math.Pow(maxValuesDifference_[index], 2)) - 2*(pointCurr.value+pointPrev.value-2*paramsZ_[index])/(r*maxValuesDifference_[index])

		case pointPrev.index < pointCurr.index:
			index = pointCurr.index
			currR = 2*delta - 4*(pointPrev.value-paramsZ_[index])/(r*maxValuesDifference_[index])
		default:
			index = pointPrev.index
			currR = 2*delta - 4*(pointPrev.value-paramsZ_[index])/(r*maxValuesDifference_[index])
		}

		if currR > maxR {
			maxR = currR
			bestPrev = pointPrev
			bestCurr = pointCurr
		}
	}
	ch <- append(res, bestPrev, bestCurr)

}

func main() {
	a := float64(3)
	b := float64(4)
	r := float64(1.1)
	eps := float64(0.12)
	c := make(chan pointTrial, 2)
	ch := make(chan []pointTrial)

	Warning := InputCheck(a, b, eps, r)
	if Warning != "OK" {
		print(Warning)
	} else {

		funcs = append(funcs, f1, f2, fi)
		trials_ = append(trials_, pointTrial{a, 0, -1}, pointTrial{b, 0, -1})
		go NewTrial(float64(a+b)/2, c)
		bestTrial := <-c
		//fmt.Println("bestTrial = ", bestTrial)
		trials_ = append(trials_, bestTrial)
		sort.Slice(trials_, func(i, j int) (less bool) { return trials_[i].x < trials_[j].x })
		maxIndex_ = bestTrial.index
		stop := false
		i := 0
		for !stop {
			i = i + 1
			fmt.Println("Прогон номер = ", i)
			wg.Add(3)
			go calculateFixedIndex()
			go calculateMaxValuesDifference()
			go calculateZ(bestTrial)
			wg.Wait()
			fmt.Println("fixedIndex_ = ", fixedIndex_)
			fmt.Println("maxValuesDifference_ = ", maxValuesDifference_)
			fmt.Println("paramsZ_ = ", paramsZ_)

			go calculateMaxR(ch)

			currInterval := <-ch
			fmt.Println("currInterval = ", currInterval)

			if math.Abs(currInterval[1].x-currInterval[0].x) < eps {
				stop = true
			} else {
				NewTrialInt(currInterval, c)
				currTrial := <-c
				fmt.Println("currTrial = ", currTrial)
				trials_ = append(trials_, currTrial)
				sort.Slice(trials_, func(i, j int) (less bool) { return trials_[i].x < trials_[j].x })
				if (currTrial.index > maxIndex_) || (currTrial.index == maxIndex_ && currTrial.value < bestTrial.value) {
					bestTrial = currTrial
					maxIndex_ = currTrial.index
				}
			}

		}

		fmt.Println(bestTrial)
	}
}
