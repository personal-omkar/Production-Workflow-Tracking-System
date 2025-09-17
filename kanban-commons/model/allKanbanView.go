package model

type AllKanbanViewTable struct {
	MachineId          string `json:"MachineId" gorm:"column:machine_id"`
	MachineName        string `json:"MachineName" gorm:"column:machine_name"`
	PartNameorKanbanNo string `json:"PartNameorKanbanNo" gorm:"column:part_details;type:jsonb"`
	TotalParts         int    `json:"TotalParts"`
}

type AllKanbanViewDetails struct {
	CompoundName  string `json:"compound_name"  gorm:"column:compound_name"`
	KanbanNo      string `json:"KanbanNo"  gorm:"column:kanban_no"`
	CustomerNotes string `json:"Notes"  gorm:"column:notes"`
}
