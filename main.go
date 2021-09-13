package main

import (
	"fmt"
	"grocery/grocery"
	"os"
)

func main() {

	configs := grocery.ReadInputFiles(os.Args[1:])

	i :=0
	for _, config := range configs {

		i++
		fmt.Printf("\n\nExample %d\n%v\n", i, config)
		t := grocery.Run(config)
		fmt.Printf("Finished at time %d\n", t)
	}
}

