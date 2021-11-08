package main

import (
	"fmt"
	"github.com/mokhtarimokhtar/goasterix"
	"github.com/mokhtarimokhtar/goasterix/uap"
)

func main() {
	dataSet := [][]byte{
		// cat 048 + cat 034
		goasterix.HexStringToByte("30 0185 f9d702 0818 3aac7308b59f093309863b75ef345535d60820032404b0928f4420e0ffd70208183aac6da878a60a77020005f06003c20200f93410f1c7566007a107e20eb54020f5ffd70208183aac6ca87ddc095e05a102586003b538f59b34510f0f2c60012c038374cd4020f5ffd70208183aac7aa8c1420f700200059f6003b73950c00464b7d025e0020f07280f6d4020f5ffd70208183aac67a8551808320e2900386001b838b81d18a113620820022201e8cc06402060ffd70208183aac67a834ef08520e2f00076001b93926e4187298160820072e0135b8554020e0f9d70208183aac730892630c26098539514e18754b3e082006f20297238e4020a0ff1608183aac6f483d9a0a6205a90190e0600ac906c10a2efbd340ffd70208183aac77a867d10dde020000556003bf44c1783cf40b620820021603e89ff14020f5ff1608183aac6e483e9a0a2d45a801a0e0580ac5061d0a76f9c040ffd70208183aac7aa8831c0fbe0298003c6001b139d2631885541208200306018703d3402020220014f60818023aac7c10944000460094000000"),
		// cat 030 ARTAS
		goasterix.HexStringToByte("1e00f3afbbf317f1300883040070a8bcf3ff07070723f0a8800713feb7022b0389038b140704012c080811580000001e7004f04aa004b0012400544e49413531313206c84c45424c48454c584d413332300101a5389075c71ca0afbbf317f130088304002aa8bcf3ff04040447fda703f7d2008f0df705280528140700000008171158000000087002f0c3c00528012d006955414c3931202007314c4c42474b4557524842373757a290f3541339c60820afbbf31101300883040335a8bcf3ff0b0b0b2be9a9b5fffefffa0fff08c008c01d0e070000001484115800000200700400ffffffffffffffff344045df7df76021d3"),
		// cat 001 + cat 002
		goasterix.HexStringToByte("01005cf52208329800bb224db58001ee3b10c9813f896e0075068801c60946b0940141d5f0081075229801c90d93b06c015aa530c18155815800752298010f0505b54c01ab9e84818154007507088801c803b9b81c014ab2d40c108202000cf4083202b83aac9722"),
		// cat 255
		goasterix.HexStringToByte("ff000ae008833aad7358"),
		// cat 048 with bds
		goasterix.HexStringToByte("30 003a fff702 0836 429b52 a0 94c70181 0913 02d0 6002b7 490d01 38a178cf4220 02e79a5d27a00c0060a3280030a4000040 063a 0743ce5b 40 20f5"),
	}

	// change User Application Profile STR by UAP ARTAS V7.0
	uap.Profiles[30] = uap.Cat030ArtasV62

	for _, data := range dataSet {
		w,_ := goasterix.NewWrapperDataBlock()
		_, err := w.Decode(data)  // data contains a set of DataBlocks
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
					fmt.Println(record)  // it displays its items in Hexadecimal
				}
			}
		}
	}
}
