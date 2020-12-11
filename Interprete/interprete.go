package Interprete

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"../Metodos"
	"../Structs"
)

func Interpreter(linea string, disco *[27]Structs.Disco) {
	linea = ruta(linea)
	linea = comentario(linea)
	comando := strings.Split(linea, " ")
	switch ejecutar := comando[0]; strings.ToLower(ejecutar) {
	case strings.ToLower("exec"):
		fmt.Println("Comando RMDISK")
		var path string = ""
		for i := 1; i <= len(comando)-1; i++ {
			if comando[i] != "" {
				path = Exec(comando[i])
			}
		}
		if path != "" {
			archivo(path, disco)
		} else {
			fmt.Println("Se requiere de una ruta")
		}
	case strings.ToLower("mkdisk"):
		fmt.Println("Comando MKDISK")
		mkdisk := Structs.Mkdisk{}
		for i := 1; i <= len(comando)-1; i++ {
			if comando[i] != "" {
				Mkdisk(comando[i], &mkdisk)
			}
		}
		if mkdisk.Path != "" && mkdisk.Size > 0 {
			if mkdisk.Unit == 0 {
				mkdisk.Unit = 'm'
			}

			if mkdisk.Fit == 0 {
				mkdisk.Fit = 'W'
			}
			if strings.ToLower(string(mkdisk.Unit)) == "k" || strings.ToLower(string(mkdisk.Unit)) == "m" {
				if strings.ToUpper(string(mkdisk.Fit)) == "B" || strings.ToUpper(string(mkdisk.Fit)) == "F" || strings.ToUpper(string(mkdisk.Fit)) == "W" {
					Metodos.CreateFile(mkdisk)
				} else {
					fmt.Println("Ajuste no renocido")
				}
			} else {
				fmt.Println("Unidades no renocidas")
			}

		} else {
			fmt.Println("Error: Hace falta uno de los siguientes atributos obligatorios path o size")
		}
	case strings.ToLower("rmdisk"):
		fmt.Println("Comando RMDISK")
		var pathEliminar string
		for i := 1; i <= len(comando)-1; i++ {
			if comando[i] != "" {
				pathEliminar = Rmdisk(comando[i])
			}
		}
		Metodos.DeleteDisk(pathEliminar, disco)
	case strings.ToLower("fdisk"):
		fmt.Println("Comando FDISK")
		fdisk := Structs.Fdisk{}
		for i := 1; i <= len(comando)-1; i++ {
			if comando[i] != "" {
				Fdisk(comando[i], &fdisk)
			}
		}
		if fdisk.Size > 0 && fdisk.Path != "" && fdisk.Name != "" {
			if fdisk.Type == 0 {
				fdisk.Type = 'p'
			}
			if fdisk.Fit == 0 {
				fdisk.Fit = 'w'
			}
			if fdisk.Unit == 0 {
				fdisk.Unit = 'k'
			}
			if strings.ToLower(string(fdisk.Type)) != "l" {
				if fdisk.Add != 0 && fdisk.Delete == "" {
					Metodos.CreatePartition(fdisk)
				} else if fdisk.Add == 0 && fdisk.Delete != "" {
					Metodos.CreatePartition(fdisk)
				} else if fdisk.Add == 0 && fdisk.Delete == "" && (strings.ToLower(string(fdisk.Fit)) == "b" || strings.ToLower(string(fdisk.Fit)) == "f" || strings.ToLower(string(fdisk.Fit)) == "w") {
					Metodos.CreatePartition(fdisk)
				} else {
					fmt.Println("Error: no se puede ejucutar add y delete al mismo tiempo")
				}
			} else {
				fmt.Println("No se peude hacer particiones logicas")
			}

		} else {
			fmt.Println("Error: Hace falta uno de los siguientes atributos obligatorios path, name o size")
		}
	case strings.ToLower("mount"):
		fmt.Println("Comando MOUNT")
		montar := Structs.Montar{}
		if len(comando) > 1 {
			for i := 1; i <= len(comando)-1; i++ {
				if comando[i] != "" {
					Mount(comando[i], &montar)
				}
			}
			if montar.Name != "" {
				Metodos.Montar(montar, disco)
			}
		} else {
			fmt.Println("Comando Mount solo")
		}
	case strings.ToLower("unmount"):
		fmt.Println("Comando UNMOUNT")
		unMount := Structs.UnMount{}
		for i := 1; i <= len(comando)-1; i++ {
			if comando[i] != "" {
				Unmount(comando[i], &unMount)
			}
		}
		if unMount.Identificador != "" {
			Metodos.Desmontar(unMount.Identificador, disco)
		} else {
			fmt.Println("No ingreso id")
		}
	case strings.ToLower("pause"):
		fmt.Println("Comando PAUSE")
		fmt.Println("Presione cualquie letra para continuar...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	case strings.ToLower("rep"):
		fmt.Println("Comando REP")
		rep := Structs.Rep{}
		for i := 1; i < len(comando); i++ {
			if comando[i] != "" {
				Rep(comando[i], &rep)
			}
		}
		if rep.Path != "" && rep.Name != "" && rep.Identificador != "" {
			Metodos.Rep(rep, disco)
		} else {
			fmt.Println("Error: Hace falta uno de los siguientes atributos obligatorios path, name o identificador")
		}
	case "":
	default:
		fmt.Println("Comando no reconocido")
	}
	fmt.Println("-----------------------------------------")
	/* mostrar(disco)
	fmt.Println("-----------------------------------------") */
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

func Mount(linea string, montar *Structs.Montar) {
	comando := strings.Split(linea, "->")
	switch ajecutar := comando[0]; strings.ToLower(ajecutar) {
	case strings.ToLower("-path"):
		fmt.Println("Atributo Path: " + ajecutar)
		fmt.Println("Valor Atributo: " + strings.ReplaceAll(comando[1], "\u0022", ""))
		montar.Path = comando[1]
	case strings.ToLower("-name"):
		fmt.Println("Atributo Name: " + ajecutar)
		fmt.Println("Valor Atributo: " + strings.ReplaceAll(comando[1], "\u0022", ""))
		montar.Name = comando[1]
	default:
		fmt.Println("Error: Este atributo " + ajecutar + " no existe en el comando mount")
	}
}

func Unmount(linea string, umount *Structs.UnMount) {
	comando := strings.Split(linea, "->")
	ejecutar := comando[0]
	if strings.Contains(strings.ToLower(ejecutar), strings.ToLower("-id")) {
		r, _ := regexp.Compile("^[vV][dD][a-zA-ZñÑ][0-9]+$")
		fmt.Println(r.FindString(comando[1]))
		umount.Identificador = r.FindString(comando[1])
	} else {
		fmt.Println("Error: Este atributo " + ejecutar + " no existe en el comando unmount")
	}
}

func Rep(linea string, rep *Structs.Rep) {
	comando := strings.Split(linea, "->")
	switch ejecutar := comando[0]; strings.ToLower(ejecutar) {
	case strings.ToLower("-path"):
		fmt.Println("Atributo Path: " + ejecutar)
		fmt.Println("Valor Atributo: " + strings.ReplaceAll(comando[1], "\u0022", ""))
		rep.Path = comando[1]
	case strings.ToLower("-name"):
		fmt.Println("Atributo Name: " + ejecutar)
		fmt.Println("Valor Atributo: " + strings.ReplaceAll(comando[1], "\u0022", ""))
		rep.Name = comando[1]
	case strings.ToLower("-id"):
		r, _ := regexp.Compile("^[vV][dD][a-zA-ZñÑ][0-9]+$")
		fmt.Println(r.FindString(comando[1]))
		rep.Identificador = r.FindString(comando[1])
	default:
		fmt.Println("Error: Este atributo " + ejecutar + " no existe en el comando Rep")
	}
}

func comentario(linea string) string {
	r, _ := regexp.Compile("#.*")
	if r.MatchString(linea) {
		fmt.Println("COMENTARIO: ", r.FindString(linea))
	}
	r = regexp.MustCompile("#.*")
	return r.ReplaceAllString(linea, "")
}

func ruta(linea string) string {
	r, _ := regexp.Compile("\"[^\"]*\"")
	if r.MatchString(linea) {
		nuevo := strings.ReplaceAll(r.FindString(linea), " ", "_")
		r = regexp.MustCompile("\"[^\"]*\"")
		linea = r.ReplaceAllString(linea, nuevo)
	}
	return linea
}

func archivo(path string, disco *[27]Structs.Disco) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Print(err)
	}

	str := string(b)

	fmt.Println(str)

	var separarComandos []string = strings.Split(str, "\n")

	for i := 0; i <= len(separarComandos)-1; i++ {
		Interpreter(separarComandos[i], disco)
	}
}

func mostrar(disco *[27]Structs.Disco) {
	for i := 0; i < len(disco); i++ {
		if disco[i].Status == 1 {
			fmt.Println(disco[i].Path)
			for j := 0; j < len(disco[i].Particiones); j++ {
				if disco[i].Particiones[j].Status == 1 {
					fmt.Println(j+1, ". Id: ", disco[i].Particiones[j].Identificador, ". Name: ", disco[i].Particiones[j].Name)
				}
			}
		}
	}
}

func validarFit(fit string) string {
	if strings.ToLower(fit) == "bf" {
		return "b"
	} else if strings.ToLower(fit) == "ff" {
		return "f"
	} else if strings.ToLower(fit) == "wf" {
		return "w"
	}
	return ""
}
