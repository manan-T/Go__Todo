package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type CmdFlags struct {
	Add       string
	Del       int
	Edit      string
	Toggle    int
	List      bool
	ExportCSV string
	ImportCSV string
	DueDate   string // New flag to specify due date
	Search    string // New flag to search tasks
	Sort      string // New flag to sort tasks
}

func NewCmdFlags() *CmdFlags {
	cf := CmdFlags{}

	flag.StringVar(&cf.Add, "add", "", "Add a new todo specify title")
	flag.StringVar(&cf.Edit, "edit", "", "Edit a todo by index & specify a new title. id:new_title")
	flag.IntVar(&cf.Del, "del", -1, "Specify a todo by index to delete")
	flag.IntVar(&cf.Toggle, "toggle", -1, "Specify a todo by index to toggle")
	flag.BoolVar(&cf.List, "list", false, "List all todos")
	flag.StringVar(&cf.ExportCSV, "export-csv", "", "Export todos to the specified CSV file")
	flag.StringVar(&cf.ImportCSV, "import-csv", "", "Import todos from the specified CSV file")
	flag.StringVar(&cf.DueDate, "due", "", "Set due date for the todo (format: RFC1123)")
	flag.StringVar(&cf.Search, "search", "", "Search todos by keyword")
	flag.StringVar(&cf.Sort, "sort", "", "Sort todos by 'created', 'due', or 'completed'")

	flag.Parse()

	return &cf
}

func (cf *CmdFlags) Execute(todos *Todos, storage *Storage[Todos]) {
	switch {
	case cf.List:
		todos.print()

	case cf.Add != "":
		// Parse due date if provided
		var dueDate *time.Time
		if cf.DueDate != "" {
			parsedDate, err := time.Parse(time.RFC1123, cf.DueDate)
			if err != nil {
				fmt.Println("Error parsing due date:", err)
				os.Exit(1)
			}
			dueDate = &parsedDate
		}
		todos.add(cf.Add, dueDate)

	case cf.Edit != "":
		parts := strings.SplitN(cf.Edit, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Error, invalid format for edit. Please use id:new_title")
			os.Exit(1)
		}

		index, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Println("Error: invalid index for edit")
			os.Exit(1)
		}

		todos.edit(index, parts[1])

	case cf.Toggle != -1:
		todos.toggle(cf.Toggle)

	case cf.Del != -1:
		todos.delete(cf.Del)

	case cf.ExportCSV != "":
		err := todos.exportToCSV(cf.ExportCSV)
		if err != nil {
			fmt.Printf("Error exporting to CSV: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Todos exported to %s successfully\n", cf.ExportCSV)

	case cf.ImportCSV != "":
		err := todos.importFromCSV(cf.ImportCSV)
		if err != nil {
			fmt.Printf("Error importing from CSV: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Todos imported from %s successfully\n", cf.ImportCSV)

	case cf.Search != "":
		result := todos.search(cf.Search)
		result.print()

	case cf.Sort != "":
		err := todos.sort(cf.Sort)
		if err != nil {
			fmt.Printf("Error sorting todos: %v\n", err)
			os.Exit(1)
		}
		todos.print()

	default:
		fmt.Println("Invalid command")
	}
}
