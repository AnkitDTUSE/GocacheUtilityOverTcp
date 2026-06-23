package server

import (
	"encoding/csv"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
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

func writeDb() error {
	db, err := os.OpenFile("db.csv", os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("error while writing DB: ", err)
		return err
	}

	recordWriter := csv.NewWriter(db)
	for key, value := range DataStorage {
		record := []string{
			string(key), string(value),
		}
		recordWriter.Write(record)
		recordWriter.Flush()
	}

	return nil
}

func Start(port int, connectionType string) error {

	LoadData() // loading data when server starts

	listener, err := net.Listen(connectionType, fmt.Sprintf(":%v", port))
	if err != nil {
		fmt.Println("error while starting server")
	}

	defer listener.Close()

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	tickerWriteDb := time.NewTicker(5 * time.Second)
	defer tickerWriteDb.Stop()

	tickerCompact := time.NewTicker(11 * time.Second)
	defer tickerCompact.Stop()

	go func() {

		for {
			conn, err := listener.Accept()

			if err != nil {
				fmt.Println("error while listening to the address ", err)
			}

			go CacheUtil(&DataStorage, conn)

		}
	}()

	for {
		select {
		case <-tickerWriteDb.C:
			fmt.Println("auto save functionality")
			writeDb()

		case <-tickerCompact.C:
			fmt.Println("compacting Db")
			Compaction(&DataStorage)
		case <-sigChan:
			fmt.Println("auto save triggered by Interupt")
			writeDb()
			Compaction(&DataStorage)
			fmt.Println("Server stopped safely")
			return nil
		}

	}

}
