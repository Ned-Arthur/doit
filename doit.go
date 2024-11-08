package main

import (
	"fmt"
	"os"			// Provides cmdline args
	"encoding/csv"	// CSV interface
	"strconv"		// Converting strings
	"text/tabwriter"// Display tabular data
)

/*
	Command format:
	doit add "Name" -> Add a task
	doit list		-> List incomplete tasks
	"    "    -a	-> List all tasks, regardless of completeness
*/

func main() {
	// Deal with our arguments
	args := os.Args[1:]	// The first element is our path, so ignore it

	if len(args) == 0 {
		// testing stuff here


		writer := tabwriter.NewWriter(os.Stdout, 0, 2, 4, ' ', 0)

		fmt.Fprintln(writer, "1\ttable test\tfalse")
		fmt.Fprintln(writer, "2\tanother test but longer\ttrue")

		writer.Flush()

		return
	}

	switch args[0] {
	case "add":
		add(args[1])
	case "list":
		list()
	default:
		panic("wtf you want bro?")
	}
}

/* Add a task to the todo list */
func add(taskName string) {
	fmt.Println("Adding task", taskName)

	// Open the data file in append mode
	file, err := os.OpenFile("data.csv", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Write to file
	writer := csv.NewWriter(file)
	defer writer.Flush()

	id := 1		// Placeholder value until we can calculate it

	// Assume new tasks are incomplete
	row := []string{strconv.Itoa(id), taskName, "false"}
	err = writer.Write(row)
	if err != nil {
		panic(err)
	}
}

func list() {
	// Open file
	file, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Read the data
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1	// Allow for a variable number of fields
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	// Print said data
	// tabwidth = 2, padding = 4, padded w/ spaces
	writer := tabwriter.NewWriter(os.Stdout, 0, 2, 4, ' ', 0)

	for _, row := range data {
		fmt.Fprintln(writer, row[0] + "\t" + row[1] + "\t" + row[2])
	}
	writer.Flush()

}

