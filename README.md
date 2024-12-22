# CLI To-Do Application in Go

A command-line application for managing to-do lists. Built in Go, this application allows users to create, manage, and manipulate tasks efficiently. It includes features such as reminders, due dates, CSV import/export, search and sort capabilities, and plans for future collaborative features.

---

## Features

### Core Functionalities
- **Add Tasks**: Add new tasks with optional due dates.
- **View Tasks**: Display tasks in a neatly formatted table.
- **Mark Completion**: Toggle the completion status of tasks.
- **Edit Tasks**: Update task titles or due dates.
- **Delete Tasks**: Remove tasks by index.

### Advanced Features
- **CSV Integration**:
  - Export tasks to a CSV file.
  - Import tasks from a CSV file.
- **Reminders**:
  - Set due dates for tasks.
  - Get system notifications for upcoming deadlines using [`beeep`](https://github.com/gen2brain/beeep).
- **Search and Sort**:
  - Search tasks by keywords.
  - Sort tasks by created date, due date, or completion status.
  
---

## Usage

### Installation
1. Clone this repository:
   ```bash
   git clone <repository_url>
   cd <repository_name>
   
2. Install dependencies:
   ```bash
   go mod tidy
   
3. Run the application:
   ```bash
   go run . [flags]
   
### Command-Line Flags

| Flag              | Description                                                     | Example                                                |
|-------------------|-----------------------------------------------------------------|--------------------------------------------------------|
| `--add`           | Add a new task.                                                | `--add "Finish project"`                               |
| `--due`           | Specify a due date for a task (RFC1123 format).                | `--add "Buy groceries" --due "Tue, 26 Dec 2024 18:00 EST"` |
| `--list`          | List all tasks.                                                | `--list`                                               |
| `--edit`          | Edit a task by index.                                          | `--edit 0:Updated_Title`                               |
| `--toggle`        | Toggle completion status of a task by index.                   | `--toggle 1`                                           |
| `--del`           | Delete a task by index.                                        | `--del 2`                                              |
| `--export-csv`    | Export tasks to a specified CSV file.                          | `--export-csv todos.csv`                               |
| `--import-csv`    | Import tasks from a specified CSV file.                        | `--import-csv todos.csv`                               |
| `--search`        | Search tasks by keywords.                                      | `--search "project"`                                   |
| `--sort`          | Sort tasks by a field (`created`, `due`, `completed`).         | `--sort "due"`                                         |


### Future Functionalities
### **Collaboration Features**
- Share and sync Todo lists with multiple users using a REST API.
- Data persistence and storage in a cloud database (PostgreSQL or Firebase).
- RESTful endpoints for managing shared Todo lists.




