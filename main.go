package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const todoFile = ".godo.json"

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

// Todo represents a single todo item
type Todo struct {
	Text      string    `json:"text"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// TodoList represents the collection of todos
type TodoList struct {
	Todos []Todo `json:"todos"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "add", "a":
		handleAdd()
	case "list", "l":
		handleList()
	case "done", "d":
		handleDone()
	case "remove", "rm":
		handleRemove()
	case "clean", "c":
		handleClean()
	case "show":
		handleShow()
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	usage := `godo - Directory-level todo manager

Usage:
  godo add <text>        Add a new todo (alias: a)
  godo list              List active todos (alias: l)
  godo list --all        List all todos including completed (alias: l --all)
  godo done <id>         Mark todo as complete (alias: d)
  godo remove <id>       Remove an active todo (alias: rm)
  godo clean             Remove all completed todos (alias: c)
  godo clean <id>        Remove a specific completed todo (alias: c)
  godo show              Show active todos (for auto-display)
  godo help              Show this help message

Examples:
  godo a "Implement user authentication"
  godo l
  godo d 1
  godo c`
	fmt.Println(usage)
}

// Load todos from the .godo.json file in current directory
func loadTodos() (*TodoList, error) {
	file, err := os.ReadFile(todoFile)
	if os.IsNotExist(err) {
		return &TodoList{Todos: []Todo{}}, nil
	}
	if err != nil {
		return nil, err
	}

	var todoList TodoList
	if err := json.Unmarshal(file, &todoList); err != nil {
		return nil, err
	}

	return &todoList, nil
}

// Save todos to the .godo.json file in current directory
func saveTodos(todoList *TodoList) error {
	data, err := json.MarshalIndent(todoList, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(todoFile, data, 0644)
}

func handleAdd() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Please provide todo text")
		fmt.Println("Usage: godo add <text>")
		os.Exit(1)
	}

	text := strings.Join(os.Args[2:], " ")

	todoList, err := loadTodos()
	if err != nil {
		fmt.Printf("Error loading todos: %v\n", err)
		os.Exit(1)
	}

	newTodo := Todo{
		Text:      text,
		Completed: false,
		CreatedAt: time.Now(),
	}

	todoList.Todos = append(todoList.Todos, newTodo)

	if err := saveTodos(todoList); err != nil {
		fmt.Printf("Error saving todos: %v\n", err)
		os.Exit(1)
	}

	// Count active todos to show the position
	activeCount := 0
	for _, todo := range todoList.Todos {
		if !todo.Completed {
			activeCount++
		}
	}

	fmt.Printf("%s[Added]%s Todo #%d: %s\n", colorGreen, colorReset, activeCount, newTodo.Text)
}

func handleList() {
	todoList, err := loadTodos()
	if err != nil {
		fmt.Printf("Error loading todos: %v\n", err)
		os.Exit(1)
	}

	showAll := len(os.Args) > 2 && (os.Args[2] == "--all" || os.Args[2] == "-a")

	activeTodos := []Todo{}
	completedTodos := []Todo{}

	for _, todo := range todoList.Todos {
		if todo.Completed {
			completedTodos = append(completedTodos, todo)
		} else {
			activeTodos = append(activeTodos, todo)
		}
	}

	if len(activeTodos) == 0 && len(completedTodos) == 0 {
		fmt.Println("No todos yet. Add one with: godo add <text>")
		return
	}

	if len(activeTodos) > 0 {
		fmt.Printf("%sActive Todos:%s\n", colorCyan, colorReset)
		for i, todo := range activeTodos {
			fmt.Printf("  %s[%d]%s %s\n", colorYellow, i+1, colorReset, todo.Text)
		}
	} else {
		fmt.Printf("%sNo active todos.%s\n", colorGray, colorReset)
	}

	if showAll && len(completedTodos) > 0 {
		fmt.Printf("\n%sCompleted Todos:%s\n", colorGray, colorReset)
		for i, todo := range completedTodos {
			fmt.Printf("  %s[%d]%s %s%s%s\n", colorGray, i+1, colorReset, colorGray, todo.Text, colorReset)
		}
	}
}

func handleDone() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Please provide todo ID")
		fmt.Println("Usage: godo done <id>")
		os.Exit(1)
	}

	position, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Error: Invalid todo ID '%s'\n", os.Args[2])
		os.Exit(1)
	}

	if position < 1 {
		fmt.Println("Error: Todo ID must be 1 or greater")
		os.Exit(1)
	}

	todoList, err := loadTodos()
	if err != nil {
		fmt.Printf("Error loading todos: %v\n", err)
		os.Exit(1)
	}

	// Find the Nth active todo
	activeIndex := 0
	var targetTodo *Todo
	for i := range todoList.Todos {
		if !todoList.Todos[i].Completed {
			activeIndex++
			if activeIndex == position {
				targetTodo = &todoList.Todos[i]
				break
			}
		}
	}

	if targetTodo == nil {
		fmt.Printf("Error: Todo #%d not found\n", position)
		os.Exit(1)
	}

	targetTodo.Completed = true
	fmt.Printf("%s[Completed]%s Todo #%d: %s\n", colorGreen, colorReset, position, targetTodo.Text)

	if err := saveTodos(todoList); err != nil {
		fmt.Printf("Error saving todos: %v\n", err)
		os.Exit(1)
	}
}

func handleRemove() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Please provide todo ID")
		fmt.Println("Usage: godo remove <id>")
		os.Exit(1)
	}

	position, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Error: Invalid todo ID '%s'\n", os.Args[2])
		os.Exit(1)
	}

	if position < 1 {
		fmt.Println("Error: Todo ID must be 1 or greater")
		os.Exit(1)
	}

	todoList, err := loadTodos()
	if err != nil {
		fmt.Printf("Error loading todos: %v\n", err)
		os.Exit(1)
	}

	// Find the Nth active todo
	activeIndex := 0
	targetIndex := -1
	for i := range todoList.Todos {
		if !todoList.Todos[i].Completed {
			activeIndex++
			if activeIndex == position {
				targetIndex = i
				break
			}
		}
	}

	if targetIndex == -1 {
		fmt.Printf("Error: Todo #%d not found\n", position)
		os.Exit(1)
	}

	removedTodo := todoList.Todos[targetIndex]
	todoList.Todos = append(todoList.Todos[:targetIndex], todoList.Todos[targetIndex+1:]...)

	if err := saveTodos(todoList); err != nil {
		fmt.Printf("Error saving todos: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s[Removed]%s Todo #%d: %s\n", colorRed, colorReset, position, removedTodo.Text)
}

func handleClean() {
	todoList, err := loadTodos()
	if err != nil {
		fmt.Printf("Error loading todos: %v\n", err)
		os.Exit(1)
	}

	// If no ID provided, remove all completed todos
	if len(os.Args) < 3 {
		completedCount := 0
		newTodos := []Todo{}
		for _, todo := range todoList.Todos {
			if !todo.Completed {
				newTodos = append(newTodos, todo)
			} else {
				completedCount++
			}
		}

		if completedCount == 0 {
			fmt.Println("No completed todos to clean")
			return
		}

		todoList.Todos = newTodos

		if err := saveTodos(todoList); err != nil {
			fmt.Printf("Error saving todos: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s[Cleaned]%s Removed %d completed todo(s)\n", colorGreen, colorReset, completedCount)
		return
	}

	// Remove specific completed todo by position
	position, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Error: Invalid todo ID '%s'\n", os.Args[2])
		os.Exit(1)
	}

	if position < 1 {
		fmt.Println("Error: Todo ID must be 1 or greater")
		os.Exit(1)
	}

	// Find the Nth completed todo
	completedIndex := 0
	targetIndex := -1
	for i := range todoList.Todos {
		if todoList.Todos[i].Completed {
			completedIndex++
			if completedIndex == position {
				targetIndex = i
				break
			}
		}
	}

	if targetIndex == -1 {
		fmt.Printf("Error: Completed todo #%d not found\n", position)
		os.Exit(1)
	}

	removedTodo := todoList.Todos[targetIndex]
	todoList.Todos = append(todoList.Todos[:targetIndex], todoList.Todos[targetIndex+1:]...)

	if err := saveTodos(todoList); err != nil {
		fmt.Printf("Error saving todos: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s[Removed]%s Completed todo #%d: %s\n", colorRed, colorReset, position, removedTodo.Text)
}

func handleShow() {
	todoList, err := loadTodos()
	if err != nil {
		// Silently fail if there's an error (for auto-display)
		return
	}

	activeTodos := []Todo{}
	for _, todo := range todoList.Todos {
		if !todo.Completed {
			activeTodos = append(activeTodos, todo)
		}
	}

	if len(activeTodos) == 0 {
		return // Don't show anything if no active todos
	}

	fmt.Printf("%sTodos:%s\n", colorBlue, colorReset)
	for i, todo := range activeTodos {
		fmt.Printf("  %s[%d]%s %s\n", colorYellow, i+1, colorReset, todo.Text)
	}
}
