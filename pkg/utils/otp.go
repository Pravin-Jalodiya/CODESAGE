package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"
)

// In-memory OTP store (for demonstration)
var otpStore = struct {
	sync.RWMutex
	data   map[string]string
	expiry map[string]time.Time
}{
	data:   make(map[string]string),
	expiry: make(map[string]time.Time),
}

// GenerateOTP generates a 6-digit OTP
func GenerateOTP() (string, error) {
	var otp string
	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp += fmt.Sprintf("%d", num)
	}
	return otp, nil
}

// SaveOTP stores the OTP for a specific email with a 5-minute expiry
func SaveOTP(email, otp string) {
	otpStore.Lock()
	defer otpStore.Unlock()
	otpStore.data[email] = otp
	otpStore.expiry[email] = time.Now().Add(5 * time.Minute)
}

// ValidateOTP validates the OTP for the given email. If valid, it deletes the OTP.
func ValidateOTP(email, otp string) bool {
	otpStore.RLock()
	savedOtp, exists := otpStore.data[email]
	expiryTime := otpStore.expiry[email]
	otpStore.RUnlock()
	if !exists || savedOtp != otp || time.Now().After(expiryTime) {
		return false
	}
	// Remove OTP after successful validation
	otpStore.Lock()
	delete(otpStore.data, email)
	delete(otpStore.expiry, email)
	otpStore.Unlock()
	return true
}
