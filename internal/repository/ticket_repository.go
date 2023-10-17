package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/arynskiii/help_desk/models"
	"github.com/jmoiron/sqlx"
)

type TicketMySQL struct {
	db *sqlx.DB
}

func NewTicketMySQL(db sqlx.DB) *TicketMySQL {
	return &TicketMySQL{
		db: &db,
	}
}

func (c *TicketMySQL) CreateCategory(ctx context.Context, category models.Category) (int64, error) {
	var query = "INSERT INTO category (title,user_id,description)  VALUES (?,?,?)"
	res, err := c.db.ExecContext(ctx, query, category.Title, category.UserId, category.Description)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *TicketMySQL) DeleteCategory(ctx context.Context, categoryId int64) error {

	query := "DELETE from  category where id=?"
	_, err := c.db.ExecContext(ctx, query, categoryId)
	if err != nil {
		return err
	}
	return nil
}

func (c *TicketMySQL) GetCategoryByID(ctx context.Context, categoryId int64) (models.Category, error) {
	var category models.Category
	query := "SELECT id,title,description,user_id from category where id=?"
	row := c.db.QueryRowContext(ctx, query, categoryId)
	if err := row.Scan(&category.Id, &category.Title, &category.Description, &category.UserId); err != nil {
		return category, err
	}
	return category, nil
}
func (c *TicketMySQL) CreateTicket(ctx context.Context, ticket models.Ticket) (int64, error) {
	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %v", err)
	}

	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("failed to rollback transaction: %v", rollbackErr)
		}
	}()

	query := "INSERT INTO ticket (title, description, deadline, user_id, category_id) VALUES (?, ?, ?, ?, ?)"
	res, err := tx.ExecContext(ctx, query, ticket.Title, ticket.Description, ticket.Deadline, ticket.UserId, ticket.CategoryId)
	if err != nil {
		return 0, fmt.Errorf("failed to execute ticket insertion: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get the ticket ID: %v", err)
	}

	history := models.History{
		TicketId:   int64(id),
		Deadline:   *ticket.Deadline,
		SenderId:   int64(ticket.UserId),
		ReceiverId: 1,
		StateId:    2,
	}
	if err := c.CreateHistory(ctx, history, tx); err != nil {

		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return id, nil
}

func (c *TicketMySQL) CreateHistory(ctx context.Context, history models.History, tx *sqlx.Tx) error {
	query := `
		INSERT INTO history_ticket (ticket_id, state_id,  deadline, receiver_id, sender_id)
		VALUES (?, ?,  ?, ?, ?)
	`
	_, err := tx.ExecContext(ctx, query, history.TicketId, history.StateId, history.Deadline, history.ReceiverId, history.SenderId)
	if err != nil {
		return fmt.Errorf("failed to execute history insertion: %v", err)
	}
	return nil
}

func (c *TicketMySQL) ShowAllTicket(ctx context.Context) ([]models.Ticket, error) {
	query := "SELECT * FROM tickets"
	rows, err := c.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var tickets []models.Ticket
	for rows.Next() {
		ticket := new(models.Ticket)
		if err := rows.Scan(&ticket.Id, &ticket.CategoryId, &ticket.UserId, &ticket.Title, &ticket.Description, &ticket.Deadline, &ticket.CreateAt); err != nil {
			return nil, err
		}

		tickets = append(tickets, *ticket)
	}
	return tickets, nil
}
func (c *TicketMySQL) UpdateHistoryState(ctx context.Context, stateId int64, historyId int64) error {
	query := "UPDATE history_ticket SET state_id=? WHERE id=?"
	_, err := c.db.ExecContext(ctx, query, stateId, historyId)
	if err != nil {
		return err
	}
	return nil
}

func (c *TicketMySQL) CreateDocAttachment(ctx context.Context, doc models.Attachments) ([]int64, error) {

	var insertedIDs []int64

	query := "INSERT INTO ticket_attachments (id, ticket_id, attachment_path) VALUES (?, ?, ?)"

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	for _, attachmentPath := range *doc.Attachment_path {

		res, err := tx.ExecContext(ctx, query, doc.Id, doc.TicketId, attachmentPath)
		if err != nil {
			return nil, err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}

		insertedIDs = append(insertedIDs, id)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return insertedIDs, nil
}

func (c *TicketMySQL) GetTicketByID(ctx context.Context, ticketId int64) (models.Ticket, error) {
	var ticket models.Ticket
	query := `
	SELECT
		t.id,
		t.category_id,
		t.user_id,
		t.title,
		t.description,
		CONVERT_TZ(t.deadline, '+00:00', '+06:00') AS deadline,
		CONVERT_TZ(t.create_at, '+00:00', '+06:00') AS create_at
	FROM
		ticket t
	WHERE
		t.id = ?;
	`

	row := c.db.QueryRowContext(ctx, query, ticketId)

	var deadlineStr string
	var createAtStr string
	if err := row.Scan(&ticket.Id, &ticket.CategoryId, &ticket.UserId, &ticket.Title, &ticket.Description, &deadlineStr, &createAtStr); err != nil {
		return ticket, err
	}

	if deadlineStr != "" {
		deadline, err := time.Parse("2006-01-02 15:04:05", deadlineStr)
		if err != nil {
			return ticket, err
		}
		ticket.Deadline = &deadline
	}

	if createAtStr != "" {
		createAt, err := time.Parse("2006-01-02 15:04:05", createAtStr)
		if err != nil {
			return ticket, err
		}
		ticket.CreateAt = createAt
	}
	var err error
	ticket.Attachment_path, err = c.GetAttachmentsPathByTicketID(ctx, ticket.Id)
	if err != nil {
		return ticket, err
	}
	return ticket, nil
}

func (c *TicketMySQL) GetAttachmentsPathByTicketID(ctx context.Context, ticketID int) ([]string, error) {
	var res []string
	query := `SELECT attachment_path FROM ticket_attachments WHERE ticket_id=? `

	rows, err := c.db.QueryContext(ctx, query, ticketID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var attachmentPath string
		if err := rows.Scan(&attachmentPath); err != nil {
			return nil, err
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		res = append(res, attachmentPath)
	}
	return res, nil
}
