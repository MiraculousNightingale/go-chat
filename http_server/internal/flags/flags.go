package flags

import "flag"

const (
	NameAddress    string = "address"
	DefaultAddress string = "0.0.0.0"
	DescAddress    string = "address on which the HTTP server is listening for requests"

	NamePort    string = "port"
	DefaultPort string = "8080"
	DescPort    string = "port on which the HTTP server is listening for requests"
)

type flagsHTTP struct {
	Address string
	Port    string
}

func ParseHTTP() flagsHTTP {
	f := flagsHTTP{}

	flag.StringVar(&f.Address, NameAddress, DefaultAddress, DescAddress)
	flag.StringVar(&f.Port, NamePort, DefaultPort, DescPort)

	flag.Parse()

	return f
}
