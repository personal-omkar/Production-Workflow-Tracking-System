package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-dao/db"
)

// RawQuery
func RawQuery(w http.ResponseWriter, r *http.Request) {
	var data model.RawQuery
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {

		if len(data.Type) > 0 {
			switch data.Type {
			case "SystemLog":
				var result []*model.SystemLog
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error
				if err != nil {
					slog.Error("RawQuery: SystemLog Record not found - " + err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to find record")
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))
			case "OrderDetails":
				var result []*model.OrderDetails
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error
				if err != nil {
					slog.Error("RawQuery: OrderDetails Record not found - " + err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to find record")
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))
			case "UserRoles":
				var result []*model.UserRoles
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error
				if err != nil {
					slog.Error("RawQuery: UserRoles Record not found - " + err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to find record")
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))
			case "Vendors":
				var result []*model.Vendors
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error
				if err != nil {
					slog.Error("RawQuery: Vendors Record not found - " + err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to find record")
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))
			case "KbExtension":
				var result []*model.KbExtension
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error
				if err != nil {
					slog.Error("RawQuery: KbExtension Record not found - " + err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to find record")
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))

			case "KbData":
				var result []*model.KbData
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error
				if err != nil {
					slog.Error("RawQuery: KbData Record not found - " + err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to find record")
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))

			case "KbRoot":
				var result []*model.KbRoot
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error
				if err != nil {
					slog.Error("RawQuery: KbRoot Record not found - " + err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to find record")
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))

			case "Compounds":
				var result []*model.Compounds
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error
				if err != nil {
					slog.Error("RawQuery: Compounds Record not found - " + err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to find record")
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))

			case "Stage":
				var result []*model.Stage
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error
				if err != nil {
					slog.Error("RawQuery: Stage Record not found - " + err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to find record")
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))

			case "Recipe":
				var result []*model.Recipe
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error
				if err != nil {
					slog.Error("RawQuery: Recipe Record not found - " + err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to find record")
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))

			case "ProdLine":
				var result []*model.ProdLine
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error
				if err != nil {
					slog.Error("RawQuery: Recipe Record not found - " + err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to find record")
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))

			case "AllKanbanViewDetails":
				var result []*model.AllKanbanViewDetails
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error
				if err != nil {
					slog.Error("RawQuery: AllKanbanViewDetails Record not found - " + err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to find record")
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))
			case "Operator":
				var result []*model.Operator
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error
				if err != nil {
					slog.Error("RawQuery: Operator Record not found - " + err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to find record")
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))
			case "DashboardStats":
				var result model.DashboardStatsResponse
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error
				if err != nil {
					slog.Error("RawQuery: DashboardStats failed - " + err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: unable to fetch stats")
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))

			case "InProgressOrderCount":
				var result model.DualCount
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error
				if err != nil {
					slog.Error("RawQuery: InProgressOrderCount query failed - " + err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: query failed")
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))
			case "MonthlyVendorOrderStatus":
				var result []model.VendorOrderStatus
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error
				if err != nil {
					slog.Error("RawQuery: MonthlyVendorOrderStatus query failed - " + err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]string{"error": "Fail: query failed"})
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))
			case "RawMaterial":
				var result []model.RawMaterial
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error
				if err != nil {
					slog.Error("RawQuery: RawMaterial query failed - " + err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]string{"error": "Fail: query failed"})
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))
			case "RecipeToStage":
				var result []*model.RecipeToStage
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error
				if err != nil {
					slog.Error("RawQuery: RecipeToStage Record not found - " + err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to find record")
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))
			case "ChemicalTypes":
				var result []*model.ChemicalTypes
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error
				if err != nil {
					slog.Error("RawQuery: ChemicalTypes Record not found - " + err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to find record")
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))
			case "ProdLineToRecipe":
				var result []*model.ProdLineToRecipe
				err := db.DBInstance.Raw(data.Query).Scan(&result).Error

				if err != nil {
					slog.Error("RawQuery: ProdLineToRecipe Record not found - " + err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to find record")
					return
				}
				bdata, _ := json.Marshal(result)
				utils.SetResponse(w, http.StatusOK, string(bdata))
			}
		}

	} else {
		slog.Error("RawQuery: Invalid Request Type - " + err.Error())
		utils.SetResponse(w, http.StatusBadRequest, "Fail: Invalid Request Type")
	}
}
