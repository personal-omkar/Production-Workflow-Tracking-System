package orders

import (
	"strconv"
	"strings"

	m "irpl.com/kanban-commons/model"
	u "irpl.com/kanban-commons/utils"
)

type OrderDetails struct {
	KbDataId          int             `json:"ID"`
	CompoundName      string          `json:"CompoundName"`
	CellNo            string          `json:"CellNo"`
	DemandDateTime    string          `json:"DemandDateTime"`
	NoOFLots          int             `json:"NoOFLots"`
	Status            string          `json:"Status"`
	OrderId           string          `josn:"orderid"`
	ProductionProcess []m.ProdProcess `josn:"productionprocess"`
}

func (o *OrderDetails) Build() string {
	var html strings.Builder

	html.WriteString(o.buildCSS())
	html.WriteString(`
			<!--html-->
				<div>
					<div class="Order-Heading d-flex">
						<i class="fas fa-arrow-left order-back-button" onclick="window.location.href='/vendor-orders'"></i>
						<span class="Order-title" id="chartTitle">Back</span>
					</div>
					<div class="d-flex">
						<div class="col-4">
							<h3>KANBAN Status</h3>
							<div class="p-4 block mb-4">
								` + BuildTimeLineItem(o.ProductionProcess, o.Status) + `
							</div>
						</div>
						<div class="col-8 px-4 ">
							<h3> Cell Number: ` + o.CellNo + `</h3>
							<div class="py-4 block mb-4">
								<h4 class="px-3">Order Details</h4>
								` + o.buildOrderDetails() + `
							<div>
						</div>
					</div>
				</div>
			<!--!html-->
		`)
	html.WriteString(o.buildJS())
	return html.String()
}
func BuildTimeLineItem(p []m.ProdProcess, currentStatus string) string {
	var html strings.Builder
	// Status mapping (only base statuses)
	statusMap := map[string]int{
		"creating":            1,
		"pending":             2,
		"approved":            2,
		"reject":              2,
		"InProductionLine":    3,
		"InProductionProcess": 4,
		"quality":             5,
		"dispatched":          6,
		"dispatch":            7,
	}

	var (
		statusKey     string
		longestPrefix int
	)

	for key := range statusMap {
		if strings.HasPrefix(currentStatus, key) && len(key) > longestPrefix {
			statusKey = key
			longestPrefix = len(key)
		}
	}

	// If no match found, use default
	if statusKey == "" {
		statusKey = currentStatus
	}

	// Handle special case for "reject"
	if statusKey == "reject" {
		html.WriteString(`
		<div class="tl-item">
			<div class="tl-dot dot-complete"></div>
			<div class="tl-content">
				<div class="">Created</div>
				<div class="tl-date text-muted mt-1 complete-decs">You have created an order</div>
			</div>
		</div>
		`)
		html.WriteString(`
		<div class="tl-item">
			<div class="tl-dot dot-complete"></div>
			<div class="tl-content">
				<div class="">` + u.StatusMap[statusKey] + `</div>
				<div class="tl-date text-muted mt-1 complete-decs">Your order is rejected.</div>
			</div>
		</div>
		`)
		return html.String()
	}

	// Handle other statuses
	for _, value := range p {
		active, css := "", ""
		statusInt, _ := strconv.Atoi(value.Status)
		if statusInt > 0 && statusInt <= statusMap[statusKey] {
			if value.Name == "Pending" && currentStatus == "pending" {
				active = "dot-complete-pending"
			} else {
				active = "dot-complete"
			}
			css = "complete-decs"
		}

		html.WriteString(`
		<!--html-->
		<div class="tl-item">
			<div class="tl-dot ` + active + `"></div>
			<div class="tl-content">
				<div class="">` + value.Name + `</div>
				<div class="tl-date text-muted mt-1 ` + css + `">` + value.Description + `</div>
			</div>
		</div>
		<!--!html-->
		`)
	}

	return html.String()
}

// buildCSS provides CSS for the page
func (o *OrderDetails) buildCSS() string {
	return u.JoinStr(`
		<style>
		/* css */

		.p-4 {
			padding: 1.5rem !important
		}

		.block {
			background: #fff;
			border-width: 0;
			border-radius: .30rem;
			box-shadow: 0 1px 3px rgba(0, 0, 0, .05);
			margin-bottom: 1.5rem;
			max-height : 75vh;
			overflow-y : auto;
		}

		.mb-4,
		.my-4 {
			margin-bottom: 1.5rem !important
		}

		.tl-item {
			border-radius: 3px;
			position: relative;
			display: -ms-flexbox;
			display: flex
		}

		.tl-item>* {
			padding: 10px
		}

		.tl-item:last-child .tl-dot:after {
			display: none
		}


		.tl-item:last-child .tl-dot:after {
			display: none
		}

		.tl-item.active .tl-dot:before {
			border-color: #448bff;
			box-shadow: 0 0 0 4px rgba(68, 139, 255, .2)
		}

		.tl-dot {
			position: relative;
			border-color: rgba(160, 175, 185, .15)
		}

		.tl-dot:after,
		.tl-dot:before {
			content: '';
			position: absolute;
			border-color: inherit;
			border-width: 2px;
			border-style: solid;
			border-radius: 50%;
			width: 10px;
			height: 10px;
			top: 15px;
			left: 50%;
			transform: translateX(-50%)
		}

		.tl-dot:after {
			width: 0;
			height: auto;
			top: 25px;
			bottom: -15px;
			border-right-width: 0;
			border-top-width: 0;
			border-bottom-width: 0;
			border-radius: 0
		}


		.tl-dot {
			position: relative;
			border-color: rgba(160, 175, 185, .15)
		}
		.dot-complete:before,
		.dot-complete:after {
			background-color : #871A83;
			content: '';
			position: absolute;
			border-color: inherit;
			border-width: 2px;
			border-style: solid;
			border-radius: 50%;
			width: 20px;
			height: 20px;
			top: 15px;
			left: 50%;
			transform: translateX(-50%)
		}
		.dot-complete-pending:before{
			background-color : #871A83;
			content: '';
			position: absolute;
			border-color: inherit;
			border-width: 2px;
			border-style: solid;
			border-radius: 50%;
			width: 20px;
			height: 20px;
			top: 15px;
			left: 50%;
			transform: translateX(-50%)
		}

		.tl-dot:after {
			width: 0;
			height: auto;
			top: 25px;
			bottom: -15px;
			border-right-width: 0;
			border-top-width: 0;
			border-bottom-width: 0;
			border-radius: 0
		}

		.tl-content p:last-child {
			margin-bottom: 0
		}

		.tl-date {
			font-size: .85em;
			margin-top: 2px;
		}
		.complete-decs{
			color : #a542a0 !important; 
		}

		.order-details:nth-child(even) {
			background-color: #F1F1F1; 
		}
		.Order-Heading{ display: flex; align-items: center; margin-bottom: 20px; font-size: 18px; }
		.Order-title { font-weight: bold; font-size: 18px; }
		.order-back-button { color: #007bff; margin-right: 10px; cursor: pointer; transition: color 0.3s ease; z-index: 1000}		
		/* !css */
		</style>
	`)
}

// buildOrderDetails generates the order details dynamically
func (o *OrderDetails) buildOrderDetails() string {
	var Status string
	statusMap := map[string]string{
		"creating":            "1",
		"pending":             "2",
		"approved":            "2",
		"reject":              "2",
		"InProductionLine":    "3",
		"InProductionProcess": "4",
		"quality":             "5",
		"dispatch":            "7",
	}

	// Check if the status is present in the statusMap
	if _, exists := statusMap[o.Status]; !exists {
		Status = o.Status
	} else {
		Status = o.mapStatusToName(statusMap[o.Status])
	}

	details := []struct {
		Label string
		Value string
	}{
		{"Part Name", o.CompoundName},
		{"Cell Number", o.CellNo},
		{"Demand Date/Time", u.FormatStringDate(o.DemandDateTime, "date")},
		{"Numbers of Lot", strconv.Itoa(o.NoOFLots)},
		{"Status", Status}, // Use the dynamically set Status
	}

	var html strings.Builder
	for _, detail := range details {
		html.WriteString("<div class=\"py-2 order-details d-flex justify-content-between\">")
		html.WriteString("<div class=\"px-3 col-3\">" + detail.Label + ": </div>")
		html.WriteString("<div class=\"px-5 col-9\">" + detail.Value + "</div>")
		html.WriteString("</div>")
	}

	return html.String()
}

// mapStatusToName maps the status string to its corresponding name
func (o *OrderDetails) mapStatusToName(status string) string {
	for _, process := range o.ProductionProcess {
		if process.Status == status {
			return process.Name
		}
	}
	return status
}

// buildJS provides JavaScript for the page
func (o *OrderDetails) buildJS() string {
	return u.JoinStr(`
		<script>
		//js
			
		//!js
		</script>
	`)
}
