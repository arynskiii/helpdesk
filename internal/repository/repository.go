package repository

import (
	"context"

	"github.com/arynskiii/help_desk/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(ctx context.Context, user models.User) (int64, error)
	GetUser(name string) (models.User, error)
	SaveTokens(name string, token string) error
	GetUserByToken(token string) (models.User, error)
}

type Ticket interface {
	CreateCategory(ctx context.Context, category models.Category) (int64, error)
	GetCategoryByID(ctx context.Context, categoryId int64) (models.Category, error)
	DeleteCategory(ctx context.Context, categoryId int64) error
	CreateTicket(ctx context.Context, ticket models.Ticket) (int64, error)
	UpdateHistoryState(ctx context.Context, stateId int64, historyId int64) error
	ShowAllTicket(ctx context.Context) ([]models.Ticket, error)
	GetTicketByID(ctx context.Context, ticketId int64) (models.Ticket, error)
	CreateDocAttachment(ctx context.Context, doc models.Attachments) ([]int64, error)
}

type Repository struct {
	Authorization
	Ticket
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthMySQL(*db),
		Ticket:        NewTicketMySQL(*db),
	}
}
