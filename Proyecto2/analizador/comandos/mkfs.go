package comandos

import (
	"fmt"
	"os"
	"strings"
)

func DoMkfs(comando string, lista *Lista) {
	sinsalto := strings.TrimRight(comando, "\n")
	coman := strings.Split(sinsalto, " ")
	lista_simple := lista
	var bandera_error bool = false
	var bandera_id bool = false
	var valor_id string = ""
	var valor_tipo string = "full"
	var valor_fs string = "2fs"
	for i := 0; i < len(coman); i++ {
		param := strings.Split(coman[i], "=")
		for j := 0; j < len(param); j++ {
			if param[j] == "-id" || param[j] == "-ID" {
				bandera_id = true
				valor_id = param[j+1]
			} else if param[j] == "-type" || param[j] == "-TYPE" {
				valor_tipo = param[j+1]
			}
		}
	}
	if !bandera_id {
		bandera_error = true
		fmt.Println("Error: El parametro -id es obligatorio")
	}
	if !bandera_error {
		ejecutar_mkfs(valor_id, valor_tipo, valor_fs, lista_simple)
	}
}

func ejecutar_mkfs(pid string, ptipo string, pfs string, lista *Lista) {

	encontrado := getParticionMontada(pid, lista)
	if len(encontrado.Path) > 0 {
		file, _ := os.OpenFile(encontrado.Path, os.O_RDWR, 0777)
		defer file.Close()
		if encontrado.Part_type == 'P' {
			fmt.Println("MKFS con exito")
		} else if encontrado.Part_type == 'E' {
			fmt.Println("MKFS con exito")
		} else {
			fmt.Println("MKFS con exito")
		}
	} else {
		fmt.Println("No se encontro la particion")
	}
}
