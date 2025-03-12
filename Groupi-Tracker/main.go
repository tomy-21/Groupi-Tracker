package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Driver struct {
	GivenName       string `json:"givenName"`
	FamilyName      string `json:"familyName"`
	PermanentNumber string `json:"permanentNumber"`
	Nationality     string `json:"nationality"`
	ImageURL        string
}

var favoriteDrivers []Driver

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

//  result page

type Race struct {
	RaceName string `json:"raceName"`
	Circuit  struct {
		CircuitName string `json:"circuitName"`
		Location    struct {
			Country string `json:"country"`
		} `json:"Location"`
	} `json:"Circuit"`
	Date  string `json:"date"`
	Round string `json:"round"`
}

type RaceTableResponse struct {
	MRData struct {
		RaceTable struct {
			Races []Race `json:"Races"`
		} `json:"RaceTable"`
	} `json:"MRData"`
}
type FastestLap struct {
	Rank string `json:"rank"`
	Lap  string `json:"lap"`
	Time struct {
		Time string `json:"time"`
	} `json:"Time"`
	AverageSpeed struct {
		Units string `json:"units"`
		Speed string `json:"speed"`
	} `json:"AverageSpeed"`
}

type RaceResult struct {
	Position string `json:"position"`
	Driver   struct {
		GivenName  string `json:"givenName"`
		FamilyName string `json:"familyName"`
	} `json:"Driver"`
	FastestLap  *FastestLap `json:"FastestLap,omitempty"`
	Constructor struct {
		Name string `json:"name"`
	} `json:"Constructor"`
	Points string `json:"points"`
}

type RaceResultResponse struct {
	MRData struct {
		RaceTable struct {
			Races []struct {
				RaceName string       `json:"raceName"`
				Results  []RaceResult `json:"Results"`
			} `json:"Races"`
		} `json:"RaceTable"`
	} `json:"MRData"`
}

// result handler

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	year := r.URL.Query().Get("year")
	if year == "" {
		year = "2024"
	}

	// Récupérer le numéro de la page (par défaut 1)
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Nombre de GP par page
	limit := 10
	offset := (page - 1) * limit

	// Récupérer les données de l'API
	apiURL := fmt.Sprintf("http://ergast.com/api/f1/%s.json", year)
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des résultats", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture des données", http.StatusInternalServerError)
		return
	}

	var result RaceTableResponse
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Erreur lors du parsing des données", http.StatusInternalServerError)
		return
	}

	// Gérer la pagination
	totalRaces := len(result.MRData.RaceTable.Races)
	totalPages := (totalRaces + limit - 1) / limit // Calcul du nombre total de pages

	// Découpage des GP affichés sur la page actuelle
	var paginatedRaces []Race
	if offset < totalRaces {
		end := offset + limit
		if end > totalRaces {
			end = totalRaces
		}
		paginatedRaces = result.MRData.RaceTable.Races[offset:end]
	}

	// Génération des années disponibles (1950 -> Année actuelle)
	years := []string{}
	currentYear := time.Now().Year()
	for y := currentYear; y >= 1950; y-- {
		years = append(years, fmt.Sprintf("%d", y))
	}

	// Struct des données envoyées au template
	data := struct {
		Years        []string
		SelectedYear string
		Races        []Race
		CurrentPage  int
		TotalPages   int
	}{
		Years:        years,
		SelectedYear: year,
		Races:        paginatedRaces,
		CurrentPage:  page,
		TotalPages:   totalPages,
	}

	tmpl.ExecuteTemplate(w, "results", data)
}

// Handler pour les détails du Grand Prix
func grandPrixDetailsHandler(w http.ResponseWriter, r *http.Request) {
	year := r.URL.Query().Get("year")
	round := r.URL.Query().Get("round")
	if year == "" || round == "" {
		http.Error(w, "Année ou manche manquante", http.StatusBadRequest)
		return
	}

	apiURL := fmt.Sprintf("http://ergast.com/api/f1/%s/%s/results.json", year, round)
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des détails du Grand Prix", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture des données", http.StatusInternalServerError)
		return
	}

	var result RaceResultResponse
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Erreur lors du parsing des données", http.StatusInternalServerError)
		return
	}

	if len(result.MRData.RaceTable.Races) == 0 {
		http.Error(w, "Aucun résultat trouvé", http.StatusNotFound)
		return
	}

	race := result.MRData.RaceTable.Races[0]

	data := struct {
		RaceName string
		Results  []RaceResult
	}{
		RaceName: race.RaceName,
		Results:  race.Results,
	}

	tmpl.ExecuteTemplate(w, "grandprixdetails", data)
}

var funcMap = template.FuncMap{
	"add": func(a, b int) int { return a + b },
	"sub": func(a, b int) int { return a - b },
}
var tmpl = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))

// Template global

// Handler pour la page des pilotes
func driversHandler(w http.ResponseWriter, r *http.Request) {
	year := r.URL.Query().Get("year")
	if year == "" {
		year = "2024" // Année par défaut
	}

	nationality := r.URL.Query().Get("nationality")
	number := r.URL.Query().Get("number")
	name := r.URL.Query().Get("name")

	apiURL := fmt.Sprintf("http://ergast.com/api/f1/%s/drivers.json", year)

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

	if result.MRData.DriverTable.Drivers == nil {
		http.Error(w, "Aucun pilote trouvé", http.StatusNotFound)
		return
	}

	// Liste des années disponibles
	years := []string{}
	currentYear := time.Now().Year()
	for y := currentYear; y >= 1950; y-- {
		years = append(years, fmt.Sprintf("%d", y))
	}

	// Liste des nationalités disponibles
	nationalitiesMap := make(map[string]bool)
	for _, driver := range result.MRData.DriverTable.Drivers {
		nationalitiesMap[driver.Nationality] = true
	}
	nationalities := []string{"All"}
	for key := range nationalitiesMap {
		nationalities = append(nationalities, key)
	}

	// Filtrage des pilotes
	var filteredDrivers []Driver
	for _, driver := range result.MRData.DriverTable.Drivers {
		if (nationality == "All" || nationality == "" || driver.Nationality == nationality) &&
			(number == "" || driver.PermanentNumber == number) &&
			(name == "" || containsIgnoreCase(driver.GivenName+" "+driver.FamilyName, name)) {
			filteredDrivers = append(filteredDrivers, driver)
		}
	}

	data := struct {
		Years               []string
		Drivers             []Driver
		Nationalities       []string
		SelectedYear        string
		SelectedNationality string
		SelectedNumber      string
		SelectedName        string
	}{
		Years:               years,
		Drivers:             filteredDrivers,
		Nationalities:       nationalities,
		SelectedYear:        year,
		SelectedNationality: nationality,
		SelectedNumber:      number,
		SelectedName:        name,
	}

	tmpl.ExecuteTemplate(w, "drivers", data)
}

// Fonction pour comparer les noms sans tenir compte de la casse
func containsIgnoreCase(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}

// fonction d'ajout en favorie

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
		"albert_park":    "/assets/images/circuits/albert_park.webp",
		"nurburgring":    "/assets/images/circuits/nurburgring.jpg",
		"bahrain":        "/assets/images/circuits/bahrein.jpg",
		"jeddah":         "/assets/images/circuits/jeddah.jpg",
		"baku":           "/assets/images/circuits/baku.jpeg",
		"shanghai":       "/assets/images/circuits/shanghai.webp",
		"miami":          "/assets/images/circuits/miami.jpg",
		"imola":          "/assets/images/circuits/imola.jpg",
		"monaco":         "/assets/images/circuits/monaco.webp",
		"indianapolis":   "/assets/images/circuits/indianapolis.webp",
		"villeneuve":     "/assets/images/circuits/villeneuve.jpg",
		"mugello":        "/assets/images/circuits/mugello.jpg",
		"catalunya":      "/assets/images/circuits/catalunya.png",
		"red_bull_ring":  "/assets/images/circuits/red_bull_ring.jpg",
		"silverstone":    "/assets/images/circuits/silverstone.webp",
		"hungaroring":    "/assets/images/circuits/hungaroring.jpg",
		"spa":            "/assets/images/circuits/spa.jpg",
		"zandvoort":      "/assets/images/circuits/zandvoort.webp",
		"monza":          "/assets/images/circuits/monza.avif",
		"marina_bay":     "/assets/images/circuits/marina_bay.jpg",
		"suzuka":         "/assets/images/circuits/suzuka.jpg",
		"losail":         "/assets/images/circuits/losail.jpg",
		"americas":       "/assets/images/circuits/austin.jpg",
		"rodriguez":      "/assets/images/circuits/rodriguez.avif",
		"interlagos":     "/assets/images/circuits/interlagos.avif",
		"vegas":          "/assets/images/circuits/vegas.jpg",
		"hockenheimring": "/assets/images/circuits/hockenheimring.jpg",
		"yas_marina":     "/assets/images/circuits/yas_marina.webp",
		"istanbul":       "/assets/images/circuits/istanbul.webp",
		"portimao":       "/assets/images/circuits/portimao.webp",
		"ricard":         "/assets/images/circuits/ricard.jpg",
		"sochi":          "/assets/images/circuits/sochi.avif",
		"adelaide":       "/assets/images/circuits/adelaide.jpg",
		"phoenix":        "/assets/images/circuits/phoenix.jpeg",
		"sepang":         "/assets/images/circuits/sepang.jpg",
		"bremgarten":     "/assets/images/circuits/bremgarten.jpg",
		"reims":          "/assets/images/circuits/reims.webp",
		"pedralbes":      "/assets/images/circuits/pedralbes.jpg",
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

//fonctionnalité Favorites

func addFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer les données du formulaire
	r.ParseForm()
	permanentNumber := r.FormValue("permanentNumber")
	givenName := r.FormValue("givenName")
	familyName := r.FormValue("familyName")

	// Vérifier si le pilote est déjà dans les favoris
	for _, driver := range favoriteDrivers {
		if driver.PermanentNumber == permanentNumber {
			http.Redirect(w, r, "/favorites", http.StatusSeeOther)
			return
		}
	}

	// Ajouter le pilote
	newDriver := Driver{
		GivenName:       givenName,
		FamilyName:      familyName,
		PermanentNumber: permanentNumber,
	}
	favoriteDrivers = append(favoriteDrivers, newDriver)
	fmt.Println("Ajout du favori :", permanentNumber, givenName, familyName)

	http.Redirect(w, r, "/drivers", http.StatusSeeOther) // Redirige vers la page des pilotes
}

func removeFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	permanentNumber := r.FormValue("permanentNumber")

	// Filtrer la liste des favoris
	var updatedFavorites []Driver
	for _, driver := range favoriteDrivers {
		if driver.PermanentNumber != permanentNumber {
			updatedFavorites = append(updatedFavorites, driver)
		}
	}
	favoriteDrivers = updatedFavorites

	http.Redirect(w, r, "/favorites", http.StatusSeeOther)
}

//retirer un pilote des favorie

// Handler pour la page d'accueil
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "base", nil)
}

// Handler pour la page "Pilotes"
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "about.html", nil)
}

func favoritesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "favorites", favoriteDrivers)
}

func main() {
	// Servir les fichiers statiques
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/circuits", circuitsHandler)
	http.HandleFunc("/results", resultsHandler)
	http.HandleFunc("/results/details", grandPrixDetailsHandler)
	http.HandleFunc("/favorites", favoritesHandler)
	http.HandleFunc("/remove_favorite", removeFavoriteHandler)
	http.HandleFunc("/drivers", driversHandler)
	// Lancer le serveur
	port := ":8080"
	log.Println("Serveur démarré sur http://localhost" + port)
	http.HandleFunc("/add_favorite", addFavoriteHandler)

	log.Fatal(http.ListenAndServe(port, nil))
}
