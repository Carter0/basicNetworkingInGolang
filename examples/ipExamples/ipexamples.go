package main

import (
	"math/big"
	"net"
	"strconv"
)

//IP4toInt -> converts an ip address into a
func IP4toInt(IPv4Addr string) int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(net.ParseIP(IPv4Addr).To4())
	return IPv4Int.Int64()
}

func main() {
	// //Figure out how to convert an ip address into binary.
	// s1 := net.ParseIP("192.0.2.33")
	// if s1 != nil {
	// 	fmt.Printf("%#v\n", s1)
	// }
	// s2 := net.ParseIP("2001:db8:8714:3a90::12")
	// if s2 != nil {
	// 	fmt.Printf("%#v\n", s2)
	// }

	// for _, item := range s2 {
	// 	var i64 int64
	// 	i64 = int64(item)
	// 	println(strconv.FormatInt(i64, 2))
	// }

	//Here is an ip to test
	network := IP4toInt("192.0.2.0")

	//Here is a subnet mask to test.
	subnet := IP4toInt("255.255.255.0")

	//The and of the network to route to and the subnet.
	result := network & subnet
	println(strconv.FormatInt(result, 2))

	//The and of something located the network.
	result2 := IP4toInt("192.0.2.33")
	println(strconv.FormatInt(result2&subnet, 2))

}
