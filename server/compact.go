package server

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
)

var (
	mu sync.RWMutex
)

func Compaction(mpp *map[string]string) {
	mu.Lock()
	// Closing the append handle before truncating
	// Reopening file specifically for truncation (notice we used the same db varible here)
	// because even after closing the file the varible remains (as it only storing the pointer to the file)

	db, err := os.OpenFile(
		"db.csv",
		os.O_RDWR|os.O_TRUNC,
		0644,
	)

	if err != nil {
		mu.Unlock()
		fmt.Println("error opening db for compact:", err)
		return 
	}

	recordWriter := csv.NewWriter(db)
	for key, value := range *mpp {

		row := []string{
			fmt.Sprint(key),
			fmt.Sprint(value),
		}

		err := recordWriter.Write(row)
		if err != nil {
			mu.Unlock()
			fmt.Println("error while compact writing:", err)
			db.Close()
			return 
		}
	}

	recordWriter.Flush()

	// Close compact handle
	db.Close()

	mu.Unlock()

	fmt.Println("Compaction done")
}
