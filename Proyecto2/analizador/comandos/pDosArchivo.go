package comandos

import (
	"fmt"
	"math"
	"os"
	"time"
	"unsafe"
)

const SIZE_BLOCKS int = 64

func CalNoEstructuras(size_particion int, part_type byte, fs string) int {
	//superbloque := Superbloque{}
	//inodos := Inodos{}
	//ebr := EBR{}
	var sizeSuperBloque int64 = int64(unsafe.Sizeof(Superbloque{}))
	var sizeInodos int64 = int64(unsafe.Sizeof(Inodos{}))
	var sizeEbr int64 = int64(unsafe.Sizeof(EBR{}))
	switch part_type {
	case 'P':
		num := math.Floor((float64(size_particion) - float64(sizeSuperBloque)) / float64((int64(4) + sizeInodos + int64((3 * SIZE_BLOCKS)))))
		return int(num)
	case 'L': //Logica : Al tamaño total de la partición se le resta el tamaño del ebr
		numL := math.Floor((float64(size_particion) - float64(sizeSuperBloque) - float64(sizeEbr)) / float64((int64(4) + sizeInodos + int64((3 * SIZE_BLOCKS)))))
		return int(numL)
	}
	return 0
}

func FormatearParticion(file *os.File, particion Particion, tipo string, fs string) {
	var no_struct int = CalNoEstructuras(int(particion.Part_size), particion.Part_type, fs)
	//inodos := Inodos{}
	superbloque := Superbloque{}
	var sizeSuperBloque int64 = int64(unsafe.Sizeof(superbloque))
	var sizeInodos int64 = int64(unsafe.Sizeof(Inodos{}))
	if no_struct > 0 {

		var tempTime [19]byte
		tiempo := time.Now()
		time := tiempo.String()

		for t := 0; t < 19; t++ {
			tempTime[t] = time[t]
		}

		var inicio_byte int = int(particion.Part_start)
		//superbloque := Superbloque{}
		superbloque.S_inodes_count = int64(no_struct)
		superbloque.S_blocks_count = (int64(no_struct) * 3)
		superbloque.S_free_blocks_count = (int64(no_struct) * 3)
		superbloque.S_free_inodes_count = int64(no_struct)
		superbloque.S_mtime = tempTime
		superbloque.S_mnt_count = 1
		superbloque.S_magic = 0xEF53
		superbloque.S_inode_size = sizeInodos
		superbloque.S_block_size = int64(SIZE_BLOCKS)
		superbloque.S_firts_ino = 0
		superbloque.S_first_blo = 0
		var fileSystem int = 2
		switch fileSystem {
		case 2:
			superbloque.S_filesystem_type = 2
			superbloque.S_bm_inode_start = (int64(inicio_byte) + sizeSuperBloque)
			superbloque.S_bm_block_start = (int64(inicio_byte) + sizeSuperBloque + int64(no_struct))
			superbloque.S_inode_start = (int64(inicio_byte) + sizeSuperBloque + int64(no_struct) + int64((3 * no_struct)))
			superbloque.S_block_start = (int64(inicio_byte) + sizeSuperBloque + int64(no_struct) + int64((3 * no_struct)) + (int64(no_struct) * sizeInodos))
		}
		//Crear carpeta raiz
		CrearCarpetaRaiz(file, superbloque, inicio_byte)
		//Crear Archivo users.txt
		fmt.Println("Se ha formateado la partición correctamente")
	} else {
		fmt.Println("Error: La particion o el -fs no es correcto")
	}
}

func CrearCarpetaRaiz(file *os.File, superbloque Superbloque, inicio_byte int) {
	var no_struct int = int(superbloque.S_inodes_count)
	//Bitmap inodos
	var bitmap_inodos []byte
	for i := 0; i < no_struct; i++ {
		bitmap_inodos[i] = '0'
	}

	//Bitmap bloques
	var bitmap_blocks []byte
	for i := 0; i < (no_struct * 3); i++ {
		bitmap_blocks[i] = '0'
	}

	var tempTime [19]byte
	tiempo := time.Now()
	time := tiempo.String()

	for t := 0; t < 19; t++ {
		tempTime[t] = time[t]
	}

	//Crear carpeta raiz
	inodos := Inodos{}
	inodos.I_uid = 1
	inodos.I_gid = 1
	inodos.I_size = 0
	inodos.I_atime = tempTime
	inodos.I_ctime = tempTime
	inodos.I_mtime = tempTime
	inodos.I_block = 0
	inodos.I_type = '0' //Carpeta
	inodos.I_perm = 000
	//Ocupamos el bloque 0 en el bitmap inodos
	bitmap_inodos[0] = '1'
	superbloque.S_free_blocks_count--

	//Crear Nodo apuntador--------------------
}
