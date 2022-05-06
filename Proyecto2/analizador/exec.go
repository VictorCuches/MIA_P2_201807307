package analizador

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func DoExec(texto string) {
	var valor_path string
	sinsalto := strings.TrimRight(texto, "\n")
	coman := strings.Split(sinsalto, " ")
	for i := 0; i < len(coman); i++ {
		param := strings.Split(coman[i], "=")
		for j := 0; j < len(param); j++ {
			if param[j] == "-path" || param[j] == "-PATH" {
				valor_path = param[j+1]
			}
		}
	}
	comillaDer := strings.TrimRight(valor_path, "\"")
	comillaIzq := strings.TrimLeft(comillaDer, "\"")
	ruta := comillaIzq
	Archivo, err := os.Open(ruta)
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(Archivo)
	fileScanner.Split(bufio.ScanLines)
	var lineas []string
	for fileScanner.Scan() {
		if fileScanner.Text() != "" {
			lineas = append(lineas, fileScanner.Text())
		}
	}
	Archivo.Close()
	for _, line := range lineas {
		palabras := strings.Split(line, " ")
		pInicio := palabras[0]
		for i := 0; i < len(pInicio); i++ {
			if pInicio[0] != 35 {
				arrLine := strings.Split(line, " ")
				Comandos_indentify(arrLine[0], arrLine, line)
				break
			}
		}
	}
}
