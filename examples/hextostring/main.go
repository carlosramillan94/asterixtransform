package main

import (
	"fmt"
	"github.com/mokhtarimokhtar/goasterix"
	"github.com/mokhtarimokhtar/goasterix/util"
)

func main() {
	dataSet := []string{
		"30001ce11f0104eb1a73134b000000dc9ee4521ae377a5a004402304",
		/*// cat 048 + cat 034
		"30 0185 f9d702 0818 3aac7308b59f093309863b75ef345535d60820032404b0928f4420e0ffd70208183aac6da878a60a77020005f06003c20200f93410f1c7566007a107e20eb54020f5ffd70208183aac6ca87ddc095e05a102586003b538f59b34510f0f2c60012c038374cd4020f5ffd70208183aac7aa8c1420f700200059f6003b73950c00464b7d025e0020f07280f6d4020f5ffd70208183aac67a8551808320e2900386001b838b81d18a113620820022201e8cc06402060ffd70208183aac67a834ef08520e2f00076001b93926e4187298160820072e0135b8554020e0f9d70208183aac730892630c26098539514e18754b3e082006f20297238e4020a0ff1608183aac6f483d9a0a6205a90190e0600ac906c10a2efbd340ffd70208183aac77a867d10dde020000556003bf44c1783cf40b620820021603e89ff14020f5ff1608183aac6e483e9a0a2d45a801a0e0580ac5061d0a76f9c040ffd70208183aac7aa8831c0fbe0298003c6001b139d2631885541208200306018703d3402020220014f60818023aac7c10944000460094000000",
		// cat 030 STR
		"1e01e2bfff016008844800e03aace4010e0a0ce01cbcd00e02570256009af969f4007017093434568704e172c32c603fff81604800c03aace4090e0e00ea51fec50e00180018ff6f00d2f4003d2f490408193a069c18a6515a08203ffb81604802303aace4090e0e00e938ffae0e001b001b012c0025f02f4904081939ee7c18a50c2a08203fff01604802383aace4010e0200e9c5ff880e0234023804f3048bf8feb937090c3964e05161b4db5c603ffb01604800423aace4010e0e1fefadccd40e003a003aff49fe8df02f4904396cdc18a10f1608203ffb804805783aace4090e0e00f202c8ed0e0020002000aeffb8f02f490408183fff01604800ac3aace4010e0e2ff37ed5430e0006000afec0ff18f8ff292f49043926e41872981608203ffa4800863aace4010e0e17f357d8010e000800080096ff09f02f49043fff01604807e23aace4010e0e2cf326dbd00e00380037ff2effe8f400472f490438b77d18a1121a08203ffb016048054e3aace4010e0e2bf2d8dce30e002000200051022df02f4904a764e93b5df60c48203ffb01604803c43aace4010e0f10f34cd8d50e00100010feee017af02f29043aac1c18d2843a08203fff01604801a03aace4010e0e2ff35ed8770e001e001cfe32ff45f4007a2f49043932051873511a0820",
		// cat 001 + cat 002
		"01005cf52208329800bb224db58001ee3b10c9813f896e0075068801c60946b0940141d5f0081075229801c90d93b06c015aa530c18155815800752298010f0505b54c01ab9e84818154007507088801c803b9b81c014ab2d40c108202000cf4083202b83aac9722",
		// cat 255
		"ff000ae008833aad7358",
		// cat 048 with bds
		"30 003a fff702 0836 429b52 a0 94c70181 0913 02d0 6002b7 490d01 38a178cf4220 02e79a5d27a00c0060a3280030a4000040 063a 0743ce5b 40 20f5",
		*/}

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
