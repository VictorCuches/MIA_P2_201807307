package comandos

import (
	"fmt"
	"os"
	"strings"
)

func Recorrido_mkdir(comando string, lista *Lista) {
	sinsalto := strings.TrimRight(comando, "\n")
	coman := strings.Split(sinsalto, " ")
	lista_simple := lista
	var bandera_error bool = false
	var bandera_path bool = false
	var valor_path string = ""
	var valor_p bool = false
	for i := 0; i < len(coman); i++ {
		param := strings.Split(coman[i], "=")
		for j := 0; j < len(param); j++ {
			if param[j] == "-path" || param[j] == "-PATH" {
				bandera_path = true
				valor_path = param[j+1]
			} else if param[j] == "-p" || param[j] == "-P" {
				valor_p = true
			}
		}
	}
	if !bandera_path {
		bandera_error = true
		fmt.Println("Error: El parametro -path es obligatorio")
	}
	if !bandera_error {
		ejecutar_mkdir(valor_path, valor_p, lista_simple)
	}
}

func ejecutar_mkdir(ppath string, valor_p bool, lista *Lista) {
	comillaDer := strings.TrimRight(ppath, "\"")
	comillaIzq := strings.TrimLeft(comillaDer, "\"")
	auxPath := comillaIzq

	file, _ := os.OpenFile(auxPath, os.O_RDWR, 0777)
	defer file.Close()

	fmt.Println("La carpeta se ha creado correctamente")
}
