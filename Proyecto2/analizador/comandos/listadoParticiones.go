package comandos

import (
	"fmt"
	"strconv"
)

type Nodo struct {
	siguiente     *Nodo
	n_particiones *Nparticiones
}

type Lista struct {
	primero  *Nodo
	ultimo   *Nodo
	contador int
}

func New_Nodo(n_particiones *Nparticiones) *Nodo {
	return &Nodo{nil, n_particiones}
}

func New_Lista() *Lista {
	return &Lista{nil, nil, 0}
}

func Insertar(n_particiones *Nparticiones, lista *Lista) {
	var nuevo *Nodo = New_Nodo(n_particiones)

	if lista.primero == nil {
		lista.primero = nuevo
		lista.ultimo = nuevo
		lista.contador += 1
	} else {
		lista.ultimo.siguiente = nuevo
		lista.ultimo = lista.ultimo.siguiente
		lista.contador += 1
	}
}

func BuscarNumero(ppath string, pname string, lista *Lista) int {
	aux := lista.primero
	var retorno int = 1
	for aux != nil {
		if (ppath == aux.n_particiones.Path) && (pname == aux.n_particiones.Name) {
			return -1
		} else {
			if ppath == aux.n_particiones.Path {
				return int(aux.n_particiones.Num)
			} else if retorno <= int(aux.n_particiones.Num) {
				retorno++
			}
		}
		aux = aux.siguiente
	}
	return retorno
}

func BuscarLetra(ppath string, pname string, lista *Lista) int {
	var retorno int = 'a'
	aux := lista.primero
	for aux != nil {
		if (ppath == aux.n_particiones.Path) && (retorno == int(aux.n_particiones.Letra)) {
			retorno++
		}
		aux = aux.siguiente
	}
	return retorno
}

func mostrarLista(lista *Lista) {
	fmt.Println("---------------------------")
	fmt.Println("|   Lista de particiones  |")
	fmt.Println("---------------------------")
	aux := lista.primero
	for aux != nil {
		var aux_l byte = byte(aux.n_particiones.Letra)
		auxLetra := string([]byte{aux_l})
		auxNum := strconv.Itoa(aux.n_particiones.Num)
		fmt.Println("   ", aux.n_particiones.Name, "   ", "07"+auxNum+auxLetra)
		fmt.Println("---------------------------")
		aux = aux.siguiente
	}
}

func GetDireccion(id string, lista *Lista) string {
	aux := lista.primero
	for aux != nil {
		var tempID string = "07"
		var aux_l byte = byte(aux.n_particiones.Letra)
		auxLetra := string([]byte{aux_l})
		auxNum := strconv.Itoa(aux.n_particiones.Num)
		tempID += auxNum + auxLetra
		if id == tempID {
			return aux.n_particiones.Path
		}
		aux = aux.siguiente
	}
	return "null"
}
func getParticionMontada(id string, lista *Lista) Nparticiones {
	aux := lista.primero
	for aux != nil {
		var tempID string = "07"
		var aux_l byte = byte(aux.n_particiones.Letra)
		auxLetra := string([]byte{aux_l})
		auxNum := strconv.Itoa(aux.n_particiones.Num)
		tempID += auxNum + auxLetra
		if tempID == id {
			return *aux.n_particiones
		}
		aux = aux.siguiente
	}
	var temp *Nparticiones = New_Nparticiones("", "", 0, 0, 0)
	return *temp
}
