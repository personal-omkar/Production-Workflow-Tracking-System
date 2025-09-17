package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-dao/dao"

	m "irpl.com/kanban-commons/model"
)

// Create STAGE Entry
func CreateStage(w http.ResponseWriter, r *http.Request) {
	var data m.Stage
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		err := dao.CreateStage(data)
		if err != nil {
			sysLog := m.SystemLog{
				Message:     "AddStage : Fail to add stage " + data.Name,
				MessageType: "ERROR",
				IsCritical:  false,
				CreatedBy:   data.CreatedBy,
			}
			utils.CreateSystemLogInternal(sysLog)
			slog.Error("Record creation failed - " + err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to create record")
			return
		}
		sysLog := m.SystemLog{
			Message:     "AddStage : Successfully added new stage " + data.Name,
			MessageType: "SUCCESS",
			IsCritical:  false,
			CreatedBy:   data.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		utils.SetResponse(w, http.StatusOK, "Successfully created record")
	} else {
		sysLog := m.SystemLog{
			Message:     "AddStage : Fail to add stage",
			MessageType: "ERROR",
			IsCritical:  false,
			CreatedBy:   data.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		slog.Error("Record creation failed - " + err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to create record")
	}
}

// Update Stage Entry
func UpdateExistingStage(w http.ResponseWriter, r *http.Request) {
	var data m.Stage
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		err := dao.UpdateExistingStage(&data)
		if err != nil {
			sysLog := m.SystemLog{
				Message:     "UpdateStage : Fail to update stage " + data.Name,
				MessageType: "ERROR",
				IsCritical:  false,
				CreatedBy:   data.CreatedBy,
			}
			utils.CreateSystemLogInternal(sysLog)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update record")
			slog.Error("Record updation failed - " + err.Error())
			return
		}
		sysLog := m.SystemLog{
			Message:     "UpdateStage : Successfully updated stage " + data.Name,
			MessageType: "SUCCESS",
			IsCritical:  false,
			CreatedBy:   data.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		utils.SetResponse(w, http.StatusOK, "Successfully updated record")
	} else {
		sysLog := m.SystemLog{
			Message:     "UpdateStage : Fail to update stage ",
			MessageType: "ERROR",
			IsCritical:  false,
			CreatedBy:   data.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to update record")
		slog.Error("Record updation failed - " + err.Error())
	}
}

// Delete UserRole Entry
func DeleteStageByID(w http.ResponseWriter, r *http.Request) {
	var data m.Stage
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		err := dao.DeleteStageByID(uint(data.ID))
		if err != nil {
			sysLog := m.SystemLog{
				Message:     "DeleteStage : Fail to delete stage ",
				MessageType: "ERROR",
				IsCritical:  false,
				CreatedBy:   data.CreatedBy,
			}
			utils.CreateSystemLogInternal(sysLog)
			slog.Error("Record deletion failed", "id", data.ID, "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to delete stage")
			return
		}
		sysLog := m.SystemLog{
			Message:     "DeleteStage : Successfully deleted stage",
			MessageType: "SUCCESS",
			IsCritical:  false,
			CreatedBy:   data.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		utils.SetResponse(w, http.StatusOK, "Successfully deleted stage")
	} else {
		sysLog := m.SystemLog{
			Message:     "DeleteStage : Fail to delete stage ",
			MessageType: "ERROR",
			IsCritical:  false,
			CreatedBy:   data.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to delete stage")
		slog.Error("Record deletion failed - " + err.Error())
	}
}

// Get UserRole Entry by parameters
func GetAllStage(w http.ResponseWriter, r *http.Request) {
	stage, err := dao.GetAllStage()
	if err != nil {
		slog.Error("Failed to get all stages", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to get all stages")
		return
	}
	allStages, err := json.Marshal(stage)
	if err != nil {
		slog.Error("Failed to get all stages", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to get all stages")
		return
	}
	utils.SetResponse(w, http.StatusOK, string(allStages))
}

// Get User by Parameter
func GetStagesByParam(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	stage, err := dao.GetStagesByParam(key, value)
	if err != nil {
		slog.Error("Stage not found", "key", key, "value", value, "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to find record")
		return
	}
	stages, _ := json.Marshal(stage)
	utils.SetResponse(w, http.StatusOK, string(stages))
}

// func GetStagesByHeader(w http.ResponseWriter, r *http.Request) {
// 	var data m.Stage
// 	decoder := json.NewDecoder(r.Body)

// 	if err := decoder.Decode(&data); err == nil {
// 		stages, err := dao.GetStagesByHeader(data.Headers)
// 		if err != nil {
// 			slog.Error("Failed to get data", "error", err.Error())
// 			utils.SetResponse(w, http.StatusInternalServerError, "Failed to get stages")
// 			return
// 		}
// 		marshalStages, _ := json.Marshal(stages)
// 		utils.SetResponse(w, http.StatusOK, string(marshalStages))
// 	} else {
// 		slog.Error("Failed to get stages - " + err.Error())
// 		utils.SetResponse(w, http.StatusInternalServerError, "Failed to get stages")
// 	}
// }
