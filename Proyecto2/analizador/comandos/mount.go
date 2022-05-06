package comandos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

func DoMount(comando string, lista *Lista) {
	sinsalto := strings.TrimRight(comando, "\n")
	coman := strings.Split(sinsalto, " ")
	lista_simple := lista
	var bandera_error bool = false
	var bandera_path bool = false
	var bandera_name bool = false
	var valor_path string = ""
	var valor_name string = ""
	for i := 0; i < len(coman); i++ {
		param := strings.Split(coman[i], "=")
		for j := 0; j < len(param); j++ {
			if param[j] == "-path" || param[j] == "-PATH" {
				bandera_path = true
				valor_path = param[j+1]
			} else if param[j] == "-name" || param[j] == "-NAME" {
				bandera_name = true
				valor_name = param[j+1]
			}
		}
	}
	if !bandera_path {
		bandera_error = true
		fmt.Println("Error: El parametro -path es obligatorio")
	}
	if !bandera_name {
		bandera_error = true
		fmt.Println("Error: El parametro -name es obligatorio")
	}
	if !bandera_error {
		buildMount(valor_path, valor_name, lista_simple)
	}
}
func buildMount(ppath string, pname string, lista *Lista) {
	var aux_path string
	var aux_nombre string
	comillaDer := strings.TrimRight(ppath, "\"")
	comillaIzq := strings.TrimLeft(comillaDer, "\"")
	aux_path = comillaIzq

	comillaDerN := strings.TrimRight(pname, "\"")
	comillaIzqN := strings.TrimLeft(comillaDerN, "\"")
	aux_nombre = comillaIzqN

	indexP := buscarParticion_P_E(ppath, pname)
	if indexP != -1 {
		mbr := MBR{}

		file, _ := os.OpenFile(aux_path, os.O_RDWR, 0777)
		defer file.Close()

		//nos posicionamos al inicio del archivo usando la funcion Seek
		file.Seek(0, 0)

		//obtenemor el size del MBR para empezar a leer desde ahi
		var sizeMbr int64 = int64(unsafe.Sizeof(mbr))

		file.Seek(0, 0)
		dataControl := leerBytes(file, int(sizeMbr))
		bufferControl := bytes.NewBuffer(dataControl)
		err := binary.Read(bufferControl, binary.BigEndian, &mbr)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		mbr.Mbr_partition[indexP].Part_status = '2'
		temp_type := mbr.Mbr_partition[indexP].Part_type
		/*temp_part := mbr.Mbr_partition[indexP]*/

		//nos posicionamos al inicio del archivo usando la funcion Seek
		file.Seek(0, 0)

		//Escribimos struct de mbr
		var bufferControlW bytes.Buffer
		binary.Write(&bufferControlW, binary.BigEndian, &mbr)
		escribirBytes(file, bufferControlW.Bytes())
		file.Close()
		//var lista *Lista = New_Lista()
		num := BuscarNumero(aux_path, aux_nombre, lista)
		//fmt.Println(num)
		if num == -1 {
			fmt.Println("Error: La particion ya esta montada")
		} else {
			letra := BuscarLetra(aux_path, aux_nombre, lista)
			var aux_l byte = byte(letra)
			auxLetra := string([]byte{aux_l})
			id := "07"
			id += strconv.Itoa(num) + auxLetra
			//fmt.Println(id)
			var n *Nparticiones = New_Nparticiones(aux_path, aux_nombre, aux_l, num, temp_type /*,temp_part*/)
			Insertar(n, lista)
			fmt.Println("Particion montada con exito")
			mostrarLista(lista)
		}
	} else {
		indexP := buscarParticion_L(aux_path, aux_nombre)
		if indexP != -1 {
			file, _ := os.OpenFile(aux_path, os.O_RDWR, 0777)
			defer file.Close()
			ebr := EBR{}
			file.Seek(int64(indexP), 0)
			var sizeEbr int64 = int64(unsafe.Sizeof(ebr))
			dataControl := leerBytes(file, int(sizeEbr))
			bufferControl := bytes.NewBuffer(dataControl)
			err := binary.Read(bufferControl, binary.BigEndian, &ebr)
			if err != nil {
				log.Fatal("binary.Read failed", err)
			}
			ebr.Part_status = '2'
			/*temp_part := ebr*/
			file.Seek(int64(indexP), 0)

			//Escribimos struct de mbr
			var bufferControlW bytes.Buffer
			binary.Write(&bufferControlW, binary.BigEndian, &ebr)
			escribirBytes(file, bufferControlW.Bytes())
			file.Close()

			letra := BuscarLetra(aux_path, aux_nombre, lista)
			if letra == -1 {
				fmt.Println("Error: la particion ya esta montada")
			} else {
				num := BuscarNumero(aux_path, aux_nombre, lista)
				var aux_l byte = byte(letra)
				auxLetra := string([]byte{aux_l})
				id := "07"
				id += strconv.Itoa(num) + auxLetra
				//fmt.Println(id)
				var n *Nparticiones = New_Nparticiones(aux_path, aux_nombre, aux_l, num, 'l' /*temp_part*/)
				Insertar(n, lista)
				fmt.Println("Particion montada con exito")
				mostrarLista(lista)
			}
		} else {
			fmt.Println("No se encontro la particion a montar")
		}
	}
}

func buscarParticion_P_E(ppath string, pname string) int {
	var aux_path string
	var aux_nombre string
	comillaDer := strings.TrimRight(ppath, "\"")
	comillaIzq := strings.TrimLeft(comillaDer, "\"")
	aux_path = comillaIzq

	comillaDerN := strings.TrimRight(pname, "\"")
	comillaIzqN := strings.TrimLeft(comillaDerN, "\"")
	aux_nombre = comillaIzqN

	mbr := MBR{}

	file, _ := os.OpenFile(aux_path, os.O_RDWR, 0777)
	defer file.Close()

	//nos posicionamos al inicio del archivo usando la funcion Seek
	file.Seek(0, 0)

	//obtenemor el size del MBR para empezar a leer desde ahi
	var sizeMbr int64 = int64(unsafe.Sizeof(mbr))

	file.Seek(0, 0)
	dataControl := leerBytes(file, int(sizeMbr))
	bufferControl := bytes.NewBuffer(dataControl)
	err := binary.Read(bufferControl, binary.BigEndian, &mbr)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	var byteName [16]byte

	for l := 0; l < 16; l++ {
		byteName[l] = '0'
	}

	for i := 0; i < len(aux_nombre); i++ {
		byteName[i] = aux_nombre[i]
	}

	for i := 0; i < 4; i++ {
		if mbr.Mbr_partition[i].Part_status != '1' {
			if mbr.Mbr_partition[i].Part_name == byteName {
				return i
			}
		}
	}
	return -1
}

func buscarParticion_L(ppath string, pname string) int {
	var aux_path string
	var aux_nombre string
	comillaDer := strings.TrimRight(ppath, "\"")
	comillaIzq := strings.TrimLeft(comillaDer, "\"")
	aux_path = comillaIzq

	comillaDerN := strings.TrimRight(pname, "\"")
	comillaIzqN := strings.TrimLeft(comillaDerN, "\"")
	aux_nombre = comillaIzqN

	mbr := MBR{}

	file, _ := os.OpenFile(aux_path, os.O_RDWR, 0777)
	defer file.Close()

	var extendida int = -1

	//nos posicionamos al inicio del archivo usando la funcion Seek
	file.Seek(0, 0)

	//obtenemor el size del MBR para empezar a leer desde ahi
	var sizeMbr int64 = int64(unsafe.Sizeof(mbr))

	file.Seek(0, 0)
	dataControl := leerBytes(file, int(sizeMbr))
	bufferControl := bytes.NewBuffer(dataControl)
	err := binary.Read(bufferControl, binary.BigEndian, &mbr)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	for i := 0; i < 4; i++ {
		if mbr.Mbr_partition[i].Part_type == 'E' {
			extendida = i
			break
		}
	}

	var byteName [16]byte

	for l := 0; l < 16; l++ {
		byteName[l] = '0'
	}

	for i := 0; i < len(aux_nombre); i++ {
		byteName[i] = aux_nombre[i]
	}

	if extendida != -1 {
		file.Seek(mbr.Mbr_partition[extendida].Part_start, 0)
		ebr := EBR{}

		file.Seek(mbr.Mbr_partition[extendida].Part_start, 0)
		var sizeEbr int64 = int64(unsafe.Sizeof(ebr))
		dataControl := leerBytes(file, int(sizeEbr))
		bufferControl := bytes.NewBuffer(dataControl)
		err := binary.Read(bufferControl, binary.BigEndian, &ebr)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		offset, _ := file.Seek(0, os.SEEK_CUR)

		for {
			if sizeEbr != 0 && (offset < (mbr.Mbr_partition[extendida].Part_size + mbr.Mbr_partition[extendida].Part_start)) {

				if ebr.Part_name == byteName {
					return (int(offset) - int(sizeEbr))
				}
			}
		}
	}
	file.Close()
	return -1
}
