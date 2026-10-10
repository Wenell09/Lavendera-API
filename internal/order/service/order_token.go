package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// Clean charset for tracking token without ambiguous characters (0, O, 1, I).
// Simple for customers to read, remember, and type without confusion.
const trackingCharset = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"

// Charset for order number.
const orderNumberCharset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateTrackingToken generates a simple, human-friendly 6-character tracking token.
// Example: TRK-8K9N2P
func GenerateTrackingToken() (string, error) {
	const length = 6
	bytes := make([]byte, length)
	charsetLen := big.NewInt(int64(len(trackingCharset)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		bytes[i] = trackingCharset[num.Int64()]
	}

	return fmt.Sprintf("TRK-%s", string(bytes)), nil
}

// GenerateOrderNumber generates a unique order number formatted with the date and random suffix.
// Example: ORD-20261008-7K2M
func GenerateOrderNumber() (string, error) {
	const length = 4
	bytes := make([]byte, length)
	charsetLen := big.NewInt(int64(len(orderNumberCharset)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		bytes[i] = orderNumberCharset[num.Int64()]
	}

	dateStr := time.Now().Format("20060102")
	return fmt.Sprintf("ORD-%s-%s", dateStr, string(bytes)), nil
}

// GenerateReferenceNumber generates a unique payment reference number.
// Example: TRX-20261008-7K2M
func GenerateReferenceNumber() (string, error) {
	const length = 6
	bytes := make([]byte, length)
	charsetLen := big.NewInt(int64(len(orderNumberCharset)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		bytes[i] = orderNumberCharset[num.Int64()]
	}

	dateStr := time.Now().Format("20060102")
	return fmt.Sprintf("TRX-%s-%s", dateStr, string(bytes)), nil
}
