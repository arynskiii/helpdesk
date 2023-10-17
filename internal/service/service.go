package service

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/arynskiii/help_desk/internal/repository"
	"github.com/arynskiii/help_desk/models"
)

type Authorization interface {
	CreateUser(ctx context.Context, user models.User) (int64, error)
	GenerateToken(name string, password string) (models.User, error)

	GetUserByToken(token string) (models.User, error)
}

type Ticket interface {
	CreateCategory(ctx context.Context, category models.Category) (int64, error)
	GetCategoryByID(ctx context.Context, categoryId int64) (models.Category, error)
	DeleteCategory(ctx context.Context, categoryId int64, userId int64) error
	CreateTicket(ctx context.Context, ticket models.Ticket) (int64, error)
	UpdateHistoryState(ctx context.Context, stateId int64, historyId int64) error
	ShowallTicket(ctx context.Context) ([]models.Ticket, error)
	CreateDocAttachment(ctx context.Context, doc models.Attachments) ([]int64, error)
	PostFiles(files []multipart.File, subrootFolder string, targetURL string, fileNames []string) ([]string, error)
	DownloadFiles(dbPaths []string, url string) ([]io.ReadCloser, []error)
	GetTicketByID(ctx context.Context, ticketId int64) (models.Ticket, error)
}

type Service struct {
	Authorization
	Ticket
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Ticket:        NewCategoryService(repo.Ticket),
	}
}
