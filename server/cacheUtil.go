package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
)

var (
	mutex sync.RWMutex
)

func CacheUtil(mpp *map[string]string, conn net.Conn) {
	defer conn.Close()

	connAddress := conn.RemoteAddr()
	fmt.Printf("[connected] Cleint address: %s\n", connAddress.String())
	fmt.Printf("Connection type: %s\n", connAddress.Network())

	reader := bufio.NewReader(conn)

	for {
		request, err := reader.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				fmt.Printf("[Disconnected] client: %s\n", connAddress.String())
				return
			}
			fmt.Println("error reading request:", err)
			return
		}

		var data map[string]string // to store the unmarshalled data out of the Req JSON

		err = json.Unmarshal([]byte(request), &data)
		if err != nil {
			fmt.Println("error decoding request:", err)
			return
		}

		switch data["cmd"] {
		case "SET":
			if len(data) != 3 {
				fmt.Println("format to set value is SET <KEY> <VALUE>")
				continue
			}

			mutex.Lock()
			(*mpp)[data["key"]] = data["value"]
			mutex.Unlock()
		case "GET":
			mutex.RLock()
			value, ok := (*mpp)[data["key"]]
			if !ok {
				fmt.Fprintln(conn, "enter valid key")
			}
			mutex.RUnlock()
			fmt.Fprintln(conn, value)
		default:
			fmt.Fprintln(conn, "Invalid entry...retry")
		}
	}
}
