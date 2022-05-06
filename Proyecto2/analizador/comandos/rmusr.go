package comandos

import (
	"fmt"
	"strings"
)

func Recorrido_rmusr(comando string) {
	sinsalto := strings.TrimRight(comando, "\n")
	coman := strings.Split(sinsalto, " ")
	var bandera_error bool = false
	var bandera_usuario bool = false
	var valor_usuario string = ""
	for i := 0; i < len(coman); i++ {
		param := strings.Split(coman[i], "=")
		for j := 0; j < len(param); j++ {
			if param[j] == "-usuario" || param[j] == "-USUARIO" {
				bandera_usuario = true
				valor_usuario = param[j+1]
			}
		}
	}
	if !bandera_usuario {
		bandera_error = true
		fmt.Println("Error: El parametro -name es obligatorio")
	}
	if !bandera_error {
		ejecutar_rmusr(valor_usuario)
	}
}

func ejecutar_rmusr(pusuario string) {
	fmt.Println("Comando rmusr reconocido correctamente")
}
