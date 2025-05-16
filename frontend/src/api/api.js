import axios from 'axios';

// Create axios instance with base URL from environment variable
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add response interceptor for consistent error handling
api.interceptors.response.use(
  response => response.data,
  error => {
    console.error('API Error:', error);
    // If the backend is unavailable or returns an error, handle gracefully
    return Promise.reject(error);
  }
);

// Coupon API functions
export const couponsApi = {
  // Get all coupons
  getAllCoupons: async () => {
    try {
      const response = await api.get('/coupons');
      return response;
    } catch (error) {
      console.log('getAllCoupons error:', error);
      // Return an empty array instead of throwing when the API fails
      return [];
    }
  },

  // Get coupon by ID
  getCouponById: async (id) => {
    try {
      const response = await api.get(`/coupons/${id}`);
      return response;
    } catch (error) {
      console.log(`getCouponById(${id}) error:`, error);
      return null;
    }
  },

  // Create new coupon
  createCoupon: async (couponData) => {
    const response = await api.post('/coupons', couponData);
    return response;
  },

  // Update coupon
  updateCoupon: async (id, couponData) => {
    const response = await api.put(`/coupons/${id}`, couponData);
    return response;
  },

  // Delete coupon
  deleteCoupon: async (id) => {
    const response = await api.delete(`/coupons/${id}`);
    return response;
  },

  // Validate coupon
  validateCoupon: async (validationData) => {
    try {
      const response = await api.post('/coupons/validate', validationData);
      return response;
    } catch (error) {
      console.log('validateCoupon error:', error);
      // Return a default invalid response
      return {
        is_valid: false,
        reason: 'Failed to validate coupon. Please try again later.'
      };
    }
  },

  // Apply coupon
  applyCoupon: async (applyData) => {
    try {
      const response = await api.post('/coupons/apply', applyData);
      return response;
    } catch (error) {
      console.log('applyCoupon error:', error);
      // Return a default error response
      return {
        is_valid: false,
        reason: 'Failed to apply coupon. Please try again later.'
      };
    }
  },

  // Get applicable coupons
  getApplicableCoupons: async (cartData) => {
    try {
      const response = await api.get('/coupons/applicable', { params: cartData });
      return response;
    } catch (error) {
      console.log('getApplicableCoupons error:', error);
      // Return an empty list
      return { applicable_coupons: [] };
    }
  },
};

export default api; 