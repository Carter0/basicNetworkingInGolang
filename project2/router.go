package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/tidwall/gjson"
)

// rtData contains a copy of the original message sent to the router
// as well as necessary data for the routing table at the top level.
type rtData struct {
	relationshipType string
	Msg              []byte
	Conn             net.Conn
}

// networkTuple represents the list of paths for our network.
type networkTuple struct {
	ip      string //ip to send data through
	netMask string //subet mask for the ip
}

// message represents a message from a neighbor to the router.
type message struct {
	Msg  interface{} `json:"msg"`
	Src  string      `json:"src"`
	Dst  string      `json:"dst"`
	Type string      `json:"type"`
}

var routingtable = map[networkTuple][]*rtData{}
var mutex sync.Mutex

//IP4toInt converts an ip address into a binary sequence
func IP4toInt(IPv4Addr string) int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(net.ParseIP(IPv4Addr).To4())
	return IPv4Int.Int64()
}

//CreateRTData creates the routing table information
func createRTData(conn net.Conn, m message, tempType string) rtData {
	message, err := json.Marshal(m.Msg)
	if err != nil {
		panic(err)
	}
	return rtData{tempType, message, conn}
}

func updateLogic(jsonMsg []byte, conn net.Conn, m message) {
	//This right here is the update logic
	tempIP := gjson.GetBytes(jsonMsg, "network")
	tempSubnet := gjson.GetBytes(jsonMsg, "netmask")
	tempTuple := networkTuple{tempIP.String(), tempSubnet.String()}
	tempRoute := createRTData(conn, m, m.Type)
	mutex.Lock()
	if val, ok := routingtable[tempTuple]; ok {
		val = append(val, &tempRoute)
	} else {
		rtArray := []*rtData{&tempRoute} //Create a pointer array and add the tempRoute pointer
		routingtable[tempTuple] = rtArray
	}
	for key, value := range routingtable {
		println("The key is: ")
		println(key.ip)
		println(key.netMask)
		for _, rtdata := range value {
			temp := *rtdata
			println(temp.relationshipType)
		}

	}
	mutex.Unlock()
}

func handleConnection(conn net.Conn) {
	for {
		var mes message

		err := json.NewDecoder(conn).Decode(&mes)
		if err != nil {
			log.Fatal("error decoding message ", err)
		}

		jsonMsg, err := json.Marshal(mes.Msg)
		if err != nil {
			panic(err)
		}

		switch mes.Type {
		case "update":
			//Either add new info to the routing table
			//Or append data to the routing table values.
			println("Adding info to routing table.")
			updateLogic(jsonMsg, conn, mes)
		case "dump":
			fmt.Println("Dump logic here.")
		case "data":
			fmt.Println("Data logic here.")
		}
	}
}

func main() {

	args := os.Args
	ip := make([]string, len(args))
	port := make([]string, len(args))
	for i, network := range args[1:] {
		split := strings.Split(network, "-")
		ip[i] = split[0]
		port[i] = split[1]

		fmt.Println("Socket connection")

		//Open sockets and start listening.
		conn, err := net.Dial("unixpacket", "./"+ip[i])
		if err != nil {
			panic(err)
		}

		fmt.Println("Starting goroutines")
		go handleConnection(conn)
	}

	//TODO make a loop through the hashmap. If statements for the key, value pairs.

	for key, value := range routingtable {
		print("Network name" + key.ip)
		print("Network mask" + key.netMask)
		for _, rtInfo := range value {
			println(rtInfo.relationshipType)
		}
	}

	//TODO, what the hell is this doing lol.
	select {}

	//TODO send update messages to the other neighbors
	//In update message, change the src you are given to the destination
	//you are given. Use the neighbors source to send info over tcp
	//You are supposed to keep a running connection to your neighbors.
}
