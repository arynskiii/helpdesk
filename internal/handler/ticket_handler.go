package handler

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/arynskiii/help_desk/models"
	"github.com/arynskiii/help_desk/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/vincent-petithory/dataurl"
)

func (h *Handler) CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.BindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	idCtx, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get userCtx"})
		return
	}
	userID, ok := idCtx.(int64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get userCtx"})
		return
	}
	category.UserId = userID
	id, err := h.services.Ticket.CreateCategory(c, category)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"category_id": id})
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	categoryId, err := strconv.Atoi(c.Request.URL.Query().Get("id"))
	if err != nil {
		c.JSON(400, err)
		return
	}
	value, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get userCtx"})
		return
	}
	userID, ok := value.(int64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user ID"})
		return
	}

	err = h.services.Ticket.DeleteCategory(c, int64(categoryId), userID)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": "deleted"})
}

func (h *Handler) CreateTicket(c *gin.Context) {
	var ticket models.Ticket
	if err := c.BindJSON(&ticket); err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	idCtx, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get userCtx"})
		return
	}
	userID, ok := idCtx.(int64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get userCtx"})
		return
	}
	ticket.UserId = userID
	id, err := h.services.Ticket.CreateTicket(c, ticket)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create ticket"})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
	c.Redirect(http.StatusSeeOther, "/allticket")
}
func (h *Handler) Changestate(c *gin.Context) {
	historyId, err := strconv.Atoi(c.Request.URL.Query().Get("historyId"))
	if err != nil {
		c.JSON(500, "Uncorrect ID")
		return
	}
	ct, ok := c.Get(roleCtx)
	if !ok {
		c.JSON(400, "fafafsa")
		return
	}
	role, ok := ct.(string)
	if role == "admin" {
		if err := h.services.Ticket.UpdateHistoryState(c, 5, int64(historyId)); err != nil {
			c.JSON(500, err)
			return
		}
	} else {
		c.JSON(500, "U tebya net prav")
	}
}
func (h *Handler) ShowallTicket(c *gin.Context) {
	h.services.Ticket.ShowallTicket(c)
}

func (h *Handler) DocumentSend(c *gin.Context) {
	var doc models.Attachments
	ticketId, err := strconv.Atoi(c.Request.URL.Query().Get("ticketid"))
	if err != nil {
		logger.Error(err)
		c.JSON(400, gin.H{"error": "Failed to parse Ticket ID"})
		return
	}

	err = c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	imgFiles, ok := c.Request.MultipartForm.File["idImage"]
	if !ok || len(imgFiles) == 0 {
		fmt.Println("Error: Image files not found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image files not found"})
		return
	}

	var fileNames []string
	for _, imgHandler := range imgFiles {
		fileNames = append(fileNames, imgHandler.Filename)
	}

	var imgReaders []multipart.File
	for _, imgHandler := range imgFiles {
		imgFile, err := imgHandler.Open()
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file"})
			return
		}
		defer imgFile.Close()
		imgReaders = append(imgReaders, imgFile)
	}

	fpaths, err := h.services.Ticket.PostFiles(imgReaders, "resume/id", "http://10.255.140.4:9090/sendfile", fileNames)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image files"})
		return
	}

	doc.TicketId = ticketId
	doc.Attachment_path = &fpaths
	_, err = h.services.Ticket.CreateDocAttachment(c, doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document attachment"})
		return
	}

	c.JSON(200, fpaths)
}
func (h *Handler) getTicketId(c *gin.Context) {
	var ticketID struct {
		ID int64 `uri:"id"`
	}
	if err := c.BindUri(&ticketID); err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	logger.Info(ticketID)
	ticket, err := h.services.Ticket.GetTicketByID(c, ticketID.ID)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket not found"})
		return
	}

	logger.Info(ticket)

	if ticket.Attachment_path != nil {
		readers, errors := h.services.Ticket.DownloadFiles(ticket.Attachment_path, "http://10.255.140.4:9090/clientdownload")

		// Handle errors from DownloadFiles
		for _, err := range errors {
			if err != nil {
				logger.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download attachment"})
				return
			}
		}
		defer func() {
			for _, reader := range readers {
				reader.Close()
			}
		}()

		var attachmentURLs []string
		for _, reader := range readers {
			data, err := ioutil.ReadAll(reader)
			if err != nil {
				logger.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read attachment data"})
				return
			}
			dataURL := dataurl.EncodeBytes(data)
			attachmentURLs = append(attachmentURLs, dataURL)
		}

		ticket.Attachment_path = attachmentURLs
	}

	c.JSON(http.StatusOK, ticket)
}
