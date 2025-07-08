package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	//"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbpool *pgxpool.Pool

func initDB() {
	var err error
	// Connection string to PostgreSQL (replace with your actual credentials)
	connString := "postgres://postgres:postgres@localhost:5400/test"

	// Create the connection pool
	dbpool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	log.Printf("Database connection pool established")
}

func main() {
	// Initialize database connection
	initDB()

	r := chi.NewRouter()
	// r.Use(middleware.Logger)// Logging API requests i/o

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Get("/getScripMaster", ExecuteScripMaster)
	r.Get("/getParticipantMaster", ExecuteParticipantMaster)
	r.Get("/getUserData", ExecuteUserQuery)
	r.Get("/processOrderRequest", ExecuteAddUpdateOrderReqeustSTP)

	http.ListenAndServe("10.10.198.204:3000", r)

	defer dbpool.Close()
}

func ExecuteScripMaster(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	// Retrieve the parameters from the URL
	//tokenNumber1 := chi.URLParam(r, "tokenid")
	tokenNumber1 := r.URL.Query().Get("tokenid")
	marketsegmentid1 := r.URL.Query().Get("marketsegmentid")
	// log.Println("TokenNumber : ", tokenNumber1)

	if tokenNumber1 == " " && marketsegmentid1 == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}

	// Query to get the person by ID
	var instrumentType int
	var err error
	var err1 error
	var err2 error

	// Convert the ID to an integer
	tokenNumber, err1 := strconv.Atoi(tokenNumber1)
	marketsegmentid, err2 := strconv.Atoi(marketsegmentid1)

	if err1 != nil || err2 != nil {
		// log.Println("error : ", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = dbpool.QueryRow(context.Background(), "SELECT nInstrumentType FROM tbl_ScripMaster WHERE nToken = ($1) AND nMarketSegmentId=($2)", tokenNumber, marketsegmentid).Scan(&instrumentType)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Token not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Query error: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instrumentType)
}

func ExecuteParticipantMaster(w http.ResponseWriter, r *http.Request) {

	// Retrieve the parameters from the URL
	participantCodeId := r.URL.Query().Get("participantcodeid")
	segmentidparam := r.URL.Query().Get("segmentid")

	if participantCodeId == " " && segmentidparam == " " {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}

	// Query to get the person by ID
	var alphaChar string
	var err error

	// Convert the ID to an integer
	segmentID, err := strconv.Atoi(segmentidparam)
	if err != nil {
		// log.Println("error : ", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = dbpool.QueryRow(context.Background(), "SELECT sAlphaChar FROM tbl_ParticipantMaster WHERE sParticipantCode = ($1) and nMarketSegmentId = ($2)", participantCodeId, segmentID).Scan(&alphaChar)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "sParticipantCode,nMarketSegmentId not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Query error: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alphaChar)
}

func ExecuteUserQuery(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	// Retrieve the parameters from the URL
	dealerid := r.URL.Query().Get("dealerid")
	if dealerid == " " {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}

	// Query to get the data by ID
	var dealercategory int
	var err error

	err = dbpool.QueryRow(context.Background(), "SELECT ndealercategory FROM tbl_usermaster UM JOIN tbl_useralias UA ON UM.sdealercode=UA.sdealercode JOIN tbl_userparticipantcode UPC ON UPC.sdealercode=UM.sdealercode WHERE UM.sdealercode=($1)", dealerid).Scan(&dealercategory)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "data not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Query error: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dealercategory)
}

func ExecuteAddUpdateOrderReqeustSTP(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	// Retrieve the parameters from the URL
	dealerid := r.URL.Query().Get("dealerid")
	tokenparam := r.URL.Query().Get("tokenid")
	participantparam := r.URL.Query().Get("participantcodeid")
	segmentidparam := r.URL.Query().Get("segmentid")

	segmentid, err1 := strconv.Atoi(segmentidparam)
	tokenid, err2 := strconv.Atoi(tokenparam)
	// participantCodeId, err3 := strconv.Atoi(participantparam)//kunalj

	if dealerid == " " {
		http.Error(w, "Missing dealerID", http.StatusBadRequest)
		return
	}
	if err1 != nil && err2 != nil { //kunalj&& err3 != nil
		// log.Println("error : ", err1)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	// Query to get the data by ID
	var dealercategory string
	dealercategory = "success"
	var err error
	//kunalj | added _ to capture data
	_, err = dbpool.Exec(context.Background(), "CALL stp_InsertAndUpdateOrderBook($1,$2,$3,$4)", dealerid, tokenid, segmentid, participantparam) //kunalj
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "data not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Query error: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dealercategory)
}
