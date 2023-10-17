package models

import "time"

type Category struct {
	Id          int64     `json:"id"`
	UserId      int64     `json:"userId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreateAt    time.Time `json:"createAt"`
}

//type CourseCategory struct {
//	ID            *int64          `json:"id"`
//	ParentID      *int64          `json:"parentID"`
//	NameKz        *string         `json:"namekz"`
//	NameRu        *string         `json:"nameru"`
//	NameEn        *string         `json:"nameen"`
//	CreatorID     int64           `json:"creatorID"`
//	CreateDate    time.Time       `json:"createDate"`
//	ModifyDate    time.Time       `json:"modifyDate"`
//	IsActive      int             `json:"isActive"`
//	ModifierID    *int64          `json:"modifierID"`
//	DescriptionKz *string         `json:"descriptionkz"`
//	DescriptionRu *string         `json:"descriptionru"`
//	DescriptionEn *string         `json:"descriptionen"`
//	OrganizerID   int64           `json:"organizerID"`
//	Organizer     CourseOrganizer `json:"organizer"`
//	//TreeTable үшін
//	Leaf bool    `json:"leaf"`
//	Key  *string `json:"key"`
//}
