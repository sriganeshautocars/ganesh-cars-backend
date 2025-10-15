package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Car struct {
	ID                 int       `json:"id" db:"id"`
	Thumbnail          string    `json:"thumbnail" db:"thumbnail"`
	Brand              string    `json:"brand" db:"brand"`
	Name               string    `json:"name" db:"name"`
	Variant            string    `json:"variant" db:"variant"`
	KMDriven           int       `json:"km_driven" db:"km_driven"`
	FuelType           string    `json:"fuel_type" db:"fuel_type"`
	BodyType           string    `json:"body_type" db:"body_type"`
	TransmissionType   string    `json:"transmission_type" db:"transmission_type"`
	Price              float64   `json:"price" db:"price"`
	Location           string    `json:"location" db:"location"`
	Insurance          string    `json:"insurance" db:"insurance"`
	NoOfSeats          int       `json:"no_of_seats" db:"no_of_seats"`
	RegNumber          string    `json:"reg_number" db:"reg_number"`
	Ownership          int       `json:"ownership" db:"ownership"`
	EngineDisplacement int       `json:"engine_displacement" db:"engine_displacement"`
	HighwayMileage     float64   `json:"highway_mileage" db:"highway_mileage"`
	MakeYear           int       `json:"make_year" db:"make_year"`
	RegYear            int       `json:"reg_year" db:"reg_year"`
	Features           string    `json:"features" db:"features"`
	Specifications     string    `json:"specifications" db:"specifications"`
	Images             string    `json:"images" db:"images"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type User struct {
	ID       int    `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// jsonToDBColumn maps JSON keys from the request to database column names.
// These keys should match the `json:"..."` tags in the Car struct for fields that can be updated.
var jsonToDBColumn = map[string]string{
	"thumbnail":           "thumbnail",
	"brand":               "brand",
	"name":                "name",
	"variant":             "variant",
	"km_driven":           "km_driven",
	"fuel_type":           "fuel_type",
	"body_type":           "body_type",
	"transmission_type":   "transmission_type",
	"price":               "price",
	"location":            "location",
	"insurance":           "insurance",
	"no_of_seats":         "no_of_seats",
	"reg_number":          "reg_number",
	"ownership":           "ownership",
	"engine_displacement": "engine_displacement",
	"highway_mileage":     "highway_mileage",
	"make_year":           "make_year",
	"reg_year":            "reg_year",
	"specifications":      "specifications",
	"features":            "features",
	"images":              "images",
}

// var supabaseClient *supabase.Client
var conn *pgx.Conn
var err error

func main() {
	// Load environment variables
	godotenv.Load()

	// Initialize Supabase client
	// supabaseURL := os.Getenv("SUPABASE_URL")
	// supabaseKey := os.Getenv("SUPABASE_KEY")
	// supabaseClient = supabase.CreateClient(supabaseURL, supabaseKey)
	conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())
	fmt.Println("Connected to the database", conn)

	// Set up routes
	router := mux.NewRouter()
	router.HandleFunc("/cars", createCarHandler).Methods("POST")
	router.HandleFunc("/cars", getAllCarsHandler).Methods("GET")
	router.HandleFunc("/cars/{id}", getCarHandler).Methods("GET")
	router.HandleFunc("/cars/{id}", updateCarHandler).Methods("PUT")
	router.HandleFunc("/cars/{id}", deleteCarHandler).Methods("DELETE")
	router.HandleFunc("/login", loginHandler).Methods("POST")

	// CORS configuration
	// In production, you should replace "http://localhost:3000" with your actual frontend domain.
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000", "http://localhost:5173"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	allowCredentials := handlers.AllowCredentials()

	// Start server
	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders, allowCredentials)(router)))

	// var version string
	// if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
	// 	log.Fatalf("Query failed: %v", err)
	// }

	// log.Println("Connected to:", version)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	// Fetch user from DB
	var userID int
	var hashedPassword string
	sqlStatement := `
		SELECT id, password
		FROM users
		WHERE email = $1
	`
	err := conn.QueryRow(context.Background(), sqlStatement, req.Email).Scan(&userID, &hashedPassword)
	fmt.Println("Error while login", err)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   req.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate token: %v", err), http.StatusInternalServerError)
		return
	}

	// Response with token
	resp := map[string]interface{}{
		"token": tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func createCarHandler(w http.ResponseWriter, r *http.Request) {
	var car Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	// ID, CreatedAt, UpdatedAt are handled by the database
	sqlStatement := `
		INSERT INTO cars (thumbnail, brand, name, variant, km_driven, fuel_type, body_type, transmission_type, price, location, insurance, no_of_seats, reg_number, ownership, engine_displacement, highway_mileage, make_year, reg_year, specifications, features, images)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,$21)
		RETURNING id, created_at, updated_at`

	err := conn.QueryRow(context.Background(), sqlStatement,
		car.Thumbnail, car.Brand, car.Name, car.Variant, car.KMDriven, car.FuelType, car.BodyType, car.TransmissionType, car.Price, car.Location, car.Insurance, car.NoOfSeats, car.RegNumber, car.Ownership, car.EngineDisplacement, car.HighwayMileage, car.MakeYear, car.RegYear, car.Specifications, car.Features, car.Images,
	).Scan(&car.ID, &car.CreatedAt, &car.UpdatedAt)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create car: %v", err), http.StatusInternalServerError)
		log.Printf("Failed to create car: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(car)
}

func getAllCarsHandler(w http.ResponseWriter, r *http.Request) {
	sqlStatement := `
		SELECT id, thumbnail, brand, name, variant, km_driven, fuel_type, body_type, transmission_type, 
		       price, location, insurance, no_of_seats, reg_number, ownership, 
		       engine_displacement, highway_mileage, make_year, reg_year,
		       specifications, features, images, created_at, updated_at
		FROM cars`
	rows, err := conn.Query(context.Background(), sqlStatement)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to query cars: %v\n", err)
		return
	}
	defer rows.Close()

	var cars []Car
	for rows.Next() {
		var car Car
		if err := rows.Scan(
			&car.ID, &car.Thumbnail, &car.Brand, &car.Name, &car.Variant, &car.KMDriven,
			&car.FuelType, &car.BodyType, &car.TransmissionType, &car.Price,
			&car.Location, &car.Insurance, &car.NoOfSeats, &car.RegNumber,
			&car.Ownership, &car.EngineDisplacement, &car.HighwayMileage,
			&car.MakeYear, &car.RegYear, &car.Specifications, &car.Features,
			&car.Images, &car.CreatedAt, &car.UpdatedAt,
		); err != nil {
			http.Error(w, fmt.Sprintf("Row scan error: %v", err), http.StatusInternalServerError)
			log.Printf("Row scan failed: %v\n", err)
			return
		}
		cars = append(cars, car)
	}

	if rows.Err() != nil {
		http.Error(w, fmt.Sprintf("Row iteration error: %v", rows.Err()), http.StatusInternalServerError)
		log.Printf("Row iteration error: %v\n", rows.Err())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

func getCarHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	var car Car
	sqlStatement := `
		SELECT id, thumbnail, brand, name, variant, km_driven, fuel_type, body_type, transmission_type, 
		       price, location, insurance, no_of_seats, reg_number, ownership, 
		       engine_displacement, highway_mileage, make_year, reg_year, 
		       specifications, features, images, created_at, updated_at 
		FROM cars WHERE id = $1`

	err = conn.QueryRow(context.Background(), sqlStatement, id).Scan(
		&car.ID, &car.Thumbnail, &car.Brand, &car.Name, &car.Variant, &car.KMDriven,
		&car.FuelType, &car.BodyType, &car.TransmissionType, &car.Price,
		&car.Location, &car.Insurance, &car.NoOfSeats, &car.RegNumber,
		&car.Ownership, &car.EngineDisplacement, &car.HighwayMileage,
		&car.MakeYear, &car.RegYear, &car.Specifications, &car.Features,
		&car.Images, &car.CreatedAt, &car.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Car not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to get car: %v", err), http.StatusInternalServerError)
			log.Printf("Failed to get car with ID %d: %v\n", id, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(car)
}

func updateCarHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	if len(updates) == 0 {
		http.Error(w, "No fields provided for update", http.StatusBadRequest)
		return
	}

	var setClauses []string
	var args []interface{}
	argID := 1

	for jsonKey, value := range updates {
		dbColumn, ok := jsonToDBColumn[jsonKey]
		if !ok {
			// Optionally, you could return an error for unknown fields,
			// or simply ignore them as done here.
			log.Printf("Unknown field in update request: %s", jsonKey)
			continue
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbColumn, argID))
		args = append(args, value)
		argID++
	}

	if len(setClauses) == 0 {
		http.Error(w, "No valid fields provided for update", http.StatusBadRequest)
		return
	}

	// Always update the updated_at timestamp
	setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")

	sqlStatement := fmt.Sprintf("UPDATE cars SET %s WHERE id = $%d RETURNING *", strings.Join(setClauses, ", "), argID)
	args = append(args, id)

	var updatedCar Car
	err = conn.QueryRow(context.Background(), sqlStatement, args...).Scan(
		&updatedCar.ID, &updatedCar.Thumbnail, &updatedCar.Brand, &updatedCar.Name, &updatedCar.Variant, &updatedCar.KMDriven,
		&updatedCar.FuelType, &updatedCar.BodyType, &updatedCar.TransmissionType, &updatedCar.Price,
		&updatedCar.Location, &updatedCar.Insurance, &updatedCar.NoOfSeats, &updatedCar.RegNumber,
		&updatedCar.Ownership, &updatedCar.EngineDisplacement, &updatedCar.HighwayMileage,
		&updatedCar.MakeYear, &updatedCar.RegYear, &updatedCar.Specifications, &updatedCar.Features,
		&updatedCar.Images, &updatedCar.CreatedAt, &updatedCar.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Car not found or no update occurred", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to update car: %v", err), http.StatusInternalServerError)
			log.Printf("Failed to update car with ID %d: %v. SQL: %s, Args: %v", id, err, sqlStatement, args)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCar)
}

func deleteCarHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	sqlStatement := `DELETE FROM cars WHERE id = $1`
	commandTag, err := conn.Exec(context.Background(), sqlStatement, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete car: %v", err), http.StatusInternalServerError)
		log.Printf("Failed to delete car with ID %d: %v\n", id, err)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
