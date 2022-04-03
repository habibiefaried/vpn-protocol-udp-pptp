package main

import (
	"flag"
	"fmt"
	"github.com/songgao/water"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"os/exec"
)

const (
	MTU = "1300"
)

var (
	localIP  = flag.String("local", "", "Local tun interface IP/MASK like 192.168.3.3⁄24")
	remoteIP = flag.String("remote", "", "Remote server (external) IP like 8.8.8.8")
	port     = flag.Int("port", 4321, "UDP port for communication")
)

func runIP(args ...string) {
	cmd := exec.Command("/sbin/ip", args...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("Running command error: ", err)
	}
	log.Printf("%s\n", stdoutStderr)
}

func main() {
	flag.Parse()
	if "" == *localIP {
		flag.Usage()
		log.Fatalln("\nlocal ip is not specified")
	}
	if "" == *remoteIP {
		flag.Usage()
		log.Fatalln("\nremote server is not specified")
	}

	iface, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		log.Fatal(err)
	}

	if nil != err {
		log.Fatalln("Unable to allocate TUN interface:", err)
	}

	log.Println("Interface allocated:", iface.Name())
	// set interface parameters
	runIP("link", "set", "dev", iface.Name(), "mtu", MTU)
	runIP("addr", "add", *localIP, "dev", iface.Name())
	runIP("link", "set", "dev", iface.Name(), "up")

	remoteAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%v", *remoteIP, *port))
	if err != nil {
		log.Fatal(err)
	}

	s, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%v", *port))
	if err != nil {
		log.Fatal(err)
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		log.Fatal(err)
	}

	defer connection.Close()
	go func() {
		buf := []byte{}
		for {
			n, addr, err := connection.ReadFromUDP(buf)
			header, _ := ipv4.ParseHeader(buf)
			log.Printf("Received %d bytes from %v: %+v\n", n, addr, header)
			if err != nil || n == 0 {
				log.Println("Error: ", err)
				continue
			}
			// write to TUN interface
			iface.Write(buf)
		}
	}()

	packet := []byte{}
	for {
		plen, err := iface.Read(packet)
		if err != nil {
			break
		}
		header, _ := ipv4.ParseHeader(packet)
		log.Printf("Sending %d bytes to remote: %+v (%+v)\n", plen, header, err)
		connection.WriteToUDP(packet, remoteAddr)
	}
}
