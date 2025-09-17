package model

import "time"

type KbTransaction struct {
	Id                int       `json:"ID" gorm:"id"`
	ProdProcessId     int       `json:"ProdProcessId" gorm:"prod_process_id"`
	Status            string    `json:"Status" gorm:"status"` //-->> it is an order of production processes for that line
	JobId             int       `json:"JobId" gorm:"job_id"`
	KbRootId          int       `json:"KbRootId" gorm:"kb_root_id"`
	ProdProcessLineID int       `json:"ProdProcessLineID" gorm:"prod_process_line_iD"`
	StartedOn         time.Time `json:"StartedOn" gorm:"started_on"`
	CompletedOn       time.Time `json:"CompletedOn" gorm:"completed_on"`
	CreatedBy         string    `json:"CreatedBy" gorm:"created_by"`
	CreatedOn         time.Time `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy        string    `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn        time.Time `json:"ModifiedOn" gorm:"modified_on"`
	Operator          string    `json:"Operator" gorm:"operator"`
}
type SimpleCount struct {
	InProgress int `json:"in_progress"`
	Scheduled  int `json:"scheduled"`
}
type DualCount struct {
	ScheduledCount  int `json:"scheduled_count"`
	InProgressCount int `json:"in_progress_count"`
}
