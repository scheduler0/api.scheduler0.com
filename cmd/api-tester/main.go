package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	APIKeyHeader    = "x-api-key"
	APISecretHeader = "x-secret-key"
	AccountIDHeader = "x-account-id"
)

type Credential struct {
	ID          uint64    `json:"id"`
	AccountID   uint64    `json:"accountId"`
	APIKey      string    `json:"apiKey"`
	APISecret   string    `json:"apiSecret"`
	DateCreated time.Time `json:"dateCreated"`
}

type CredentialResponse struct {
	Success bool       `json:"success"`
	Data    Credential `json:"data"`
}

type PaginatedCredentialsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Total       int          `json:"total"`
		Offset      int          `json:"offset"`
		Limit       int          `json:"limit"`
		Credentials []Credential `json:"credentials"`
	} `json:"data"`
}

type Project struct {
	ID          uint64    `json:"id"`
	AccountID   uint64    `json:"accountId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DateCreated time.Time `json:"dateCreated"`
}

type ProjectResponse struct {
	Success bool    `json:"success"`
	Data    Project `json:"data"`
}

type PaginatedProjectsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Total    int       `json:"total"`
		Offset   int       `json:"offset"`
		Limit    int       `json:"limit"`
		Projects []Project `json:"projects"`
	} `json:"data"`
}

type Job struct {
	ID           uint64    `json:"id"`
	AccountID    uint64    `json:"accountId"`
	ProjectID    uint64    `json:"projectId"`
	Description  string    `json:"description"`
	ExecutorID   uint64    `json:"executorId"`
	Data         string    `json:"data"`
	Spec         string    `json:"spec"`
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
	Timezone     string    `json:"timezone"`
	DateCreated  time.Time `json:"dateCreated"`
	DateModified time.Time `json:"dateModified"`
	DateDeleted  time.Time `json:"dateDeleted"`
	CreatedBy    string    `json:"createdBy"`
	ModifiedBy   string    `json:"modifiedBy"`
	DeletedBy    string    `json:"deletedBy"`
}

type JobResponse struct {
	Success bool `json:"success"`
	Data    Job  `json:"data"`
}

type PaginatedJobsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Total  int   `json:"total"`
		Offset int   `json:"offset"`
		Limit  int   `json:"limit"`
		Jobs   []Job `json:"jobs"`
	} `json:"data"`
}

type AsyncTask struct {
	ID          uint64    `json:"id"`
	RequestID   string    `json:"requestId"`
	Input       string    `json:"input"`
	Output      string    `json:"output"`
	Service     string    `json:"service"`
	State       int       `json:"state"`
	DateCreated time.Time `json:"dateCreated"`
}

type AsyncTaskResponse struct {
	Success bool      `json:"success"`
	Data    AsyncTask `json:"data"`
}

func main() {
	apiKey := flag.String("api-key", "", "API Key")
	apiSecret := flag.String("api-secret", "", "API Secret")
	host := flag.String("host", "", "API Host (e.g., http://localhost:3002)")
	accountID := flag.String("account-id", "", "Account ID")
	flag.Parse()

	if *apiKey == "" || *apiSecret == "" || *host == "" || *accountID == "" {
		fmt.Println("All flags are required: -api-key, -api-secret, -host, -account-id")
		os.Exit(1)
	}

	// Test Credentials CRUD
	fmt.Println("=== Testing Credentials CRUD ===")
	// Create a credential
	fmt.Println("\nCreating a credential...")
	credential, err := createCredential(*host, *apiKey, *apiSecret, *accountID)
	if err != nil {
		fmt.Printf("Failed to create credential: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created credential with ID: %d\n", credential.ID)

	// Get all credentials
	fmt.Println("\nGetting all credentials...")
	credentials, err := getAllCredentials(*host, *apiKey, *apiSecret, *accountID, 10, 0)
	if err != nil {
		fmt.Printf("Failed to get credentials: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d credentials\n", len(credentials))

	// Update the credential
	fmt.Println("\nUpdating the credential...")
	err = updateCredential(*host, *apiKey, *apiSecret, *accountID, credential.ID)
	if err != nil {
		fmt.Printf("Failed to update credential: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully updated credential")

	// Delete the credential
	fmt.Println("\nDeleting the credential...")
	err = deleteCredential(*host, *apiKey, *apiSecret, *accountID, credential.ID)
	if err != nil {
		fmt.Printf("Failed to delete credential: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully deleted credential")

	// Test Projects CRUD
	fmt.Println("\n=== Testing Projects CRUD ===")
	// Create a project
	fmt.Println("\nCreating a project...")
	project, err := createProject(*host, *apiKey, *apiSecret, *accountID, "Test Project", "This is a test project")
	if err != nil {
		fmt.Printf("Failed to create project: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created project with ID: %s\n", project.ID)

	// Get all projects
	fmt.Println("\nGetting all projects...")
	projects, err := getAllProjects(*host, *apiKey, *apiSecret, *accountID, 10, 0)
	if err != nil {
		fmt.Printf("Failed to get projects: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d projects\n", len(projects))

	// Update the project
	fmt.Println("\nUpdating the project...")
	err = updateProject(*host, *apiKey, *apiSecret, *accountID, project.ID, "Updated Test Project")
	if err != nil {
		fmt.Printf("Failed to update project: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully updated project")

	// Delete the project
	fmt.Println("\nDeleting the project...")
	err = deleteProject(*host, *apiKey, *apiSecret, *accountID, project.ID)
	if err != nil {
		fmt.Printf("Failed to delete project: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully deleted project")

	// Test Jobs CRUD
	fmt.Println("\n=== Testing Jobs CRUD ===")
	// Create a project first (needed for job creation)
	fmt.Println("\nCreating a project for job...")
	project, err = createProject(*host, *apiKey, *apiSecret, *accountID, "Job Test Project", "Project for testing jobs")
	if err != nil {
		fmt.Printf("Failed to create project: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created project with ID: %d\n", project.ID)

	// Create a job
	fmt.Println("\nCreating a job...")
	job, err := createJob(*host, *apiKey, *apiSecret, *accountID, project.ID, "Test Job", "This is a test job", "@every 1m")
	if err != nil {
		fmt.Printf("Failed to create job: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created job with ID: %d\n", job.ID)

	// Get all jobs
	fmt.Println("\nGetting all jobs...")
	jobs, err := getAllJobs(*host, *apiKey, *apiSecret, *accountID, 10, 0)
	if err != nil {
		fmt.Printf("Failed to get jobs: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d jobs\n", len(jobs))

	// Update the job
	fmt.Println("\nUpdating the job...")
	err = updateJob(*host, *apiKey, *apiSecret, *accountID, job.ID, "Updated Test Job")
	if err != nil {
		fmt.Printf("Failed to update job: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully updated job")

	// Delete the job
	fmt.Println("\nDeleting the job...")
	err = deleteJob(*host, *apiKey, *apiSecret, *accountID, job.ID)
	if err != nil {
		fmt.Printf("Failed to delete job: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully deleted job")

	// Clean up - delete the test project
	fmt.Println("\nCleaning up - deleting test project...")
	err = deleteProject(*host, *apiKey, *apiSecret, *accountID, project.ID)
	if err != nil {
		fmt.Printf("Failed to delete test project: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully deleted test project")
}

func createCredential(host, apiKey, apiSecret, accountID string) (*Credential, error) {
	url := fmt.Sprintf("%s/api/v1/credentials", host)
	reqBody := map[string]string{}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(APIKeyHeader, apiKey)
	req.Header.Set(APISecretHeader, apiSecret)
	req.Header.Set(AccountIDHeader, accountID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create credential: %s", string(body))
	}

	var response CredentialResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func getAllCredentials(host, apiKey, apiSecret, accountID string, limit, offset int) ([]Credential, error) {
	url := fmt.Sprintf("%s/api/v1/credentials?limit=%d&offset=%d", host, limit, offset)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(APIKeyHeader, apiKey)
	req.Header.Set(APISecretHeader, apiSecret)
	req.Header.Set(AccountIDHeader, accountID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get credentials: %s", string(body))
	}

	var response PaginatedCredentialsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data.Credentials, nil
}

func updateCredential(host, apiKey, apiSecret, accountID string, credentialID uint64) error {
	url := fmt.Sprintf("%s/api/v1/credentials/%d", host, credentialID)
	reqBody := map[string]bool{
		"archived": true,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(APIKeyHeader, apiKey)
	req.Header.Set(APISecretHeader, apiSecret)
	req.Header.Set(AccountIDHeader, accountID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update credential: %s", string(body))
	}

	return nil
}

func deleteCredential(host, apiKey, apiSecret, accountID string, credentialID uint64) error {
	url := fmt.Sprintf("%s/api/v1/credentials/%d", host, credentialID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set(APIKeyHeader, apiKey)
	req.Header.Set(APISecretHeader, apiSecret)
	req.Header.Set(AccountIDHeader, accountID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete credential: %s", string(body))
	}

	return nil
}

func createProject(host, apiKey, apiSecret, accountID, name, description string) (*Project, error) {
	url := fmt.Sprintf("%s/api/v1/projects", host)
	reqBody := map[string]string{
		"name":        name,
		"description": description,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(APIKeyHeader, apiKey)
	req.Header.Set(APISecretHeader, apiSecret)
	req.Header.Set(AccountIDHeader, accountID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create project: %s", string(body))
	}

	var response ProjectResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func getAllProjects(host, apiKey, apiSecret, accountID string, limit, offset int) ([]Project, error) {
	url := fmt.Sprintf("%s/api/v1/projects?limit=%d&offset=%d", host, limit, offset)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(APIKeyHeader, apiKey)
	req.Header.Set(APISecretHeader, apiSecret)
	req.Header.Set(AccountIDHeader, accountID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get projects: %s", string(body))
	}

	var response PaginatedProjectsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data.Projects, nil
}

func updateProject(host, apiKey, apiSecret, accountID string, projectID uint64, description string) error {
	url := fmt.Sprintf("%s/api/v1/projects/%d", host, projectID)
	reqBody := map[string]string{
		"description": description,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(APIKeyHeader, apiKey)
	req.Header.Set(APISecretHeader, apiSecret)
	req.Header.Set(AccountIDHeader, accountID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update project: %s", string(body))
	}

	return nil
}

func deleteProject(host, apiKey, apiSecret, accountID string, projectID uint64) error {
	url := fmt.Sprintf("%s/api/v1/projects/%d", host, projectID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set(APIKeyHeader, apiKey)
	req.Header.Set(APISecretHeader, apiSecret)
	req.Header.Set(AccountIDHeader, accountID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete project: %s", string(body))
	}

	return nil
}

func createJob(host, apiKey, apiSecret, accountID string, projectID uint64, description, callbackUrl, spec string) (*Job, error) {
	url := fmt.Sprintf("%s/api/v1/jobs", host)
	reqBody := map[string]interface{}{
		"description": description,
		"callbackUrl": callbackUrl,
		"spec":        spec,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(APIKeyHeader, apiKey)
	req.Header.Set(APISecretHeader, apiSecret)
	req.Header.Set(AccountIDHeader, accountID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create job: %s", string(body))
	}

	// Get the Location header for async task status
	location := resp.Header.Get("Location")
	if location == "" {
		return nil, fmt.Errorf("no Location header in response")
	}

	// Extract request ID from Location header
	// Location format: /async-tasks/{requestId}
	parts := strings.Split(location, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid Location header format: %s", location)
	}
	requestID := parts[2]

	// Poll for async task completion
	fmt.Println("Waiting for job creation to complete...")
	for i := 0; i < 30; i++ { // Try for 30 seconds
		time.Sleep(time.Second)
		task, err := getAsyncTask(host, apiKey, apiSecret, accountID, requestID)
		if err != nil {
			return nil, fmt.Errorf("failed to check async task status: %v", err)
		}

		switch task.State {
		case 2: // Success
			var job Job
			if err := json.Unmarshal([]byte(task.Output), &job); err != nil {
				return nil, fmt.Errorf("failed to parse job from async task output: %v", err)
			}
			return &job, nil
		case 3: // Fail
			return nil, fmt.Errorf("job creation failed: %s", task.Output)
		case 0, 1: // Not Started or In Progress
			continue
		default:
			return nil, fmt.Errorf("unknown async task state: %d", task.State)
		}
	}

	return nil, fmt.Errorf("job creation timed out after 30 seconds")
}

func getAsyncTask(host, apiKey, apiSecret, accountID, requestID string) (*AsyncTask, error) {
	url := fmt.Sprintf("%s/api/v1/async-tasks/%s", host, requestID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(APIKeyHeader, apiKey)
	req.Header.Set(APISecretHeader, apiSecret)
	req.Header.Set(AccountIDHeader, accountID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get async task: %s", string(body))
	}

	var response AsyncTaskResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func getAllJobs(host, apiKey, apiSecret, accountID string, limit, offset int) ([]Job, error) {
	url := fmt.Sprintf("%s/api/v1/jobs?limit=%d&offset=%d", host, limit, offset)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(APIKeyHeader, apiKey)
	req.Header.Set(APISecretHeader, apiSecret)
	req.Header.Set(AccountIDHeader, accountID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get jobs: %s", string(body))
	}

	var response PaginatedJobsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data.Jobs, nil
}

func updateJob(host, apiKey, apiSecret, accountID string, jobID uint64, description string) error {
	url := fmt.Sprintf("%s/api/v1/jobs/%d", host, jobID)
	reqBody := map[string]string{
		"description": description,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(APIKeyHeader, apiKey)
	req.Header.Set(APISecretHeader, apiSecret)
	req.Header.Set(AccountIDHeader, accountID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update job: %s", string(body))
	}

	return nil
}

func deleteJob(host, apiKey, apiSecret, accountID string, jobID uint64) error {
	url := fmt.Sprintf("%s/api/v1/jobs/%d", host, jobID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set(APIKeyHeader, apiKey)
	req.Header.Set(APISecretHeader, apiSecret)
	req.Header.Set(AccountIDHeader, accountID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete job: %s", string(body))
	}

	return nil
}
