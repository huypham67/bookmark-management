package jwtutils

import (
	"fmt"
)

// Provider encapsulates JWT token generation and validation capabilities.
// It loads RSA keys from configured paths and provides token generator/validator instances.
type Provider interface {
	Generator() TokenGenerator
	Validator() TokenValidator
}

type jwtProvider struct {
	generator TokenGenerator
	validator TokenValidator
}

// NewProvider initializes a new JWT provider from environment configuration.
// It loads configuration from environment variables, validates it, loads RSA keys,
// and creates TokenGenerator and TokenValidator instances.
//
// Parameters:
//   - envPrefix: Environment variable prefix (e.g., "" for default, "STAGING_" for staging)
//
// Returns an error if:
// - Configuration loading fails
// - Configuration validation fails
// - Key files cannot be read or are invalid
// - Key parsing fails
func NewProvider(envPrefix string) (Provider, error) {
	cfg, err := LoadJWTConfig(envPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to load jwt config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid jwt config: %w", err)
	}

	// Load RSA keys from configured paths
	privateKey, err := LoadRSAPrivateKeyFromFile(cfg.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	publicKey, err := LoadRSAPublicKeyFromFile(cfg.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load public key: %w", err)
	}

	// Create TokenGenerator
	generator, err := NewTokenGenerator(privateKey, cfg.Issuer, cfg.Audience, cfg.ExpirationDuration())
	if err != nil {
		return nil, fmt.Errorf("failed to create token generator: %w", err)
	}

	// Create TokenValidator
	validator, err := NewTokenValidator(publicKey, cfg.Issuer, cfg.Audience)
	if err != nil {
		return nil, fmt.Errorf("failed to create token validator: %w", err)
	}

	return &jwtProvider{
		generator: generator,
		validator: validator,
	}, nil
}

// Generator returns the TokenGenerator instance from the provider.
func (p *jwtProvider) Generator() TokenGenerator {
	return p.generator
}

// Validator returns the TokenValidator instance from the provider.
func (p *jwtProvider) Validator() TokenValidator {
	return p.validator
}
