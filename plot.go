package main

import (
		"github.com/knodos/kmodel/graphics"
        "math/rand"
)

func main() {
        rand.Seed(int64(0))
        n := 100
        X, Y := randomTriples(n)

        graphics.DotsPaint(X,Y,"X","Y","bubble.png")
}

// randomTriples returns some random x, y, z triples
// with some interesting kind of trend.
func randomTriples(n int) ([]float64, []float64) {
        datax := make([]float64, n)
        datay := make([]float64, n)
        for i := range datay {
            datax[i] = rand.Float64()
            datay[i] = rand.Float64()  
        }
        return datax, datay
}