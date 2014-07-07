package renta

import (
	"fmt"
	"github.com/knodos/kmath"
	"math"
)

type Persona struct {
	// Renta es la renta personal. Generamos un valor aleatorio
	// dentro de una distribución centrada en Rmed, con una dispersión que
	// acerque su coeficiente Gini al real (0.34 en 2011, para España).
	// La distribución de rentas es cercana a log-normal.
	Renta float64
	
	// Impuesto sobre la renta 
	Impuesto float64
}

const (
	// Rmin es la renta mínima, un valor cercano al umbral de pobreza, lo minimo
	// necesario para vivir dignamente.
	Rmin float64 = 6000.0

	// Rmed es la renta media por persona. Usamos la estadística de 2011
	// http://www.ine.es/jaxi/tabla.do?type=pcaxis&path=/t25/p453/provi/l0/&file=01001.px
	Rmed float64 = 9326.0
	
	// IRmax es el porcentaje máximo de IRPF 
	Imax float64 = 60.0
	
	// N es la cantidad de personas que queremos modelar.
	N int = 1000
)

func Main() {
	
	H := make([]Persona, N)

	// Asigna ingresos (anuales) a cada persona
	DistribuyeRenta(H)	
	Datos(H)
	
	IrpfTramos2014(H)
	Datos(H)
	
	IrpfTramos2015(H)
	Datos(H)
	
	IrpfProgresivo(H)
	Datos(H)
}

// Fuente AEAT (Estadisticas IRPF, redimiento neto 2011)
//
// Negativo y Cero	1,15
// <1500	5,80
// 1,5 - 6 13,63
// 6 - 12 18,86
// 12 - 21 26,54
// 21 - 30	15,60
// 30 - 60	14,90
// 60 - 150 3,14
// 150 - 601 0,35
// >601  0,03

func DistribuyeRenta(H []Persona) {

    ri := 500
    rf := 300000
     
    r := (rf-ri) / N
    
    for i := 0; i < N; i++ {
        H[i].Renta = float64(ri)
        ri += r
    }
}

//func RentaCheck(H []Persona) float64 {
//}

func GiniRenta(H []Persona) float64 {
	in := make([]float64, len(H))

	for i := 0; i < len(H); i++ {
		in[i] = H[i].Renta - H[i].Impuesto
	}
	return kmath.Gini(in)
}

func Recaudacion(H []Persona) float64 {
	
    c := 0.0

	for i := 0; i < N; i++ {
		c += H[i].Impuesto
	}
	return c
}

func Cobertura (H []Persona, lev float64) (float64,float64) {

    var i int
    
    // Crea un array que podamos ordenar
	h := make([]float64, N)

	for i = 0; i < N; i++ {
		h[i] = H[i].Renta - H[i].Impuesto
	}
	return kmath.Cover(h,lev)
}


func Datos(H []Persona) {

    ad, au := Cobertura(H,Rmin)   
    gi := GiniRenta(H)
    r := Recaudacion(H)
    
    fmt.Printf("%f, %f, %f, %f\n",ad,au,gi,r)
}

func IrpfProgresivo(H []Persona) {
    
    coef := (47500-6000) / 25.00
     
    for i := 0; i < N; i++ {
        pc := H[i].Renta / coef + 5
        if pc > Imax {
            pc = Imax
        }
        //fmt.Printf("%f %f\n",H[i].Renta,pc)
        H[i].Impuesto = H[i].Renta * pc / 100
    }
}

// Fuente: http://elpais.com/elpais/2014/06/20/media/1403284152_965185.html
// 17707 : 24.75%
// 33007 : 30%
// 53407 : 40%
// 120000: 47%
// 175000: 49%
// 300000: 51%
// > : 52%
func IrpfTramos2014(H []Persona) {

    tr := [...]float64{17707, 33007, 53407, 120000, 175000, 300000, math.MaxFloat64 }
    pc := [...]float64{24.75, 30, 40, 47, 49, 51, 52 }

    for i := 0; i < N; i++ {
        for j:=0;j<len(tr); j++ {
            if H[i].Renta <= tr[j] {
                H[i].Impuesto = H[i].Renta * pc[j] / 100.0
                break 
            }
        }
        
	}
}

// 12450: 20%
// 20200: 25%
// 35200: 31%
// 60000: 39%
// > 47%
func IrpfTramos2015(H []Persona) {

    tr := [...]float64{ 12450, 20200, 35200, 60000, math.MaxFloat64 }
    pc := [...]float64{ 20,25,31,39,47 }

    for i := 0; i < N; i++ {
        for j:=0;j<len(tr); j++ {
            if H[i].Renta <= tr[j] {
                H[i].Impuesto = H[i].Renta * pc[j] / 100.0
                break 
            }
        }
        
	}
}



