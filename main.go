package main

import (
	"fmt"
	"log"
	i "taskmaster/internal"
)

func main() {
	client, err := i.CreateSupabaseClient()
	if err != nil {
		log.Fatalf("Error creating supabase client: %v", err)
	}
	user, err := client.GetUserById(1)
	if err != nil {
		log.Fatalf("Error getting all users: %v", err)
	}
	fmt.Println(*user)
}
