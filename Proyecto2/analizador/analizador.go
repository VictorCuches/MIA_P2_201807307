package analizador

import (
	"Proyecto2/analizador/comandos"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Start_Analizer() {

	intContador := 0
	for { // simulacion de un while infinito

		EncabezadoInfo()
		// leyendo entrada por consola
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		entrada := scanner.Text()

		arrEntrada := strings.Split(entrada, " ")

		if strings.ToUpper(arrEntrada[0]) == "SALIR" {
			fmt.Println("qchau :D")
			break
		}

		fmt.Println("\n-------------------------------------------")
		Comandos_indentify(arrEntrada[0], arrEntrada, entrada)
		fmt.Println("---------------------------------------------")
		fmt.Println(" ")

		intContador++

	}
}

var bandera_login bool = false
var lista *comandos.Lista = comandos.New_Lista()
var listaLogin *comandos.ListaL = comandos.New_ListaL()
var listaFiles *comandos.ListaF = comandos.New_ListaF()

// Identifica cada comando
func Comandos_indentify(command string, params []string, entrada string) {
	// fmt.Println("Comando -> ", command)
	// fmt.Println("**************************************")
	// fmt.Println(params)
	// fmt.Println(params[1])
	// fmt.Println(params[2])

	if strings.ToUpper(command) == "MKDISK" {
		fmt.Println("============ MKDISK ============")
		// enviando los parametros
		//mkdisk.DoMkdisk(params)
		comandos.DoMkdisk(params)

	} else if strings.ToUpper(command) == "RMDISK" {
		fmt.Println("============ RMKDISK ============")
		comandos.DoRmdisk(params)

	} else if strings.ToUpper(command) == "FDISK" {
		fmt.Println("============ FDISK ============")
		comandos.DoFdisk(entrada)
	} else if strings.ToUpper(command) == "MOUNT" {
		fmt.Println("============ MOUNT ============")
		comandos.DoMount(entrada, lista)
	} else if strings.ToUpper(command) == "MKFS" {
		fmt.Println("============ MKFS ============")
		comandos.DoMkfs(entrada, lista)
	} else if strings.ToUpper(command) == "LOGIN" {
		fmt.Println("============ LOGIN ============")
		if !bandera_login {
			comandos.DoLogin(entrada, lista, listaLogin)
			bandera_login = true
		} else {
			fmt.Println("Error: ya hay un usuario logeado")
		}
	} else if strings.ToUpper(command) == "LOGOUT" {
		fmt.Println("============ LOGOUT ============")
		if !bandera_login {
			fmt.Println("No hay ninguna sesion iniciada")
		} else {
			fmt.Println("Se ha terminado la sesion correctamente")
			bandera_login = false
		}
	} else if strings.ToUpper(command) == "MKGRP" {
		fmt.Println("============ MKGRP ============")
		if !bandera_login {
			fmt.Println("Error: No se ha iniciado sesión")
		} else {
			comandos.DoMkgrp(entrada)
		}
	} else if strings.ToUpper(command) == "RMGRP" {
		fmt.Println("< < < < < RMGRP > > > > >")
	} else if strings.ToUpper(command) == "MKUSER" {
		fmt.Println("< < < < < MKUSER > > > > >")
	} else if strings.ToUpper(command) == "RMUSR" {
		fmt.Println("< < < < < RMUSR > > > > >")
	} else if strings.ToUpper(command) == "MKFILE" {
		fmt.Println("============ MKFILE ============")
		if !bandera_login {
			fmt.Println("Error: No se ha iniciado sesión")
		} else {
			comandos.DoRep(entrada, lista, listaFiles)
		}
	} else if strings.ToUpper(command) == "MKDIR" {
		fmt.Println("< < < < < MKDIR > > > > >")
	} else if strings.ToUpper(command) == "PAUSE" {
		fmt.Println("============ PAUSE ============")
		DoPause()
	} else if strings.ToUpper(command) == "EXEC" {
		fmt.Println("============ EXEC ============")
		DoExec(entrada)
	} else if strings.ToUpper(command) == "REP" {
		fmt.Println("============ REP ============")
		comandos.DoRep(entrada, lista, listaFiles)
	} else {
		fmt.Println("Comando no reconocido")
	}

}

// Muestra informacion personal
func EncabezadoInfo() {
	fmt.Println("******************************")
	fmt.Println("***** VICTOR CUCHES   P2 *****")
	fmt.Println("********* 201807307 **********")
	fmt.Println("******************************")
	fmt.Println("Ingrese comando: ")
}

func DoPause() {
	fmt.Println("Presione enter para poder continuar")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	entrada := text
	for {
		if entrada == "\n" {
			break
		}
	}
}
