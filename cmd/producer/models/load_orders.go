package models

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const filepath = "cmd/producer/orders.json"

func OredersLoad() ([]Orders, error) {

	var orders []Orders

	file, err := os.Open(filepath)
	if err != nil {
		return orders, fmt.Errorf("error with file open: %w", err)
	}
	defer file.Close()

	bfile, err := io.ReadAll(file)
	if err != nil {
		return orders, fmt.Errorf("error with read file: %w", err)
	}

	json.Unmarshal(bfile, &orders)

	return orders, nil
}
