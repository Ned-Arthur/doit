package main

import (
	"fmt"
	"os"			// Provides cmdline args and files
	"encoding/csv"	// CSV interface
	"strconv"		// Converting strings
	"text/tabwriter"// Display tabular data
	"time"			// Datetime operations
	"path/filepath" // String to filepath

	"github.com/mergestat/timediff"	// Human-readable time differences
)

/* CONSTANTS */
var HEADERS []string = []string{"ID", "Description", "Created", "Completed"}
const TIMEFORMAT = "2006-01-02 15:04:05 -0700 MST"
const CSVFILENAME = "doit.csv"
var CSVFILE string

/* HELPER FUNCTIONS */

/* Read in our data as a [][]string */
func readData() [][]string { // Open file
	file, err := os.Open(CSVFILE)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Read the data
	reader := csv.NewReader(file)
	//reader.FieldsPerRecord = -1	// Allow for a variable number of fields
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	return data
}

func displayTable(data [][]string) {
	// Display a table with specific rows
	// Automatically displays header as well

	// tabwidth = 2, padding = 4, padded w/ spaces
	writer := tabwriter.NewWriter(os.Stdout, 0, 2, 4, ' ', 0)

	fmt.Fprintln(writer, HEADERS[0] +"\t"+ HEADERS[1] +"\t"+ HEADERS[2] +"\t"+ HEADERS[3])

	for _, row := range data {
		timeValue, _ := time.Parse(TIMEFORMAT, row[2])
		timeSince := timediff.TimeDiff(timeValue)

		fmt.Fprintln(writer, row[0] +"\t"+ row[1] +"\t"+ timeSince +"\t"+ row[3])
	}
	writer.Flush()
}

/*
	Command format:
	doit add <"Name">		-> Add a task
	doit list				-> List incomplete tasks
	"    "    -a			-> List all tasks, regardless of completeness
	doit complete <taskId>	-> Change the completed flag of the indicated task
	doit delete <taskId>	-> Remove specified task from the list
	"    "      -c			-> Remove all tasks marked as complete=true		*
	doit init				-> Create the CSV file with headers
*/

/* ENTRY POINT */
func main() {
	// Setup our file location
	binPath, err := os.Executable()
	if err != nil {
		fmt.Println("!! Error getting executable path")
		panic(err)
	}
	binDir := filepath.Dir(binPath)
	CSVFILE = filepath.Join(binDir, CSVFILENAME)


	// Deal with our arguments
	args := os.Args[1:]	// The first element is our path, so ignore it

	if len(args) == 0 {
		fmt.Println("Doit is a todo list CLI app.")
		fmt.Println("Usage:")
		fmt.Println("        doit <command> [arguments]")
		fmt.Println("\nUse the -h or --help flag for more info")

		return
	}

	switch args[0] {
	case "add":
		add(args[1])
	case "list":
		if len(args) > 1 {
			list(args[1])
		} else {
			list("")
		}
	case "complete":
		complete(args[1])
	case "delete":
		deleteTask(args[1])
	case "init":
		initData()
	case "-h":
		help()
	case "--help":
		help()
	default:
		fmt.Println("Invalid command. Use -h or --help to see valid commands.")
	}
}

/* OPERATION FUNCTIONS */

/* Add a task to the todo list */ func add(taskName string) {
	// Get existing records and find largest id
	data := readData()
	maxId := 0
	for _, row := range data {
		rowId, _ := strconv.Atoi(row[0])
		if rowId > maxId {
			maxId = rowId
		}
	}

	// Open the data file in append mode
	file, err := os.OpenFile(CSVFILE, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Write to file
	writer := csv.NewWriter(file)
	defer writer.Flush()

	id := maxId + 1		// Placeholder value until we can calculate it
	currentTime := time.Now()

	// Assume new tasks are incomplete
	row := []string{strconv.Itoa(id), taskName, currentTime.Format(TIMEFORMAT), "false"}
	err = writer.Write(row)
	if err != nil {
		panic(err)
	}
}

func list(flags string) {
	data := readData()

	if len(data) == 1 {
		fmt.Println("No todos saved; nothing to display")
		fmt.Println("Add a todo with `doit add \"Task name\"`")
		return
	}

	// The function handles the headings for us
	displayTable(data[1:])
}

func complete(taskId string) {
	data := readData()
	var rowToUpdate int
	for i, row := range data {
		if row[0] == taskId {
			rowToUpdate = i
		}
	}

	if rowToUpdate == 0 {
		panic("ID doesn't match a task")
	}

	data[rowToUpdate][3] = "true"

	file, err := os.Create(CSVFILE)
	if err!= nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	err = writer.WriteAll(data)
	if err!= nil {
		panic(err)
	}

	fmt.Println("Successfully completed a task. Good Job!")
	displayTable([][]string{data[rowToUpdate]})
}

func deleteTask(arg string) {
	data := readData()

	var rowToDelete int
	for i, row := range data {
		if row[0] == arg {
			rowToDelete = i
		}
	}

	if rowToDelete == 0 {
		panic("No rows to delete")
		return
	}

	removeData := data[rowToDelete]

	data = append(data[:rowToDelete], data[rowToDelete+1:]...)

	file, err := os.Create(CSVFILE)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	err = writer.WriteAll(data)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully deleted row:")
	displayTable([][]string{removeData})
}

func initData() {
	file, err := os.Create(CSVFILE)
	if err != nil {
		panic(err)
	}
	defer file.Close()


	writer := csv.NewWriter(file)
	defer writer.Flush()
	err = writer.Write(HEADERS)
	if err != nil {
		panic(err)
	}
}

func help() {
	fmt.Println("Help for doit")
	helpString := `
Command format:
	doit add <"Name">       -> Add a task
	doit list               -> List incomplete tasks
	"    "    -a            -> List all tasks, regardless of completeness
	doit complete <taskId>  -> Change the completed flag of the indicated task
	doit delete <taskId>    -> Remove specified task from the list
	doit init               -> Create the CSV file with headers
	`
	fmt.Println(helpString)
}

