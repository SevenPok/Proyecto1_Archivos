package Structs

type Partition struct {
	Status byte
	Type   byte
	Fit    byte
	Start  int64
	Size   int64
	Name   int64
}

type EBR struct {
	Status byte
	Fit    byte
	Start  int64
	Size   int64
	Name   [16]byte
	Next   int64
}

type MBR struct {
	Size      int64
	Date      [20]byte
	Signature int64
	Fit       byte
	Partition [4]Partition
}

type Mkdisk struct {
	Path string
	Size int64
	Unit byte
	Fit  byte
}

type Fdisk struct {
	Size   int64
	Unit   byte
	Path   string
	Type   byte
	Fit    byte
	Delete string
	Name   string
	Add    int64
}