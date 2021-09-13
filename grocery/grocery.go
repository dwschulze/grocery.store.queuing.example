package grocery

import (
	"bufio"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)


func Run(config Config) int {

	var registers []Register
	for i:=1; i<=config.NumRegisters; i++ {
		processingTime := 1
		if i == config.NumRegisters {
			processingTime = 2
		}

		register := Register{}
		register.Id = i
		register.ProcessingTime = processingTime
		registers = append(registers, register)
	}

	completion := Completion{
		MaxArrivalTime:  config.MaxArrivalTime,
		Registers:       registers,
		CurrentCustomer: Customer{},
	}

	customers := make([]Customer, len(config.Customers))
	copy(customers, config.Customers)

	//	Sort customers by
	//		ArrivalTime
	//		NumItems
	//		Type (A goes before B)

	sort.Slice(customers, func(p, q int) bool {

		diff := customers[p].ArrivalTime - customers[q].ArrivalTime
		if diff < 0 {
			return true
		} else if diff > 0 {
			return false
		}

		diff = customers[p].NumItems - customers[q].NumItems
		if diff < 0 {
			return true
		} else if diff > 0 {
			return false
		}

		return customers[p].Type == "A"
	})

	t := 0
	for !completion.IsFinished(t) {

		for len(customers) > 0 && customers[0].ArrivalTime == t {

			c := customers[0]
			if c.Type == "A" {
				queueCustomerTypeA(c, registers)
			} else if c.Type == "B" {
				queueCustomerTypeB(c, registers)
			} else {
				log.Fatal("Unknown Customer.Type: " + c.Type)
			}

			customers = customers[1:]
		}

		for i := range registers {
			(&registers[i]).Process()
		}

		t++
	}

	return t
}


func queueCustomerTypeA(c Customer, registers []Register) {

	shortestLine := math.MaxInt32

	for _, r := range registers {
		if r.NumberCustomers() < shortestLine {
			shortestLine = r.NumberCustomers()
		}
	}

	foundRegister := false
	for i, r := range registers {
		if r.NumberCustomers() == shortestLine {
			r.CustomerQueue = append(r.CustomerQueue, c)
			foundRegister = true
			log.Printf("  time: %d, register: %d, customer: %v\n", c.ArrivalTime, r.Id, c)
			registers[i].CustomerQueue = append(registers[i].CustomerQueue, c)
			break
		}
	}

	if !foundRegister {
		log.Fatal("Could not find register")
	}
}


func queueCustomerTypeB(c Customer, registers []Register) {

	minItems := math.MaxInt32
	var register *Register = nil

	for i, r := range registers {

		if r.GetNumberItemsLastCustomerInLine() < 1 {
			register = &registers[i]
			break
		} else {
			if r.GetNumberItemsLastCustomerInLine() < minItems {
				minItems = r.GetNumberItemsLastCustomerInLine()
				register = &registers[i]
			}
		}
	}

	if nil == register {
		log.Fatalln("Could not find register for Customer.Type B")
	}

	register.CustomerQueue = append(register.CustomerQueue, c)
	log.Printf("  time: %d, register: %d, customer: %v\n", c.ArrivalTime, register.Id, c)
}


func ReadInputFiles(args []string) []Config {

	if len(args) < 1 {
		log.Fatal("No input files given")
	}

	var configs []Config
	for _, fileName := range args {

		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		config := Config{}
		maxArrivalTime := 0
		scanner := bufio.NewScanner(file)
		i := 0
		for scanner.Scan() {

			text := scanner.Text()
			if (0 == i) {
				config.NumRegisters, err = strconv.Atoi(text)
				if err != nil {
					log.Fatal(err)
				}
				i++
			} else {
				parts := strings.Split(text, " ")
				if len(parts) != 3 {
					log.Fatal("The line " + text + " from file " + fileName + " must have 3 space separated characters.")
				}
				customer := Customer{}
				customer.Type = parts[0]
				customer.ArrivalTime, err = strconv.Atoi(parts[1])
				if err != nil {
					log.Fatal(err)
				}
				customer.NumItems, err = strconv.Atoi(parts[2])
				if err != nil {
					log.Fatal(err)
				}

				if customer.ArrivalTime > maxArrivalTime {
					maxArrivalTime = customer.ArrivalTime
				}

				config.Customers = append(config.Customers, customer)
			}
		}

		config.MaxArrivalTime = maxArrivalTime
		configs = append(configs, config)

	}

	return configs
}


