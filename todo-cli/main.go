package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Todo struct {
	ID          int        `json:"id"`
	Task        string     `json:"task"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type TodoList struct {
	Todos []Todo `json:"todos"`
}

func (tl *TodoList) save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return err
	}

	todoDir := filepath.Join(homeDir, ".todo-cli")
	if err := os.MkdirAll(todoDir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(todoDir, "todos.json"), data, 0644)
}

func loadTodoList() (*TodoList, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	todoFile := filepath.Join(homeDir, ".todo-cli", "todos.json")
	if _, err := os.Stat(todoFile); os.IsNotExist(err) {
		return &TodoList{}, nil
	}

	data, err := os.ReadFile(todoFile)
	if err != nil {
		return nil, err
	}

	var todoList TodoList
	if err := json.Unmarshal(data, &todoList); err != nil {
		return nil, err
	}

	return &todoList, nil
}

func main() {
	todoList, err := loadTodoList()
	if err != nil {
		fmt.Println("Error loading todo list:", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a task description")
			return
		}
		task := strings.Join(os.Args[2:], " ")
		todo := Todo{
			ID:        len(todoList.Todos) + 1,
			Task:      task,
			Completed: false,
			CreatedAt: time.Now(),
		}
		todoList.Todos = append(todoList.Todos, todo)
		if err := todoList.save(); err != nil {
			fmt.Println("Error saving todo:", err)
			return
		}
		fmt.Printf("Added todo: %s\n", task)

	case "list":
		if len(todoList.Todos) == 0 {
			fmt.Println("No todos found")
			return
		}
		for _, todo := range todoList.Todos {
			status := " "
			if todo.Completed {
				status = "✓"
			}
			completedInfo := ""
			if todo.Completed && todo.CompletedAt != nil {
				completedInfo = fmt.Sprintf(" (completed at: %s)", todo.CompletedAt.Format("2006-01-02 15:04:05"))
			}
			fmt.Printf("[%s] %d. %s%s\n", status, todo.ID, todo.Task, completedInfo)
		}

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a todo ID")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid todo ID")
			return
		}
		for i := range todoList.Todos {
			if todoList.Todos[i].ID == id {
				todoList.Todos[i].Completed = true
				completedTime := time.Now()
				todoList.Todos[i].CompletedAt = &completedTime
				if err := todoList.save(); err != nil {
					fmt.Println("Error saving todo:", err)
					return
				}
				fmt.Printf("Marked todo %d as completed at %s\n", id, completedTime.Format("2006-01-02 15:04:05"))
				return
			}
		}
		fmt.Println("Todo not found")

	case "remove":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a todo ID")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid todo ID")
			return
		}
		for i := range todoList.Todos {
			if todoList.Todos[i].ID == id {
				todoList.Todos = append(todoList.Todos[:i], todoList.Todos[i+1:]...)
				if err := todoList.save(); err != nil {
					fmt.Println("Error saving todo:", err)
					return
				}
				fmt.Printf("Removed todo %d\n", id)
				return
			}
		}
		fmt.Println("Todo not found")

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  todo add <task>    - Add a new todo")
	fmt.Println("  todo list          - List all todos")
	fmt.Println("  todo done <id>     - Mark a todo as completed")
	fmt.Println("  todo remove <id>   - Remove a todo")
}
