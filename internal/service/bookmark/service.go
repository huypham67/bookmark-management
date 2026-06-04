package bookmark

import (
	"context"
	"errors"

	bookmarkDTO "github.com/huypham67/bookmark-service/internal/dto/bookmark"
	"github.com/huypham67/bookmark-service/internal/model"
	"github.com/huypham67/bookmark-service/internal/repository/bookmark"
	"github.com/huypham67/bookmark-service/pkg/utils"
)

var (
	ErrInternalServerError = errors.New("internal server error")
)

const bookmarkCodeLength = 6

// PaginationResult holds pagination metadata.
type PaginationResult struct {
	Total int64
}

// Service defines the interface for bookmark operations.
//
//go:generate mockery --name=Service --output=./mocks --outpkg=mocks --filename=mock_service.go
type Service interface {
	Create(ctx context.Context, userID string, req bookmarkDTO.CreateBookmarkRequest) (*model.Bookmark, error)
	List(ctx context.Context, userID string, page, limit int64, sort string) ([]*model.Bookmark, *PaginationResult, error)
}

type service struct {
	bookmarkRepo  bookmark.Repository
	codeGenerator utils.CodeGenerator
}

// NewService creates a new instance of the bookmark service with the provided dependencies.
func NewService(
	bookmarkRepo bookmark.Repository,
	codeGenerator utils.CodeGenerator,
) Service {
	return &service{
		bookmarkRepo:  bookmarkRepo,
		codeGenerator: codeGenerator,
	}
}
