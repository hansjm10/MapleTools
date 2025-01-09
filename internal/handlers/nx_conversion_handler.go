package handlers

import (
	"embed"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"

	"MapleTools/internal/models"
	"MapleTools/internal/utils"
)

const fixedMesoAmount = 100000000.0

//go:embed templates/*
var templatesFS embed.FS

func FormHandler(w http.ResponseWriter, r *http.Request) {
	data := models.NXConversionTemplate{}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			data.Error = "Error parsing form data."
			renderTemplate(w, data)
			return
		}

		nxFor100MStr := r.FormValue("nxFor100M")
		numberOfItemsStr := r.FormValue("numberOfItems")
		nxPerItemStr := r.FormValue("nxPerItem")

		// Preserve values
		data.NXFor100M = nxFor100MStr
		data.NumberOfItems = numberOfItemsStr
		data.NXPerItem = nxPerItemStr

		// Convert input
		nxFor100M, err1 := strconv.ParseFloat(nxFor100MStr, 64)
		numberOfItems, err2 := strconv.ParseFloat(numberOfItemsStr, 64)
		nxPerItem, err3 := strconv.ParseFloat(nxPerItemStr, 64)

		if err1 != nil || err2 != nil || err3 != nil || nxFor100M == 0 {
			data.Error = "Invalid input. Please enter valid numbers."
			renderTemplate(w, data)
			return
		}

		// Perform calculations
		totalNX := numberOfItems * nxPerItem
		adjustedTotalNX := totalNX * 1.01
		roundedTotalNX := math.Ceil(adjustedTotalNX/100) * 100
		requiredMeso := (roundedTotalNX / nxFor100M) * fixedMesoAmount

		// Fill the data
		data.TotalNX = utils.FormatNumber(totalNX)
		data.AdjustedTotalNX = utils.FormatNumber(adjustedTotalNX)
		data.RoundedTotalNX = utils.FormatNumber(roundedTotalNX)
		data.ConversionRate = fmt.Sprintf("%.2f", nxFor100M)
		data.TotalMeso = utils.FormatNumber(requiredMeso)
	}

	renderTemplate(w, data)
}

func renderTemplate(w http.ResponseWriter, data models.NXConversionTemplate) {
	// Use ParseFS with the embedded file system and the path to the template.
	tmpl, err := template.ParseFS(templatesFS, "templates/form.gohtml")
	if err != nil {
		http.Error(w, "Error parsing template.", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template.", http.StatusInternalServerError)
	}
}
