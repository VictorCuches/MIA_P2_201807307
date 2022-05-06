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

func DoFdisk(comando string) {
	sinsalto := strings.TrimRight(comando, "\n")
	coman := strings.Split(sinsalto, " ")
	var bandera_error bool = false
	var bandera_size bool = false
	var bandera_path bool = false
	var bandera_name bool = false
	var valor_size int = 0
	var valor_unit string = "k"
	var valor_path string = ""
	var valor_type string = "p"
	var valor_fit string = "wf"
	var valor_name string = ""
	for i := 0; i < len(coman); i++ {
		param := strings.Split(coman[i], "=")
		for j := 0; j < len(param); j++ {
			if param[j] == "-size" || param[j] == "-SIZE" {
				sv, _ := strconv.Atoi(param[j+1])
				valor_size = sv
				if valor_size <= 0 {
					bandera_error = true
					fmt.Println("Los valores de size no pueden ser menores o iguales a cero")
				} else {
					bandera_size = true
				}
			} else if param[j] == "-unit" || param[j] == "-UNIT" {
				valor_unit = param[j+1]
			} else if param[j] == "-path" || param[j] == "-PATH" {
				bandera_path = true
				valor_path = param[j+1]
			} else if param[j] == "-type" || param[j] == "-TYPE" {
				valor_type = param[j+1]
			} else if param[j] == "-fit" || param[j] == "-FIT" {
				valor_fit = param[j+1]
			} else if param[j] == "-name" || param[j] == "-NAME" {
				bandera_name = true
				valor_name = param[j+1]
			}
		}
	}
	if !bandera_size {
		bandera_error = true
		fmt.Println("Error: El parametro -size es obligatorio")
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
		crear_particiones(valor_size, valor_unit, valor_path, valor_type, valor_fit, valor_name)
	}
}

func crear_particiones(psize int, punit string, ppath string, ptype string, pfit string, pname string) {
	if ptype == "p" || ptype == "P" {
		if archivoExiste(ppath) {
			crearParticionPrimaria(psize, punit, ppath, pfit, pname, "principal")
		} else {
			fmt.Println("Error: No se encontro el disco")
		}
	} else if ptype == "e" || ptype == "E" {
		if archivoExiste(ppath) {
			crearParticionExtendida(psize, punit, ppath, pfit, pname, "principal")
		} else {
			fmt.Println("Error: No se encontro el disco")
		}
	} else if ptype == "l" || ptype == "L" {
		if archivoExiste(ppath) {
			crearParticionLogica(psize, punit, ppath, pfit, pname, "principal")
		} else {
			fmt.Println("Error: No se encontro el disco")
		}
	}
}

func crearParticionPrimaria(psize int, punit string, ppath string, pfit string, pnombre string, archivo string) {

	var auxFit byte = 0
	//var auxUnit byte = 0
	var auxPath string
	var pname string
	var size_bytes int = 1024

	comillaDer := strings.TrimRight(ppath, "\"")
	comillaIzq := strings.TrimLeft(comillaDer, "\"")
	auxPath = comillaIzq

	comillaDerN := strings.TrimRight(pnombre, "\"")
	comillaIzqN := strings.TrimLeft(comillaDerN, "\"")
	pname = comillaIzqN

	var temporal int8 = '1'
	s := &temporal
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s)

	if pfit == "bf" || pfit == "BF" {
		auxFit = 'B'
	} else if pfit == "ff" || pfit == "FF" {
		auxFit = 'F'
	} else if pfit == "wf" || pfit == "WF" {
		auxFit = 'W'
	}

	if punit == "b" || punit == "B" {
		//auxUnit = 'b'
		size_bytes = psize
	} else if punit == "k" || punit == "K" {
		//auxUnit = 'k'
		size_bytes = psize * 1024
	} else if punit == "m" || punit == "M" {
		//auxUnit = 'm'
		size_bytes = psize * 1024 * 1024
	}

	mbr := MBR{}

	file, _ := os.OpenFile(auxPath, os.O_RDWR, 0777)
	defer file.Close()

	var flagParticion bool = false //Flag para ver si hay una particion disponible
	var numParticion int = 0       //Que numero de particion es

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
		if mbr.Mbr_partition[i].Part_start == -1 || mbr.Mbr_partition[i].Part_status == '1' && mbr.Mbr_partition[i].Part_size >= int64(size_bytes) {
			flagParticion = true
			numParticion = i
			break
		}
	}

	if flagParticion {
		//verificar el espacio libre del disco
		var espacioUsado int = 0
		for i := 0; i < 4; i++ {
			if mbr.Mbr_partition[i].Part_status != '1' {
				espacioUsado += int(mbr.Mbr_partition[i].Part_size)
			}
		}

		if archivo == "principal" {
			fmt.Println("Espacio Disponible: ", (mbr.Mbr_tamano - int64(espacioUsado)), " Bytes")
			fmt.Println("Espacio Requerido: ", size_bytes, " Bytes")
		}

		//verificar que haya espacio suficiente para crear la particion
		if (mbr.Mbr_tamano - int64(espacioUsado)) >= int64(size_bytes) {
			if !existeParticion(auxPath, pname) {
				if mbr.Dsk_fit == 'F' { //FIRST FIT
					mbr.Mbr_partition[numParticion].Part_type = 'P'
					mbr.Mbr_partition[numParticion].Part_fit = auxFit
					//start
					if numParticion == 0 {
						mbr.Mbr_partition[numParticion].Part_start = sizeMbr
					} else {
						mbr.Mbr_partition[numParticion].Part_start = mbr.Mbr_partition[numParticion-1].Part_start + mbr.Mbr_partition[numParticion-1].Part_size
					}
					mbr.Mbr_partition[numParticion].Part_size = int64(size_bytes)
					mbr.Mbr_partition[numParticion].Part_status = '0'
					for n := 0; n < len(pname); n++ {
						mbr.Mbr_partition[numParticion].Part_name[n] = pname[n]
					}
					//se guarde de nuevo el MBR
					//nos posicionamos al inicio del archivo usando la funcion Seek
					file.Seek(0, 0)

					//Escribimos struct de mbr
					var bufferControl bytes.Buffer
					binary.Write(&bufferControl, binary.BigEndian, &mbr)
					escribirBytes(file, bufferControl.Bytes())

					//se guardan los bytes de la particion
					file.Seek(mbr.Mbr_partition[numParticion].Part_start, 0)
					for i := 0; i < size_bytes; i++ {
						escribirBytes(file, binario.Bytes())
					}

					if archivo == "principal" {
						fmt.Println("Particion primaria creada con exito")
					}
				} else if mbr.Dsk_fit == 'B' { //BEST FIT
					var bestIndex int = numParticion
					for i := 0; i < 4; i++ {
						if mbr.Mbr_partition[i].Part_start == -1 || (mbr.Mbr_partition[i].Part_status == '1' && mbr.Mbr_partition[i].Part_size >= int64(size_bytes)) {
							if i != numParticion {
								if mbr.Mbr_partition[bestIndex].Part_size > mbr.Mbr_partition[i].Part_size {
									bestIndex = i
									break
								}
							}
						}
					}
					mbr.Mbr_partition[bestIndex].Part_type = 'P'
					mbr.Mbr_partition[bestIndex].Part_fit = auxFit
					//start
					if bestIndex == 0 {
						mbr.Mbr_partition[bestIndex].Part_start = sizeMbr
					} else {
						mbr.Mbr_partition[bestIndex].Part_start = mbr.Mbr_partition[bestIndex-1].Part_start + mbr.Mbr_partition[bestIndex-1].Part_size
					}
					mbr.Mbr_partition[bestIndex].Part_size = int64(size_bytes)
					mbr.Mbr_partition[bestIndex].Part_status = '0'
					for n := 0; n < len(pname); n++ {
						mbr.Mbr_partition[bestIndex].Part_name[n] = pname[n]
					}
					//se guarda de nuevo el MBR
					//nos posicionamos al inicio del archivo usando la funcion Seek
					file.Seek(0, 0)

					//Escribimos struct de mbr
					var bufferControl bytes.Buffer
					binary.Write(&bufferControl, binary.BigEndian, &mbr)
					escribirBytes(file, bufferControl.Bytes())

					//se guardan los bytes de la particion
					file.Seek(mbr.Mbr_partition[bestIndex].Part_start, 0)
					for i := 0; i < size_bytes; i++ {
						escribirBytes(file, binario.Bytes())
					}

					if archivo == "principal" {
						fmt.Println("Particion primaria creada con exito")
					}
				} else if mbr.Dsk_fit == 'W' { //WORST FIT
					var worstIndex int = numParticion
					for i := 0; i < 4; i++ {
						if mbr.Mbr_partition[i].Part_start == -1 || (mbr.Mbr_partition[i].Part_status == '1' && mbr.Mbr_partition[i].Part_size >= int64(size_bytes)) {
							if i != numParticion {
								if mbr.Mbr_partition[worstIndex].Part_size < mbr.Mbr_partition[i].Part_size {
									worstIndex = i
									break
								}
							}
						}
					}
					mbr.Mbr_partition[worstIndex].Part_type = 'P'
					mbr.Mbr_partition[worstIndex].Part_fit = auxFit
					//start
					if worstIndex == 0 {
						mbr.Mbr_partition[worstIndex].Part_start = sizeMbr
					} else {
						mbr.Mbr_partition[worstIndex].Part_start = mbr.Mbr_partition[worstIndex-1].Part_start + mbr.Mbr_partition[worstIndex-1].Part_size
					}
					mbr.Mbr_partition[worstIndex].Part_size = int64(size_bytes)
					mbr.Mbr_partition[worstIndex].Part_status = '0'
					for n := 0; n < len(pname); n++ {
						mbr.Mbr_partition[worstIndex].Part_name[n] = pname[n]
					}
					//se guarda de nuevo el MBR
					//nos posicionamos al inicio del archivo usando la funcion Seek
					file.Seek(0, 0)

					//Escribimos struct de mbr
					var bufferControl bytes.Buffer
					binary.Write(&bufferControl, binary.BigEndian, &mbr)
					escribirBytes(file, bufferControl.Bytes())

					//se guardan los bytes de la particion
					file.Seek(mbr.Mbr_partition[worstIndex].Part_start, 0)
					for i := 0; i < size_bytes; i++ {
						escribirBytes(file, binario.Bytes())
					}

					if archivo == "principal" {
						fmt.Println("Particion primaria creada con exito")
					}
				}
			} else {
				fmt.Println("Error: ya existe una particion con ese nombre")
			}
		} else {
			fmt.Println("Error: la particion a crear excede el espacio libre")
		}
	} else {
		fmt.Println("Error: Ya existen 4 particiones, no se puede crear otra")
		fmt.Println("Elimine alguna para poder crear una")
	}
}

func crearParticionExtendida(psize int, punit string, ppath string, pfit string, pnombre string, archivo string) {

	var auxFit byte = 0
	//var auxUnit byte = 0
	var auxPath string
	var pname string
	var size_bytes int = 1024

	comillaDer := strings.TrimRight(ppath, "\"")
	comillaIzq := strings.TrimLeft(comillaDer, "\"")
	auxPath = comillaIzq

	comillaDerN := strings.TrimRight(pnombre, "\"")
	comillaIzqN := strings.TrimLeft(comillaDerN, "\"")
	pname = comillaIzqN

	var temporal int8 = '1'
	s := &temporal
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s)

	if pfit == "bf" || pfit == "BF" {
		auxFit = 'B'
	} else if pfit == "ff" || pfit == "FF" {
		auxFit = 'F'
	} else if pfit == "wf" || pfit == "WF" {
		auxFit = 'W'
	}

	if punit == "b" || punit == "B" {
		//auxUnit = 'b'
		size_bytes = psize
	} else if punit == "k" || punit == "K" {
		//auxUnit = 'k'
		size_bytes = psize * 1024
	} else if punit == "m" || punit == "M" {
		//auxUnit = 'm'
		size_bytes = psize * 1024 * 1024
	}

	mbr := MBR{}

	file, _ := os.OpenFile(auxPath, os.O_RDWR, 0777)
	defer file.Close()

	var flagParticion bool = false //Flag para ver si hay una particion disponible
	var flagExtendida bool = false //Flag para ver si hay una particion extendida
	var numParticion int = 0       //Que numero de particion es

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
			flagExtendida = true
			break
		}
	}

	if !flagExtendida {
		//Verificar si existe una particion disponible
		for i := 0; i < 4; i++ {
			if mbr.Mbr_partition[i].Part_start == -1 || (mbr.Mbr_partition[i].Part_status == '1' && mbr.Mbr_partition[i].Part_size >= int64(size_bytes)) {
				flagParticion = true
				numParticion = i
				break
			}
		}
		if flagParticion {
			//Verificar el espacio libre del disco
			var espacioUsado = 0
			for i := 0; i < 4; i++ {
				if mbr.Mbr_partition[i].Part_status != '1' {
					espacioUsado += int(mbr.Mbr_partition[i].Part_size)
				}
			}
			if archivo == "principal" {
				fmt.Println("Espacio Disponible: ", (mbr.Mbr_tamano - int64(espacioUsado)), " Bytes")
				fmt.Println("Espacio Requerido: ", size_bytes, " Bytes")
			}
			//Verificar que haya espacio suficiente para crear la particion
			if (mbr.Mbr_tamano - int64(espacioUsado)) >= int64(size_bytes) {
				if !(existeParticion(auxPath, pname)) {
					if mbr.Dsk_fit == 'F' { //FIRST FIT
						mbr.Mbr_partition[numParticion].Part_type = 'E'
						mbr.Mbr_partition[numParticion].Part_fit = auxFit
						//start
						if numParticion == 0 {
							mbr.Mbr_partition[numParticion].Part_start = sizeMbr
						} else {
							mbr.Mbr_partition[numParticion].Part_start = mbr.Mbr_partition[numParticion-1].Part_start + mbr.Mbr_partition[numParticion-1].Part_size
						}
						mbr.Mbr_partition[numParticion].Part_size = int64(size_bytes)
						mbr.Mbr_partition[numParticion].Part_status = '0'
						for n := 0; n < len(pname); n++ {
							mbr.Mbr_partition[numParticion].Part_name[n] = pname[n]
						}
						//se guarda de nuevo el MBR
						//nos posicionamos al inicio del archivo usando la funcion Seek
						file.Seek(0, 0)

						//Escribimos struct de mbr
						var bufferControl bytes.Buffer
						binary.Write(&bufferControl, binary.BigEndian, &mbr)
						escribirBytes(file, bufferControl.Bytes())

						//se guardan los bytes de la particion extendida
						file.Seek(mbr.Mbr_partition[numParticion].Part_start, 0)

						ebr := EBR{}

						ebr.Part_fit = auxFit
						ebr.Part_status = '0'
						ebr.Part_start = mbr.Mbr_partition[numParticion].Part_start
						ebr.Part_size = 0
						ebr.Part_next = -1
						for n := 0; n < 16; n++ {
							ebr.Part_name[n] = '0'
						}

						//obtenemor el size del MBR para empezar a leer desde ahi
						var sizeEbr int64 = int64(unsafe.Sizeof(ebr))

						//Escribimos struct de ebr
						var bufferControlEBR bytes.Buffer
						binary.Write(&bufferControlEBR, binary.BigEndian, &ebr)
						escribirBytes(file, bufferControlEBR.Bytes())

						//se guardan los bytes de la particion
						for i := 0; i < size_bytes-int(sizeEbr); i++ {
							escribirBytes(file, binario.Bytes())
						}

						if archivo == "principal" {
							fmt.Println("Particion extendida creada con exito")
						}
					} else if mbr.Dsk_fit == 'B' { //BEST FIT
						var bestIndex int = numParticion
						for i := 0; i < 4; i++ {
							if mbr.Mbr_partition[i].Part_start == -1 || (mbr.Mbr_partition[i].Part_status == '1' && mbr.Mbr_partition[i].Part_size >= int64(size_bytes)) {
								if i != numParticion {
									if mbr.Mbr_partition[bestIndex].Part_size > mbr.Mbr_partition[i].Part_size {
										bestIndex = i
										break
									}
								}
							}
						}
						mbr.Mbr_partition[bestIndex].Part_type = 'E'
						mbr.Mbr_partition[bestIndex].Part_fit = auxFit
						//start
						if bestIndex == 0 {
							mbr.Mbr_partition[bestIndex].Part_start = sizeMbr
						} else {
							mbr.Mbr_partition[bestIndex].Part_start = mbr.Mbr_partition[bestIndex-1].Part_start + mbr.Mbr_partition[bestIndex-1].Part_size
						}
						mbr.Mbr_partition[bestIndex].Part_size = int64(size_bytes)
						mbr.Mbr_partition[bestIndex].Part_status = '0'
						for n := 0; n < len(pname); n++ {
							mbr.Mbr_partition[bestIndex].Part_name[n] = pname[n]
						}
						//se guarda de nuevo el MBR
						//nos posicionamos al inicio del archivo usando la funcion Seek
						file.Seek(0, 0)

						//Escribimos struct de mbr
						var bufferControl bytes.Buffer
						binary.Write(&bufferControl, binary.BigEndian, &mbr)
						escribirBytes(file, bufferControl.Bytes())

						//se guardan los bytes de la particion extendida
						file.Seek(mbr.Mbr_partition[bestIndex].Part_start, 0)

						ebr := EBR{}

						ebr.Part_fit = auxFit
						ebr.Part_status = '0'
						ebr.Part_start = mbr.Mbr_partition[bestIndex].Part_start
						ebr.Part_size = 0
						ebr.Part_next = -1
						for n := 0; n < 16; n++ {
							ebr.Part_name[n] = '0'
						}

						//obtenemor el size del MBR para empezar a leer desde ahi
						var sizeEbr int64 = int64(unsafe.Sizeof(ebr))

						//Escribimos struct de ebr
						var bufferControlEBR bytes.Buffer
						binary.Write(&bufferControlEBR, binary.BigEndian, &ebr)
						escribirBytes(file, bufferControlEBR.Bytes())

						//se guardan los bytes de la particion
						for i := 0; i < size_bytes-int(sizeEbr); i++ {
							escribirBytes(file, binario.Bytes())
						}

						if archivo == "principal" {
							fmt.Println("Particion extendida creada con exito")
						}
					} else if mbr.Dsk_fit == 'W' { //WORST FIT
						var worstIndex int = numParticion
						for i := 0; i < 4; i++ {
							if mbr.Mbr_partition[i].Part_start == -1 || (mbr.Mbr_partition[i].Part_status == '1' && mbr.Mbr_partition[i].Part_size >= int64(size_bytes)) {
								if i != numParticion {
									if mbr.Mbr_partition[worstIndex].Part_size < mbr.Mbr_partition[i].Part_size {
										worstIndex = i
										break
									}
								}
							}
						}
						mbr.Mbr_partition[worstIndex].Part_type = 'E'
						mbr.Mbr_partition[worstIndex].Part_fit = auxFit
						//start
						if worstIndex == 0 {
							mbr.Mbr_partition[worstIndex].Part_start = sizeMbr
						} else {
							mbr.Mbr_partition[worstIndex].Part_start = mbr.Mbr_partition[worstIndex-1].Part_start + mbr.Mbr_partition[worstIndex-1].Part_size
						}
						mbr.Mbr_partition[worstIndex].Part_size = int64(size_bytes)
						mbr.Mbr_partition[worstIndex].Part_status = '0'
						for n := 0; n < len(pname); n++ {
							mbr.Mbr_partition[worstIndex].Part_name[n] = pname[n]
						}

						//se guarda de nuevo el MBR
						//nos posicionamos al inicio del archivo usando la funcion Seek
						file.Seek(0, 0)

						//Escribimos struct de mbr
						var bufferControl bytes.Buffer
						binary.Write(&bufferControl, binary.BigEndian, &mbr)
						escribirBytes(file, bufferControl.Bytes())

						//se guardan los bytes de la particion extendida
						file.Seek(mbr.Mbr_partition[worstIndex].Part_start, 0)

						ebr := EBR{}

						ebr.Part_fit = auxFit
						ebr.Part_status = '0'
						ebr.Part_start = mbr.Mbr_partition[worstIndex].Part_start
						ebr.Part_size = 0
						ebr.Part_next = -1
						for n := 0; n < 16; n++ {
							ebr.Part_name[n] = '0'
						}

						//obtenemor el size del MBR para empezar a leer desde ahi
						var sizeEbr int64 = int64(unsafe.Sizeof(ebr))

						//Escribimos struct de ebr
						var bufferControlEBR bytes.Buffer
						binary.Write(&bufferControlEBR, binary.BigEndian, &ebr)
						escribirBytes(file, bufferControlEBR.Bytes())

						//se guardan los bytes de la particion
						for i := 0; i < size_bytes-int(sizeEbr); i++ {
							escribirBytes(file, binario.Bytes())
						}

						if archivo == "principal" {
							fmt.Println("Particion extendida creada con exito")
						}
					}
				} else {
					fmt.Println("Error: ya existe una particion con ese nombre")
				}
			} else {
				fmt.Println("Error: la particion a crear excede el espacio libre")
			}
		} else {
			fmt.Println("Error: Ya existen 4 particiones, no se puede crear otra")
			fmt.Println("Elimine alguna para poder crear una")
		}
	} else {
		fmt.Println("Error: ya existe una particion extendida en este disco")
	}
}

func crearParticionLogica(psize int, punit string, ppath string, pfit string, pnombre string, archivo string) {

	var auxFit byte = 0
	//var auxUnit byte = 0
	var auxPath string
	var size_bytes int = 1024
	var pname string

	comillaDer := strings.TrimRight(ppath, "\"")
	comillaIzq := strings.TrimLeft(comillaDer, "\"")
	auxPath = comillaIzq

	comillaDerN := strings.TrimRight(pnombre, "\"")
	comillaIzqN := strings.TrimLeft(comillaDerN, "\"")
	pname = comillaIzqN

	var temporal int8 = '1'
	s := &temporal
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s)

	if pfit == "bf" || pfit == "BF" {
		auxFit = 'B'
	} else if pfit == "ff" || pfit == "FF" {
		auxFit = 'F'
	} else if pfit == "wf" || pfit == "WF" {
		auxFit = 'W'
	}

	if punit == "b" || punit == "B" {
		//auxUnit = 'b'
		size_bytes = psize
	} else if punit == "k" || punit == "K" {
		//auxUnit = 'k'
		size_bytes = psize * 1024
	} else if punit == "m" || punit == "M" {
		//auxUnit = 'm'
		size_bytes = psize * 1024 * 1024
	}

	mbr := MBR{}

	file, _ := os.OpenFile(auxPath, os.O_RDWR, 0777)
	defer file.Close()

	var numExtendida int = -1

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
	//Verificar si existe una particion extendida
	for i := 0; i < 4; i++ {
		if mbr.Mbr_partition[i].Part_type == 'E' {
			numExtendida = i
			break
		}
	}

	if !existeParticion(auxPath, pname) {
		if numExtendida != -1 {

			ebr := EBR{}

			cont := mbr.Mbr_partition[numExtendida].Part_start
			file.Seek(cont, 0)
			var sizeEbr int64 = int64(unsafe.Sizeof(ebr))
			dataControl := leerBytes(file, int(sizeEbr))
			bufferControl := bytes.NewBuffer(dataControl)
			err := binary.Read(bufferControl, binary.BigEndian, &ebr)
			if err != nil {
				log.Fatal("binary.Read failed", err)
			}

			if ebr.Part_size == 0 { //Si es la primera
				if mbr.Mbr_partition[numExtendida].Part_size < int64(size_bytes) {
					if archivo == "principal" {
						fmt.Println("Error: la particion logica a crear excede el espacio disponible de la particion extendida")
					}
				} else {
					ebr.Part_status = '0'
					ebr.Part_fit = auxFit
					offset, _ := file.Seek(0, os.SEEK_CUR) //ftell en go
					ebr.Part_start = offset - sizeEbr      //Para regresar al inicio de la extendida
					ebr.Part_size = int64(size_bytes)
					ebr.Part_next = -1
					for n := 0; n < len(pname); n++ {
						ebr.Part_name[n] = pname[n]
					}
					file.Seek(mbr.Mbr_partition[numExtendida].Part_start, 0)
					//Escribimos struct de ebr
					var bufferControlEBR bytes.Buffer
					binary.Write(&bufferControlEBR, binary.BigEndian, &ebr)
					escribirBytes(file, bufferControlEBR.Bytes())

					if archivo == "principal" {
						fmt.Println("Particion logica creada con exito")
					}
				}
			} else {
				offset, _ := file.Seek(0, os.SEEK_CUR) //ftell en go
				for {
					if (ebr.Part_next != -1) && (offset < mbr.Mbr_partition[numExtendida].Part_size+mbr.Mbr_partition[numExtendida].Part_start) {
						file.Seek(ebr.Part_next, 0)
						dataControlebr := leerBytes(file, int(sizeEbr))
						bufferControlerebr := bytes.NewBuffer(dataControlebr)
						err := binary.Read(bufferControlerebr, binary.BigEndian, &ebr)
						if err != nil {
							log.Fatal("binary.Read failed", err)
						}
					} else {
						break
					}
				}
				//var espacioNecesario int
				espacioNecesario := ebr.Part_start + ebr.Part_size + int64(size_bytes)

				if espacioNecesario <= (mbr.Mbr_partition[numExtendida].Part_size + mbr.Mbr_partition[numExtendida].Part_start) {
					ebr.Part_next = ebr.Part_start + ebr.Part_size
					//Escribimos el next del ultimo EBR
					file.Seek(offset-sizeEbr, 0)
					//Escribimos struct de ebr
					/*var bufferControlEBR bytes.Buffer
					binary.Write(&bufferControlEBR, binary.BigEndian, &ebr)
					escribirBytes(file, bufferControlEBR.Bytes())*/
					//Escribimos el nuevo EBR
					file.Seek(ebr.Part_start+ebr.Part_size, 0)
					ebr.Part_status = 0
					ebr.Part_fit = auxFit
					ebr.Part_start = offset
					ebr.Part_size = int64(size_bytes)
					ebr.Part_next = -1
					for i := 0; i < 16; i++ {
						ebr.Part_name[i] = '0'
					}
					for n := 0; n < len(pname); n++ {
						ebr.Part_name[n] = pname[n]
					}
					//file.Seek(mbr.Mbr_partition[numExtendida].Part_start, 0)
					file.Seek(ebr.Part_start+ebr.Part_size, 0)
					//Escribimos struct de ebr
					var bufferNuevo bytes.Buffer
					binary.Write(&bufferNuevo, binary.BigEndian, &ebr)
					escribirBytes(file, bufferNuevo.Bytes())

					if archivo == "principal" {
						fmt.Println("Particion logica creada con exito")
					}
				} else {
					fmt.Println("Error: la particion logica a crear excede el")
					fmt.Println("espacio disponible de la particion extendida")
				}
			}
		} else {
			fmt.Println("Error: se necesita una particion extendida donde guardar la logica")
		}
	} else {
		fmt.Println("Error: ya existe una particion con ese nombre")
	}
}

func archivoExiste(ruta string) bool {
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		return false
	}
	return true
}

func leerBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func existeParticion(ppath string, pname string) bool {
	var extendida int = -1
	comillaDer := strings.TrimRight(ppath, "\"")
	comillaIzq := strings.TrimLeft(comillaDer, "\"")
	auxPath := comillaIzq
	file, _ := os.OpenFile(auxPath, os.O_RDWR, 0777)
	defer file.Close()

	mbr := MBR{}

	file.Seek(0, 0)
	var sizeMbr int64 = int64(unsafe.Sizeof(mbr))
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

	for i := 0; i < len(pname); i++ {
		byteName[i] = pname[i]
	}

	for i := 0; i < 4; i++ {
		if mbr.Mbr_partition[i].Part_name == byteName {
			defer file.Close()
			return true
		} else if mbr.Mbr_partition[i].Part_type == 'E' {
			extendida = i
		}
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
					defer file.Close()
					return true
				}
				if ebr.Part_next == -1 {
					defer file.Close()
					return false
				}
			}
		}
	}

	defer file.Close()
	return false
}
