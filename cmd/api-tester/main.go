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
	"testing"
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

type JobCreateRequest struct {
	ProjectID   uint64 `json:"projectId"`
	Description string `json:"description"`
	CallbackUrl string `json:"callbackUrl"`
	Spec        string `json:"spec"`
	Timezone    string `json:"timezone"`
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

	// Create multiple jobs
	fmt.Println("\nCreating multiple jobs...")
	jobs, err := createJobs(*host, *apiKey, *apiSecret, *accountID, []JobCreateRequest{
		{
			ProjectID:   project.ID,
			Description: "Test Job 1",
			CallbackUrl: "This is test job 1",
			Spec:        "@every 1m",
			Timezone:    "UTC",
		},
		{
			ProjectID:   project.ID,
			Description: "Test Job 2",
			CallbackUrl: "This is test job 2",
			Spec:        "@every 2m",
			Timezone:    "UTC",
		},
	})
	if err != nil {
		fmt.Printf("Failed to create jobs: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created %d jobs\n", len(jobs))

	// Get all jobs
	fmt.Println("\nGetting all jobs...")
	allJobs, err := getAllJobs(*host, *apiKey, *apiSecret, *accountID, project.ID, 10, 0)
	if err != nil {
		fmt.Printf("Failed to get jobs: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d jobs\n", len(allJobs))

	// Update the first job
	fmt.Println("\nUpdating the first job...")
	err = updateJob(*host, *apiKey, *apiSecret, *accountID, jobs[0].ID, "Updated Test Job 1")
	if err != nil {
		fmt.Printf("Failed to update job: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully updated job")

	// Delete all jobs
	fmt.Println("\nDeleting all jobs...")
	for _, job := range jobs {
		err = deleteJob(*host, *apiKey, *apiSecret, *accountID, job.ID)
		if err != nil {
			fmt.Printf("Failed to delete job %d: %v\n", job.ID, err)
			os.Exit(1)
		}
		fmt.Printf("Successfully deleted job %d\n", job.ID)
	}

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

func createJobs(host, apiKey, apiSecret, accountID string, jobs []JobCreateRequest) ([]Job, error) {
	url := fmt.Sprintf("%s/api/v1/jobs", host)
	jsonBody, err := json.Marshal(jobs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal jobs: %v", err)
	}

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

	if resp.StatusCode != http.StatusAccepted {
		fmt.Println(string(body))
		return nil, fmt.Errorf("failed to create jobs: %s", string(body))
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
	fmt.Println("Waiting for jobs creation to complete...")
	for i := 0; i < 30; i++ { // Try for 30 seconds
		time.Sleep(time.Second)
		task, err := getAsyncTask(host, apiKey, apiSecret, accountID, requestID)
		if err != nil {
			return nil, fmt.Errorf("failed to check async task status: %v", err)
		}

		switch task.State {
		case 2: // Success
			var createdJobs []Job
			if err := json.Unmarshal([]byte(task.Output), &createdJobs); err != nil {
				return nil, fmt.Errorf("failed to parse jobs from async task output: %v", err)
			}
			return createdJobs, nil
		case 3: // Fail
			return nil, fmt.Errorf("jobs creation failed: %s", task.Output)
		case 0, 1: // Not Started or In Progress
			continue
		default:
			return nil, fmt.Errorf("unknown async task state: %d", task.State)
		}
	}

	return nil, fmt.Errorf("jobs creation timed out after 30 seconds")
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

func getAllJobs(host, apiKey, apiSecret, accountID string, projectID uint64, limit, offset int) ([]Job, error) {
	url := fmt.Sprintf("%s/api/v1/jobs?projectId=%d&limit=%d&offset=%d", host, projectID, limit, offset)
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

func testGetAllCredentials(t *testing.T, client *http.Client, baseURL string, accountId string) {
	// Test pagination
	testPagination(t, client, baseURL, accountId)

	// Test ordering
	testOrdering(t, client, baseURL, accountId)
}

func testOrdering(t *testing.T, client *http.Client, baseURL string, accountId string) {
	orderByFields := []string{"date_created", "date_modified", "created_by", "modified_by", "deleted_by"}
	directions := []string{"asc", "desc"}

	for _, field := range orderByFields {
		for _, direction := range directions {
			url := fmt.Sprintf("%s/api/v1/credentials?accountId=%s&orderBy=%s&orderByDirection=%s", baseURL, accountId, field, direction)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status code %d, got %d for ordering by %s %s", http.StatusOK, resp.StatusCode, field, direction)
				continue
			}

			var response map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				t.Errorf("Failed to decode response: %v", err)
				continue
			}

			// Verify response structure
			if _, ok := response["data"]; !ok {
				t.Errorf("Response missing 'data' field for ordering by %s %s", field, direction)
			}
			if _, ok := response["pagination"]; !ok {
				t.Errorf("Response missing 'pagination' field for ordering by %s %s", field, direction)
			}

			// Verify ordering (basic check)
			data, ok := response["data"].([]interface{})
			if !ok {
				t.Errorf("Data field is not an array for ordering by %s %s", field, direction)
				continue
			}

			if len(data) > 1 {
				// Get the first two items to verify ordering
				item1, ok1 := data[0].(map[string]interface{})
				item2, ok2 := data[1].(map[string]interface{})

				if ok1 && ok2 {
					val1, exists1 := item1[field]
					val2, exists2 := item2[field]

					if exists1 && exists2 {
						// Convert to comparable values
						var comparable1, comparable2 interface{}
						switch v1 := val1.(type) {
						case string:
							comparable1 = v1
							comparable2 = val2.(string)
						case float64:
							comparable1 = v1
							comparable2 = val2.(float64)
						}

						// Verify ordering
						if direction == "asc" {
							if comparable1.(float64) > comparable2.(float64) {
								t.Errorf("Incorrect ascending order for field %s", field)
							}
						} else {
							if comparable1.(float64) < comparable2.(float64) {
								t.Errorf("Incorrect descending order for field %s", field)
							}
						}
					}
				}
			}
		}
	}
}

func testPagination(t *testing.T, client *http.Client, baseURL string, accountId string) {
	// Test with default pagination
	url := fmt.Sprintf("%s/api/v1/credentials?accountId=%s", baseURL, accountId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		return
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify response structure
	if _, ok := response["data"]; !ok {
		t.Error("Response missing 'data' field")
	}
	if _, ok := response["pagination"]; !ok {
		t.Error("Response missing 'pagination' field")
	}

	// Test with custom pagination
	url = fmt.Sprintf("%s/api/v1/credentials?accountId=%s&limit=5&offset=0", baseURL, accountId)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		return
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify pagination parameters
	pagination, ok := response["pagination"].(map[string]interface{})
	if !ok {
		t.Error("Pagination field is not an object")
		return
	}

	if limit, ok := pagination["limit"].(float64); !ok || limit != 5 {
		t.Errorf("Expected limit to be 5, got %v", pagination["limit"])
	}

	if offset, ok := pagination["offset"].(float64); !ok || offset != 0 {
		t.Errorf("Expected offset to be 0, got %v", pagination["offset"])
	}
}
