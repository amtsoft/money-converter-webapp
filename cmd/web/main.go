package main

import (
	"fmt"
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

    // 2. Parse templates
    tmpl := template.Must(template.ParseFiles("cmd/web/templates/index.html"))

    // 3. Register routes
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl.Execute(w, nil)
    })

    http.HandleFunc("/convert", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Redirect(w, r, "/", http.StatusSeeOther)
            return
        }

        // Parse form
        amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
        from := r.FormValue("from")
        to := r.FormValue("to")

        // Fetch live rate using the apiKey
        rate, err := converter.GetRate(apiKey, from, to)
        if err != nil {
            http.Error(w, "Failed to fetch exchange rate", http.StatusInternalServerError)
            return
        }

        // Calculate and execute template
        result := converter.Convert(converter.ConversionRequest{Amount: amount}, rate)
        tmpl.Execute(w, PageData{Result: result})
    })

    // 4. Start server
    fmt.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}