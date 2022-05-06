package comandos

import (
	"fmt"
	"strings"
)

func Recorrido_mkusr(comando string) {
	sinsalto := strings.TrimRight(comando, "\n")
	coman := strings.Split(sinsalto, " ")
	var bandera_error bool = false
	var bandera_usuario bool = false
	var bandera_pwd bool = false
	var bandera_grp bool = false
	var valor_usuario string = ""
	var valor_pwd string = ""
	var valor_grp string = ""
	for i := 0; i < len(coman); i++ {
		param := strings.Split(coman[i], "=")
		for j := 0; j < len(param); j++ {
			if param[j] == "-usuario" || param[j] == "-USUARIO" {
				bandera_usuario = true
				valor_usuario = param[j+1]
			} else if param[j] == "-pwd" || param[j] == "-PWD" {
				bandera_pwd = true
				valor_pwd = param[j+1]
			} else if param[j] == "-grp" || param[j] == "-GRP" {
				bandera_grp = true
				valor_grp = param[j+1]
			}
		}
	}
	if !bandera_usuario {
		bandera_error = true
		fmt.Println("Error: El parametro -usuario es obligatorio")
	}
	if !bandera_pwd {
		bandera_error = true
		fmt.Println("Error: El parametro -pwd es obligatorio")
	}
	if !bandera_grp {
		bandera_error = true
		fmt.Println("Error: El parametro -grp es obligatorio")
	}
	if !bandera_error {
		ejecutar_mkusr(valor_usuario, valor_pwd, valor_grp)
	}
}

func ejecutar_mkusr(pusuario string, ppwd string, pgrp string) {
	fmt.Println("Comando mkusr reconocido correctamente")
}
