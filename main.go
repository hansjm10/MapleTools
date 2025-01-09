package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

// Fixed meso amount: 100,000,000 meso
const fixedMesoAmount = 100000000.0

// TemplateData holds the data passed to the HTML template
type TemplateData struct {
	NXFor100M       string
	NumberOfItems   string
	NXPerItem       string
	TotalNX         string
	AdjustedTotalNX string
	RoundedTotalNX  string
	ConversionRate  string
	Calculation     string
	TotalMeso       string
	Error           string
}

var formTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Meso to NX Calculator</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .container { max-width: 500px; margin: auto; }
        input[type=text], input[type=number] {
            width: 100%; padding: 8px 12px; margin: 6px 0 12px 0; box-sizing: border-box;
        }
        input[type=submit] {
            background-color: #4CAF50; color: white; padding: 10px 16px; border: none; cursor: pointer; width: 100%;
        }
        input[type=submit]:hover { background-color: #45a049; }
        .result, .error {
            background-color: #f2f2f2; padding: 15px; margin-top: 20px; border-radius: 5px;
        }
        .error { background-color: #ffe6e6; color: #cc0000; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Meso to NX Calculator</h1>
        <form method="POST" action="/">
            <label for="nxFor100M">NX received for 100M meso:</label><br>
            <input type="number" id="nxFor100M" name="nxFor100M" placeholder="e.g., 805" value="{{.NXFor100M}}" required><br>
            
            <label for="numberOfItems">Number of items:</label><br>
            <input type="number" id="numberOfItems" name="numberOfItems" placeholder="e.g., 5" value="{{.NumberOfItems}}" required><br>
            
            <label for="nxPerItem">NX cost per item:</label><br>
            <input type="number" id="nxPerItem" name="nxPerItem" placeholder="e.g., 1300" value="{{.NXPerItem}}" required><br>
            
            <input type="submit" value="Calculate">
        </form>
        
        {{if .Error}}
            <div class="error">
                <p>{{.Error}}</p>
            </div>
        {{end}}

        {{if .TotalMeso}}
            <div class="result">
                <h2>Calculation Result</h2>
                <p><strong>Total NX needed for {{.NumberOfItems}} item(s) at {{.NXPerItem}} NX each:</strong> {{.TotalNX}} NX</p>
                <p><strong>After applying 1% fee:</strong> {{.AdjustedTotalNX}} NX</p>
                <p><strong>Rounded up to nearest 100 intervals:</strong> {{.RoundedTotalNX}} NX</p>
                <p><strong>Conversion rate:</strong> {{.ConversionRate}} NX per 100M meso</p>
                <p><strong>Calculation:</strong> (({{.TotalNX}} NX × 1.01) rounded to {{.RoundedTotalNX}} NX) / {{.NXFor100M}} NX × 100,000,000 meso = {{.TotalMeso}} meso</p>
                <p><strong>Total Meso Required:</strong> {{.TotalMeso}}</p>
            </div>
        {{end}}
    </div>
</body>
</html>
`

func main() {
	http.HandleFunc("/", formHandler)
	fmt.Println("Server starting at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// formHandler handles both GET and POST requests
func formHandler(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{}

	if r.Method == http.MethodPost {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			data.Error = "Error parsing form data."
			renderTemplate(w, data)
			return
		}

		// Retrieve form values
		nxFor100MStr := r.FormValue("nxFor100M")
		numberOfItemsStr := r.FormValue("numberOfItems")
		nxPerItemStr := r.FormValue("nxPerItem")

		// Preserve user input in case of error
		data.NXFor100M = nxFor100MStr
		data.NumberOfItems = numberOfItemsStr
		data.NXPerItem = nxPerItemStr

		// Convert strings to floats
		nxFor100M, err1 := strconv.ParseFloat(nxFor100MStr, 64)
		numberOfItems, err2 := strconv.ParseFloat(numberOfItemsStr, 64)
		nxPerItem, err3 := strconv.ParseFloat(nxPerItemStr, 64)

		if err1 != nil || err2 != nil || err3 != nil || nxFor100M == 0 {
			data.Error = "Invalid input. Please enter valid numbers. Ensure that 'NX received for 100M meso' is not zero."
			renderTemplate(w, data)
			return
		}

		// Calculations
		totalNX := numberOfItems * nxPerItem
		adjustedTotalNX := totalNX * 1.01 // Apply 1% fee
		// Round up to nearest 100 intervals
		roundedTotalNX := math.Ceil(adjustedTotalNX/100) * 100
		requiredMeso := (roundedTotalNX / nxFor100M) * fixedMesoAmount

		// Format numbers and fill TemplateData
		data.TotalNX = formatNumber(totalNX)
		data.AdjustedTotalNX = formatNumber(adjustedTotalNX)
		data.RoundedTotalNX = formatNumber(roundedTotalNX)
		data.ConversionRate = fmt.Sprintf("%.2f", nxFor100M)
		data.Calculation = fmt.Sprintf("((%.2f NX × 1.01) rounded to %s NX) / %.2f NX × 100,000,000 meso = %.2f meso",
			totalNX, data.RoundedTotalNX, nxFor100M, requiredMeso)
		data.TotalMeso = formatNumber(requiredMeso)
	}

	renderTemplate(w, data)
}

func renderTemplate(w http.ResponseWriter, data TemplateData) {
	tmpl, err := template.New("form").Parse(formTemplate)
	if err != nil {
		http.Error(w, "Error parsing template.", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template.", http.StatusInternalServerError)
	}
}

// formatNumber formats large numbers into K, M, B, etc.
func formatNumber(n float64) string {
	switch {
	case n >= 1e12:
		return fmt.Sprintf("%.2fT", n/1e12)
	case n >= 1e9:
		return fmt.Sprintf("%.2fB", n/1e9)
	case n >= 1e6:
		return fmt.Sprintf("%.2fM", n/1e6)
	case n >= 1e3:
		return fmt.Sprintf("%.2fK", n/1e3)
	default:
		return fmt.Sprintf("%.0f", math.Round(n))
	}
}
