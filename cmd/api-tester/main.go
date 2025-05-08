package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Credential struct {
	ID          uint64    `json:"id"`
	AccountID   string    `json:"accountId"`
	APIKey      string    `json:"apiKey"`
	APISecret   string    `json:"apiSecret"`
	DateCreated time.Time `json:"dateCreated"`
}

type CredentialResponse struct {
	Success bool       `json:"success"`
	Data    Credential `json:"data"`
}

type Project struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"accountId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DateCreated time.Time `json:"dateCreated"`
}

type ProjectResponse struct {
	Success bool    `json:"success"`
	Data    Project `json:"data"`
}

type Job struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"accountId"`
	ProjectID   string    `json:"projectId"`
	Description string    `json:"description"`
	ExecutorID  string    `json:"executorId"`
	Data        string    `json:"data"`
	Spec        string    `json:"spec"`
	StartDate   string    `json:"startDate"`
	EndDate     string    `json:"endDate"`
	Timezone    string    `json:"timezone"`
	DateCreated time.Time `json:"dateCreated"`
}

type JobResponse struct {
	Success bool `json:"success"`
	Data    Job  `json:"data"`
}

type Execution struct {
	ID                    uint64    `json:"id"`
	AccountID             string    `json:"accountId"`
	UniqueID              string    `json:"uniqueId"`
	State                 string    `json:"state"`
	NodeID                string    `json:"nodeId"`
	JobID                 string    `json:"jobId"`
	LastExecutionDatetime string    `json:"lastExecutionDatetime"`
	NextExecutionDatetime string    `json:"nextExecutionDatetime"`
	JobQueueVersion       int       `json:"jobQueueVersion"`
	ExecutionVersion      int       `json:"executionVersion"`
	Logs                  string    `json:"logs"`
	DateCreated           time.Time `json:"dateCreated"`
}

type ExecutionResponse struct {
	Success bool      `json:"success"`
	Data    Execution `json:"data"`
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

	fmt.Println("apiKey: ", *apiKey)
	fmt.Println("apiSecret: ", *apiSecret)
	fmt.Println("host: ", *host)
	fmt.Println("accountID: ", *accountID)

	if *apiKey == "" || *apiSecret == "" || *host == "" || *accountID == "" {
		fmt.Println("All flags are required: -api-key, -api-secret, -host, -account-id")
		os.Exit(1)
	}

	// Create a credential
	credential, err := createCredential(*host, *apiKey, *apiSecret, *accountID)
	if err != nil {
		fmt.Printf("Failed to create credential: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created credential with API Key: %s\n", credential.APIKey)

	// Create a project
	project, err := createProject(*host, credential.APIKey, credential.APISecret, *accountID)
	if err != nil {
		fmt.Printf("Failed to create project: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created project with ID: %s\n", project.ID)

	// Create a job
	job, err := createJob(*host, credential.APIKey, credential.APISecret, *accountID, project.ID)
	if err != nil {
		fmt.Printf("Failed to create job: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created job with ID: %s\n", job.ID)

	// Get executions
	executions, err := getExecutions(*host, credential.APIKey, credential.APISecret, *accountID)
	if err != nil {
		fmt.Printf("Failed to get executions: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Retrieved %d executions\n", len(executions))

	// Get async task status if available
	if job.ID != "" {
		task, err := getAsyncTask(*host, credential.APIKey, credential.APISecret, *accountID, job.ID)
		if err != nil {
			fmt.Printf("Failed to get async task: %v\n", err)
		} else {
			fmt.Printf("Async task state: %d\n", task.State)
		}
	}
}

func createCredential(host, apiKey, apiSecret, accountID string) (*Credential, error) {
	url := fmt.Sprintf("%s/api/v1/credentials", host)
	reqBody := map[string]string{
		"accountId": accountID,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("x-api-secret", apiSecret)
	req.Header.Set("x-account-id", accountID)

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

func createProject(host, apiKey, apiSecret, accountID string) (*Project, error) {
	url := fmt.Sprintf("%s/api/v1/projects", host)
	reqBody := map[string]string{
		"name":        "Test Project",
		"description": "Project created by API tester",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	fmt.Println("apiKey: ", apiKey)
	fmt.Println("apiSecret: ", apiSecret)
	fmt.Println("accountID: ", accountID)
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("x-api-secret", apiSecret)
	req.Header.Set("x-account-id", accountID)

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

func createJob(host, apiKey, apiSecret, accountID, projectID string) (*Job, error) {
	url := fmt.Sprintf("%s/api/v1/jobs", host)
	reqBody := map[string]string{
		"description": "Test Job",
		"timezone":    "UTC",
		"startDate":   time.Now().Add(time.Hour).Format(time.RFC3339),
		"endDate":     time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"spec":        "@every 1h",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("x-api-secret", apiSecret)
	req.Header.Set("x-account-id", accountID)

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

	var response JobResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func getExecutions(host, apiKey, apiSecret, accountID string) ([]Execution, error) {
	url := fmt.Sprintf("%s/api/v1/executions", host)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("x-api-secret", apiSecret)
	req.Header.Set("x-account-id", accountID)

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
		return nil, fmt.Errorf("failed to get executions: %s", string(body))
	}

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			Total      int         `json:"total"`
			Offset     int         `json:"offset"`
			Limit      int         `json:"limit"`
			Executions []Execution `json:"executions"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data.Executions, nil
}

func getAsyncTask(host, apiKey, apiSecret, accountID, requestID string) (*AsyncTask, error) {
	url := fmt.Sprintf("%s/api/v1/async-tasks/%s", host, requestID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("x-api-secret", apiSecret)
	req.Header.Set("x-account-id", accountID)

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
