package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	// Import your internal logic package
	"github.com/mewteebee/money-converter/internal/converter"
)

// PageData holds the data to be rendered in our HTML template
type PageData struct {
	Result float64
}

func main() {
	// Parse the template file
	tmpl := template.Must(template.ParseFiles("cmd/web/templates/index.html"))

	// Handle GET requests to show the form
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	// Handle POST requests for conversion logic
	http.HandleFunc("/convert", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Retrieve form values
		amountStr := r.FormValue("amount")
		from := r.FormValue("from")
		to := r.FormValue("to")

		amount, _ := strconv.ParseFloat(amountStr, 64)

		// Call the logic from internal/converter/logic.go
		result := converter.Convert(converter.ConversionRequest{
			Amount:       amount,
			FromCurrency: from,
			ToCurrency:   to,
		})

		// Execute template with the result
		tmpl.Execute(w, PageData{Result: result})
	})

	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}