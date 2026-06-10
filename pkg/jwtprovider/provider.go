package jwtprovider

import (
	"fmt"

	"github.com/huypham67/bookmark-service-monolithic/pkg/jwt"
)

// Provider encapsulates JWT token generation and validation capabilities.
type Provider interface {
	Generator() jwt.TokenGenerator
	Validator() jwt.TokenValidator
}

type jwtProvider struct {
	generator jwt.TokenGenerator
	validator jwt.TokenValidator
}

// NewProvider initializes a new JWT provider from environment configuration.
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
	generator, err := jwt.NewTokenGenerator(privateKey, cfg.Issuer, cfg.Audience, cfg.ExpirationDuration())
	if err != nil {
		return nil, fmt.Errorf("failed to create token generator: %w", err)
	}

	// Create TokenValidator
	validator, err := jwt.NewTokenValidator(publicKey, cfg.Issuer, cfg.Audience)
	if err != nil {
		return nil, fmt.Errorf("failed to create token validator: %w", err)
	}

	return &jwtProvider{
		generator: generator,
		validator: validator,
	}, nil
}

// Generator returns the TokenGenerator instance from the provider.
func (p *jwtProvider) Generator() jwt.TokenGenerator {
	return p.generator
}

// Validator returns the TokenValidator instance from the provider.
func (p *jwtProvider) Validator() jwt.TokenValidator {
	return p.validator
}
