package comandos

import (
	"fmt"
	"strconv"
	"strings"
)

func DoMkfile(comando string, lista *ListaF) {
	sinsalto := strings.TrimRight(comando, "\n")
	coman := strings.Split(sinsalto, " ")
	listaSimple := lista
	var bandera_error bool = false
	var bandera_path bool = false
	var valor_path string = ""
	var valor_r bool = false
	var valor_size int = 0
	var valor_count string = "null"
	for i := 0; i < len(coman); i++ {
		param := strings.Split(coman[i], "=")
		for j := 0; j < len(param); j++ {
			if param[j] == "-path" || param[j] == "-PATH" {
				bandera_path = true
				valor_path = param[j+1]
			} else if param[j] == "-r" || param[j] == "-R" {
				valor_r = true
			} else if param[j] == "-size" || param[j] == "-SIZE" {
				sv, _ := strconv.Atoi(param[j+1])
				valor_size = sv
				if valor_size <= 0 {
					bandera_error = true
					fmt.Println("Los valores de size no pueden ser menores o iguales a cero")
				}
			} else if param[j] == "-count" || param[j] == "-COUNT" {
				valor_count = param[j+1]
			}
		}
	}
	if !bandera_path {
		bandera_error = true
		fmt.Println("Error: El parametro -path es obligatorio")
	}
	if !bandera_error {
		ejecutar_mkfile(valor_path, valor_r, valor_size, valor_count, listaSimple)
	}
}

func ejecutar_mkfile(ppath string, valor_r bool, psize int, pcount string, lista *ListaF) {

	comillaDer := strings.TrimRight(ppath, "\"")
	comillaIzq := strings.TrimLeft(comillaDer, "\"")
	auxPath := comillaIzq

	var nfile *NFiles = New_NFiles(auxPath)
	InsertarF(nfile, lista)

	fmt.Println("Archivo creado correctamente")
}
