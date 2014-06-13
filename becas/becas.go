package becas

import (
	"fmt"
	"kmath"
	"log"
	"math/rand"
	"github.com/GaryBoone/GoStats/stats"
	"math"
	"sort"
)

type Estudiante struct {
	// Renta es la renta personal del estudiante. Generamos un valor aleatorio
	// dentro de una distribución centrada en Rmed, con una dispersión que
	// acerque su coeficiente Gini al real (0.34 en 2011, para España).
	// La distribución de rentas es cercana a log-normal.
	Renta float64

	// Beca es obviamente la cantidad asignada por el algoritmo a este estudiante.
	Beca float64
}

const (
	// Rmax es la renta máxima por persona (elegimos una familia típica de 4
	// personas). Este es el umbral 2 de la tabla de becas. Por encima de este
	// umbral los estudiantes no reciben parte variable.
	Rmax float64 = 36421.0 / 4.0

	// Rmin es la renta mínima, un valor cercano al umbral de pobreza, lo minimo
	// necesario para vivir dignamente.
	Rmin float64 = 6000.0

	// Rmed es la renta media por persona. Usamos la estadística de 2011
	// http://www.ine.es/jaxi/tabla.do?type=pcaxis&path=/t25/p453/provi/l0/&file=01001.px
	Rmed float64 = 9326.0

	// Ce es el coste de los estudios. Se toma un valor fijo indicativo.
	Ce float64 = 3600.0

	// C es la dotación presupuestaria para la parte variable de becas, por persona
	Cp float64 = 2300.0
)

func Main() {

	var i int

	// N es la cantidad de estudiantes que queremos modelar.
	N := 1000
	H := make([]Estudiante, N)

	// C es la cantidad total de dinero que vamos a repartir.
	C := Cp * float64(N)

	// Asigna ingresos (anuales) a cada estudiante.
	for i = 0; i < N; i++ {
		H[i].Renta = math.Exp(rand.NormFloat64()*0.6+1)*2870
		if H[i].Renta < 0.0 {
			H[i].Renta = 0
		}
	}
	fmt.Println("Coeficiente Gini inicial", GiniTotal(H))

	// El dinero necesario para que todos puedan estudiar
	Cn := 0.0
	for i := 0; i < N; i++ {
		c := H[i].Renta - Ce - Rmin
		if c < 0 {
			Cn -= c
		}
	}
	fmt.Println("Dotación necesaria", Cn)
	fmt.Println("Dotación presupuestada", C)
	Stats(H)

	// Reparto de becas segun el algoritmo oficial
	println(AlgoritmoOficial(C, H))
	fmt.Println("Gini algoritmo oficial", GiniTotal(H))
	Stats(H)
	
	// Reparto de becas segun el algoritmo oficial
	AlgoritmoDelPantano(C, H)
	fmt.Println("Gini algoritmo pantanoso", GiniTotal(H))
	Stats(H)
	
	// Print(H)
}

func Print(H []Estudiante) {

    for i:=0;i<len(H); i++ {
        fmt.Println(H[i].Renta,H[i].Beca,H[i].Renta+H[i].Beca)
    }
}

func GiniRenta(H []Estudiante) float64 {
	in := make([]float64, len(H))

	for i := 0; i < len(H); i++ {
		in[i] = H[i].Renta
	}
	return kmath.Gini(in)
}

func GiniTotal(H []Estudiante) float64 {
	in := make([]float64, len(H))

	for i := 0; i < len(H); i++ {
		in[i] = H[i].Renta + H[i].Beca
	}
	return kmath.Gini(in)
}

func Stats(H []Estudiante) {
	in := make([]float64, len(H))

	for i := 0; i < len(H); i++ {
		in[i] = H[i].Renta + H[i].Beca
	}
	log.Println("StdDev",stats.StatsPopulationStandardDeviation(in))
	log.Println("Mean",stats.StatsMean(in))
}

// AlgoritmoOficial reparte C entre H según la fórmula oficial de las becas
// para el curso 2013-2014, con alguna simplificación: no tiene en cuenta 
// las notas, y Rmax es fijo.
func AlgoritmoOficial(C float64, H []Estudiante) int {

	// Cantidad de estudiantes que no superan el umbral 2 (Rmax) -> S
	N := len(H)
	S := 0
	for i := 0; i < N; i++ {
		if H[i].Renta < Rmax {
			S++
		}
	}

	// Reparto de la cantidad fija de 60 euros
	Ci := C // usado para comparar luego
	C -= 60 * float64(S)

	// Calcula sumatorio de k = K
	Nm := 1.0
	K := 0.0
	for i := 0; i < N; i++ {
		if H[i].Renta >= Rmax {
			continue
		}

		K += Nm * (1 - H[i].Renta/Rmax)
	}

	// Calcula el importe por becario
	for i := 0; i < N; i++ {
		if H[i].Renta >= Rmax {
			continue
		}

		k := Nm * (1 - H[i].Renta/Rmax)
		H[i].Beca = C*k/K + 60
	}

	// Verifica que la suma da el total
	Cv := 0.0
	for i := 0; i < N; i++ {
		Cv += H[i].Beca
	}

	if math.Abs(Cv-Ci) > 0.001 {
		log.Fatal("AlgoritmoOficial erroneo")
	}

	return S
}

// AlgoritmoDelPantano reparte C entre H de forma que el coeficiente Gini
// sea mínimo.
func AlgoritmoDelPantano(C float64, H []Estudiante) {

    c := C
    d := 0.0
    diff := 0.0
    level := 0.0
    var i int
    
	// Calcula el nivel del agua
	N := len(H)
	
	// Crea un array que podamos ordenar
	h := make([]float64, len(H))

	for i = 0; i < N; i++ {
		h[i] = H[i].Renta
		H[i].Beca = 0
	}
	
	// Ordenamos las rentas de menor a mayor
	sort.Float64s(h)

    // Averiguamos el nivel que podemos llenar hasta acabar con C
	for i=1; i < N; i++ {
		
		// Rectángulo: altura * anchura
		diff = (h[i] - h[i-1]) * float64(i)	
		
	    if c<diff {
	        break
	    }
		    
		c -= diff
		level = h[i]
	}
	// El resto lo repartimos
	level = h[i] + c/float64(i)
	log.Println("Nivel del agua de renta",level,"resto",c)
	
	c = C
	for i=0; i<N; i++ {
		d = level - H[i].Renta
		
		if d<=0 {
		    continue
		}
		
		if c<d {
		    H[i].Beca = c
		    break
		}
		
		H[i].Beca = d
		c -= d
	}
	
	// Comprueba que la suma de becas es (casi) igual a C
	Cv := 0.0
	for i = 0; i <N; i++ {
		Cv += H[i].Beca
	}
	
	if math.Abs(Cv-C) > 0.001 {
		log.Fatal("AlgoritmoDelPantano erroneo. Diferencia ",Cv-C)
	}
}
