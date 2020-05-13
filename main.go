package main

import (
	"bufio"
	"fmt"
	"github.com/google/gopacket/pcap"
	cmap "github.com/orcaman/concurrent-map"
	"log"
	"net"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type Interface struct {
	Display string
	Iface   net.Interface
}

func windowsInterfaces() []Interface {
	pcapInterfaces, _ := pcap.FindAllDevs()

	interfaces := make([]Interface, len(pcapInterfaces))
	for i, iface := range pcapInterfaces {
		interfaces[i] = Interface{
			Display: iface.Description + " => " + iface.Name,
			Iface: net.Interface{
				Index:        100 + i,
				MTU:          1500,
				Name:         iface.Name,
				HardwareAddr: nil,
				Flags:        net.Flags(iface.Flags),
			},
		}
	}

	return interfaces
}

func findInterfaces() []Interface {
	if runtime.GOOS == "windows" {
		return windowsInterfaces()
	}

	netInterfaces, _ := net.Interfaces()

	interfaces := make([]Interface, len(netInterfaces))
	for i, iface := range netInterfaces {
		interfaces[i] = Interface{
			Display: iface.Name,
			Iface:   iface,
		}
	}

	return interfaces
}

func main() {
	interfaces := findInterfaces()

	interfacesByIndex := make(map[int]Interface)
	index := 0

	for _, iface := range interfaces {
		interfacesByIndex[index] = iface
		index = index + 1
	}

	printInterfaces(interfacesByIndex)

	ifaceI := selectInterface(interfacesByIndex, "Please select first interface to be bridged")
	ifaceII := selectInterface(interfacesByIndex, "Please select second interface to be bridged")

	if (ifaceI.HardwareAddr != nil && ifaceI.HardwareAddr.String() == ifaceII.HardwareAddr.String()) || (ifaceI.HardwareAddr == nil && ifaceI.Name == ifaceII.Name) {
		log.Fatal("You need to select different interfaces to bridge")
	}

	macConnections := cmap.New()

	go Bridge(ifaceI, ifaceII, &macConnections)
	Bridge(ifaceII, ifaceI, &macConnections)
}

func printInterfaces(interfacesByIndex map[int]Interface) {
	keys := make([]int, 0)
	for key := range interfacesByIndex {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	for _, key := range keys {
		fmt.Printf("[ %d ] %s => %s\n", key, interfacesByIndex[key].Display, interfacesByIndex[key].Iface.HardwareAddr)
	}
}

func selectInterface(interfacesByIndex map[int]Interface, message string) *net.Interface {

	for {
		input := askUserInput(message)

		key, err := strconv.Atoi(input)

		if err != nil {
			key = -1
		}

		iface, ok := interfacesByIndex[key]

		if ok {
			return &iface.Iface
		} else {
			fmt.Println("Wrong selection!")
		}
	}
}

func askUserInput(message string) string {
	fmt.Println(message)

	reader := bufio.NewReader(os.Stdin)
	sentence, err := reader.ReadString('\n')

	if err != nil {
		log.Fatalf("Error occurred while reading user input: %s", err.Error())
	}

	sentence = strings.Replace(sentence, "\r\n", "", -1)

	return strings.TrimSuffix(sentence, "\n")
}
