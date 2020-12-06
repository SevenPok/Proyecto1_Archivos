package Interprete

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"../Structs"
)

func Interpreter(linea string) {
	comando := strings.Split(linea, " ")
	switch ejecutar := comando[0]; strings.ToLower(ejecutar) {
	case strings.ToLower("exec"):
		fmt.Println("Comando RMDISK")
	case strings.ToLower("mkdisk"):
		fmt.Println("Comando MKDISK")
		mkdisk := Structs.Mkdisk{}
		for i := 1; i <= len(comando)-1; i++ {
			Mkdisk(comando[i], &mkdisk)
		}
		if mkdisk.Path != "" && mkdisk.Size > 0 {
			fmt.Println("COMANDO EJECUTADO")
		} else {
			fmt.Println("Error: Hace falta uno de los siguientes atributos obligatorios path o size")
		}
	case strings.ToLower("rmdisk"):
		fmt.Println("Comando RMDISK")
		fmt.Println(Rmdisk(comando[1]))

	case strings.ToLower("fdisk"):
		fmt.Println("Comando FDISK")
		fdisk := Structs.Fdisk{}
		for i := 1; i <= len(comando)-1; i++ {
			Fdisk(comando[i], &fdisk)
		}
		if fdisk.Size > 0 && fdisk.Path != "" && fdisk.Name != "" {
			fmt.Println("COMANDO EJECUTADO")
		} else if fdisk.Path != "" && fdisk.Name != "" {
			fmt.Println("COMANDO EJECUTADO")
		} else {
			fmt.Println("Error: Hace falta uno de los siguientes atributos obligatorios path, name o size")
		}
	case strings.ToLower("mount"):
		fmt.Println("Comando MOUNT")
		if len(comando) > 1 {
			for i := 1; i <= len(comando)-1; i++ {
				Mount(comando[i])
			}
		} else {
			fmt.Println("Comando Mount solo")
		}
	case strings.ToLower("unmount"):
		fmt.Println("Comando UNMOUNT")
		for i := 1; i <= len(comando)-1; i++ {
			Unmount(comando[i])
		}
	case strings.ToLower("pause"):
		fmt.Println("Comando PAUSE")
		fmt.Println("Presione cualquie letra para continuar...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	case strings.ToLower("rep"):
	case strings.ToLower("mbr"):
	case strings.ToLower("disk"):
	default:
		fmt.Println("Comando no reconocido")
	}

}

func Exec(linea string) string {
	comando := strings.Split(linea, "->")
	switch ejecutar := comando[0]; strings.ToLower(ejecutar) {
	case strings.ToLower("-path"):
		return strings.ReplaceAll(comando[1], "\u0022", "")
	default:
		fmt.Println("Error: Este atributo " + ejecutar + " no existe en el comando exec")
		return ""
	}
}

func Mkdisk(linea string, mkdisk *Structs.Mkdisk) {
	comando := strings.Split(linea, "->")
	switch ejecutar := comando[0]; strings.ToLower(ejecutar) {
	case strings.ToLower("-path"):
		mkdisk.Path = strings.ReplaceAll(comando[1], "\u0022", "")
	case strings.ToLower("-size"):
		var s string = comando[1]
		i, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			mkdisk.Size = i
		} else {
			fmt.Println("Error: Se esperaba valor numerico")
		}
	case strings.ToLower("-unit"):
		var unit []byte = []byte(comando[1])
		if strings.Compare(strings.ToLower(string(unit)), strings.ToLower("k")) == 0 || strings.Compare(strings.ToLower(string(unit)), strings.ToLower("m")) == 0 {
			mkdisk.Unit = unit[0]
		} else {
			fmt.Println("Error: Valor del atributo unit solo puede ser k o m")
		}
	default:
		fmt.Println("Error: Este atributo " + ejecutar + " no existe en el comando mkdisk")
	}
}

func Rmdisk(linea string) string {
	comando := strings.Split(linea, "->")
	switch ajecutar := comando[0]; strings.ToLower(ajecutar) {
	case strings.ToLower("-path"):
		return strings.ReplaceAll(comando[1], "\u0022", "")
	default:
		fmt.Println("Error: Este atributo " + ajecutar + " no existe en el comando rmdisk")
		return ""
	}
}

func Fdisk(linea string, fdisk *Structs.Fdisk) {
	comando := strings.Split(linea, "->")
	switch ejecutar := comando[0]; strings.ToLower(ejecutar) {
	case strings.ToLower("-path"):
		fdisk.Path = strings.ReplaceAll(comando[1], "\u0022", "")
	case strings.ToLower("-size"):
		var s string = comando[1]
		i, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			fdisk.Size = i
		} else {
			fmt.Println("Error: Se esperaba valor numerico")
		}
	case strings.ToLower("-name"):
		fdisk.Name = strings.ReplaceAll(comando[1], "\u0022", "")
	case strings.ToLower("-unit"):
		var unit []byte = []byte(comando[1])
		if strings.Compare(strings.ToLower(string(unit)), strings.ToLower("k")) == 0 || strings.Compare(strings.ToLower(string(unit)), strings.ToLower("m")) == 0 || strings.Compare(strings.ToLower(string(unit)), strings.ToLower("b")) == 0 {
			fdisk.Unit = unit[0]
		} else {
			fmt.Println("Error: Valor del atributo unit solo puede ser k, m o b")
		}
	case strings.ToLower("-type"):
		var unit []byte = []byte(strings.ToLower(comando[1]))
		if strings.Compare(strings.ToLower(string(unit)), strings.ToLower("p")) == 0 || strings.Compare(strings.ToLower(string(unit)), strings.ToLower("e")) == 0 || strings.Compare(strings.ToLower(string(unit)), strings.ToLower("l")) == 0 {
			fdisk.Type = unit[0]
		} else {
			fmt.Println("Error: Valor del atributo Type solo puede ser P, E o L")
		}
	case strings.ToLower("-fit"):
		var unit []byte = []byte(comando[1])
		if strings.Compare(strings.ToLower(string(unit)), strings.ToLower("bf")) == 0 || strings.Compare(strings.ToLower(string(unit)), strings.ToLower("ff")) == 0 || strings.Compare(strings.ToLower(string(unit)), strings.ToLower("wf")) == 0 {
			fdisk.Fit = unit[0]
		} else {
			fmt.Println("Error: Valor del atributo Fit solo puede ser bf, ff o wf")
		}
	case strings.ToLower("-delete"):
		if strings.Compare(strings.ToLower(comando[1]), strings.ToLower("fast")) == 0 || strings.Compare(strings.ToLower(comando[1]), strings.ToLower("full")) == 0 {
			fdisk.Delete = comando[1]
		} else {
			fmt.Println("Error: Valor del atributo Delete solo puede ser Fast o Full")
		}
	case strings.ToLower("-add"):
		var s string = comando[1]
		i, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			fdisk.Add = i
		} else {
			fmt.Println("Error: Se esperaba valor numerico")
		}
	default:
		fmt.Println("Error: Este atributo " + ejecutar + " no existe en el comando fdisk")
	}
}

func Mount(linea string) {
	comando := strings.Split(linea, "->")
	switch ajecutar := comando[0]; strings.ToLower(ajecutar) {
	case strings.ToLower("-path"):
		fmt.Println("Atributo Path: " + ajecutar)
		fmt.Println("Valor Atributo: " + strings.ReplaceAll(comando[1], "\u0022", ""))
	case strings.ToLower("-name"):
		fmt.Println("Atributo Name: " + ajecutar)
		fmt.Println("Valor Atributo: " + strings.ReplaceAll(comando[1], "\u0022", ""))
	default:
		fmt.Println("Error: Este atributo " + ajecutar + " no existe en el comando mount")
	}
}

func Unmount(linea string) {
	comando := strings.Split(linea, "->")
	ejecutar := comando[0]
	valor := comando[1]
	if strings.Contains(strings.ToLower(ejecutar), strings.ToLower("-id")) {
		numero := strings.Split(strings.ToLower(ejecutar), strings.ToLower("-id"))
		var s string = numero[1]
		i, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			fmt.Println("Atributo IDn: " + ejecutar)
			fmt.Println("la n en IDn: ", i)
			fmt.Println("Valor Atributo: ", valor)
		} else {
			fmt.Println("Error: Se esperaba valor numerico")
		}
	} else {
		fmt.Println("Error: Este atributo " + ejecutar + " no existe en el comando unmount")
	}
}

func Rep() {

}

func Mbr() {

}

func Disk() {

}
