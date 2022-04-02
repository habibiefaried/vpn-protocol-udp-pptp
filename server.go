package main

import (
    "log"
    "github.com/songgao/water"
    "os/exec"
)

const (
    BUFFERSIZE = 1500
    MTU        = "1300"
)



func runIP(args ...string) {
    cmd := exec.Command("/sbin/ip", args...)
    stdoutStderr, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatal("Running command error: ",err)
    }
    log.Printf("%s\n", stdoutStderr)
}

func main(){
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
    runIP("addr", "add", "192.168.9.11/24", "dev", iface.Name())
    runIP("link", "set", "dev", iface.Name(), "up")    
}
