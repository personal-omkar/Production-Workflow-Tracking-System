package flowchart

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"sort"
	"strconv"
	"strings"

	m "irpl.com/kanban-commons/model"
	u "irpl.com/kanban-commons/utils"
	"irpl.com/kanban-web/services"
)

type FlowchartPage struct {
	Title     string
	Steps     []services.StepCard
	InfoCards []m.Cell
}

func NewFlowchartPage(title string, steps []services.StepCard, card []m.Cell) *FlowchartPage {
	sort.Slice(steps, func(i, j int) bool {
		order1 := steps[i].Order
		order2 := steps[j].Order
		return order1 < order2
	})
	return &FlowchartPage{
		Title:     title,
		Steps:     steps,
		InfoCards: card,
	}
}

func (f *FlowchartPage) Build() string {
	var Html strings.Builder

	sort.Slice(f.InfoCards, func(i, j int) bool {
		running1, _ := strconv.Atoi(f.InfoCards[i].KBRunningNo)
		running2, _ := strconv.Atoi(f.InfoCards[j].KBRunningNo)
		return running1 < running2
	})

	sort.Slice(f.Steps, func(i, j int) bool {
		order1 := f.Steps[i].Order
		order2 := f.Steps[j].Order
		return order1 < order2
	})

	// Convert f.Steps to a slice of pointers
	stepPointers := make([]*services.StepCard, len(f.Steps))
	for i := range f.Steps {
		stepPointers[i] = &f.Steps[i]
	}

	// HTML structure for the flowchart and info cards
	Html.WriteString(`
		<div class="main-container" style="overflow:auto; height:92vh;">
			<div class="flowchart-container">
				<div id="flowchartWrapper" class="flowchart-wrapper">
					<div class="flowchart-header">
						<i class="fas fa-arrow-left flowchart-back-button" onclick="window.location.href='/production-line'"></i>
						<span class="flowchart-title" id="chartTitle">` + f.Title + `</span>
						<button class="btn btn-primary me-1 mb-1 mx-4" id="process_forward" type="button" style="background-color:#871a83; border:none; cursor: pointer !important; z-index:999;" disabled>Forward</button>
						<div class="h3 me-1 mb-1 mx-4" style="color:#871A83">` + u.DefaultsMap["flow_chart_heading"] + `</div>
					</div>
					<div id="flowchart" class="flowchart-content">
						<svg id="connections" class="flowchart-svg">
							<defs>
								<marker id="arrowhead" markerWidth="6" markerHeight="4" refX="6" refY="2" orient="auto">
									<polygon points="0 0, 6 2, 0 4" fill="#333"></polygon>
								</marker>
							</defs>
						</svg>
	`)

	// Build flowchart steps in alternating alignment (left, right, center)
	c := 0                  //number of containers
	numberofdoublecard := 0 //number of containers having double cards
	isSingleCardRight := true
	for i := 0; i < len(f.Steps); i++ {
		step := f.Steps[i]
		stepID := fmt.Sprintf("step-%d", i)
		// Create centered row with two cards every third step if there's a pair available
		if i%3 == 0 && i+1 < len(f.Steps) {
			numberofdoublecard++           // to get numbr of container having double cards
			if numberofdoublecard%2 == 0 { //apply flow right to left for all even number of continers having double cards
				Html.WriteString(`<div class="flowchart-row center"  style="direction: rtl;" >`)
				Html.WriteString(`<div id="` + stepID + `" style="direction: ltr;">` + step.Build(stepPointers) + `</div>`)
				Html.WriteString(`<div id="step-` + fmt.Sprintf("%d", i+1) + `" style="direction: ltr;">` + f.Steps[i+1].Build(stepPointers) + `</div>`)
				Html.WriteString(`</div>`)
				i++ // Skip the next step since it is part of this row
			} else {
				Html.WriteString(`<div class="flowchart-row center">`)
				Html.WriteString(`<div id="` + stepID + `">` + step.Build(stepPointers) + `</div>`)
				Html.WriteString(`<div id="step-` + fmt.Sprintf("%d", i+1) + `">` + f.Steps[i+1].Build(stepPointers) + `</div>`)
				Html.WriteString(`</div>`)
				i++ // Skip the next step since it is part of this row
			}

		} else {
			// Single card alignment
			alignment := "left"
			if isSingleCardRight {
				alignment = "right"
			}
			Html.WriteString(`<div class="flowchart-row ` + alignment + `">`)
			Html.WriteString(`<div id="` + stepID + `">` + step.Build(stepPointers) + `</div>`)
			Html.WriteString(`</div>`)
			isSingleCardRight = !isSingleCardRight
		}
		c++

	}
	Html.WriteString(`
					</div>
				</div>
			</div>
			
			<!-- Process Card Section (Right Side, 3 Col Width) -->
			<div class="process-container">
			<div class="section-header">In Process</div>
			<div class="scroll-div">
		`)

	sort.Slice(f.InfoCards, func(i, j int) bool {
		return f.InfoCards[i].ProductionProcessLineOrder > f.InfoCards[j].ProductionProcessLineOrder
	})

	for _, value := range f.InfoCards {
		ProdProcessLineID, _ := strconv.Atoi(value.ProdProcessID)
		if ProdProcessLineID == 2 {
			continue
		}
		if ProdProcessLineID > 1 {
			card := &services.Card{
				Header: `
					<div class="card-header" style="display: flex; justify-content: space-between; align-items: center;" Root="` + value.KRId + `">
						<span>Cell Name ` + value.CellNumber + `</span>
						<i class="fas fa-ellipsis-v"></i>
					</div>
				`,
				Body: `
					<div class="card-body">
						<div class="row">
							<div class="col-6 row-label">Compound Code</div>
							<div class="col-6 row-label text-end">Lot No.</div>
							<div class="col-6 row-data">` + value.CompoundName + `</div>
							<div class="col-6 row-data text-end">` + value.LotNo + `</div>
						</div>
						<div class="row">
							<div class="col-6 row-label">Demand Date/Time</div>
							<div class="col-6 row-label text-end">MFG. Date/Time</div>
							<div class="col-6 row-data">` + u.FormatStringDate(value.DemandDateTime, "date") + `</div>
							<div class="col-6 row-data text-end">` + u.FormatStringDate(value.MfgDateTime, "date-time") + `</div>
						</div>
					</div>`,
				Width: "line-up-card",
				Style: "box-shadow: 2px 2px 5px rgba(0, 0, 0, 0.1); font-size: 10px;",
			}
			Html.WriteString(card.Build())
		}
	}
	Html.WriteString(`
				</div>
				</div>
				<!-- Lineup Card Section (Right Side, 3 Col Width) -->
				<div class="lineup-container">
				<div class="section-header">Line Up</div>
				<div class="scroll-div">			
			`)
	for _, value := range f.InfoCards {
		ProdProcessLineID, _ := strconv.Atoi(value.ProdProcessID)
		if ProdProcessLineID == 1 {
			kanbanNo := ""
			if value.KanbanNo.Valid {
				kanbanNo = value.KanbanNo.String
			}
			runningno, _ := strconv.Atoi(value.KBRunningNo)
			if runningno >= 1 && runningno <= 4 {
				// Serialize value to JSON for data-data attribute
				valueJSON, err := json.Marshal(value)
				if err != nil {
					log.Println("Error serializing value:", err)
					continue
				}

				// Escape JSON string to prevent HTML parsing issues
				escapedValueJSON := html.EscapeString(string(valueJSON))

				card := &services.Card{
					Header: `
						<div class="card-header" style="display: flex; justify-content: space-between; align-items: center;" data-data="` + escapedValueJSON + `">
							<span>Cell Name ` + value.CellNumber + `</span>
							<i class="fas fa-ellipsis-v"></i>
						</div>
					`,
					Body: `
						<div class="card-body">
							<div class="row">
								<div class="col-6 row-label">Compound Code</div>
								<div class="col-6 row-label text-end">Kanban No.</div>
								<div class="col-6 row-data">` + value.CompoundName + `</div>
								<div class="col-6 row-data text-end">` + kanbanNo + `</div>
							</div>
							<div class="row">
								<div class="col-6 row-label">Demand Date/Time</div>
								<div class="col-6 row-label text-end">MFG. Date/Time</div>
								<div class="col-6 row-data">` + u.FormatStringDate(value.DemandDateTime, "date") + `</div>
								<div class="col-6 row-data text-end">` + u.FormatStringDate(value.MfgDateTime, "date-time") + `</div>
							</div>
						</div>`,
					Width: "line-up-card",
					Style: "box-shadow: 2px 2px 5px rgba(0, 0, 0, 0.1); font-size: 10px;",
				}
				Html.WriteString(card.Build())
			}

		}
	}

	Html.WriteString(`
			</div>
			</div>
		</div>
	`)

	Html.WriteString(`
	<style>
	/*css*/
		.main-container { display: flex; width: 100%; height:100vh;position: relative; ;background-image: url('data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20"><circle cx="2" cy="2" r="1.5" fill="%23999" /></svg>');
        background-size: 20px 20px;  }
		.flowchart-container { flex: 6; padding-right: 20px; padding-left: 10px; }
		.process-container { flex: 3; display: flex; flex-direction: column; padding-right: 20px; font-size: 18px; height:100%; }
		.lineup-container { flex: 3; display: flex; flex-direction: column; font-size: 18px; height:100%; padding-bottom:100px; }

		.card-row { display: flex; justify-content: space-between; gap: 20px; }
		.info-card { flex: 1; padding: 10px; border: 1px solid #ddd; border-radius: 8px; background-color: #f9f9f9; box-shadow: 2px 2px 5px rgba(0, 0, 0, 0.1); }
		.info-card-header { font-weight: bold; font-size: 16px; margin-bottom: 10px; }
		.info-card-item { display: flex; justify-content: space-between; font-size: 14px; padding: 4px 0; background-color: #e8f4fb; border-radius: 4px; }

		.flowchart-wrapper { display: flex; flex-direction: column; align-items: flex-start; max-width: 80%; }
		.flowchart-header{ display: flex; align-items: center; margin-bottom: 20px; font-size: 18px; }
		.section-header{ display: flex; align-items: center; margin-bottom: 20px; font-size: 18px; font-style : bold ; }
		.flowchart-back-button { color: #007bff; margin-right: 10px; cursor: pointer; transition: color 0.3s ease; z-index: 1000}
		.flowchart-back-button:hover { color: #0056b3; }
		.flowchart-title { font-weight: bold; font-size: 18px; }
		.flowchart-content { display: flex; flex-direction: column; align-items: flex-start; width: 100%; margin-top: 10px }
		.flowchart-tooltip { position: absolute; padding: 10px; background-color: #333; color: #fff; border-radius: 5px; font-size: 14px; opacity: 0; pointer-events: none; transition: opacity 0.2s ease-in-out; }
		.flowchart-step-card { position: relative; width: 120px; height: 130px; padding: 15px; border: 1px solid #ddd; border-radius: 8px; text-align: center; background-color: #fff; font-weight: bold; box-shadow: 2px 2px 5px rgba(0, 0, 0, 0.1); cursor: pointer; transition: box-shadow 0.3s ease, transform 0.3s ease; }
		.order-badge {
			position: absolute;
			top: 5px; 
			right: 5px; 
			color: #871A83;
			font-size: 12px; 
			padding: 2px 6px; 
			border-radius: 50%; 
			z-index: 1; 
			text-align: center;
		}
		.bottom-text {
			position: absolute;
			bottom: 10px; 
			left: 50%; 
			transform: translateX(-50%); 
			font-size: 14px; 
			color: #871A83;
		}
		.flowchart-step-card:hover { box-shadow: 0px 8px 15px rgba(0, 0, 0, 0.3); transform: translateY(-5px); }
		.flowchart-step-card.active { box-shadow: 0px 0px 20px rgba(135, 26, 131, 0.8);}
		.flowchart-step-card.select { box-shadow: 0px 0px 30px rgba(40, 160, 255, 0.8); transform: translateY(-5px); }
		.line-up-card.select { box-shadow: 0px 0px 30px rgba(40, 160, 255, 0.8); transform: translateY(-5px); }
		.flowchart-row { display: flex; gap: 20px; margin-bottom: 20px; width: 100%; }
		.right { justify-content: flex-end; width: 100%; }
		.left { justify-content: flex-start; width: 100%; }
		.center { justify-content: center; }
		.flowchart-svg { position: absolute; top: 0; left: 0; overflow: visible; }
		.flowchart-arrow-line { fill: none; stroke: #333; stroke-width: 0.8; marker-end: url(#arrowhead); }
		.icon img { width: 24px; height: 24px; vertical-align: middle; }
		.line-up-card {
			box-shadow: 2px 2px 5px rgba(0, 0, 0, 0.1);
			font-size: 14px;
			margin-bottom: 15px; /* Adds gap between cards */
			border-radius: 8px;
			background-color: #fff;
			width: 85%;
		}

		/* Reduce height of card content */
		.line-up-card .card-body {
			padding: 8px;
		}

		.line-up-card .row {
			margin-bottom: 8px; /* Adds gap between rows inside each card */
			margin-left: 0;
			margin-right: 0;
		}

		.card-body .row:nth-child(even) {
			background-color: #ffffff;
		}

		.card-body .row:nth-child(odd) {
			background-color: #e0f7ff;
		}

		.scroll-div{
			width: 85%;
			padding-right: 0;   /* No padding to prevent pushing of content */
    		margin: 0;          
			/* overflow-y: auto; */
			&::-webkit-scrollbar-thumb:hover {
				background: #a8a8a8 !important;
			}
			&::-webkit-scrollbar-track {
				-webkit-box-shadow: inset 0 0 6px rgba(83, 83, 83, 0.07);
				background-color: #f1f1f1;
			}
			&::-webkit-scrollbar {
				width: 7px;
				background-color: #f1f1f1;
			}
			&::-webkit-scrollbar-thumb {
				background-color: #c1c1c1;
			}
		}


		#tooltip {
            position: relative;
            display: inline-block;
            cursor: pointer;
			
        }

        #tooltip:hover::after {
            content: attr(data-tooltip);
            position: absolute;
            background-color: #871a83;
            color:#ffffff;
            padding: 8px;
            border-radius: 5px;
            white-space: pre-line; /* Enables multi-line content */
            top: 100%; /* Show below the element */
            left: 100%;
            transform: translateX(-50%);
            z-index: 1000;
            width: 250px;
            text-align: left;
        }
		/*!css*/
	</style>
	`)

	js := `
		<script>
		//js
			document.addEventListener("DOMContentLoaded", function () {
				let isUpdating = false; // Flag to track if an update is in progress
    			let updateInterval;

				function updateFlowchart() {
					if (isUpdating) return;
			        isUpdating = true;

					const urlParams = new URLSearchParams(window.location.search);
					const line = urlParams.get('line');
					const apiUrl = '/update-flowchart?line=' + line; 

					fetch(apiUrl)
						.then(response => response.json())
						.then(data => {
							const container = document.querySelector('.main-container');
							const parent = container.parentNode;
							parent.querySelectorAll('style, script').forEach((el) => el.remove());
							container.innerHTML = '';
							container.outerHTML = data.content; 
						addCardClickListeners();
						firstCardChecked = false;					
						checkFirstCard();					
						reinitializeConnections();
						})
						.catch(error => console.error('Error fetching data:', error))
						.finally(() => {
							isUpdating = false;
						});
				}

				function startUpdateInterval() {
    			    if (updateInterval) clearInterval(updateInterval);
        			updateInterval = setInterval(updateFlowchart, 10000);
    			}

    			function stopUpdateInterval() {
        			if (updateInterval) clearInterval(updateInterval);
    			}

				function highlightCardsByRoot(rootId) {
					let Processcards = document.querySelectorAll(".process-container .line-up-card");
					Processcards.forEach(c => c.classList.remove("select"));
					Processcards.forEach(c => {
						let cHeader = c.querySelector(".card .card-header");
						if (cHeader) {
							let cardRootId = cHeader.getAttribute("root");
							if (cardRootId == rootId) {
								c.classList.add("select");
							}
						}
					});
				}

				function addCardClickListeners() {
					const cards = document.querySelectorAll('.flowchart-step-card');
					const processForwardButton = document.getElementById('process_forward');
					let selectedCardIndex = null;

					// Get unique orders from all cards
					const orders = Array.from(cards).map(card => parseInt(card.getAttribute('data-order'), 10));
					const uniqueOrders = [...new Set(orders)].sort((a, b) => a - b);

					cards.forEach((card, index) => {
						if (card.getAttribute('data-status') === '1') {
							card.classList.add('active');
						}

						card.addEventListener('click', (event) => {
							// Only allow clicks on cards with data-status = '1'
							if (card.getAttribute('data-status') !== '1') return;

							// Deselect all cards and select the clicked one
							cards.forEach(c => c.classList.remove('select'));
							card.classList.add('select');
							let dataData = card.getAttribute("data-data");
							
							if (dataData) {
								try {
									let parsedData = JSON.parse(dataData);
									let kbRootId = parsedData.KbRootId; // Assuming KbRootId is inside the JSON

									if (kbRootId) {
										highlightCardsByRoot(kbRootId);
									}
								} catch (error) {
									console.error("Invalid JSON in data-data attribute:", error);
								}
							}
							selectedCardIndex = index;

							// Get the current order and the next order
							const currentOrder = parseInt(card.getAttribute('data-order'), 10);
							const currentOrderIndex = uniqueOrders.indexOf(currentOrder);
							const nextOrder = uniqueOrders[currentOrderIndex + 1];

							if (nextOrder) {
								// Find the next card by order
								const nextCard = Array.from(cards).find(
									c => parseInt(c.getAttribute('data-order'), 10) === nextOrder
								);

								if (nextCard && nextCard.getAttribute('data-status') === '0') {
									processForwardButton.disabled = false;
									processForwardButton.innerHTML = "Forward";
								} else {
									processForwardButton.disabled = true;
								}
							} else {
								// If there's no next order, handle "Packing"
								processForwardButton.innerHTML = "Packing";
								processForwardButton.disabled = false;
							}

							event.stopPropagation();
						});
					});

					processForwardButton.addEventListener('click', () => {
						if (selectedCardIndex !== null) {
							stopUpdateInterval();
							const currentCard = cards[selectedCardIndex];
							const currentOrder = parseInt(currentCard.getAttribute('data-order'), 10);
							const currentCardData = JSON.parse(currentCard.getAttribute('data-data'));
							const kbRootId = currentCardData.KbRootId;

							// Get the next order
							const currentOrderIndex = uniqueOrders.indexOf(currentOrder);
							const nextOrder = uniqueOrders[currentOrderIndex + 1];

							if (nextOrder) {
								const nextCard = Array.from(cards).find(
									c => parseInt(c.getAttribute('data-order'), 10) === nextOrder
								);
								const nextCardData = JSON.parse(nextCard.getAttribute('data-data'));
								const data = {
									KbRootId: kbRootId,
									ProdProcessLineId: nextCardData.ProdProcessLineId
								};

								// Send fetch request for the next card
								fetch("/create-new-KbTransaction", {
									method: 'POST',
									headers: { 'Content-Type': 'application/json' },
									body: JSON.stringify(data)
								}).then(response => {
									if (!response.ok) throw new Error("Server error");
									updateFlowchart();
									startUpdateInterval();
								});
							} else {
								const data = {
									KbRootId: kbRootId,
									Status: "Packing"
								};

								fetch("/create-new-KbTransaction", {
									method: 'POST',
									headers: { 'Content-Type': 'application/json' },
									body: JSON.stringify(data)
								}).then(response => {
									if (!response.ok) throw new Error("Server error");
									updateFlowchart();
									startUpdateInterval();
								});
							}
						}
					});

					// Reset selection and disable button on outside click
					document.addEventListener('click', () => {
						if (processForwardButton) {
							processForwardButton.disabled = true;
							processForwardButton.innerHTML = "Forward";
						}
						cards.forEach(c => c.classList.remove('select'));
						let Processcards = document.querySelectorAll(".process-container .line-up-card");
						Processcards.forEach(c => c.classList.remove('select'));

						selectedCardIndex = null;
					});
				}

				let firstCardChecked = false; // Add this flag
				function checkFirstCard() {
					if (firstCardChecked) return; // Prevent duplicate execution
					firstCardChecked = true;
					const firstCard = document.querySelector('#step-0 .flowchart-step-card');
					if (firstCard) {
						if (!firstCard.classList.contains('active')) {
							const firstCardData = JSON.parse(firstCard.getAttribute('data-data'));
							const lineUpCard = $(".lineup-container .line-up-card")[0];
							if (lineUpCard) {
								const cardHeader = $(lineUpCard).find(".card-header")[0];
								if (cardHeader) {
									const dataData = $(cardHeader).attr("data-data");
									if (dataData) {
										const parsedData = JSON.parse(dataData);
										const data = {
											KbRootId: parseInt(parsedData.krid, 10),
											ProdProcessLineId: firstCardData.ProdProcessLineId,
											Status: "Line-Up"
										};
										// Perform the fetch request
										fetch("/create-new-KbTransaction", {
											method: 'POST',
											headers: {
												'Content-Type': 'application/json',
											},
											body: JSON.stringify(data)
										})
										.then(response => {
											if (!response.ok) {
												throw new Error("Server error");
											}
											const urlParams = new URLSearchParams(window.location.search);
											const line = urlParams.get('line');
											const IDdata = {
												ID : parseInt(line,10)
											};
											fetch("/update-running-number",{
												method: 'POST',
												headers: {
												'Content-Type': 'application/json',
												},
												body: JSON.stringify(IDdata)
											})
											.then(response => {
												if (!response.ok) {
													throw new Error("Server error");
												}
												updateFlowchart(); // Ensure this function updates the flowchart
											})
										})
									}
								}
							}
						}
					}
				}


			function reinitializeConnections() {
				const connections = document.getElementById("connections");
				//connections.innerHTML = "";
				function addConnection(startElem, endElem) {
					const startRect = startElem.getBoundingClientRect();
					const endRect = endElem.getBoundingClientRect();
					const svgRect = connections.getBoundingClientRect();
					
					const startX = startRect.left + startRect.width / 2 - svgRect.left;
					const startY = startRect.bottom - svgRect.top;
					const endX = endRect.left + endRect.width / 2 - svgRect.left;
					const endY = endRect.top - svgRect.top;


					const path = document.createElementNS("http://www.w3.org/2000/svg", "path");
					path.setAttribute("class", "flowchart-arrow-line");
					path.setAttribute("marker-end", "url(#arrowhead)");

					if (startX === endX) {
						path.setAttribute("d", "M " + startX + " " + startY + " L " + endX + " " + endY);
					} else {
						const midY = startY + (endY - startY) / 2;
						path.setAttribute("d", "M " + startX + " " + startY + " L " + startX + " " + midY + " L " + endX + " " + midY + " L " + endX + " " + endY);
					}
					connections.appendChild(path);
				}
				
				// Define custom connection order
				const connectionOrder = [
					{ start: 0, end: 1 },
					{ start: 1, end: 2 },
					{ start: 2, end: 3 },
					{ start: 3, end: 4 },
					{ start: 4, end: 5 },
					{ start: 5, end: 6 },
					{ start: 6, end: 7 },
					{ start: 7, end: 8 },
					{ start: 8, end: 9 },
					{ start: 9, end: 10 },
					{ start: 10, end: 11 }
					// Continue adding more as needed
				];

				// Add connections based on the custom order
				connectionOrder.forEach(({ start, end }) => {
					const startElem = document.getElementById("step-" + start);
					const endElem = document.getElementById("step-" + end);
					if (startElem && endElem) {
						addConnection(startElem, endElem);
					}
				});
			}
			updateFlowchart();
			startUpdateInterval();
		});


		//!js
		</script>
	`
	Html.WriteString(js)

	return Html.String()
}
