package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	baseURL       = "http://localhost:8080"
	walletAddress = "0x7C1409F8280144B3e3FBD3F977A097f492b4289e"
)

func main() {
	fmt.Println("🧪 Pass Chain Smoke Tests")
	fmt.Println("=" + strings.Repeat("=", 50))

	passed := 0
	failed := 0

	// Test 1: Health check
	if testHealthCheck() {
		passed++
	} else {
		failed++
	}

	// Test 2: Get user
	if testGetUser() {
		passed++
	} else {
		failed++
	}

	// Test 3: List credentials
	if testListCredentials() {
		passed++
	} else {
		failed++
	}

	// Test 4: Get stats
	if testGetStats() {
		passed++
	} else {
		failed++
	}

	// Results
	fmt.Println()
	fmt.Println(strings.Repeat("=", 52))
	fmt.Printf("Results: %d passed, %d failed\n", passed, failed)
	if failed > 0 {
		os.Exit(1)
	}
}

func testHealthCheck() bool {
	fmt.Print("✓ Testing /health... ")
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		fmt.Printf("❌ FAIL: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("❌ FAIL: status %d\n", resp.StatusCode)
		return false
	}

	fmt.Println("✅ PASS")
	return true
}

func testGetUser() bool {
	fmt.Print("✓ Testing GET /api/v1/me... ")
	
	req, _ := http.NewRequest("GET", baseURL+"/api/v1/me", nil)
	req.Header.Set("X-Wallet-Address", walletAddress)
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("❌ FAIL: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ FAIL: status %d, body: %s\n", resp.StatusCode, string(body))
		return false
	}

	var user map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&user)
	
	if user["walletAddress"] != walletAddress {
		fmt.Printf("❌ FAIL: wallet address mismatch\n")
		return false
	}

	fmt.Println("✅ PASS")
	return true
}

func testListCredentials() bool {
	fmt.Print("✓ Testing GET /api/v1/credentials... ")
	
	req, _ := http.NewRequest("GET", baseURL+"/api/v1/credentials", nil)
	req.Header.Set("X-Wallet-Address", walletAddress)
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("❌ FAIL: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ FAIL: status %d, body: %s\n", resp.StatusCode, string(body))
		return false
	}

	fmt.Println("✅ PASS")
	return true
}

func testGetStats() bool {
	fmt.Print("✓ Testing GET /api/v1/stats... ")
	
	req, _ := http.NewRequest("GET", baseURL+"/api/v1/stats", nil)
	req.Header.Set("X-Wallet-Address", walletAddress)
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("❌ FAIL: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ FAIL: status %d, body: %s\n", resp.StatusCode, string(body))
		return false
	}

	var stats map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&stats)
	
	if _, ok := stats["totalCredentials"]; !ok {
		fmt.Printf("❌ FAIL: missing totalCredentials\n")
		return false
	}

	fmt.Println("✅ PASS")
	return true
}

