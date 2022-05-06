package comandos

import (
	"fmt"
	"strings"
)

func DoLogin(comando string, lista *Lista, listaLogin *ListaL) {
	sinsalto := strings.TrimRight(comando, "\n")
	coman := strings.Split(sinsalto, " ")
	lista_simple := lista
	lista_simpleL := listaLogin
	var bandera_error bool = false
	var bandera_usuario bool = false
	var bandera_password bool = false
	var bandera_id bool = false
	var valor_usuario string = ""
	var valor_password string = ""
	var valor_id string = ""
	for i := 0; i < len(coman); i++ {
		param := strings.Split(coman[i], "=")
		for j := 0; j < len(param); j++ {
			if param[j] == "-usuario" || param[j] == "-USUARIO" {
				valor_usuario = param[j+1]
				bandera_usuario = true
			} else if param[j] == "-password" || param[j] == "-PASSWORD" {
				valor_password = param[j+1]
				bandera_password = true
			} else if param[j] == "-id" || param[j] == "-ID" {
				valor_id = param[j+1]
				bandera_id = true
			}
		}
	}
	if !bandera_usuario {
		bandera_error = true
		fmt.Println("Error: El parametro -usuario es obligatorio")
	}
	if !bandera_password {
		bandera_error = true
		fmt.Println("Error: El parametro -password es obligatorio")
	}
	if !bandera_id {
		bandera_error = true
		fmt.Println("Error: El parametro -id es obligatorio")
	}
	if !bandera_error {
		crear_usuario(valor_usuario, valor_password, valor_id, lista_simple, lista_simpleL)
	}
}

func crear_usuario(pusuario string, ppassword string, pid string, lista *Lista, listaLogin *ListaL) {

	encontrado := getParticionMontada(pid, lista)
	if len(encontrado.Path) > 0 {
		var nl *NLogin = New_NLogin(1, pusuario, pid, "1", 1)
		InsertarL(nl, listaLogin)
		fmt.Println("Se ha iniciado la sesion con el usuario:", pusuario)
	} else {
		fmt.Println("No se encontro la particion")
	}
}
