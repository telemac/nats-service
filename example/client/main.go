package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	fmt.Println("NATS Greeting Service Client")
	fmt.Println("============================")

	// Test the hello endpoint
	fmt.Println("\n1. Testing 'hello' endpoint:")
	testHello(nc, "World")

	// Test hello with different languages
	fmt.Println("\n2. Testing 'hello' with different languages:")
	languages := []struct {
		lang string
		name string
	}{
		{"es", "Mundo"},
		{"fr", "Monde"},
		{"de", "Welt"},
		{"ja", "世界"},
		{"zh", "世界"},
	}

	for _, test := range languages {
		testHelloWithLanguage(nc, test.name, test.lang)
	}

	// Test the goodbye endpoint
	fmt.Println("\n3. Testing 'goodbye' endpoint:")
	testGoodbye(nc, "Friend", "See you later")

	// Test with empty name
	fmt.Println("\n4. Testing with empty name (should default to 'World'/'Friend'):")
	testHello(nc, "")
	testGoodbye(nc, "", "")

	// Test with custom timeout
	fmt.Println("\n5. Testing with custom timeout:")
	testHelloWithTimeout(nc, "Timeout Test", 2*time.Second)

	fmt.Println("\nAll tests completed!")
}

func testHello(nc *nats.Conn, name string) {
	// Prepare request
	request := map[string]interface{}{
		"name": name,
	}

	// Send request
	msg, err := nc.Request("greeter-service.greet.hello", encodeJSON(request), 2*time.Second)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		log.Printf("Failed to parse response: %v", err)
		return
	}

	// Print result
	fmt.Printf("  Request: Hello %s\n", name)
	fmt.Printf("  Response: %s\n", response["greeting"])
	fmt.Printf("  Timestamp: %s\n", response["timestamp"])
}

func testHelloWithLanguage(nc *nats.Conn, name, language string) {
	// Prepare request
	request := map[string]interface{}{
		"name":     name,
		"language": language,
	}

	// Send request
	msg, err := nc.Request("greeter-service.greet.hello", encodeJSON(request), 2*time.Second)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		log.Printf("Failed to parse response: %v", err)
		return
	}

	// Print result
	fmt.Printf("  Request: Hello %s (language: %s)\n", name, language)
	fmt.Printf("  Response: %s\n", response["greeting"])
}

func testGoodbye(nc *nats.Conn, name, farewell string) {
	// Prepare request
	request := map[string]interface{}{
		"name":     name,
		"farewell": farewell,
	}

	// Send request
	msg, err := nc.Request("greeter-service.greet.goodbye", encodeJSON(request), 2*time.Second)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		log.Printf("Failed to parse response: %v", err)
		return
	}

	// Print result
	fmt.Printf("  Request: Farewell %s (%s)\n", name, farewell)
	fmt.Printf("  Response: %s\n", response["message"])
	fmt.Printf("  Timestamp: %s\n", response["timestamp"])
}

func testHelloWithTimeout(nc *nats.Conn, name string, timeout time.Duration) {
	// Prepare request
	request := map[string]interface{}{
		"name": name,
	}

	fmt.Printf("  Sending request with %v timeout...\n", timeout)

	// Send request with custom timeout
	start := time.Now()
	msg, err := nc.Request("greeter-service.greet.hello", encodeJSON(request), timeout)
	elapsed := time.Since(start)

	if err != nil {
		log.Printf("  Request failed after %v: %v", elapsed, err)
		return
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		log.Printf("  Failed to parse response: %v", err)
		return
	}

	// Print result
	fmt.Printf("  Response received in %v: %s\n", elapsed, response["greeting"])
}

// encodeJSON helper to encode JSON and handle errors
func encodeJSON(data interface{}) []byte {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode JSON: %v", err)
		return []byte("{}")
	}
	return jsonData
}
