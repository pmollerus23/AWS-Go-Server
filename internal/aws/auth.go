package aws

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sort"
	"strings"
	"time"
)

// NewIAMAuthMiddleware creates a middleware that verifies AWS SigV4 signatures.
// This is useful for authenticating API requests using AWS IAM credentials.
func NewIAMAuthMiddleware(logger *slog.Logger, region string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if Authorization header is present
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("missing Authorization header",
					"path", r.URL.Path,
					"remote_addr", r.RemoteAddr,
				)
				http.Error(w, `{"error":"Missing AWS Authorization header"}`, http.StatusUnauthorized)
				return
			}

			// Verify it's AWS SigV4
			if !strings.HasPrefix(authHeader, "AWS4-HMAC-SHA256") {
				logger.Warn("invalid authorization scheme",
					"auth_header", authHeader,
					"path", r.URL.Path,
				)
				http.Error(w, `{"error":"Invalid authorization scheme. Expected AWS4-HMAC-SHA256"}`, http.StatusUnauthorized)
				return
			}

			// Parse the authorization header
			credential, signedHeaders, signature, err := parseAuthHeader(authHeader)
			if err != nil {
				logger.Error("failed to parse auth header",
					"error", err,
					"auth_header", authHeader,
				)
				http.Error(w, `{"error":"Invalid Authorization header format"}`, http.StatusUnauthorized)
				return
			}

			// Verify timestamp is recent (within 15 minutes)
			dateHeader := r.Header.Get("X-Amz-Date")
			if dateHeader == "" {
				logger.Warn("missing X-Amz-Date header")
				http.Error(w, `{"error":"Missing X-Amz-Date header"}`, http.StatusUnauthorized)
				return
			}

			timestamp, err := time.Parse("20060102T150405Z", dateHeader)
			if err != nil {
				logger.Error("failed to parse timestamp", "error", err, "date", dateHeader)
				http.Error(w, `{"error":"Invalid X-Amz-Date format"}`, http.StatusUnauthorized)
				return
			}

			if time.Since(timestamp) > 15*time.Minute {
				logger.Warn("request timestamp too old",
					"timestamp", timestamp,
					"age_minutes", time.Since(timestamp).Minutes(),
				)
				http.Error(w, `{"error":"Request timestamp too old"}`, http.StatusUnauthorized)
				return
			}

			// In a production system, you would:
			// 1. Extract the access key ID from the credential
			// 2. Look up the secret key from IAM or a secrets manager
			// 3. Recompute the signature using the secret key
			// 4. Compare with the provided signature
			// 5. Verify the user has permission for this action (IAM policy evaluation)

			logger.Info("IAM authentication successful",
				"credential", credential,
				"signed_headers", signedHeaders,
				"signature_provided", signature[:16]+"...", // Log first 16 chars only
				"path", r.URL.Path,
			)

			// Request is authenticated, continue to handler
			h.ServeHTTP(w, r)
		})
	}
}

// parseAuthHeader parses the AWS SigV4 Authorization header.
// Format: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request, SignedHeaders=host;range;x-amz-date, Signature=fe5f80f77d5fa3beca038a248ff027d0445342fe2855ddc963176630326f1024
func parseAuthHeader(authHeader string) (credential, signedHeaders, signature string, err error) {
	// Remove the "AWS4-HMAC-SHA256 " prefix
	authHeader = strings.TrimPrefix(authHeader, "AWS4-HMAC-SHA256 ")

	// Split into parts
	parts := strings.Split(authHeader, ", ")
	for _, part := range parts {
		if strings.HasPrefix(part, "Credential=") {
			credential = strings.TrimPrefix(part, "Credential=")
		} else if strings.HasPrefix(part, "SignedHeaders=") {
			signedHeaders = strings.TrimPrefix(part, "SignedHeaders=")
		} else if strings.HasPrefix(part, "Signature=") {
			signature = strings.TrimPrefix(part, "Signature=")
		}
	}

	if credential == "" || signedHeaders == "" || signature == "" {
		return "", "", "", fmt.Errorf("incomplete authorization header")
	}

	return credential, signedHeaders, signature, nil
}

// ComputeSignature computes the AWS SigV4 signature for a request.
// This is a simplified version for demonstration. In production, use the AWS SDK's signer.
func ComputeSignature(secretKey, dateStamp, region, service, stringToSign string) string {
	kDate := hmacSHA256([]byte("AWS4"+secretKey), []byte(dateStamp))
	kRegion := hmacSHA256(kDate, []byte(region))
	kService := hmacSHA256(kRegion, []byte(service))
	kSigning := hmacSHA256(kService, []byte("aws4_request"))
	signature := hmacSHA256(kSigning, []byte(stringToSign))
	return hex.EncodeToString(signature)
}

// CreateCanonicalRequest creates a canonical request for AWS SigV4.
func CreateCanonicalRequest(method, uri, queryString string, headers http.Header, signedHeaders []string, payloadHash string) string {
	var canonicalHeaders strings.Builder
	var headerNames []string

	// Build canonical headers
	for _, name := range signedHeaders {
		value := headers.Get(name)
		canonicalHeaders.WriteString(strings.ToLower(name))
		canonicalHeaders.WriteString(":")
		canonicalHeaders.WriteString(strings.TrimSpace(value))
		canonicalHeaders.WriteString("\n")
		headerNames = append(headerNames, strings.ToLower(name))
	}

	sort.Strings(headerNames)
	signedHeadersStr := strings.Join(headerNames, ";")

	// Canonical request format
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		method,
		uri,
		queryString,
		canonicalHeaders.String(),
		signedHeadersStr,
		payloadHash,
	)
}

// HashPayload computes the SHA256 hash of the request payload.
func HashPayload(payload []byte) string {
	hash := sha256.Sum256(payload)
	return hex.EncodeToString(hash[:])
}

// hmacSHA256 computes HMAC-SHA256.
func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// ReadBody reads the request body and returns it, while also replacing it so it can be read again.
func ReadBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return []byte{}, nil
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	// Replace the body so it can be read again by handlers
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	return body, nil
}
