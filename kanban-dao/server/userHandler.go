package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-dao/dao"

	m "irpl.com/kanban-commons/model"
)

// Create User
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var data m.User
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		user, _ := dao.GetUserWithEmail(data.Email)
		if user.ID != 0 {
			// create log
			sysLog := m.SystemLog{
				Message:     "AddUser : Fail to add new user " + data.Username,
				MessageType: "ERROR",
				IsCritical:  false,
				CreatedBy:   data.CreatedBy,
			}
			utils.CreateSystemLogInternal(sysLog)
			utils.SetResponse(w, http.StatusForbidden, "User with this email already exists")
			return
		}
		roleName, _ := dao.GetUserRoleByRoleID(int(data.RoleID))
		if roleName.RoleName == "Customer" || roleName.RoleName == "Operator" {
			vendor, _ := dao.GetVendorByParam("vendor_code", data.VendorsCode)
			if len(vendor) <= 0 {
				sysLog := m.SystemLog{
					Message:     "AddUser : Fail to add user, Invalid vendor code for " + data.Username,
					MessageType: "ERROR",
					IsCritical:  false,
					CreatedBy:   data.CreatedBy,
				}
				utils.CreateSystemLogInternal(sysLog)
				utils.SetResponse(w, http.StatusNotAcceptable, "Invalid vendor code")
				return
			}
		}

		// Create the user
		err = dao.CreateNewOrUpdateExistingUser(&data)
		if err != nil {
			sysLog := m.SystemLog{
				Message:     "AddUser : Fail to add new user " + data.Username,
				MessageType: "ERROR",
				IsCritical:  false,
				CreatedBy:   data.CreatedBy,
			}
			utils.CreateSystemLogInternal(sysLog)
			slog.Error("Record creation failed", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to create record")
			return
		}

		// Create user-to-role relation
		var usertorole m.UserToRole
		usertorole.UserId = int(data.ID)
		usertorole.UserRoleID = roleName.ID
		usertorole.CreatedOn = time.Now()
		usertorole.CreatedBy = data.CreatedBy
		err = dao.CreateNewOrUpdateExistingUserToRole(&usertorole)
		if err != nil {
			sysLog := m.SystemLog{
				Message:     "AddUser : Fail to add user-role for " + data.Username,
				MessageType: "ERROR",
				IsCritical:  false,
				CreatedBy:   data.CreatedBy,
			}
			utils.CreateSystemLogInternal(sysLog)
			slog.Error("User-role relation creation failed", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to create user-role relation")
			return
		}

		// Create user-to-vendor relation for customers and operators
		if roleName.RoleName == "Customer" || roleName.RoleName == "Operator" {
			vendor, _ := dao.GetVendorByParam("vendor_code", data.VendorsCode)

			user, _ := dao.GetUserWithEmail(data.Email) // Fetch updated user with ID
			var userTovendor m.UserToVendor
			userTovendor.VendorId = vendor[0].ID
			userTovendor.UserId = int(user.ID)
			userTovendor.CreatedOn = time.Now()

			err = dao.CreateNewOrUpdateVendorToUser(&userTovendor)
			if err != nil {
				sysLog := m.SystemLog{
					Message:     "AddUser : Fail to create vendor link for " + data.Username,
					MessageType: "ERROR",
					IsCritical:  false,
					CreatedBy:   data.CreatedBy,
				}
				utils.CreateSystemLogInternal(sysLog)
				slog.Error("User-vendor relation creation failed", "error", err.Error())
				utils.SetResponse(w, http.StatusInternalServerError, "Failed to create vendor relation")
				return
			}
		}

		// Success log
		sysLog := m.SystemLog{
			Message:     "AddUser : Successfully added new user " + data.Username,
			MessageType: "SUCCESS",
			IsCritical:  false,
			CreatedBy:   data.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		utils.SetResponse(w, http.StatusOK, "successfully created record")
	} else {
		sysLog := m.SystemLog{
			Message:     "AddUser : Fail to added new user",
			MessageType: "ERROR",
			IsCritical:  false,
			CreatedBy:   data.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		slog.Error("Record creation failed", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "failed to create record")
	}
}

// Update User
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var data m.User
	var userToVendor m.UserToVendor
	var userTOrole m.UserToRole
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		err := dao.CreateNewOrUpdateExistingUser(&data)
		if err != nil {
			// create log
			sysLog := m.SystemLog{
				Message:     "UpdateUser : Fail to updated user" + data.Username,
				MessageType: "ERROR",
				IsCritical:  false,
				CreatedBy:   data.CreatedBy,
			}
			utils.CreateSystemLogInternal(sysLog)
			slog.Error("Record update failed", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update record")
			return
		}
		role, _ := dao.GetUserToRolesByParam("userid", strconv.Itoa(int(data.ID)))
		if len(role) <= 0 {
			userTOrole.UserId = int(data.ID)
			userTOrole.UserRoleID = int(data.RoleID)
			userTOrole.CreatedOn = time.Now()
			userTOrole.CreatedBy = data.ModifiedBy
			err := dao.CreateNewOrUpdateExistingUserToRole(&userTOrole)
			if err != nil {
				// create log
				sysLog := m.SystemLog{
					Message:     "UpdateUser : Fail to updated user" + data.Username,
					MessageType: "ERROR",
					IsCritical:  false,
					CreatedBy:   data.CreatedBy,
				}
				utils.CreateSystemLogInternal(sysLog)
				slog.Error("Record update failed", "error", err.Error())
				utils.SetResponse(w, http.StatusInternalServerError, "Failed to update user role")
				return
			}
		} else if role[0].UserRoleID != int(data.RoleID) {
			role[0].UserRoleID = int(data.RoleID)
			role[0].ModifiedOn = time.Now()
			role[0].ModifiedBy = data.ModifiedBy
			err := dao.CreateNewOrUpdateExistingUserToRole(&role[0])
			if err != nil {
				// create log
				sysLog := m.SystemLog{
					Message:     "UpdateUser : Fail to updated user" + data.Username,
					MessageType: "ERROR",
					IsCritical:  false,
					CreatedBy:   data.CreatedBy,
				}
				utils.CreateSystemLogInternal(sysLog)
				slog.Error("Record update failed", "error", err.Error())
				utils.SetResponse(w, http.StatusInternalServerError, "Failed to update user role")
				return
			}
		}
		roleNae, _ := dao.GetUserRoleByRoleID(int(data.RoleID))
		if roleNae.RoleName == "Customer" {
			vendorId, _ := dao.GetVendorByParam("vendor_code", data.VendorsCode)
			if len(vendorId) <= 0 {
				// create log
				sysLog := m.SystemLog{
					Message:     "UpdateUser : Fail to updated user" + data.Username,
					MessageType: "ERROR",
					IsCritical:  false,
					CreatedBy:   data.CreatedBy,
				}
				utils.CreateSystemLogInternal(sysLog)
				utils.SetResponse(w, http.StatusNotAcceptable, "Failed to update vendor code, Invalid vendor code")
				return
			}
			rec, _ := dao.GetUserToVendorByParam("user_id", strconv.Itoa(int(data.ID)))
			if len(rec) <= 0 {
				userToVendor.UserId = int(data.ID)
				userToVendor.VendorId = vendorId[0].ID
				userToVendor.CreatedOn = time.Now()
				err := dao.CreateNewOrUpdateVendorToUser(&userToVendor)
				if err != nil {
					// create log
					sysLog := m.SystemLog{
						Message:     "UpdateUser : Fail to updated user" + data.Username,
						MessageType: "ERROR",
						IsCritical:  false,
						CreatedBy:   data.CreatedBy,
					}
					utils.CreateSystemLogInternal(sysLog)
					slog.Error("Record update failed", "error", err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Failed to update vendor code")
					return
				}
			} else if rec[0].VendorId != vendorId[0].ID {
				rec[0].VendorId = vendorId[0].ID
				rec[0].ModifiedOn = time.Now()
				err := dao.CreateNewOrUpdateVendorToUser(&rec[0])
				if err != nil {
					// create log
					sysLog := m.SystemLog{
						Message:     "UpdateUser : Fail to updated user" + data.Username,
						MessageType: "ERROR",
						IsCritical:  false,
						CreatedBy:   data.CreatedBy,
					}
					utils.CreateSystemLogInternal(sysLog)
					slog.Error("Record update failed", "error", err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Failed to update vendor code")
					return
				}
			}
		}
		// create log
		sysLog := m.SystemLog{
			Message:     "UpdateUser : Successfully updated user" + data.Username,
			MessageType: "SUCCESS",
			IsCritical:  false,
			CreatedBy:   data.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		utils.SetResponse(w, http.StatusOK, "successfully updated record")
	} else {
		// create log
		sysLog := m.SystemLog{
			Message:     "UpdateUser : Fail to updated user",
			MessageType: "ERROR",
			IsCritical:  false,
			CreatedBy:   data.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		log.Print("Error in decoding data")
		slog.Error("Record update failed", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
	}
}

// Delete User
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var ids []int64
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&ids); err == nil {
		for _, id := range ids {
			err := dao.DeleteUser(id)
			if err != nil {
				utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to delete record")
				slog.Error("Record deletion failed", "id", id, "error", err.Error())
				return
			}
		}
		utils.SetResponse(w, http.StatusOK, "Success: successfully deleted record")
	} else {
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to delete record")
		slog.Error("Record deletion failed", "error", err.Error())
	}
}

// Get User by Parameter
func GetUserByParam(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	m := map[string]interface{}{}
	m[key] = value

	users, _, err := dao.GetUsersByCriteria(1, 1, m)
	if err != nil {
		slog.Error("Record not found", "key", key, "value", value, "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to find record")
		return
	}

	bData, _ := json.Marshal(users)
	utils.SetResponse(w, http.StatusOK, string(bData))
}

// GetUserDetails retrives user details
func GetUserDetails(w http.ResponseWriter, r *http.Request) {

	userdetails, err := dao.GetUserDetails()
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(userdetails); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetUserDetailsByEmail retrives user details
func GetUserDetailsByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	userdetails, err := dao.GetUserDetailsByEmail(email)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(userdetails); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Update UserDetails
func UpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	var data m.UserManagement
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err == nil {
		userrec, err := dao.GetUserWithID(strconv.Itoa(data.UserID))
		if err != nil {
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
		} else {
			userrec.ID = uint(data.UserID)
			userrec.RoleID = uint(data.RoleId)
			userrec.Username = data.UserName
			userrec.Password = data.Password
			userrec.ModifiedOn = time.Now()
			err := dao.CreateNewOrUpdateExistingUser(&userrec)
			if err != nil {
				utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
			} else {
				utils.SetResponse(w, http.StatusOK, "Success: successfully updated user record")
			}
		}
	} else {
		slog.Error("Record updation failed", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
	}
}

// Update User status
func UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	var data m.User
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		data.ModifiedOn = time.Now()
		err := dao.UpdateExistingUser(&data)
		if err != nil {
			sysLog := m.SystemLog{
				Message:     "UpdateUser : Fail to updated user",
				MessageType: "ERROR",
				IsCritical:  false,
				CreatedBy:   data.CreatedBy,
			}
			utils.CreateSystemLogInternal(sysLog)
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
		} else {
			// create log
			sysLog := m.SystemLog{
				Message:     "UpdateUser : Successfully updated user" + data.Username,
				MessageType: "SUCCESS",
				IsCritical:  false,
				CreatedBy:   data.CreatedBy,
			}
			utils.CreateSystemLogInternal(sysLog)
			utils.SetResponse(w, http.StatusOK, "successfully updated record")
		}

	} else {
		// create log
		sysLog := m.SystemLog{
			Message:     "UpdateUser : Fail to updated user",
			MessageType: "ERROR",
			IsCritical:  false,
			CreatedBy:   data.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		slog.Error("Record update failed", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
	}
}
func GetAllUserBySearchAndPagination(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"conditions"`
	}
	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	KbRoot, PaginationResp, err := dao.GetAllUserBySearchAndPagination(requestData.Pagination, requestData.Conditions)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.UserManagement
	}
	Response.Pagination = PaginationResp
	Response.Data = KbRoot

	if err := json.NewEncoder(w).Encode(Response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
