// Package link provides link shortening and retrieval services.
// It handles URL shortening with unique code generation and URL retrieval operations.
package link

import (
	"context"

	linkDTO "github.com/huypham67/bookmark-service/internal/dto/link"
	"github.com/huypham67/bookmark-service/internal/repository/link"
	"github.com/huypham67/bookmark-service/pkg/utils"
)

const shortCodeLength = 7

// Service defines the contract for link operations.
//
//go:generate mockery --name=Service --output=./mocks --outpkg=mocks --filename=mock_service.go
type Service interface {
	ShortenURL(ctx context.Context, request linkDTO.ShortenURLRequest) (string, error)
	GetOriginalURL(ctx context.Context, code string) (string, error)
}

type service struct {
	linkRepo      link.Repository
	codeGenerator utils.CodeGenerator
}

// NewService creates a new link service with the provided repository and code generator.
func NewService(linkRepo link.Repository, codeGenerator utils.CodeGenerator) Service {
	return &service{
		linkRepo:      linkRepo,
		codeGenerator: codeGenerator,
	}
}
