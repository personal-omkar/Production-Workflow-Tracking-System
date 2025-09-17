package model

import (
	"time"

	"github.com/lib/pq"
)

type KbData struct {
	Id             int            `json:"ID" gorm:"id"`
	CompoundId     int            `json:"CompoundId" gorm:"compound_id"`
	MFGDateTime    time.Time      `json:"MFGDateTime" gorm:"mfg_date_time"`
	DemandDateTime time.Time      `json:"DemandDateTime" gorm:"demand_date_time"`
	SubmitDateTime time.Time      `json:"SubmitDateTime" gorm:"submit_date_time"`
	ExpDate        time.Time      `json:"ExpDate" gorm:"exp_date"`
	CellNo         string         `json:"CellNo" gorm:"cell_no"`
	NoOFLots       int            `json:"NoOFLots" gorm:"no_of_lots"`
	Location       string         `json:"Location" gorm:"location"`
	CreatedBy      string         `json:"CreatedBy" gorm:"created_by"`
	CreatedOn      time.Time      `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy     string         `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn     time.Time      `json:"ModifiedOn" gorm:"modified_on"`
	KbExtensionID  int            `json:"kbExtensionId" gorm:"kb_extension_id"`
	KanbanNo       pq.StringArray `json:"KanbanNo" gorm:"column:kanban_no;type:text[]"`
	Note           string         `json:"Note" gorm:"column:note"`
}

// DashboardStatsResponse represents the response structure for dashboard statistics.
type DashboardStatsResponse struct {
	Code                 int    `json:"code"`
	Message              string `json:"message"`
	DailySubmitted       int    `json:"daily_submitted"`
	MonthlySubmitted     int    `json:"monthly_submitted"`
	DailyDispatched      int    `json:"daily_dispatched"`
	MonthlyDispatched    int    `json:"monthly_dispatched"`
	DailyPercentage      int    `json:"daily_percentage"`
	MonthlyPercentage    int    `json:"monthly_percentage"`
	MonthlyRejPercentage int    `json:"monthly_rej_percentage"`
	MonthlyRejected      int    `json:"monthly_rejected"`
}
