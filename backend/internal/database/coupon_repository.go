package database

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/farmako/coupon-system/internal/models"
	"gorm.io/gorm"
)

// CouponRepository handles database operations for coupons
type CouponRepository struct {
	db    *gorm.DB
	mutex sync.RWMutex
}

// NewCouponRepository creates a new coupon repository
func NewCouponRepository(conn *DBConn) *CouponRepository {
	return &CouponRepository{
		db: conn.DB,
	}
}

// CreateCoupon creates a new coupon
func (r *CouponRepository) CreateCoupon(ctx context.Context, coupon *models.Coupon) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// With our new StringArray type, GORM will automatically handle the JSON serialization/deserialization
	// We don't need manual JSON handling anymore
	return r.db.WithContext(ctx).Create(coupon).Error
}

// GetCouponByCode fetches a coupon by its code
func (r *CouponRepository) GetCouponByCode(ctx context.Context, code string) (*models.Coupon, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var coupon models.Coupon
	if err := r.db.WithContext(ctx).Where("coupon_code = ?", code).First(&coupon).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// The StringArray fields are automatically handled by GORM now

	return &coupon, nil
}

// UpdateCoupon updates an existing coupon
func (r *CouponRepository) UpdateCoupon(ctx context.Context, coupon *models.Coupon) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// The StringArray fields are automatically handled by GORM now
	return r.db.WithContext(ctx).Save(coupon).Error
}

// DeleteCoupon deletes a coupon by ID
func (r *CouponRepository) DeleteCoupon(ctx context.Context, id uint) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.db.WithContext(ctx).Delete(&models.Coupon{}, id).Error
}

// ListCoupons returns all coupons
func (r *CouponRepository) ListCoupons(ctx context.Context) ([]models.Coupon, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var coupons []models.Coupon
	if err := r.db.WithContext(ctx).Find(&coupons).Error; err != nil {
		return nil, err
	}

	// The StringArray fields are automatically handled by GORM now

	return coupons, nil
}

// RecordCouponUsage records the usage of a coupon by a user
func (r *CouponRepository) RecordCouponUsage(ctx context.Context, couponID uint, userID, orderID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	usage := models.CouponUsage{
		CouponID: couponID,
		UserID:   userID,
		OrderID:  orderID,
		UsedAt:   time.Now(),
	}

	return r.db.WithContext(ctx).Create(&usage).Error
}

// GetCouponUsageCount returns the number of times a user has used a coupon
func (r *CouponRepository) GetCouponUsageCount(ctx context.Context, couponID uint, userID string) (int64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var count int64
	err := r.db.WithContext(ctx).Model(&models.CouponUsage{}).
		Where("coupon_id = ? AND user_id = ?", couponID, userID).
		Count(&count).Error

	return count, err
}

// GetApplicableCoupons returns all valid coupons for a given cart and user
func (r *CouponRepository) GetApplicableCoupons(ctx context.Context, req *models.ApplicableCouponsRequest) ([]models.Coupon, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var coupons []models.Coupon

	// Get active coupons that haven't expired
	if err := r.db.WithContext(ctx).
		Where("expiry_date > ?", req.Timestamp).
		Where("min_order_value <= ?", req.OrderTotal).
		Find(&coupons).Error; err != nil {
		return nil, err
	}

	// The StringArray fields are automatically handled by GORM now

	return coupons, nil
}
