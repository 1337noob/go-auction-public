package api

import (
	"encoding/json"
	"errors"
	"log"
	"main/auction/internal/application/read_model"
	"main/auction/internal/domain"
	"main/auction/module"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type HttpHandler struct {
	module *module.AuctionModule
}

func NewHttpHandler(module *module.AuctionModule) *HttpHandler {
	return &HttpHandler{
		module: module,
	}
}

type CreateAuctionRequest struct {
	LotID      string `json:"lot_id"`
	StartPrice int    `json:"start_price"`
	MinBidStep int    `json:"min_bid_step"`
	SellerID   string `json:"seller_id"`
	StartTime  string `json:"start_time"`
	Endtime    string `json:"end_time"`
	Timeout    string `json:"timeout"`
}

type CreateAuctionResponse struct {
	ID string `json:"id"`
}

type PlaceBidRequest struct {
	UserID string `json:"user_id"`
	Amount int    `json:"amount"`
}

type PlaceBidResponse struct {
	Message string `json:"message"`
}

type CancelAuctionRequest struct {
	SellerID string `json:"seller_id"`
	Reason   string `json:"reason"`
}

type CancelAuctionResponse struct {
	Message string `json:"message"`
}

type CreateLotRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id"`
}

type CreateLotResponse struct {
	ID string `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *HttpHandler) writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func (h *HttpHandler) CreateAuction(w http.ResponseWriter, r *http.Request) {
	var req CreateAuctionRequest
	reader := json.NewDecoder(r.Body)
	err := reader.Decode(&req)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.StartTime, time.Local)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid start_time format: "+err.Error())
		return
	}
	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.Endtime, time.Local)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid end_time format: "+err.Error())
		return
	}
	timeout, err := time.ParseDuration(req.Timeout)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid timeout format: "+err.Error())
		return
	}

	id, err := h.module.CreateAuction(
		r.Context(),
		req.LotID,
		req.StartPrice,
		req.MinBidStep,
		req.SellerID,
		startTime,
		endTime,
		timeout,
	)
	if err != nil {
		log.Printf("error creating auction: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to create auction: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateAuctionResponse{ID: id})
}

func (h *HttpHandler) PlaceBid(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auctionID := vars["id"]
	if auctionID == "" {
		h.writeError(w, http.StatusBadRequest, "auction ID is required")
		return
	}

	var req PlaceBidRequest
	reader := json.NewDecoder(r.Body)
	err := reader.Decode(&req)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	err = h.module.PlaceBid(r.Context(), auctionID, req.UserID, req.Amount)
	if err != nil {
		log.Printf("error placing bid: %v", err)
		statusCode := http.StatusInternalServerError
		if errors.Is(err, domain.ErrBidTooLow) ||
			errors.Is(err, domain.ErrAuctionNotStarted) ||
			errors.Is(err, domain.ErrAuctionCancelled) ||
			errors.Is(err, domain.ErrAuctionCompleted) {
			statusCode = http.StatusBadRequest
		}
		h.writeError(w, statusCode, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(PlaceBidResponse{Message: "bid accepted for processing"})
}

func (h *HttpHandler) FindAuctionById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auctionID := vars["id"]
	if auctionID == "" {
		h.writeError(w, http.StatusBadRequest, "auction ID is required")
		return
	}

	auction, err := h.module.FindAuctionByID(r.Context(), auctionID)
	if err != nil {
		log.Printf("error finding auction: %v", err)
		statusCode := http.StatusInternalServerError
		if errors.Is(err, read_model.ErrAuctionReadModelNotFound) {
			statusCode = http.StatusNotFound
		}
		h.writeError(w, statusCode, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(auction)
}

func (h *HttpHandler) CancelAuction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auctionID := vars["id"]
	if auctionID == "" {
		h.writeError(w, http.StatusBadRequest, "auction ID is required")
		return
	}

	var req CancelAuctionRequest
	reader := json.NewDecoder(r.Body)
	err := reader.Decode(&req)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	err = h.module.CancelAuction(r.Context(), auctionID, req.SellerID, req.Reason)
	if err != nil {
		log.Printf("error cancelling auction: %v", err)
		statusCode := http.StatusInternalServerError
		if errors.Is(err, domain.ErrAuctionCancelled) ||
			errors.Is(err, domain.ErrAuctionCompleted) {
			statusCode = http.StatusBadRequest
		} else if errors.Is(err, domain.ErrSellerIDInvalid) {
			statusCode = http.StatusForbidden
		}
		h.writeError(w, statusCode, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CancelAuctionResponse{Message: "auction cancelled successfully"})
}

func (h *HttpHandler) CreateLot(w http.ResponseWriter, r *http.Request) {
	var req CreateLotRequest
	reader := json.NewDecoder(r.Body)
	err := reader.Decode(&req)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	lotID, err := h.module.CreateLot(r.Context(), req.Name, req.Description, req.OwnerID)
	if err != nil {
		log.Printf("error creating lot: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to create lot: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateLotResponse{ID: lotID})
}
