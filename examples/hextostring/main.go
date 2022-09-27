package main

import (
	"fmt"
	"github.com/mokhtarimokhtar/goasterix"
	"github.com/mokhtarimokhtar/goasterix/util"
)

func main() {
	dataSet := []string{
		"01005e03354db4de31c3d7c008004500009600004000031187d5ac1e2e13e003354d20dc2e48008231c0011c007a0000004d006900000fa1000015006af9a0eb3300287434e701477acacc25e8026f0805c8f9a0eb3300287434e9042ef7ca7273e8025e0801dbf980eb3340207434ea035890cb43d6ad4b0508f9a0eb3300287434ed076746ca89d90ac9c30805a0f9a0eb3300287434ef0382eacac092ad63a8080320",
		"01005e03354db4de31c3d7c00800450000650000400003118807ac1e2e12e003354d24402cda0051b50e004600490000004d029600000fa10000080039e8eb19011001171c2a00e8eb19011001181c2c00e8eb190110012a313a00e8eb190110012a313c00e8eb190110022c313e003d413e00",
	}
	for _, data := range dataSet {
		w, _ := goasterix.NewWrapperDataBlock()
		tmp, _ := util.HexStringToByte(data)
		_, err := w.Decode(tmp) // data contains a set of DataBlocks
		if err != nil {
			fmt.Println("ERROR Wrapper: ", err)
		}

		for _, dataBlock := range w.DataBlocks {
			// dataBlock contains one datablock = CAT + LEN + RECORD(S)
			fmt.Printf("\nCategory: %v, Len: %v\n", dataBlock.Category, dataBlock.Len)
			for i, records := range dataBlock.String() {
				// records contains one or more records = N * items
				fmt.Println("Record: ", i+1)
				for _, record := range records {
					fmt.Println(record) // it displays its items in Hexadecimal
				}
			}
		}
	}
}
