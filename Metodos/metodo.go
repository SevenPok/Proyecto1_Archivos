package Metodos

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"../Structs"
)

func ReadFile(path string, mbr Structs.MBR) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var size int = int(unsafe.Sizeof(mbr))
	data := readBytes(file, size)
	buffer := bytes.NewBuffer(data)

	err = binary.Read(buffer, binary.BigEndian, &mbr)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	fmt.Println(mbr)
}

func CreateFile(mkdisk Structs.Mkdisk) {
	ensureDir(mkdisk.Path)
	file, err := os.Create(mkdisk.Path)
	defer file.Close()
	if err != nil {
		fmt.Println("No se pudo crear el archivo")
	} else {
		var otro int8 = 0
		var size int64 = 1024 * 1024
		if strings.ToLower(string(mkdisk.Unit)) == "k" {
			size = 1024
		}

		s := &otro

		var binario bytes.Buffer
		binary.Write(&binario, binary.BigEndian, s)
		writeBytes(file, binario.Bytes())

		file.Seek((mkdisk.Size-1)*size, 0)

		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, s)
		writeBytes(file, binario2.Bytes())

		file.Seek(0, 0)

		disco := Structs.MBR{Size: (mkdisk.Size - 1) * size, Fit: mkdisk.Fit}

		t := time.Now()
		fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())
		cadenita := fecha
		copy(disco.Date[:], cadenita)
		rand.Seed(time.Now().UnixNano())

		disco.Signature = int64(rand.Intn(10001))
		s1 := &disco

		var binario3 bytes.Buffer
		binary.Write(&binario3, binary.BigEndian, s1)
		writeBytes(file, binario3.Bytes())

		fmt.Println("El archivo se creo exitosamente, el tamanio del mbr es de: ", int(unsafe.Sizeof(disco)), " bytes")
		size, err := GetFileSize1(mkdisk.Path)
		if err != nil {

		}
	}
}

func DeleteDisk(path string, disco *[27]Structs.Disco) {
	fmt.Println("Presione y/n para continuar...")
	reader := bufio.NewReader(os.Stdin)
	comando, _ := reader.ReadString('\n')
	if strings.TrimSpace(comando) == "y" {
		err := os.Remove(path)
		if err != nil {
			for i := 0; i < len(disco); i++ {
				if path == disco[i].Path {
					disco[i] = Structs.Disco{}
					break
				}
			}
			fmt.Println("Error al eliminar disco: ", err)
		} else {
			fmt.Println("Disco eliminado correctamente")
		}
	}
}

func GetFileSize1(filepath string) (int64, error) {
	fi, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}
	// get the size
	return fi.Size(), nil
}

func CreatePartition(fdisk Structs.Fdisk) {
	file, err := os.OpenFile(fdisk.Path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		fmt.Println("No se encontro el disco")
	} else {
		mbr := Structs.MBR{}
		var size int = int(unsafe.Sizeof(mbr))
		mbr = readDisk(file, size, mbr)

		var sizeByte int64 = fdisk.Size * definirUnidad(string(fdisk.Unit))
		partition := Structs.Partition{Status: 1, Type: fdisk.Type, Fit: fdisk.Fit, Size: sizeByte}
		copy(partition.Name[:], fdisk.Name)

		if strings.ToLower(string(fdisk.Type)) == "e" {
			ebr := Structs.EBR{Status: 1, Fit: fdisk.Fit, Next: -1}
			copy(ebr.Name[:], fdisk.Name)
			ebr.Size = int64(unsafe.Sizeof(ebr))
			partition.Ebr = ebr
		}

		if fdisk.Add != 0 && validarNombre(mbr, string(partition.Name[:])) {
			partition.Size = fdisk.Add * definirUnidad(string(fdisk.Unit))
			mbr = add(mbr, partition)
		} else if fdisk.Delete != "" && fdisk.Add == 0 && validarNombre(mbr, string(partition.Name[:])) {

		} else if !validarNombre(mbr, string(partition.Name[:])) && soloUnaParticionExtendida(mbr, partition) && validacion(mbr, partition, fdisk) {
			mbr = ajuste(mbr, partition)
		} else {
			fmt.Println("El nombre de la particion ya existe, no se puede crear")
		}
		file.Seek(0, 0)
		var buffer bytes.Buffer
		binary.Write(&buffer, binary.BigEndian, &mbr)
		writeBytes(file, buffer.Bytes())
		Show(fdisk.Path)

	}
}

func Montar(mount Structs.Montar, disco *[27]Structs.Disco) {
	file, err := os.OpenFile(mount.Path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		fmt.Println("No se encontro el disco")
	} else {
		mbr := Structs.MBR{}
		var size int = int(unsafe.Sizeof(mbr))
		mbr = readDisk(file, size, mbr)

		partition := Structs.Partition{}
		copy(partition.Name[:], mount.Name)
		agregarPath(mount, disco)
		for i := 0; i < len(mbr.Partition); i++ {
			if partition.Name == mbr.Partition[i].Name {
				agregarMount(mount, disco)
				break
			} else if i == 3 {
				fmt.Println("No existe la particion")
			}
		}

	}
}

func Desmontar(id string, disco *[27]Structs.Disco) bool {
	letra := strings.ReplaceAll(id, "vd", "")
	r, _ := regexp.Compile("[0-9]+")
	if r.MatchString(letra) {
		r = regexp.MustCompile("[0-9]+")
		letra = r.ReplaceAllString(letra, "")
	}
	for i := 0; i < len(disco); i++ {
		if letra == disco[i].Letra {
			for j := 0; j < len(disco[i].Particiones); j++ {
				if id == disco[i].Particiones[j].Identificador && disco[i].Particiones[j].Status == 1 {
					disco[i].Particiones[j] = Structs.ParticionMontada{}
					fmt.Println("Se ha desmontado correctamente la particion")
					return true
				}
			}
			fmt.Println("No existe el identificador")
			return false
		}
	}
	fmt.Println("No existe el identificador")
	return false
}

func agregarMount(mount Structs.Montar, disco *[27]Structs.Disco) {
	for i := 0; i < len(disco); i++ {
		if mount.Path == disco[i].Path {
			for j := 0; j < len(disco[i].Particiones); j++ {
				if disco[i].Particiones[j].Status == 0 {
					disco[i].Particiones[j].Identificador = "vd" + letra(i) + strconv.Itoa(j+1)
					disco[i].Particiones[j].Path = mount.Path
					disco[i].Particiones[j].Name = mount.Name
					disco[i].Particiones[j].Status = 1
					fmt.Println("Se ha montado correctamente la particion")
					break
				}
			}
			break
		}
	}
}

func agregarPath(mount Structs.Montar, disco *[27]Structs.Disco) {
	for i := 0; i < len(disco); i++ {
		if disco[i].Path == mount.Path && mount.Name != "" {
			break
		} else if disco[i].Status == 0 {
			disco[i].Letra = letra(i)
			disco[i].Path = mount.Path
			disco[i].Status = 1
			break
		}
	}
}

func ajuste(mbr Structs.MBR, partition Structs.Partition) Structs.MBR {
	if strings.ToLower(string(partition.Fit)) == "b" {
		return firstFit(mbr, partition)
	} else if strings.ToLower(string(partition.Fit)) == "f" {
		return firstFit(mbr, partition)
	} else {
		return firstFit(mbr, partition)
	}
}

func bestFit(mbr Structs.MBR, partition Structs.Partition) Structs.MBR {

	return mbr
}

func worstFit(mbr Structs.MBR, partition Structs.Partition) Structs.MBR {

	return mbr
}

func firstFit(mbr Structs.MBR, partition Structs.Partition) Structs.MBR {
	var size int64 = int64(unsafe.Sizeof(mbr))
	for i := 0; i < len(mbr.Partition)-1; i++ {
		if mbr.Partition[i].Status == 0 {
			if mbr.Partition[i+1].Status == 0 {
				partition.Start = size
				if strings.ToLower(string(partition.Type)) == "e" {
					partition.Ebr.Start = size
				}
				mbr.Partition[i] = partition
				fmt.Println("Se creo la particion exitosamente")
				return mbr
			}
		} else if espacioEntreParticiones(mbr.Partition[i], mbr.Partition[i+1], partition.Size) {

			partition.Start = size + mbr.Partition[i].Size
			if strings.ToLower(string(partition.Type)) == "e" {
				partition.Ebr.Start = size + mbr.Partition[i].Size
			}
			mbr = insert(mbr, partition)
			fmt.Println("Se creo la particion exitosamente!!!!!")
			mbr = ordenarParticiones(mbr)
			return mbr
		}
		size = mbr.Partition[i].Start + mbr.Partition[i].Size
	}
	if mbr.Size-size > partition.Size {
		partition.Start = size
		if strings.ToLower(string(partition.Type)) == "e" {
			partition.Ebr.Start = size
		}
		mbr.Partition[3] = partition
		fmt.Println("Se creo la particion exitosamente")
	} else {
		fmt.Println("No hay espacio suficiente para la particion")
	}
	return mbr
}

func insert(mbr Structs.MBR, partition Structs.Partition) Structs.MBR {
	for i := 0; i < len(mbr.Partition); i++ {
		if mbr.Partition[i].Status == 0 {
			mbr.Partition[i] = partition
			return mbr
		}
	}
	return mbr
}

func add(mbr Structs.MBR, partition Structs.Partition) Structs.MBR {
	for i := 0; i < len(mbr.Partition)-1; i++ {
		if mbr.Partition[i].Name == partition.Name {
			size := mbr.Partition[i].Size + partition.Size
			if mbr.Partition[i+1].Status == 0 {
				if size < mbr.Size && size > 0 {
					mbr.Partition[i].Size = size
					fmt.Println("Se modifico el espacio de la particion exitosamente")
				} else {
					fmt.Println("No hay espacio suficiente para aumentar la particion")
				}
			} else {
				if partition.Size > 0 {
					if espacioEntreParticiones(mbr.Partition[i], mbr.Partition[i+1], partition.Size) {
						mbr.Partition[i].Size = size
						fmt.Println("Se modifico el espacio de la particion exitosamente")
					} else {
						fmt.Println("No hay espacio suficiente para aumentar la particion")
					}
				} else if size > 0 {
					mbr.Partition[i].Size = size
					fmt.Println("Se modifico el espacio de la particion exitosamente")
				} else {
					fmt.Println("El espacio reducido es mayor que el que tiene la particion")
				}
			}
			return mbr
		}
	}

	if mbr.Partition[3].Name == partition.Name {
		size := mbr.Partition[3].Size + partition.Size
		if size < mbr.Size && size > 0 {
			mbr.Partition[3].Size = size
			fmt.Println("Se modifico el espacio de la particion exitosamente")
		} else {
			fmt.Println("No hay espacio suficiente para aumentar la particion")
		}
	}
	return mbr
}

func valorAbsoluto(num int64) int64 {
	if num < 0 {
		num = num * -1
	}
	return num
}

func validacion(mbr Structs.MBR, partition Structs.Partition, fdisk Structs.Fdisk) bool {
	if validarParticion(mbr) && hayEspacioEnDisco(mbr, partition) {
		return true
	}
	return false
}

func validarParticion(mbr Structs.MBR) bool {
	espacio := 0
	for i := 0; i < len(mbr.Partition); i++ {
		if mbr.Partition[i].Status != 0 {
			espacio++
		}
	}
	if espacio < 4 {
		return true
	} else {
		fmt.Println("Solo se permite crear 4 particiones para este disco")
		return false
	}
}

func hayEspacioEnDisco(mbr Structs.MBR, partition Structs.Partition) bool {
	if mbr.Size < partition.Size {
		fmt.Println("El tamañio de la particion no debe ser mas grande que el disco")
	}
	return mbr.Size > partition.Size
}

func validarNombre(mbr Structs.MBR, name string) bool {
	for i := 0; i < 4; i++ {
		if string(mbr.Partition[i].Name[:]) == name {
			return true
		}
	}
	return false
}

func soloUnaParticionExtendida(mbr Structs.MBR, partition Structs.Partition) bool {
	for i := 0; i < len(mbr.Partition); i++ {
		if strings.ToLower(string(mbr.Partition[i].Type)) == "e" && strings.ToLower(string(partition.Type)) == "e" {
			fmt.Println("Solo se permite crear una particion extendida")
			return false
		}
	}
	return true
}

func ordenarParticiones(mbr Structs.MBR) Structs.MBR {
	for i := 0; i < len(mbr.Partition)-1; i++ {
		for j := 0; j < len(mbr.Partition)-1; j++ {
			if mbr.Partition[j].Start > mbr.Partition[j+1].Start {
				aux := mbr.Partition[j+1]
				mbr.Partition[j+1] = mbr.Partition[j]
				mbr.Partition[j] = aux
			}
		}
	}
	return mbr
}

func espacioEntreParticiones(part1 Structs.Partition, part2 Structs.Partition, size int64) bool {
	dif := part2.Start - (part1.Start + part1.Size)
	//sfmt.Println("La difrenecia es de: ", dif, " espacio a meter: ", size)
	return dif >= size
}

func libre(part1 Structs.Partition, part2 Structs.Partition) int64 {
	dif := part2.Start - (part1.Start + part1.Size)
	//sfmt.Println("La difrenecia es de: ", dif, " espacio a meter: ", size)
	return dif
}

func definirUnidad(unidad string) int64 {
	var size int64
	if strings.ToLower(unidad) == "b" {
		size = 1
	} else if strings.ToLower(unidad) == "k" {
		size = 1024
	} else {
		size = 1024 * 1024
	}
	return size
}

func Show(path string) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		fmt.Println("No se encontro el disco")
	} else {
		mbr := Structs.MBR{}

		var size int = int(unsafe.Sizeof(mbr))

		mbr = readDisk(file, size, mbr)

		for i := 0; i < len(mbr.Partition); i++ {
			if mbr.Partition[i].Status != 0 {
				fmt.Println((i + 1), "). type: ", string(mbr.Partition[i].Type), " size: ", mbr.Partition[i].Size, " start: ", mbr.Partition[i].Start, " name: ", string(mbr.Partition[i].Name[:]))
			}
		}
	}
}

func readDisk(file *os.File, size int, disco Structs.MBR) Structs.MBR {
	data := readBytes(file, size)

	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &disco)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	return disco
}

func writeBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}

func readBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func ensureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			panic(merr)
		}
	}
}

func Rep(rep Structs.Rep, disco *[27]Structs.Disco) {
	ruta := existe(rep, disco)
	if ruta != "" {
		ensureDir(rep.Path)
		dot := ""
		if strings.ToLower(rep.Name) == "mbr" {
			dot = graficarMBR(ruta)
			ruta = strings.Split(rep.Path, ".")[0] + ".dot"

			ext := strings.Split(rep.Path, ".")[1]

			file, err := os.Create(ruta)
			if err != nil {
				fmt.Println("No se pudo crear el archivo")
			} else {
				fmt.Fprintln(file, dot)
			}
			file.Close()
			arch := strings.Split(rep.Path, "/")
			tipo := strings.Split(rep.Path, "/")
			ruta = ""
			for i := 0; i < len(tipo)-1; i++ {

				ruta = ruta + tipo[i] + "/"

			}

			aux := strings.Split(arch[len(arch)-1], ".")[0] + ".dot"

			out, err := exec.Command("dot", "-T"+ext, ruta+aux, "-o", rep.Path).Output()
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			fmt.Println(string(out))
			fmt.Println("Se creo exitosamente el reporte")

		} else if strings.ToLower(rep.Name) == "disk" {

		}

	}
}

func graficarMBR(ruta string) string {
	linea := ""

	file, err := os.Open(ruta)
	defer file.Close()
	if err != nil {
		fmt.Println("No se encontro el disco")
	} else {
		mbr := Structs.MBR{}
		var size int = int(unsafe.Sizeof(mbr))
		mbr = readDisk(file, size, mbr)
		linea = linea + "digraph G {\n node [shape=record];\n"
		linea = linea + "mbr[label=\"MBR"

		for i := 0; i < len(mbr.Partition)-1; i++ {
			if size == int(mbr.Partition[i].Start) && mbr.Partition[i].Status == 0 {
				linea = linea + "|Libre"
			}
			if mbr.Partition[i].Status != 0 {

				if strings.ToLower(string(mbr.Partition[i].Type)) == "p" || strings.ToLower(string(mbr.Partition[i].Type)) == "l" {
					linea = linea + "|Primaria"
				} else {

					linea = linea + "|{Extendida|{EBR|Libre}}"
				}

				if espacioEntreParticiones(mbr.Partition[i], mbr.Partition[i+1], 1) {

					linea = linea + "|Libre"
				}
				size = int(mbr.Partition[i].Start + mbr.Partition[i].Size)
			}
		}

		if mbr.Partition[3].Status != 0 {
			if strings.ToLower(string(mbr.Partition[3].Type)) == "p" || strings.ToLower(string(mbr.Partition[3].Type)) == "l" {
				linea = linea + "|Primaria"
			} else {

				linea = linea + "|{Extendida|{EBR|Libre}}"
			}
		}
		aux := int(mbr.Size) - size

		if aux > 0 {
			linea = linea + "|Libre "
		}

		linea = linea + "\"] }"

	}

	return linea
}

func graficarDisk(ruta string) string {
	linea := ""

	file, err := os.Open(ruta)
	defer file.Close()
	if err != nil {
		fmt.Println("No se encontro el disco")
	} else {
		mbr := Structs.MBR{}
		var size int = int(unsafe.Sizeof(mbr))
		mbr = readDisk(file, size, mbr)
	}

	return linea
}

func existe(rep Structs.Rep, disco *[27]Structs.Disco) string {
	letra := strings.ReplaceAll(rep.Identificador, "vd", "")
	r, _ := regexp.Compile("[0-9]+")
	if r.MatchString(letra) {
		r = regexp.MustCompile("[0-9]+")
		letra = r.ReplaceAllString(letra, "")
	}
	for i := 0; i < len(disco); i++ {
		if letra == disco[i].Letra {
			for j := 0; j < len(disco[i].Particiones); j++ {
				if rep.Identificador == disco[i].Particiones[j].Identificador && disco[i].Particiones[j].Status == 1 {
					return disco[i].Path
				}
			}
			fmt.Println("No existe el identificador")
			return ""
		}
	}
	fmt.Println("No existe el identificador")
	return ""
}

func letra(i int) string {
	switch i {
	case 0:
		return "a"
	case 1:
		return "b"
	case 2:
		return "c"
	case 3:
		return "d"
	case 4:
		return "e"
	case 5:
		return "f"
	case 6:
		return "g"
	case 7:
		return "h"
	case 8:
		return "i"
	case 9:
		return "j"
	case 10:
		return "k"
	case 11:
		return "l"
	case 12:
		return "m"
	case 13:
		return "n"
	case 14:
		return "ñ"
	case 15:
		return "o"
	case 16:
		return "p"
	case 17:
		return "q"
	case 18:
		return "r"
	case 19:
		return "s"
	case 20:
		return "t"
	case 21:
		return "u"
	case 22:
		return "v"
	case 23:
		return "w"
	case 24:
		return "x"
	case 25:
		return "y"
	case 26:
		return "z"
	default:
		return ""
	}

}
