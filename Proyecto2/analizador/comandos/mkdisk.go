package comandos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func DoMkdisk(params []string) {
	ParamSize := ""
	ParamFit := ""
	ParamUnit := ""
	ParamPath := ""

	// Recorriendo cada elemento del array
	for _, item := range params {
		fmt.Println(item)
		// variables
		intPosIgual := strings.Index(item, "=")
		intLenItem := len(item)

		// datos que me interesan
		// strParametro := strings.Replace(item[0:intPosIgual+1], "-", "", 5)
		strParametro := item[0 : intPosIgual+1]
		strParametro = strings.TrimLeft(strParametro, "–")
		strParametro = strings.TrimLeft(strParametro, "-")
		strParametro = strings.TrimRight(strParametro, "=")
		strParametro = strings.TrimSpace(strParametro)

		strDato := item[intPosIgual+1 : intLenItem]

		// Guardando informacion de cada parametro
		// y verificando que no esten repetidos
		if strings.ToUpper(strParametro) == "SIZE" {
			if ParamSize == "" {
				ParamSize = strDato
			} else {
				fmt.Println("Se ha repetido el parametro SIZE")
			}
		} else if strings.ToUpper(strParametro) == "FIT" {
			if ParamFit == "" {
				ParamFit = strDato
			} else {
				fmt.Println("Se ha repetido el parametro FIT")
			}
		} else if strings.ToUpper(strParametro) == "UNIT" {
			if ParamUnit == "" {
				ParamUnit = strDato
			} else {
				fmt.Println("Se ha repetido el parametro UNIT")
			}
		} else if strings.ToUpper(strParametro) == "PATH" {
			if ParamPath == "" {
				ParamPath = strDato
			} else {
				fmt.Println("Se ha repetido el parametro PATH")
			}
		}

		//fmt.Println(intPosIgual, "///", strParametro, "///", strDato)

	}
	// Parametros opcionales
	if ParamFit == "" {
		ParamFit = "FF"
	}
	if ParamUnit == "" {
		ParamUnit = "M"
	}

	// Validaciones que no falte ningun parametro
	if ParamPath == "" && ParamSize == "" {
		fmt.Println("Los parametros PATH Y SIZE son necesarios")
	} else if ParamPath == "" {

		fmt.Println("Los parametros PATH es necesario")

	} else if ParamSize == "" {
		fmt.Println("Los parametros SIZE es necesario")

	} else if strings.ToUpper(ParamFit) != "BF" && strings.ToUpper(ParamFit) != "FF" && strings.ToUpper(ParamFit) != "WF" {
		fmt.Println("Los valores para FIT son desconocidos")

	} else if strings.ToUpper(ParamUnit) != "K" && strings.ToUpper(ParamUnit) != "M" {
		fmt.Println("Los valores para UNIT son desconocidos")

	} else {
		// convirtiendo a int el parametro de size
		intSize, _ := strconv.Atoi(ParamSize)

		buildDisco(intSize, ParamFit, ParamUnit, ParamPath)
		fmt.Println("----------------------------------------")
		fmt.Println("fit ", ParamFit)
		fmt.Println("size ", ParamSize)
		fmt.Println("path ", ParamPath)
		fmt.Println("unit ", ParamUnit)
	}

}

func buildDisco(psize int, pfit string, punit string, ppath string) {

	comillaDer := strings.TrimRight(ppath, "\"")
	comillaIzq := strings.TrimLeft(comillaDer, "\"")
	cambio_ruta := comillaIzq

	crearDirectorio(cambio_ruta)

	file, _ := os.Create(cambio_ruta)
	defer file.Close()

	prueba := MBR{}

	var temporal int8 = 0
	s := &temporal
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s)

	if punit == "k" || punit == "K" {
		prueba.Mbr_tamano = int64(psize * 1024)

		for i := 0; i < psize*1024; i++ {
			escribirBytes(file, binario.Bytes())
		}
	} else if punit == "m" || punit == "M" {
		prueba.Mbr_tamano = int64(psize * 1024 * 1024)

		for i := 0; i < psize*1024*1024; i++ {
			escribirBytes(file, binario.Bytes())
		}
	}

	tiempo := time.Now()
	time := tiempo.String()

	for t := 0; t < 19; t++ {
		prueba.Mbr_fecha_creacion[t] = time[t]
	}

	numeroRandom := rand.Intn(100-1) + 1
	prueba.Mbr_dsk_signature = int64(numeroRandom)

	var auxfit byte = 0

	if pfit == "bf" || pfit == "BF" {
		auxfit = 'B'
	} else if pfit == "ff" || pfit == "FF" {
		auxfit = 'F'
	} else if pfit == "wf" || pfit == "WF" {
		auxfit = 'W'
	}

	prueba.Dsk_fit = auxfit

	for p := 0; p < 4; p++ {
		prueba.Mbr_partition[p].Part_status = '0'
		prueba.Mbr_partition[p].Part_type = '0'
		prueba.Mbr_partition[p].Part_fit = '0'
		prueba.Mbr_partition[p].Part_size = 0
		prueba.Mbr_partition[p].Part_start = -1
		for n := 0; n < 16; n++ {
			prueba.Mbr_partition[p].Part_name[n] = '0'
		}
	}

	//nos posicionamos al inicio del archivo usando la funcion Seek
	file.Seek(0, 0)

	//Escribimos struct de mbr
	var bufferControl bytes.Buffer
	binary.Write(&bufferControl, binary.BigEndian, &prueba)
	escribirBytes(file, bufferControl.Bytes())

	//movemos el puntero a donde ira nuestra primera estructura
	file.Seek(int64(unsafe.Sizeof(prueba)), 0)

	fmt.Println("Disco creado correctamente")
	fmt.Println("Se agrego el mbr de manera correcta")

}

func crearDirectorio(path string) {
	directorio := obtener_path(path)
	if _, err := os.Stat(directorio); os.IsNotExist(err) {
		err = os.MkdirAll(directorio, 0777)
		if err != nil {
			panic(err)
		}
	}
}

func obtener_path(path string) string {
	var aux_path int
	var aux_ruta string

	for i := len(path) - 1; i >= 0; i-- {
		aux_path++
		if string(path[i]) == "/" {
			break
		}
	}

	for i := 0; i < ((len(path)) - (aux_path - 1)); i++ {
		aux_ruta += string(path[i])
	}

	return aux_ruta
}

func escribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}
