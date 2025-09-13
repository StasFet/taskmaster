package database

import (
	"os"

	"github.com/supabase-community/supabase-go"
)

type SupabaseClient struct {
	client *supabase.Client
}

func CreateSupabaseClient() (*SupabaseClient, error) {
	SUPABASE_URL := os.Getenv("SUPABASE_URL")
	SUPABASE_KEY := os.Getenv("SUPABASE_KEY")
	client, err := supabase.NewClient(SUPABASE_URL, SUPABASE_KEY, &supabase.ClientOptions{})
	if err != nil {
		return nil, err
	}
	return &SupabaseClient{client}, nil
}

func (c *SupabaseClient) GetClient() *supabase.Client {
	return c.client
}
