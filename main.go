package main

import (
	"fmt"
	"os"
)

func main() {
	// Initialize todos and storage
	todos := Todos{}
	storage := NewStorage[Todos]("todos.json")
	storage.Load(&todos)

	// Parse command-line flags
	cmdFlags := NewCmdFlags()

	// Check for export and import commands
	if cmdFlags.ExportCSV != "" {
		// Export todos to CSV file
		err := todos.exportToCSV(cmdFlags.ExportCSV)
		if err != nil {
			fmt.Println("Error exporting to CSV:", err)
			os.Exit(1)
		}
		fmt.Println("Todos exported to", cmdFlags.ExportCSV)
		return
	}

	if cmdFlags.ImportCSV != "" {
		// Import todos from CSV file
		err := todos.importFromCSV(cmdFlags.ImportCSV)
		if err != nil {
			fmt.Println("Error importing from CSV:", err)
			os.Exit(1)
		}
		fmt.Println("Todos imported from", cmdFlags.ImportCSV)
	}

	// Execute other commands (add, edit, delete, etc.)
	cmdFlags.Execute(&todos, storage)

	// Save updated todos to storage (JSON)
	storage.Save(todos)
}
