package utils

var DefaultsMap = make(map[string]string)

var StatusMap = make(map[string]string)

var KanbanPriorityColors = map[string]map[string]string{
	"": {
		"bg-color":   "#54df96",
		"text-color": "black",
	},
	"urgent": {
		"bg-color":   "#ffe98d",
		"text-color": "black",
	},
	"regular": {
		"bg-color":   "#54df96",
		"text-color": "black",
	},
	"mosturgent": {
		"bg-color":   "#fb9290",
		"text-color": "black",
	},
}
