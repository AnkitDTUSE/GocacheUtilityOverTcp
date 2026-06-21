package server

import (
	"encoding/csv"
	"fmt"
	"net"
	"os"
)

var DataStorage = make(map[string]string)

func LoadData() {
	mpp := &DataStorage
	db, _ := os.OpenFile("db.csv", os.O_RDONLY|os.O_CREATE, 0644)

	defer db.Close() 

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

func Start(port int, connectionType string) error {
	listener, err := net.Listen(connectionType, fmt.Sprintf(":%v", port))

	if err != nil {
		fmt.Println("error while starting server")
	}

	go LoadData()
	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("error while listening to the address")
			return err
		}

		go CacheUtil(&DataStorage, conn)

	}

}
