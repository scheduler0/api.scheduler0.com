package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	APIKeyHeader    = "x-api-key"
	APISecretHeader = "x-secret-key"
	AccountIDHeader = "x-account-id"
)

var (
	apiKey    *string
	apiSecret *string
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}

	// Set header names from environment variables with defaults
	apiKeyEnv := getEnvOrDefault("API_KEY", "x-api-key")
	apiSecretEnv := getEnvOrDefault("API_SECRET", "x-secret-key")

	fmt.Println("apiKeyEnv", apiKeyEnv)
	fmt.Println("apiSecretEnv", apiSecretEnv)

	apiKey = &apiKeyEnv
	apiSecret = &apiSecretEnv
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	host := flag.String("host", "", "API Host (e.g., http://localhost:3002)")
	accountID := flag.String("account-id", "", "Account ID")
	flag.Parse()

	if *apiKey == "" || *apiSecret == "" || *host == "" || *accountID == "" {
		fmt.Println("All flags are required: -api-key, -api-secret, -host, -account-id")
		os.Exit(1)
	}

	// Test Projects CRUD
	fmt.Println("=== Testing Projects CRUD ===")
	// Create a project
	fmt.Println("\nCreating a project...")
	project, err := CreateProject(*host, *apiKey, *apiSecret, *accountID, "Test Project", "This is a test project")
	if err != nil {
		fmt.Printf("Error creating project: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created project: %+v\n", project)

	// Get all projects
	fmt.Println("\nGetting all projects...")
	projects, err := GetAllProjects(*host, *apiKey, *apiSecret, *accountID, 10, 0)
	if err != nil {
		fmt.Printf("Error getting projects: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d projects\n", len(projects))

	// Update project
	fmt.Println("\nUpdating project...")
	projectID, _ := strconv.ParseInt(project.ID, 10, 64)
	err = UpdateProject(*host, *apiKey, *apiSecret, *accountID, projectID, "Updated description")
	if err != nil {
		fmt.Printf("Error updating project: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Project updated successfully")

	// Test Executors CRUD
	fmt.Println("\n=== Testing Executors CRUD ===")
	// Create an executor
	fmt.Println("\nCreating an executor...")
	executor := ExecutorCreateRequest{
		Name:           "Test Webhook Executor",
		Type:           "webhook_url",
		WebhookURL:     "http://localhost:3000/webhook",
		WebhookMethod:  "POST",
		WebhookHeaders: "Content-Type: application/json",
	}
	createdExecutor, err := CreateExecutor(*host, *apiKey, *apiSecret, *accountID, executor)
	if err != nil {
		fmt.Printf("Error creating executor: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created executor: %+v\n", createdExecutor)

	// Get all executors
	fmt.Println("\nGetting all executors...")
	executors, err := GetAllExecutors(*host, *apiKey, *apiSecret, *accountID, 10, 0)
	if err != nil {
		fmt.Printf("Error getting executors: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d executors\n", len(executors))

	// Get single executor
	fmt.Println("\nGetting single executor...")
	singleExecutor, err := GetExecutor(*host, *apiKey, *apiSecret, *accountID, fmt.Sprintf("%d", createdExecutor.ID))
	if err != nil {
		fmt.Printf("Error getting executor: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Retrieved executor: %+v\n", singleExecutor)

	// Update executor
	fmt.Println("\nUpdating executor...")
	updateRequest := ExecutorUpdateRequest{
		Name:           "Updated Webhook Executor",
		Type:           "webhook_url",
		WebhookURL:     "http://localhost:3000/webhook/updated",
		WebhookMethod:  "POST",
		WebhookHeaders: "Content-Type: application/json",
	}
	err = UpdateExecutor(*host, *apiKey, *apiSecret, *accountID, fmt.Sprintf("%d", createdExecutor.ID), updateRequest)
	if err != nil {
		fmt.Printf("Error updating executor: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Executor updated successfully")

	// Test Jobs CRUD
	fmt.Println("\n=== Testing Jobs CRUD ===")
	// Create jobs
	fmt.Println("\nCreating jobs...")
	jobs := []JobCreateRequest{
		{
			ProjectID:  projectID,
			Data:       "Test Job 1",
			Spec:       "*/5 * * * *",
			StartDate:  "2024-01-01T00:00:00Z",
			EndDate:    "2024-12-31T23:59:59Z",
			Timezone:   "UTC",
			ExecutorID: &createdExecutor.ID,
		},
		{
			ProjectID:  projectID,
			Data:       "Test Job 2",
			Spec:       "*/10 * * * *",
			StartDate:  "2024-01-01T00:00:00Z",
			EndDate:    "2024-12-31T23:59:59Z",
			Timezone:   "UTC",
			ExecutorID: &createdExecutor.ID,
		},
	}
	createdJobs, err := CreateJobs(*host, *apiKey, *apiSecret, *accountID, jobs)
	if err != nil {
		fmt.Printf("Error creating jobs: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created %d jobs\n", len(createdJobs))

	// Get all jobs
	fmt.Println("\nGetting all jobs...")
	allJobs, err := GetAllJobs(*host, *apiKey, *apiSecret, *accountID, projectID, 10, 0)
	if err != nil {
		fmt.Printf("Error getting jobs: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d jobs\n", len(allJobs))

	// Update job
	fmt.Println("\nUpdating job...")
	err = UpdateJob(*host, *apiKey, *apiSecret, *accountID, createdJobs[0].ID, "Updated job description")
	if err != nil {
		fmt.Printf("Error updating job: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Job updated successfully")

	// Delete job
	fmt.Println("\nDeleting job...")
	err = DeleteJob(*host, *apiKey, *apiSecret, *accountID, createdJobs[0].ID)
	if err != nil {
		fmt.Printf("Error deleting job: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Job deleted successfully")

	// Delete executor
	fmt.Println("\nDeleting executor...")
	err = DeleteExecutor(*host, *apiKey, *apiSecret, *accountID, fmt.Sprintf("%d", createdExecutor.ID))
	if err != nil {
		fmt.Printf("Error deleting executor: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Executor deleted successfully")

	// Delete job 2
	fmt.Println("\nDeleting job 2...")
	err = DeleteJob(*host, *apiKey, *apiSecret, *accountID, createdJobs[1].ID)
	if err != nil {
		fmt.Printf("Error deleting job 2: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Job 2 deleted successfully")

	// Delete project
	fmt.Println("\nDeleting project...")
	err = DeleteProject(*host, *apiKey, *apiSecret, *accountID, projectID)
	if err != nil {
		fmt.Printf("Error deleting project: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Project deleted successfully")
}
