package becas

import (
	"fmt"
	"github.com/knodos/kmath"
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
	
	// Nota es la nota que sirve de parámetro en el algoritmo oficial
	Nota float64
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
	
	// Nmed es la nota media de todos los estudiantes. Junto con la desviación
	// típica se usa para crear una distribución normal de notas.
	Nmed float64 = 6.0
	
	// Sd es la desviación típica de las notas de los estudiantes.
	Sd float64 = 2.0
	
	// N es la cantidad de estudiantes que queremos modelar.
	N int = 1000
)

func Main() {

	var i int
	
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
	
	// Reparto aleatorio de notas (distribución normal centrada en Nmed con 
    // stddev = Sd), limitada entre 0 y 10. 
    for i = 0; i < N; i++ {
		H[i].Nota = rand.NormFloat64()*Sd+Nmed
		if H[i].Nota< 0 {
		    H[i].Nota = 0.0
		} else if H[i].Nota>10 {
		    H[i].Nota = 10.0
		}
	}
	
	println("ANTES DE LAS BECAS")	
	Datos(H)
	
	AlgoritmoDelPantano(C, H)
	Datos(H)
	
	AlgoritmoEmpujar(C, H)
	Datos(H)
	
	AlgoritmoOficial(C, H)
	Datos(H)
	
	for i:=0; i<100; i++ {
	    AlgoritmoDelPantano(C, H) // AlgoritmoOficialConNotas(C, H)
	    Datos(H)
	
	    // Simulamos una relación sencilla entre recursos y notas, y vemos
	    // el coeficiente de Pearson.
	    Examen(H)
	    Datos(H)
	
	    // Tambien hay una relación entre notas e ingresos (a medio y largo plazo)
	    // Queda contrarrestado este hecho por el algoritmo? 	
	    // Repetimos un año más el reparto de becas y vemos
	    EvolucionRenta(H)
	    Datos(H)
	}
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

    fmt.Println("ALGORITMO OFICIAL, SIN NOTAS")
    
	// Cantidad de estudiantes que no superan el umbral 2 (Rmax) -> S
	N := len(H)
	S := 0
	for i := 0; i < N; i++ {
		if H[i].Renta < Rmax {
			S++
		}
		H[i].Beca = 0
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
		log.Fatal("  [!] Algoritmo erroneo")
	}
	
	// Cuantos estudiantes pueden estudiar
	PuedenEstudiar(H)
	Pobres(H)

	return S
}

// AlgoritmoDelPantano reparte C entre H de forma que el coeficiente Gini
// sea mínimo.
func AlgoritmoDelPantano(C float64, H []Estudiante) {

    fmt.Println("ALGORITMO DEL PANTANO")
    
    c := C
    d := 0.0
    diff := 0.0
    level := 0.0
    var i int
	
	// Crea un array que podamos ordenar
	h := make([]float64,N)

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
	}
	// El resto lo repartimos
	level = h[i-1] + c/float64(i)
	fmt.Println("  Renta mínima conseguida",level,"resto",c)
	
	c = C
	for i=0; i<N; i++ {
		d = level - H[i].Renta
		
		if d<=0 {
		    continue
		}
		
		H[i].Beca = d
		c -= d
		if c<0.1 {
		    break;
		}
	}
	
	// Comprueba que la suma de becas es (casi) igual a C
	Cv := 0.0
	for i = 0; i <N; i++ {
		Cv += H[i].Beca
	}
	
	if math.Abs(Cv-C) > 0.001 {
		log.Fatal("  [!] Algoritmo erroneo. Diferencia ",Cv-C)
	}
	
	// Comprueba que todos están por encima del minimo
	for i = 0; i <N; i++ {
	    if H[i].Renta + H[i].Beca < level {
		    log.Fatal("  [!] Algoritmo erroneo. Hay rentas no niveladas")
		}
	}
	
	// Cuantos estudiantes pueden estudiar
	PuedenEstudiar(H)
	Pobres(H)
}

// AlgoritmoOficialConNotas es una simulación de la realimentación que las notas
// pueden tener sobre la renta de los estudiantes.
//
// Se asigna la beca en función de la renta y de las notas. Aquellos estudiantes
// con notas bajas pueden no recibir beca y no tener la oportunidad de seguir
// estudiando, con lo que se genera un circulo vicioso.
func AlgoritmoOficialConNotas (C float64, H []Estudiante) {

    var i int
    
    N := len(H)

    fmt.Println("ALGORITMO OFICIAL, CON NOTAS")
	
	// Cantidad de estudiantes que no superan el umbral 2 (Rmax) -> S	
	S := 0
	for i = 0; i < N; i++ {
		if H[i].Renta < Rmax {
			S++
		}
		H[i].Beca = 0 // reset beca
	}
	
	// Calcula Nmax, la nota media del mejor 10% = 0.1 S
	Nmax := NotaMediaMejor (H, S/10)
	//log.Println("Nmax",Nmax)

	// Reparto de la cantidad fija de 60 euros
	Ci := C // usado para comparar luego
	C -= 60 * float64(S)

	// Calcula sumatorio de k = K
	K := 0.0
	for i = 0; i < N; i++ {
		if H[i].Renta >= Rmax {
			continue
		}

		K += H[i].Nota / Nmax * (1 - H[i].Renta/Rmax)
	}

	// Calcula el importe por becario
	for i = 0; i < N; i++ {
		if H[i].Renta >= Rmax {
			continue
		}

		k := H[i].Nota / Nmax * (1 - H[i].Renta/Rmax)
		H[i].Beca = C*k/K + 60
	}

	// Verifica que la suma da el total
	Cv := 0.0
	for i = 0; i < N; i++ {
		Cv += H[i].Beca
	}

	if math.Abs(Cv-Ci) > 0.001 {
		log.Fatal("  Algoritmo erroneo, por ",Cv-Ci)
	}
	
	// Cuantos estudiantes pueden estudiar
	PuedenEstudiar(H)
	Pobres(H)
}

func AlgoritmoEmpujar (C float64, H []Estudiante) {

    var i int
    
    N := len(H)

    fmt.Println("ALGORITMO EMPUJAR")
    
    // Crea un array que podamos ordenar
	h := make([]float64, N)

	for i = 0; i < N; i++ {
		h[i] = H[i].Renta
		H[i].Beca = 0
	}
	
	// Ordenamos las rentas de menor a mayor
	sort.Float64s(h)
	
	// Damos dinero a los que menos les falta para llegar a Rmin + Ce
	// Averigua que renta mínima vamos a apoyar
	c := C

	for i = N-1; i >=0; i-- {
	    d := (Rmin+Ce) - h[i]
	    if d>0 {
	        if c<d {
	            break
	        }
		    c -= d    
		}
	}
	
	if i<0 {
	    i = 0
	    //println("  Reparto para todos")
	    c = C
	} else {	    
	    //fmt.Println("  Apoyamos por encima de",h[i+1],c)
	    c = h[i+1]
	}
	
	for i=0; i < N; i++ {
	    if H[i].Renta<c {
	        continue
	    }
	    d := (Rmin+Ce) - H[i].Renta
	    if d>0 {
		    H[i].Beca = d
		}
	}

	// Verifica que la suma da el total
	Cv := 0.0
	for i = 0; i < N; i++ {
		Cv += H[i].Beca
	}

	if math.Abs(Cv-C) > 0.001 {
		fmt.Println("  [!] Algoritmo erroneo, por ",Cv-C)
	}
	
	// Cuantos estudiantes pueden estudiar
	PuedenEstudiar(H)
	Pobres(H)
}

func Pobres(H []Estudiante) int {
    var i int
    
    N := len(H)

    n := 0
    p := 0.0
	for i = 0; i < N; i++ {
		if H[i].Renta + H[i].Beca < Rmin {
		    n++
		    p += Rmin - (H[i].Renta + H[i].Beca)
		}
	}
	return n
}

func PuedenEstudiar(H []Estudiante) int {
    n := 0
    b := 0.0
	for i := 0; i < len(H); i++ {
		if H[i].Renta + H[i].Beca >= Rmin + Ce {
			n++
		}
		b += H[i].Beca
	}
	return n
}

func Cobertura (H []Estudiante, lev float64) (float64,float64) {

    var i int
    
    N := len(H)
    
    // Crea un array que podamos ordenar
	h := make([]float64, N)

	for i = 0; i < N; i++ {
		h[i] = H[i].Renta + H[i].Beca
	}
	return kmath.Cover(h,lev)
}

func NotaMediaMejor (H []Estudiante, n int) float64 {

    var i int
    
    N := len(H)
    
    // Crea un array que podamos ordenar
	h := make([]float64, N)

	for i = 0; i < N; i++ {
		h[i] = H[i].Nota
	}
	
	// Ordenamos las notas de menor a mayor
	sort.Float64s(h)
	
	// Calculamos la media de los mejores n estudiantes
	m := 0.0
    for i = N-n; i < N; i++ {
		m += h[i]
	}
	
	m /= float64(n)
	
	// Cuantos estudiantes igualan o superan la nota media ?
	for i = N-1; i >=0; i-- {
		if h[i]<m {
		    break
		}
	}
	
	fmt.Printf("  De %d estudiantes, %d superan la nota media mejor (%f)\n",N,N-i,m)
	
	return m
}

// Corr devuelve el coeficiente de correlación Pearson entre becas y notas
func Corr(H []Estudiante) (float64, error) {

    // Crea arrays 
	r := make([]float64, N)
	e := make([]float64, N)

	for i := 0; i < N; i++ {
		r[i] = H[i].Renta
	}
	
	for i := 0; i < N; i++ {
		e[i] = H[i].Nota
	}
	
	return kmath.Pearson(r,e)
}

func Datos(H []Estudiante) {

    ade, aue := Cobertura(H,Rmin+Ce)
    ad, au := Cobertura(H,Rmin)
    
    gi := GiniTotal(H)
    ne := PuedenEstudiar(H)
    np := Pobres(H)
    c,_ := Corr(H)
    
    fmt.Printf("%f, %f, %f, %f, %f, %d, %d, %f\n",ade,aue,ad,au,gi,ne,np,c)
}

// Nueva nota = nota anterior +/-2 random, +/- 1 renta
// Rmin + Ce es la referencia para la renta
func Examen(H []Estudiante) {

    for i := 0; i < N; i++ {
        r := H[i].Renta / (Rmin+Ce)
        if r>2 {
            r = 2
        }     
        k := 1 + (rand.Float64()-0.5) / 5 * 2  + (r-1) * 0.1
		H[i].Nota *= k
// fmt.Println("factor nota",k, H[i].Nota)			
	
		if H[i].Nota<0 {
		    H[i].Nota = 0
		} else if H[i].Nota > 10 {
		    H[i].Nota = 10
		}
	}
}

// Nueva renta = renta anterior -10% con nota 0 +10% con 10, más un factor
// aleatorio de 10%
func EvolucionRenta (H []Estudiante) {

    for i := 0; i < N; i++ {

        k := 1 + (H[i].Nota - 5)/50 + (rand.Float64()-0.5)/5
		H[i].Renta *= k
// fmt.Println("factor renta",k, H[i].Renta)		
		if H[i].Renta<0 {
		    H[i].Renta = 0
		} 
	}
}

