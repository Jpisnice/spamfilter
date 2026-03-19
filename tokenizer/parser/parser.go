package tokenizer

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"

	logger "github.com/Jpisnice/spamfilter/logger"
)

type Record struct {
	Label    int
	Message  string
	Filename string
}

// Parse reads up to n data rows from the CSV (after the header)
// and returns them as a slice of Record structs.
// If n <= 0, it reads all rows.
func Parse(filename string, n ...int) ([]Record, int) {
	logger.Logger.Print("Parsing file: ", filename, " rows:", n)

	file, err := os.Open(filename)
	if err != nil {
		logger.Logger.Print("Error opening file: ", err)
		return nil, 0
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header row
	header, err := reader.Read()
	if err == io.EOF {
		logger.Logger.Print("Empty CSV file")
		return nil, 0
	}
	if err != nil {
		logger.Logger.Print("Error reading header: ", err)
		return nil, 0
	}
	logger.Logger.Print("Header: ", header)

	// Determine how many rows to read; if not provided, read all (limit <= 0).
	limit := -1
	if len(n) > 0 {
		limit = n[0]
	}

	var records []Record
	rowsRead := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Logger.Print("Error reading record: ", err)
			break
		}

		// Convert first column to int label
		labelInt, err := strconv.Atoi(record[0])
		if err != nil {
			logger.Logger.Print("Error converting label to int: ", err)
			continue
		}

		r := Record{
			Label:    labelInt,
			Message:  record[1],
			Filename: record[2],
		}
		records = append(records, r)

		rowsRead++
		if limit > 0 && rowsRead >= limit {
			break
		}
	}
	return records,len(records)
}

