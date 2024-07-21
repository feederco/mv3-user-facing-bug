package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	webpush "github.com/SherClockHolmes/webpush-go"
)

const WebPushPublic = "BAzsP8FJ4nf_fPgTTv8Agj5z6WbIJMFWr7AezO3_b_zfLuCFhrzE8O1GLRvfXKQ7B4JKkElxlLBKjEszW7NuYQc"
const WebPushPrivate = "m-QQp-4SZXCGvbAInOWFeGeFiqyDO2eYsP3H-kC7NJQ"

func main() {
	cmd := os.Args[1]

	// Use this to generate new VAPID keys
	if cmd == "keygen" {
		keyGen()
		// Sends a message to a single token
	} else if cmd == "send-single" {
		token := os.Args[2]
		sendMessageToToken(token)
		// Sends a message to all tokens
	} else if cmd == "send" {
		sendMessageToAll()
	} else {
		// Start the server
		startServer()
	}
}

func keyGen() {
	privateKey, publicKey, _ := webpush.GenerateVAPIDKeys()
	log.Println("Private key", privateKey)
	log.Println("Public key", publicKey)
}

func sendMessageToAll() {
	subscriptions := readSubscriptionsFromFile()

	for _, token := range subscriptions {
		sendMessageToToken(token)
	}
}

func sendMessageToToken(token string) {
	var parsed string
	err := json.Unmarshal([]byte(token), &parsed)
	if err != nil {
		log.Fatalln("Invalid token " + token)
	}

	fmt.Println("Sending message to token", parsed)
	var webPushSubscription webpush.Subscription
	err = json.Unmarshal([]byte(parsed), &webPushSubscription)
	if err != nil {
		log.Fatalln("Invalid token")
	}

	pushData := map[string]interface{}{
		"title": "Test",
		"body":  time.Now().String(),
	}

	pushDataBytes, _ := json.Marshal(pushData)

	// Send Notification
	resp, err := webpush.SendNotification(pushDataBytes, &webPushSubscription, &webpush.Options{
		Subscriber:      "example@example.com",
		VAPIDPublicKey:  WebPushPublic,
		VAPIDPrivateKey: WebPushPrivate,
	})

	if err != nil {
		log.Println("Error sending notification", err)
		return
	}

	if resp.StatusCode == 410 {
		removeSubscriptionFromFile(token)
	}

	log.Println(resp, err)
}

func startServer() {
	http.HandleFunc("/token", handleToken)
	fmt.Println("Starting server on :1337")
	log.Fatal(http.ListenAndServe(":1337", nil))
}

func handleToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	bodyString, _ := io.ReadAll(r.Body)
	writeSubscriptionToFile(bodyString)

	fmt.Printf("Received subscription: %+v\n", bodyString)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success"}`))
}

func writeSubscriptionToFile(body []byte) {
	f, err := os.OpenFile("subscriptions.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := f.Write(body); err != nil {
		log.Fatal(err)
	}

	f.Write([]byte("\n"))
}

func removeSubscriptionFromFile(body string) {
	file, err := os.Open("subscriptions.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var buffer bytes.Buffer
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != body {
			buffer.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("subscriptions.txt", buffer.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func readSubscriptionsFromFile() []string {
	f, err := os.Open("subscriptions.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var subscriptions []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		subscriptions = append(subscriptions, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return subscriptions
}
