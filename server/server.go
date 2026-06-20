package Server

import (
	"encoding/csv"
	"fmt"
	"net"
	"os"
)

var DataStorage = make(map[string]string)

func LoadData(mpp *map[string]string) {
	db, _ := os.OpenFile("db.csv", os.O_RDONLY, 0644)

	defer db.Close() // close the file

	if dbInfo, _ := db.Stat(); dbInfo.Size() != 0 {

		csvReader := csv.NewReader(db)

		dbData, _ := csvReader.ReadAll()

		for _, row := range dbData {
			(*mpp)[row[0]] = row[1]
		}
		fmt.Println("old data retrived")
	} else {
		fmt.Println("db is Empty continue with your writing")
	}

}

func Start(port int, connectionType  string) error {
	

	return nil
}


func main() {
	listener, err := net.Listen("tcp", ":3000")

	if err != nil {
		fmt.Println("error while starting server")
	}

	go LoadData(&DataStorage)
	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("error while listening to the address")
		}

		go CacheUtil(&DataStorage, conn)

	}
}
