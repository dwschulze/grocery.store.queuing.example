package grocery

type Customer struct {
	Type              string
	ArrivalTime       int
	NumItems          int
	NumProcessedItems int
}

func (c *Customer) Process() {
 	c.NumProcessedItems++
}

type Config struct {
	NumRegisters int
	Customers []Customer
	MaxArrivalTime int
}

type Register struct {
	Id int
	ProcessingTime int
	CustomerQueue []Customer
	CurrentCustomer *Customer
	NumCalls int
}

type Completion struct {
	MaxArrivalTime int
	Registers []Register
	CurrentCustomer Customer
}


func (r Register) NumberCustomers() int {

	n := len(r.CustomerQueue)
	if r.CurrentCustomer != nil {
		n++
	}

	return n
}

func (r Register) IsEmpty() bool {
	b := r.NumberCustomers() < 1
	return b
}

func (r *Register) AddCustomer(c Customer) {

	r.CustomerQueue = append(r.CustomerQueue, c)
}

func (r *Register) GetNumberItemsLastCustomerInLine() int {

	if len(r.CustomerQueue) > 0 {
		lastCustomer := r.CustomerQueue[len(r.CustomerQueue) - 1]
		return lastCustomer.NumItems
	} else if nil != r.CurrentCustomer {
		return r.CurrentCustomer.NumItems
	}

	return 0
}

func (r *Register) Process() {

	if nil == r.CurrentCustomer {
		if len(r.CustomerQueue) > 0 {
			r.CurrentCustomer = &r.CustomerQueue[0]
			if len(r.CustomerQueue) > 1 {
				r.CustomerQueue = r.CustomerQueue[1:]
			} else {
				r.CustomerQueue = nil
			}

			//	reset this when starting to process a new customer
			r.NumCalls = 0
		} else {
			return
		}
	}

	//	Handle training register
	if r.ProcessingTime > 1 {
		r.NumCalls++
		if r.NumCalls % r.ProcessingTime != 0 {
			return
		}
	}

	r.CurrentCustomer.Process()

	if r.CurrentCustomer.NumProcessedItems >= r.CurrentCustomer.NumItems {
		r.CurrentCustomer = nil
	}

}


func (c Completion) IsFinished(time int) bool {

	if time <= c.MaxArrivalTime {
		return false
	}

	for _, register := range c.Registers {

		if !register.IsEmpty() {
			return false
		}
	}

	return true
}

