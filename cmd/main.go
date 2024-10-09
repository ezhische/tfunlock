package main

import (
	"fmt"
	"log"
	"os"

	dbclient "github.com/ezhische/tfunlock/internal/client"
	"github.com/ezhische/tfunlock/internal/config"
)

func main() {
	// Статические ключи доступа
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	client, err := dbclient.NewClient(cfg.Key, cfg.Secret, cfg.Region, cfg.Endpoint)
	if err != nil {
		log.Fatalf("client creation failed: %v", err)
	}
	switch cfg.Args[0] {
	case "show":
		show(client, cfg.TableName)
	case "showID":
		if len(cfg.Args) < 2 {
			fmt.Fprintln(os.Stderr, "Not enough arguments showID <LockID>")
			return
		}
		getItem(client, cfg.TableName, cfg.Args[1])
	case "delete":
		if len(cfg.Args) < 2 {
			fmt.Fprintln(os.Stderr, "Not enough arguments delete <LockID>")
			return
		}
		delete(client, cfg.TableName, cfg.Args[1])
	case "update":
		if len(cfg.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Not enough arguments update <LockID> <newValue>")
			return
		}
		update(client, cfg.TableName, cfg.Args[1], cfg.Args[2])
	case "":
		fmt.Fprintln(os.Stderr, "Not enough arguments Usage show|showID|delete|update")
	default:
		fmt.Fprintf(os.Stderr, "Unknown argument: %v\nUsage show|showID <LockID>|delete <LockID>|update <LockID> <newValue>\n", cfg.Args[0])

	}
}

func printTable(input map[string]string) {
	for k, v := range input {
		fmt.Println(v, k)
	}
}

func show(client *dbclient.Client, lockTable string) {
	result, err := client.FullScanTable(lockTable)
	if err != nil {
		log.Fatalf("error reading table: %v", err)
	}
	printTable(result)
}
func delete(client *dbclient.Client, lockTable, LockID string) {
	err := client.DeleteItemByLockID(lockTable, LockID)
	if err != nil {
		log.Fatalf("error deleteting: %v", err)
	}
	log.Println("Item deleted successfully!")
}

func update(client *dbclient.Client, lockTable, LockID, newValue string) {
	err := client.UpdateDigest(lockTable, LockID, newValue)
	if err != nil {
		log.Fatalf("error updating: %v", err)
	}
	log.Println("Item updated successfully!")
}

func getItem(client *dbclient.Client, lockTable, LockID string) {
	result, err := client.GetItemByLockID(lockTable, LockID)
	if err != nil {
		log.Fatalf("error reading: %v", err)
	}
	printTable(result)
}
