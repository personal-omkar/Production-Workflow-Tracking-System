package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	u "irpl.com/kanban-commons/utils"
)

type StepCard struct {
	Icon              string `json:"Icon" gorm:"icon"`
	Name              string `json:"Name" gorm:"name"`
	Description       string `json:"Description" gorm:"description"`
	Style             string
	Text              string
	Order             int       `json:"Order" gorm:"order"`
	Status            string    `json:"Status" gorm:"status"`
	CompoundName      string    `json:"CompoundName" gorm:"compound_name"`
	CellNo            string    `json:"CellNo" gorm:"cell_no"`
	KbRootId          int       `json:"KbRootId" gorm:"kb_root_id"`
	ProdProcessLineId int       `json:"ProdProcessLineId" gorm:"prod_process_line_id"`
	MFGDateTime       time.Time `json:"MfgDateTime" gorm:"mfg_date_time"`
	IsGroup           bool      `json:"is_group" gorm:"column:isgroup"`
	GroupName         string    `json:"GroupName" gorm:"group_name"`
	CreatedOn         string    `json:"created_on"`
}

func (s *StepCard) Build(cards []*StepCard) string {

	s.Style = "user-select: none;"

	cardOrderMap := make(map[int]string)
	var i int = 1
	for _, card := range cards {
		if _, exists := cardOrderMap[card.Order]; !exists {
			cardOrderMap[card.Order] = fmt.Sprintf("%d", i)
			i++
		}
	}

	// Initialize necessary maps to track active groups, their minimum orders, and status
	activeGroupData := make(map[string]*StepCard)
	groupMinOrder := make(map[string]int)
	groupActiveStatus := make(map[string]bool)

	// Pass 1 - Identify active groups and their minimum orders
	for _, card := range cards {
		if card.IsGroup {
			if card.CellNo != "" {
				// Ensure the first card with a CellNo is stored as active
				if existingCard, exists := activeGroupData[card.GroupName]; !exists || card.Order < existingCard.Order {
					activeGroupData[card.GroupName] = card
				}
			}
			if minOrder, exists := groupMinOrder[card.GroupName]; !exists || card.Order < minOrder {
				groupMinOrder[card.GroupName] = card.Order
			}
			groupActiveStatus[card.GroupName] = card.CellNo != "" // Track active status for the group
		}
	}

	// If the current card is a group, loop through all cards in the same group to find a valid card with CellNo
	if s.IsGroup {
		for _, card := range cards {
			if card.GroupName == s.GroupName && card.CellNo != "" {
				// If a card with a non-empty CellNo is found in the same group, update s card with this card's data
				s.CellNo = card.CellNo
				s.CompoundName = card.CompoundName
				s.MFGDateTime = card.MFGDateTime
				s.KbRootId = card.KbRootId
				s.Status = card.Status
				s.Style = card.Style
				s.CreatedOn = card.CreatedOn
				break
			}
		}
	}

	// Pass 2 - Update all cards in the same group with the active card's data and minimum order
	for i := range cards {
		if cards[i].IsGroup {
			// Ensure the first card gets updated with the active group's data
			if activeCard, exists := activeGroupData[cards[i].GroupName]; exists {
				// Update the card with the active card's data
				cards[i].CellNo = activeCard.CellNo
				cards[i].CompoundName = activeCard.CompoundName
				cards[i].MFGDateTime = activeCard.MFGDateTime
				cards[i].KbRootId = activeCard.KbRootId
				cards[i].Status = activeCard.Status
				cards[i].Style = activeCard.Style
				cards[i].IsGroup = activeCard.IsGroup
				cards[i].GroupName = activeCard.GroupName
				cards[i].CreatedOn = activeCard.CreatedOn
				cards[i].Style = activeCard.Style
			}
			// Ensure the minimum order is propagated to each card
			if minOrder, exists := groupMinOrder[cards[i].GroupName]; exists {
				cards[i].Order = minOrder
			}
		}
	}

	// Check if the current card is active based on its group status (including the first card if the group is active)
	isActive := s.CellNo != "" || (s.IsGroup && groupActiveStatus[s.GroupName])
	if s.IsGroup {
		if groupStatus, exists := groupActiveStatus[s.GroupName]; exists {
			isActive = groupStatus
		}
	}
	dataStatus := "0"
	if isActive {
		dataStatus = "1"
	}

	// Prepare the data for HTML rendering
	var TableMap map[string]interface{}
	datadata, _ := json.Marshal(s)
	json.Unmarshal(datadata, &TableMap)
	data := strings.ReplaceAll(string(datadata), "'", "&#39;")

	// Use the group's minimum order in the HTML
	minOrder := s.Order
	if s.IsGroup {
		if groupOrder, exists := groupMinOrder[s.GroupName]; exists {
			minOrder = groupOrder
		}
	}
	img := u.ImageFetcher{
		DirPath: "./static/" + s.Icon + "/",
	}
	// Step 9: Generate the HTML based on the card's active status
	if dataStatus == "1" {

		return u.JoinStr(`
		<div class="flowchart-step-card" id="tooltip" style="`, s.Style, `" data-toggle="tooltip" data-bs-html="true" data-placement="right" data-tooltip="Cell: `, s.CellNo, ` Compound: `, s.CompoundName, `     Started-On Date: `, s.MFGDateTime.Format("02.01.2006"), `     Started-On Time: `, s.MFGDateTime.Format("15:04:05"), `" data-status="1" data-order="`, strconv.Itoa(minOrder), `" data-data='`, data, `'>
		    <div class="order-badge">`, cardOrderMap[s.Order], `</div>
			<div class="icon" style="margin:1 !important;">
			<img src="`, img.GetImagePath(), `" alt="`, s.Text, `" style="width:80px !important; height:72px !important; display:block !important; ">
			</div>
			<div class="text" style="margin-top: 2px !important;font-size: 14px;text-align: center;">`, s.Name, `</div>
			<div class="bottom-text w-75">`, truncateGroupName(s.GroupName), `</div> 
		</div>
		`)
	}

	// If not active, return a different HTML with status "0"
	return u.JoinStr(`
	<div class="flowchart-step-card" id="tooltip" style="`, s.Style, `" data-toggle="tooltip" data-placement="right" data-tooltip="No Processing data" data-status="0" data-order="`, strconv.Itoa(minOrder), `" data-data='`, data, `'>
		<div class="order-badge">`, cardOrderMap[s.Order], `</div>
		<div class="icon" style="margin:1 !important;">
			<img src="`, img.GetImagePath(), `" alt="`, s.Text, `" style="width:80px !important; height:72px !important; display:block !important; ">
			</div>
			<div class="text" style="margin-top: 2px !important;font-size: 14px;text-align: center;">`, s.Name, `</div>
			<div class="bottom-text w-75">`, truncateGroupName(s.GroupName), `</div> 
		</div>
	`)
}

type InfoCard struct {
	Title  string
	Fields map[string]string
}

func (c *InfoCard) Build() string {
	var fieldsHTML strings.Builder
	for key, value := range c.Fields {
		fieldsHTML.WriteString(u.JoinStr(`
			<div class="info-card-item"><span>`, key, `</span><span>`, value, `</span></div>
		`))
	}

	return u.JoinStr(`
		<div class="info-card">
			<div class="info-card-header">`, c.Title, `</div>
			`, fieldsHTML.String(), `
		</div>
	`)
}

// Function to truncate GroupName if its length exceeds 6
func truncateGroupName(groupName string) string {
	if len(groupName) > 6 {
		return groupName[:6] + "..."
	}
	return groupName
}
