package internal

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"gitlab.com/slon/shad-go/gitfame/internal/exceptions"
)

func convertRecordToMap(record [4]string) map[string]interface{} {
	lines, err := strconv.Atoi(record[1])
	exceptions.Exception(err, "convertRecordToMap: could not convert number of lines")

	commits, err := strconv.Atoi(record[2])
	exceptions.Exception(err, "convertRecordToMap: could not convert number of commits")

	files, err := strconv.Atoi(record[3])
	exceptions.Exception(err, "convertRecordToMap: could not convert number of files")

	return map[string]interface{}{
		"name":    record[0],
		"lines":   lines,
		"commits": commits,
		"files":   files,
	}
}

// Tabular

func (stats *Statistics) PrintTabular() {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	printTabularHeader(writer, "Name\tLines\tCommits\tFiles")

	for _, record := range stats.OrderedData {
		printTabularRecord(writer, record[0]+"\t"+record[1]+"\t"+record[2]+"\t"+record[3])
	}

	flushTabularWriter(writer)
}

func printTabularHeader(writer *tabwriter.Writer, header string) {
	_, err := fmt.Fprintln(writer, header)
	exceptions.Exception(err, "printTabularHeader: could not print header")
}

func printTabularRecord(writer *tabwriter.Writer, record string) {
	_, err := fmt.Fprintln(writer, record)
	exceptions.Exception(err, "printTabularRecord: could not print record")
}

func flushTabularWriter(writer *tabwriter.Writer) {
	err := writer.Flush()
	exceptions.Exception(err, "flushTabularWriter: could not flush writer")
}

// CSV

func (stats *Statistics) PrintCSV() {
	writer := csv.NewWriter(os.Stdout)
	writeCSVHeader(writer, []string{"Name", "Lines", "Commits", "Files"})

	for _, record := range stats.OrderedData {
		writeCSVRecord(writer, record[:])
	}

	flushCSVWriter(writer)
}

func writeCSVHeader(writer *csv.Writer, header []string) {
	err := writer.Write(header)
	exceptions.Exception(err, "writeCSVHeader: could not write header")
}

func writeCSVRecord(writer *csv.Writer, record []string) {
	err := writer.Write(record)
	exceptions.Exception(err, "writeCSVRecord: could not write record")
}

func flushCSVWriter(writer *csv.Writer) {
	writer.Flush()
	err := writer.Error()
	exceptions.Exception(err, "flushCSVWriter: could not flush writer")
}

// JSON

func (stats *Statistics) PrintJSON() {
	var records []map[string]interface{}
	for _, record := range stats.OrderedData {
		records = append(records, convertRecordToMap(record))
	}
	printJSONData(records)
}

func printJSONData(records []map[string]interface{}) {
	jsonData, err := json.Marshal(records)
	exceptions.Exception(err, "printJSONData: could not marshal json")

	fmt.Println(string(jsonData))
}

// JSON Lines

func (stats *Statistics) PrintJSONLines() {
	for _, record := range stats.OrderedData {
		jsonLine := convertRecordToMap(record)
		printJSONLineData(jsonLine)
	}
}

func printJSONLineData(record map[string]interface{}) {
	jsonData, err := json.Marshal(record)
	exceptions.Exception(err, "printJSONLineData: could not marshal json")

	fmt.Println(string(jsonData))
}
