package models

import "time"

type History struct {
	Id         int64     `json:"id"`
	TicketId   int64     `json:"ticketId"`
	StateId    int64     `json:"stateId"`
	CreateAt   time.Time `json:"createAt"`
	Deadline   time.Time `json:"dedline"`
	ReceiverId int64     `json:"receiverId"`
	SenderId   int64     `json:"senderId"`
}

type Document struct {
	Id              int64  `json:"id"`
	TicketId        int64  `json:"ticketId"`
	Attachment_path string `json:"path"`
}

type Ticket struct {
	Id              int        `json:"id"`
	CategoryId      int        `json:"categoryId"`
	UserId          int64      `json:"userId"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Deadline        *time.Time `json:"deadline"`
	CreateAt        time.Time  `json:"createAt"`
	Attachment_path []string   `json:"path"`
}

//
//type Files struct {
//	ID               int64   `json:"id"`
//	FilePath         string  `json:"filepath"`
//	FileName         string  `json:"filename"`
//	Uuid             string  `json:"uuid"`
//	Thumb            *string `json:"thumb"`
//	OriginalFilePath *string `json:"original_file_path"`
//	// Detail           *ImageDetail `json:"detail"`
//	IsDeleted  bool       `json:"is_deleted"`
//	CreateDate *time.Time `json:"create_date"`
//	ModiDate   *time.Time `json:"modi_date"`
//}
//
//type ImageDetail struct {
//	FileName  string  `json:"filename"`
//	Author    *string `json:"author"`
//	ImgDescKz *string `json:"img_desc_kz"`
//	ImgDescRu *string `json:"img_desc_ru"`
//	ImgDescEn *string `json:"img_desc_en"`
//}
//type Test struct {
//	Image *string
//}

// HISTORY Table
// id
// ticket id
// state id
// create at
// dedline
// receiver id
// sender_id
