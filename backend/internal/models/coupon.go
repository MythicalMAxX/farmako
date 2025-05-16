package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// DiscountType represents the type of discount (percentage or fixed amount)
type DiscountType string

const (
	PercentageDiscount DiscountType = "percentage"
	FixedDiscount      DiscountType = "fixed"
)

// UsageType represents how a coupon can be used
type UsageType string

const (
	OneTime   UsageType = "one_time"
	MultiUse  UsageType = "multi_use"
	TimeBased UsageType = "time_based"
)

// StringArray is a custom type that implements database serialization methods
type StringArray []string

// Scan implements the sql.Scanner interface
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = StringArray{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		// Try to unmarshal as JSON array
		if err := json.Unmarshal(v, s); err != nil {
			// Not valid JSON, store the string directly
			*s = StringArray{string(v)}
		}
		return nil
	case string:
		// Try to unmarshal as JSON array
		if err := json.Unmarshal([]byte(v), s); err != nil {
			// Not valid JSON, store the string directly
			*s = StringArray{v}
		}
		return nil
	default:
		return errors.New("unsupported Scan source")
	}
}

// Value implements the driver.Valuer interface
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// Coupon represents a discount coupon in the system
type Coupon struct {
	ID                    uint         `json:"id" gorm:"primaryKey"`
	CouponCode            string       `json:"coupon_code" gorm:"uniqueIndex;not null"`
	ExpiryDate            time.Time    `json:"expiry_date"`
	UsageType             UsageType    `json:"usage_type"`
	ApplicableMedicineIDs StringArray  `json:"applicable_medicine_ids" gorm:"type:text"`
	ApplicableCategories  StringArray  `json:"applicable_categories" gorm:"type:text"`
	MinOrderValue         float64      `json:"min_order_value"`
	ValidTimeWindowStart  *time.Time   `json:"valid_time_window_start,omitempty"`
	ValidTimeWindowEnd    *time.Time   `json:"valid_time_window_end,omitempty"`
	TermsAndConditions    string       `json:"terms_and_conditions"`
	DiscountType          DiscountType `json:"discount_type"`
	DiscountValue         float64      `json:"discount_value"`
	MaxUsagePerUser       int          `json:"max_usage_per_user"`
	CreatedAt             time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
}

// CouponUsage tracks coupon usage by users
type CouponUsage struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	CouponID uint      `json:"coupon_id" gorm:"index"`
	UserID   string    `json:"user_id" gorm:"index"`
	UsedAt   time.Time `json:"used_at" gorm:"autoCreateTime"`
	OrderID  string    `json:"order_id"`
}

// CartItem represents an item in the user's cart
type CartItem struct {
	ID       string  `json:"id"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

// ValidationRequest represents a request to validate a coupon
type ValidationRequest struct {
	CouponCode string     `json:"coupon_code"`
	CartItems  []CartItem `json:"cart_items"`
	OrderTotal float64    `json:"order_total"`
	UserID     string     `json:"user_id"`
	Timestamp  time.Time  `json:"timestamp"`
}

// ValidationResponse represents the response to a coupon validation request
type ValidationResponse struct {
	IsValid  bool      `json:"is_valid"`
	Discount *Discount `json:"discount,omitempty"`
	Message  string    `json:"message,omitempty"`
	Reason   string    `json:"reason,omitempty"`
}

// Discount represents the discount applied to an order
type Discount struct {
	ItemsDiscount   float64 `json:"items_discount"`
	ChargesDiscount float64 `json:"charges_discount"`
	TotalDiscount   float64 `json:"total_discount"`
}

// ApplicableCouponsRequest represents a request to get applicable coupons
type ApplicableCouponsRequest struct {
	CartItems  []CartItem `json:"cart_items"`
	OrderTotal float64    `json:"order_total"`
	UserID     string     `json:"user_id"`
	Timestamp  time.Time  `json:"timestamp"`
}

// ApplicableCouponsResponse represents a response with applicable coupons
type ApplicableCouponsResponse struct {
	ApplicableCoupons []ApplicableCoupon `json:"applicable_coupons"`
}

// ApplicableCoupon is a simplified coupon representation for client response
type ApplicableCoupon struct {
	CouponCode    string       `json:"coupon_code"`
	DiscountValue float64      `json:"discount_value"`
	DiscountType  DiscountType `json:"discount_type"`
	ExpiryDate    time.Time    `json:"expiry_date"`
	Description   string       `json:"description,omitempty"`
}
