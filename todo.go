package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aquasecurity/table"
	"github.com/gen2brain/beeep"
)

type Todo struct {
	Title       string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time // can be null
	DueDate     *time.Time // new field for task due date
}

type Todos []Todo

// Add a new Todo
func (todos *Todos) add(title string, dueDate *time.Time) {
	todo := Todo{
		Title:       title,
		Completed:   false,
		CompletedAt: nil,
		CreatedAt:   time.Now(),
		DueDate:     dueDate,
	}

	*todos = append(*todos, todo)

	// If there is a due date, schedule a reminder
	if dueDate != nil {
		go todos.scheduleReminder(len(*todos) - 1) // Schedule reminder for the added task
	}
}

// Search Todos by keyword
func (todos *Todos) search(keyword string) Todos {
	var result Todos
	for _, t := range *todos {
		if strings.Contains(strings.ToLower(t.Title), strings.ToLower(keyword)) {
			result = append(result, t)
		}
	}
	return result
}

// Sort Todos based on a criterion
func (todos *Todos) sort(criterion string) error {
	switch criterion {
	case "created":
		sort.Slice(*todos, func(i, j int) bool {
			return (*todos)[i].CreatedAt.Before((*todos)[j].CreatedAt)
		})
	case "due":
		sort.Slice(*todos, func(i, j int) bool {
			if (*todos)[i].DueDate == nil && (*todos)[j].DueDate == nil {
				return false
			}
			if (*todos)[i].DueDate == nil {
				return false
			}
			if (*todos)[j].DueDate == nil {
				return true
			}
			return (*todos)[i].DueDate.Before(*(*todos)[j].DueDate)
		})
	case "completed":
		sort.Slice(*todos, func(i, j int) bool {
			return !(*todos)[i].Completed && (*todos)[j].Completed
		})
	default:
		return errors.New("invalid sort criterion: use 'created', 'due', or 'completed'")
	}
	return nil
}

// Function to schedule a reminder
func (todos *Todos) scheduleReminder(index int) {
	t := *todos

	if err := t.validateIndex(index); err != nil {
		return
	}

	// Check if the task has a due date and if it is due
	dueDate := t[index].DueDate
	if dueDate == nil || dueDate.After(time.Now()) {
		return
	}

	// If the due date has passed or is today, show the notification
	err := beeep.Notify("Reminder: "+t[index].Title, "Your task is due!", "")
	if err != nil {
		fmt.Println("Error sending notification:", err)
	}
}

// Check operations like edit, remove, or toggle
func (todos *Todos) validateIndex(index int) error {
	if index < 0 || index >= len(*todos) {
		err := errors.New("invalid index")
		fmt.Println(err)
		return err
	}

	return nil
}

// Delete a Todo by index
func (todos *Todos) delete(index int) error {
	t := *todos

	if err := t.validateIndex(index); err != nil {
		return err
	}

	*todos = append(t[:index], t[index+1:]...)
	return nil
}

// Toggle the completion state of a Todo
func (todos *Todos) toggle(index int) error {
	t := *todos

	if err := t.validateIndex(index); err != nil {
		return err
	}

	isCompleted := t[index].Completed

	if !isCompleted {
		completionTime := time.Now()
		t[index].CompletedAt = &completionTime
	}

	t[index].Completed = !isCompleted

	return nil
}

// Edit a Todo's title
func (todos *Todos) edit(index int, title string) error {
	t := *todos

	if err := t.validateIndex(index); err != nil {
		return err
	}

	t[index].Title = title
	return nil
}

// Print all Todos in a tabular format
func (todos *Todos) print() {
	table := table.New(os.Stdout)
	table.SetRowLines(false)
	table.SetHeaders("#", "Title", "Completed", "Created At", "Due Date")
	for index, t := range *todos {
		completed := "❌"
		// completedAt := ""
		dueDate := ""

		if t.Completed {
			completed = "✅"
			// if t.CompletedAt != nil {
			// 	completedAt = t.CompletedAt.Format(time.RFC1123)
			// }
		}

		if t.DueDate != nil {
			dueDate = t.DueDate.Format(time.RFC1123)
		}

		table.AddRow(strconv.Itoa(index), t.Title, completed, t.CreatedAt.Format(time.RFC1123), dueDate)
	}

	table.Render()
}

// Export Todos to a CSV file
func (todos *Todos) exportToCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	err = writer.Write([]string{"Title", "Completed", "CreatedAt", "CompletedAt"})
	if err != nil {
		return err
	}

	// Write Todo data
	for _, t := range *todos {
		completedAt := ""
		if t.CompletedAt != nil {
			completedAt = t.CompletedAt.Format(time.RFC3339)
		}
		record := []string{
			t.Title,
			strconv.FormatBool(t.Completed),
			t.CreatedAt.Format(time.RFC3339),
			completedAt,
		}
		err := writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}

// Import Todos from a CSV file
func (todos *Todos) importFromCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header row
	_, err = reader.Read()
	if err != nil {
		return err
	}

	var importedTodos Todos

	// Read each record and create a Todo
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		completed, err := strconv.ParseBool(record[1])
		if err != nil {
			return err
		}

		createdAt, err := time.Parse(time.RFC3339, record[2])
		if err != nil {
			return err
		}

		var completedAt *time.Time
		if record[3] != "" {
			parsedCompletedAt, err := time.Parse(time.RFC3339, record[3])
			if err != nil {
				return err
			}
			completedAt = &parsedCompletedAt
		}

		todo := Todo{
			Title:       record[0],
			Completed:   completed,
			CreatedAt:   createdAt,
			CompletedAt: completedAt,
		}

		importedTodos = append(importedTodos, todo)
	}

	*todos = importedTodos
	return nil
}
