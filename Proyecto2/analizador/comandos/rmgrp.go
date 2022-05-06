package comandos

import (
	"fmt"
	"strings"
)

func Recorrido_rmgrp(comando string) {
	sinsalto := strings.TrimRight(comando, "\n")
	coman := strings.Split(sinsalto, " ")
	var bandera_error bool = false
	var bandera_name bool = false
	var valor_name string = ""
	for i := 0; i < len(coman); i++ {
		param := strings.Split(coman[i], "=")
		for j := 0; j < len(param); j++ {
			if param[j] == "-name" || param[j] == "-NAME" {
				bandera_name = true
				valor_name = param[j+1]
			}
		}
	}
	if !bandera_name {
		bandera_error = true
		fmt.Println("Error: El parametro -name es obligatorio")
	}
	if !bandera_error {
		ejecutar_rmgrp(valor_name)
	}
}

func ejecutar_rmgrp(pname string) {
	fmt.Println("Comando rmgrp reconocido correctamente")
}
