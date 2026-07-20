package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/amtsoft/money-converter/internal/converter"
)

// PageData holds the data to be rendered in our HTML template
type PageData struct {
	Result float64
}

func main() {
	// 1. Load config and check keys
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY is not set in environment")
	}

	// 2. Fetch supported currency codes once at startup and cache in memory
	codes, err := converter.GetSupportedCodes(apiKey)
	if err != nil {
		log.Fatalf("Failed to fetch supported currency codes: %v", err)
	}

	// 3. Parse templates
	tmpl := template.Must(template.ParseFiles("cmd/web/templates/index.html"))

	// 4. Serve static assets (CSS, etc.) from /public
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// 5. Register routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, nil); err != nil {
			log.Printf("template execute error: %v", err)
		}
	})

	http.HandleFunc("/codes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(codes); err != nil {
			log.Printf("codes encode error: %v", err)
		}
	})

	http.HandleFunc("/convert", func(w http.ResponseWriter, r *http.Request) {
		amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
		if err != nil {
			http.Error(w, "Amount must be a valid number", http.StatusBadRequest)
			return
		}
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")

		rate, err := converter.GetRate(apiKey, from, to)
		if err != nil {
			log.Printf("GetRate error: %v", err)
			http.Error(w, "Failed to fetch exchange rate", http.StatusInternalServerError)
			return
		}

		result := converter.Convert(converter.ConversionRequest{Amount: amount}, rate)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(PageData{Result: result}); err != nil {
			log.Printf("convert encode error: %v", err)
		}
	})

	// 6. Start server
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}