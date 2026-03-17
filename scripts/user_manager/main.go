package main

import (
	"bufio"
	"fmt"
	"logger-app/db"
	"logger-app/storage"
	"os"
	"strings"

	"crypto/rand"
	"encoding/hex"
)

func generateID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func main() {
	db.InitDB()
	reader := bufio.NewReader(os.Stdin)

	if len(os.Args) < 2 {
		interactiveMenu(reader)
		return
	}

	command := os.Args[1]
	switch command {
	case "create":
		name := ""
		if len(os.Args) > 2 {
			name = strings.Join(os.Args[2:], " ")
		} else {
			name = prompt(reader, "Enter user name: ")
		}
		handleCreate(name)
	case "list":
		handleList()
	case "update":
		id := ""
		name := ""
		if len(os.Args) > 2 {
			id = os.Args[2]
		} else {
			id = prompt(reader, "Enter User ID to update: ")
		}
		if len(os.Args) > 3 {
			name = strings.Join(os.Args[3:], " ")
		} else {
			name = prompt(reader, "Enter new name: ")
		}
		handleUpdate(id, name)
	case "delete":
		id := ""
		if len(os.Args) > 2 {
			id = os.Args[2]
		} else {
			id = prompt(reader, "Enter User ID to delete: ")
		}
		handleDelete(id)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}

func interactiveMenu(reader *bufio.Reader) {
	for {
		fmt.Println("\n--- User Manager Interactive Menu ---")
		fmt.Println("1. List Users")
		fmt.Println("2. Create User")
		fmt.Println("3. Update User")
		fmt.Println("4. Delete User")
		fmt.Println("5. Or q to Exit")
		choice := prompt(reader, "Select an option (1-5): ")

		switch choice {
		case "1":
			handleList()
		case "2":
			name := prompt(reader, "Enter user name: ")
			handleCreate(name)
		case "3":
			id := prompt(reader, "Enter User ID to update: ")
			name := prompt(reader, "Enter new name: ")
			handleUpdate(id, name)
		case "4":
			id := prompt(reader, "Enter User ID to delete: ")
			confirm := prompt(reader, fmt.Sprintf("Are you sure you want to delete user %s and all their events? (y/N): ", id))
			if strings.ToLower(confirm) == "y" {
				handleDelete(id)
			} else {
				fmt.Println("Deletion cancelled.")
			}
		case "q":
			fmt.Println("Goodbye!")
			return
		case "5":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}

func prompt(reader *bufio.Reader, message string) string {
	fmt.Print(message)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func handleCreate(name string) {
	if name == "" {
		fmt.Println("Error: Name cannot be empty.")
		return
	}
	id := generateID()
	if err := storage.CreateUser(id, name); err != nil {
		fmt.Printf("Failed to create user: %v\n", err)
		return
	}
	fmt.Printf("User created successfully!\nID: %s\nName: %s\n", id, name)
}

func handleList() {
	users, err := storage.GetAllUsers()
	if err != nil {
		fmt.Printf("Failed to list users: %v\n", err)
		return
	}
	fmt.Println("\nID               | Name")
	fmt.Println("-----------------|-----------------")
	for _, u := range users {
		fmt.Printf("%s | %s\n", u.ID, u.Name)
	}
}

func handleUpdate(id, name string) {
	if id == "" || name == "" {
		fmt.Println("Error: ID and Name are required for update.")
		return
	}
	exists, _ := storage.UserExists(id)
	if !exists {
		fmt.Printf("Error: User with ID %s not found.\n", id)
		return
	}
	if err := storage.UpdateUser(id, name); err != nil {
		fmt.Printf("Failed to update user: %v\n", err)
		return
	}
	fmt.Println("User updated successfully.")
}

func handleDelete(id string) {
	if id == "" {
		fmt.Println("Error: User ID required for deletion.")
		return
	}
	exists, _ := storage.UserExists(id)
	if !exists {
		fmt.Printf("Error: User with ID %s not found.\n", id)
		return
	}
	if err := storage.DeleteUser(id); err != nil {
		fmt.Printf("Failed to delete user: %v\n", err)
		return
	}
	fmt.Println("User and all associated events deleted successfully.")
}

func printUsage() {
	fmt.Println("\nUsage:")
	fmt.Println("  manage-user                (Interactive mode)")
	fmt.Println("  manage-user list           (List all users)")
	fmt.Println("  manage-user create <name>  (Create new user)")
	fmt.Println("  manage-user update <id> <new_name>")
	fmt.Println("  manage-user delete <id>")
}
