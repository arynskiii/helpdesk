package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/arynskiii/help_desk/internal/repository"
	"github.com/arynskiii/help_desk/models"

	"github.com/twinj/uuid"
)

type TicketService struct {
	repo repository.Ticket
}

func NewCategoryService(repo repository.Ticket) *TicketService {
	return &TicketService{
		repo: repo,
	}
}

func (s *TicketService) CreateCategory(ctx context.Context, category models.Category) (int64, error) {

	return s.repo.CreateCategory(ctx, category)
}
func (s *TicketService) GetCategoryByID(ctx context.Context, categoryId int64) (models.Category, error) {
	return s.repo.GetCategoryByID(ctx, categoryId)
}

func (s *TicketService) DeleteCategory(ctx context.Context, categoryID int64, userId int64) error {
	category, err := s.repo.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return err
	}
	if category.UserId != userId {
		return fmt.Errorf("NET PRAV!")
	}
	return s.repo.DeleteCategory(ctx, categoryID)
}

func (s *TicketService) CreateTicket(ctx context.Context, ticket models.Ticket) (int64, error) {
	return s.repo.CreateTicket(ctx, ticket)
}

func (s *TicketService) UpdateHistoryState(ctx context.Context, stateId int64, historyId int64) error {
	return s.repo.UpdateHistoryState(ctx, stateId, historyId)
}

func (s *TicketService) ShowallTicket(ctx context.Context) ([]models.Ticket, error) {
	return s.repo.ShowAllTicket(ctx)
}
func (s *TicketService) CreateDocAttachment(ctx context.Context, doc models.Attachments) ([]int64, error) {
	return s.repo.CreateDocAttachment(ctx, doc)
}

func (s *TicketService) GetTicketByID(ctx context.Context, ticketId int64) (models.Ticket, error) {
	return s.repo.GetTicketByID(ctx, ticketId)
}
func (s *TicketService) PostFiles(files []multipart.File, subrootFolder string, targetURL string, fileNames []string) ([]string, error) {
	fmt.Println("HTTP client post files started . . . ")

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for i, file := range files {
		err := bodyWriter.WriteField("subrootfolder", subrootFolder)
		if err != nil {
			fmt.Println("Error creating form field")
			return nil, err
		}

		fileWriter, err := bodyWriter.CreateFormFile("file", "upload_"+uuid.NewV4().String()+filepath.Ext(fileNames[i]))
		if err != nil {
			fmt.Println("Error writing to buffer")
			return nil, err
		}

		// Copy file content to the request body
		if _, err := io.Copy(fileWriter, file); err != nil {
			return nil, err
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetURL, contentType, bodyBuf)
	if err != nil {
		fmt.Printf("File endpoint error: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(resp.Status)
	fmt.Println(string(respBody))
	fmt.Println("HTTP client post files called and successfully terminated . . .")

	var uploadedPaths []string
	responseString := string(respBody)
	for range files {
		uploadedPaths = append(uploadedPaths, responseString)
	}

	return uploadedPaths, nil
}

func (s *TicketService) DownloadFiles(dbPaths []string, url string) ([]io.ReadCloser, []error) {
	fmt.Println("URL:>", url)
	type PostObj struct {
		FilePath string `json:"filePath"`
	}

	var bodies []io.ReadCloser
	var errors []error

	for _, dbPath := range dbPaths {
		p := &PostObj{FilePath: dbPath}
		jsonStr, err := json.Marshal(p)
		if err != nil {
			fmt.Printf("download error: %v\n", err)
			errors = append(errors, err)
			bodies = append(bodies, nil)
			continue
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		if err != nil {
			fmt.Printf("download error: %v\n", err)
			errors = append(errors, err)
			bodies = append(bodies, nil)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("download error: %v\n", err)
			errors = append(errors, err)
			bodies = append(bodies, nil)
			continue
		}
		bodies = append(bodies, resp.Body)
		errors = append(errors, nil)
	}

	return bodies, errors
}
