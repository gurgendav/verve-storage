package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/gurgendav/verve-storage/models"
	"github.com/gurgendav/verve-storage/pkg/gredis"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const FileReSyncInterval = 30 * time.Minute
const PromotionsFileName = "promotions.csv"

func startFileWatcher(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	loadData(PromotionsFileName)
	wg.Done()

	ticker := time.NewTicker(FileReSyncInterval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			wg.Add(1)
			loadData(PromotionsFileName)
			wg.Done()
		}
	}
}

func loadData(filename string) {
	start := time.Now()
	fmt.Printf("Started loading data from file: %s\n", filename)

	err := gredis.DropDatabase()
	if err != nil {
		println("Unable to drop database")
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		price, err := strconv.ParseFloat(record[1], 64)

		r := models.Promotion{
			ID:             record[0],
			Price:          price,
			ExpirationDate: record[2],
		}

		if err = gredis.Set(r.ID, r, 0); err != nil {
			fmt.Printf("Error setting promotion %s: %v\n", r.ID, err)
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Finished loading data from file: %s, in: %s\n", filename, elapsed)
}

func main() {
	// Create a context to manage the shutdown process
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a wait group to ensure no file processing is happening before shutdown
	wg := &sync.WaitGroup{}

	// Register OS signals to handle graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go startFileWatcher(ctx, wg)

	select {
	case sig := <-signalCh:
		fmt.Printf("Received signal: %v. Shutting down...\n", sig)
	case <-ctx.Done():
		fmt.Println("Context canceled. Shutting down...")
	}

	cancel()

	wg.Wait()

	fmt.Println("Shutdown completed. Exiting...")
}
