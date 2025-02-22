package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Driver struct {
	GivenName       string `json:"givenName"`
	FamilyName      string `json:"familyName"`
	PermanentNumber string `json:"permanentNumber"`
	Nationality     string `json:"nationality"`
	ImageURL        string
}

// Structure pour parser la réponse de l'API Ergast
type ErgastResponse struct {
	MRData struct {
		DriverTable struct {
			Drivers []Driver `json:"Drivers"`
		} `json:"DriverTable"`
	} `json:"MRData"`
}

type Circuit struct {
	CircuitID   string `json:"circuitId"`
	CircuitName string `json:"circuitName"`
	Location    struct {
		Locality string `json:"locality"`
		Country  string `json:"country"`
	} `json:"Location"`
	URL      string `json:"url"`
	ImageURL string
}

type ErgastCircuitResponse struct {
	MRData struct {
		CircuitTable struct {
			Circuits []Circuit `json:"Circuits"`
		} `json:"CircuitTable"`
	} `json:"MRData"`
}

// Template global
var tmpl = template.Must(template.ParseGlob("templates/*.html"))

// Handler pour la page des pilotes
func driversHandler(w http.ResponseWriter, r *http.Request) {
	year := r.URL.Query().Get("year")
	if year == "" {
		year = "2024" // Année par défaut
	}

	nationality := r.URL.Query().Get("nationality") // Récupère la nationalité sélectionnée

	apiURL := "http://ergast.com/api/f1/" + year + "/drivers.json"

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des pilotes", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture de la réponse", http.StatusInternalServerError)
		return
	}

	var result ErgastResponse
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Erreur lors du parsing des données", http.StatusInternalServerError)
		return
	}

	// Liste des années disponibles
	years := []string{}
	currentYear := time.Now().Year()
	for y := currentYear; y >= 1950; y-- {
		years = append(years, fmt.Sprintf("%d", y))
	}

	// Liste des nationalités uniques disponibles
	nationalitiesMap := make(map[string]bool)
	for _, driver := range result.MRData.DriverTable.Drivers {
		nationalitiesMap[driver.Nationality] = true
	}
	nationalities := []string{}
	for key := range nationalitiesMap {
		nationalities = append(nationalities, key)
	}

	// Filtrage des pilotes par nationalité
	var filteredDrivers []Driver
	for _, driver := range result.MRData.DriverTable.Drivers {
		if nationality == "All" || nationality == "" || driver.Nationality == nationality {
			filteredDrivers = append(filteredDrivers, driver)
		}
	}

	// Préparer les données pour le template
	data := struct {
		Years               []string
		Drivers             []Driver
		Nationalities       []string
		SelectedYear        string
		SelectedNationality string
	}{
		Years:               years,
		Drivers:             filteredDrivers,
		Nationalities:       nationalities,
		SelectedYear:        year,
		SelectedNationality: nationality,
	}

	// Exécuter le template avec les données filtrées
	tmpl.ExecuteTemplate(w, "drivers", data)
}

// Handler pour la page circuits

func circuitsHandler(w http.ResponseWriter, r *http.Request) {
	// Récupération de l'année depuis les paramètres de requête
	year := r.URL.Query().Get("year")
	if year == "" {
		year = "2024" // Année par défaut
	}

	apiURL := fmt.Sprintf("http://ergast.com/api/f1/%s/circuits.json", year)

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des circuits", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture de la réponse", http.StatusInternalServerError)
		return
	}

	var result ErgastCircuitResponse
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Erreur lors du parsing des données", http.StatusInternalServerError)
		return
	}

	// Associer les images des circuits
	circuitImages := map[string]string{
		"albert_park":   "/assets/images/circuits/melbourne.jpg",
		"bahrain":       "/assets/images/circuits/bahrein.jpg",
		"jeddah":        "/assets/images/circuits/jeddah.jpg",
		"baku":          "/assets/images/circuits/baku.jpeg",
		"shanghai":      "/assets/images/shanghai.jpg",
		"miami":         "/assets/images/miami.jpg",
		"imola":         "/assets/images/imola.jpg",
		"monaco":        "/assets/images/monaco.jpg",
		"montreal":      "/assets/images/montreal.jpg",
		"barcelona":     "/assets/images/barcelona.jpg",
		"red_bull_ring": "/assets/images/red_bull_ring.jpg",
		"silverstone":   "/assets/images/silverstone.jpg",
		"hungaroring":   "/assets/images/hungaroring.jpg",
		"spa":           "/assets/images/spa.jpg",
		"zandvoort":     "/assets/images/zandvoort.jpg",
		"monza":         "/assets/images/monza.jpg",
		"marina_bay":    "/assets/images/marina_bay.jpg",
		"suzuka":        "/assets/images/suzuka.jpg",
		"losail":        "/assets/images/losail.jpg",
		"americas":      "/assets/images/circuits/austin.jpg",
		"mexico_city":   "/assets/images/mexico_city.jpg",
		"sao_paulo":     "/assets/images/sao_paulo.jpg",
		"las_vegas":     "/assets/images/las_vegas.jpg",
		"yas_marina":    "/assets/images/yas_marina.jpg",
	}

	for i, circuit := range result.MRData.CircuitTable.Circuits {
		if img, ok := circuitImages[circuit.CircuitID]; ok {
			result.MRData.CircuitTable.Circuits[i].ImageURL = img
		} else {
			result.MRData.CircuitTable.Circuits[i].ImageURL = "/assets/images/default.jpg"
		}
	}

	// Générer une liste d'années de 1950 à aujourd'hui
	years := []string{}
	currentYear := time.Now().Year()
	for y := currentYear; y >= 1950; y-- {
		years = append(years, fmt.Sprintf("%d", y))
	}

	data := struct {
		Years        []string
		Circuits     []Circuit
		SelectedYear string
	}{
		Years:        years,
		Circuits:     result.MRData.CircuitTable.Circuits,
		SelectedYear: year,
	}

	tmpl.ExecuteTemplate(w, "circuits", data)
}

// Handler pour la page d'accueil
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "base", nil)
}

// Handler pour la page "Pilotes"
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "about.html", nil)
}

// Handler pour la page "Design"
func collectionHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "collection.html", nil)
}

// Handler pour la page "Circuits"
func favoritesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "favorites.html", nil)
}

func main() {
	// Servir les fichiers statiques
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/circuits", circuitsHandler)
	http.HandleFunc("/collection", collectionHandler)
	http.HandleFunc("/favorites", favoritesHandler)
	http.HandleFunc("/drivers", driversHandler) // Nouvelle route pour récupérer les pilotes

	// Lancer le serveur
	port := ":8080"
	log.Println("Serveur démarré sur http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
