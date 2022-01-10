package router

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

// Start starts router.
func Start(r *gin.Engine, port string) {
	log.Println("router.Start: starting on port", port)
	printAddresses(port)
	log.Fatal("router.Start(ERROR):", r.Run(":"+port))
}

// PrintAddresses prints listening addresses.
func printAddresses(port string) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("router.PrintAddresses(ERROR):", err)
		return
	}

	message := "\nNavigate to one of the following addresses:\n"
	for i := range addr {
		ip := net.ParseIP(strings.Split(addr[i].String(), "/")[0])
		if ip != nil && ip.To4() != nil {
			message = message + "\thttp://" + ip.String() + ":" + port + "\n"
		}
	}
	fmt.Println(message)
}
