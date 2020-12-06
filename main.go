package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"./Interprete"
)

func main() {
	menu()
}

func menu() {
	finalizar := 0
	fmt.Println("Bienvenido a la consola de comandos... ('x' para finalizar)")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter comands: ")
	comando, _ := reader.ReadString('\n')

	if comando == "x\n" {
		finalizar = 1
	} else {
		if comando != "" {
			Interprete.Interpreter(strings.TrimSpace(comando))
		}
	}
	fmt.Println("-----------------------------------------")
	for finalizar != 1 {
		fmt.Print("Enter comands: ")
		reader := bufio.NewReader(os.Stdin)
		comando, _ := reader.ReadString('\n')
		if comando == "x\n" {
			finalizar = 1
		} else {
			if comando != "" {
				Interprete.Interpreter(strings.TrimSpace(comando))
			}
		}
		fmt.Println("-----------------------------------------")
	}
}
