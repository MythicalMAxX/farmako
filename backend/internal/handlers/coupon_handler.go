package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/farmako/coupon-system/internal/models"
	"github.com/farmako/coupon-system/internal/services"
	"github.com/go-chi/chi/v5"
)

// CouponHandler handles HTTP requests for coupons
type CouponHandler struct {
	service *services.CouponService
}

// NewCouponHandler creates a new coupon handler
func NewCouponHandler(service *services.CouponService) *CouponHandler {
	return &CouponHandler{
		service: service,
	}
}

// RegisterRoutes registers routes for the coupon handler
func (h *CouponHandler) RegisterRoutes(r chi.Router) {
	r.Route("/coupons", func(r chi.Router) {
		r.Get("/", h.ListCoupons)
		r.Post("/", h.CreateCoupon)
		r.Get("/{id}", h.GetCoupon)
		r.Put("/{id}", h.UpdateCoupon)
		r.Delete("/{id}", h.DeleteCoupon)
		r.Get("/applicable", h.GetApplicableCoupons)
		r.Post("/validate", h.ValidateCoupon)
		r.Post("/apply", h.ApplyCoupon)
	})
}

// Helper functions for JSON responses
func respondJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

// CreateCoupon creates a new coupon
func (h *CouponHandler) CreateCoupon(w http.ResponseWriter, r *http.Request) {
	var coupon models.Coupon

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&coupon); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Ensure valid expiry date
	if coupon.ExpiryDate.IsZero() || coupon.ExpiryDate.Before(time.Now()) {
		// Set default expiry date to 1 year from now if not provided or invalid
		coupon.ExpiryDate = time.Now().AddDate(1, 0, 0)
	}

	if err := h.service.CreateCoupon(r.Context(), &coupon); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, coupon)
}

// GetCoupon gets a coupon by ID
func (h *CouponHandler) GetCoupon(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid coupon ID")
		return
	}

	// For simplicity, we'll list all coupons and filter
	coupons, err := h.service.ListCoupons(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, coupon := range coupons {
		if coupon.ID == uint(id) {
			respondJSON(w, http.StatusOK, coupon)
			return
		}
	}

	respondError(w, http.StatusNotFound, "Coupon not found")
}

// UpdateCoupon updates an existing coupon
func (h *CouponHandler) UpdateCoupon(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid coupon ID")
		return
	}

	var coupon models.Coupon
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&coupon); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Ensure ID matches route parameter
	coupon.ID = uint(id)

	if err := h.service.UpdateCoupon(r.Context(), &coupon); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, coupon)
}

// DeleteCoupon deletes a coupon by ID
func (h *CouponHandler) DeleteCoupon(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid coupon ID")
		return
	}

	if err := h.service.DeleteCoupon(r.Context(), uint(id)); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListCoupons lists all coupons
func (h *CouponHandler) ListCoupons(w http.ResponseWriter, r *http.Request) {
	coupons, err := h.service.ListCoupons(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, coupons)
}

// ValidateCoupon validates a coupon against a cart
func (h *CouponHandler) ValidateCoupon(w http.ResponseWriter, r *http.Request) {
	var req models.ValidationRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Set timestamp if not provided
	if req.Timestamp.IsZero() {
		req.Timestamp = time.Now()
	}

	// Basic validation
	if req.CouponCode == "" {
		respondError(w, http.StatusBadRequest, "Coupon code is required")
		return
	}

	response, err := h.service.ValidateCoupon(r.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrCouponNotFound) {
			respondError(w, http.StatusNotFound, "Coupon not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, response)
}

// ApplyCoupon applies a coupon to an order
func (h *CouponHandler) ApplyCoupon(w http.ResponseWriter, r *http.Request) {
	var req models.ValidationRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Set timestamp if not provided
	if req.Timestamp.IsZero() {
		req.Timestamp = time.Now()
	}

	// Basic validation
	if req.CouponCode == "" {
		respondError(w, http.StatusBadRequest, "Coupon code is required")
		return
	}

	if req.UserID == "" {
		respondError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	// Generate a fake order ID for demonstration
	orderID := "order_" + strconv.FormatInt(time.Now().Unix(), 10)

	response, err := h.service.ApplyCoupon(r.Context(), &req, orderID)
	if err != nil {
		if errors.Is(err, services.ErrCouponNotFound) {
			respondError(w, http.StatusNotFound, "Coupon not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, response)
}

// GetApplicableCoupons returns applicable coupons for a cart
func (h *CouponHandler) GetApplicableCoupons(w http.ResponseWriter, r *http.Request) {
	var req models.ApplicableCouponsRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Set timestamp if not provided
	if req.Timestamp.IsZero() {
		req.Timestamp = time.Now()
	}

	response, err := h.service.GetApplicableCoupons(r.Context(), &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, response)
}
