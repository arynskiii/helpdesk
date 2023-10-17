package models

type Attachments struct {
	Id              int       `json:"id"`
	TicketId        int       `json:"id"`
	Attachment_path *[]string `json:"path"`
}
