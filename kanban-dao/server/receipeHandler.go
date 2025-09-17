package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-dao/dao"
)

// Create Recipe Entry
func CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var data []model.Recipe

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err == nil {
		var recipe model.Recipe
		recipe.CompoundCode = data[0].CompoundCode
		recipe.CompoundName = data[0].CompoundName
		recipe.CreatedOn = time.Now()
		recipe.CreatedBy = data[0].CreatedBy
		recipe.BaseQty = data[0].BaseQty
		err := dao.CreateNewOrUpdateExistingRecipe(&recipe)
		if err != nil {
			slog.Error("Recipe creation failed - " + err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Recipe already exist")
			return
		}

		var prodLineToRecipe model.ProdLineToRecipe
		prodLineToRecipe.ProdLineId = data[0].ProdLineId
		prodLineToRecipe.RecipeId = recipe.Id
		recipe.CreatedOn = time.Now()
		recipe.CreatedBy = data[0].CreatedBy
		err = dao.CreateNewOrUpdateExistingProdLineToRecipe(&prodLineToRecipe)
		if err != nil {
			slog.Error("Recipe creation failed - " + err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "prod line to recipe relation already exist")
			return
		}

		recipes, err := dao.GetRecipeByParam("Compound_Code", recipe.CompoundCode)
		if len(recipes) > 0 {
			for _, rts := range data {

				var recipeToStage model.RecipeToStage
				recipeToStage.RecipeId = recipes[0].Id
				recipeToStage.StageId = int(rts.StageId)
				recipeToStage.Data = rts.Data
				recipeToStage.CreatedBy = rts.CreatedBy
				err := dao.CreateRecipeToStage(&recipeToStage)
				if err != nil {
					slog.Error("Recipe to stage creation failed - " + err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Recipe to stage failed")
					return
				}
			}
			utils.SetResponse(w, http.StatusOK, "Success: successfully created record")
		} else {
			slog.Error("Recipe creation failed - " + err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Recipe not found")
			return
		}
	} else {
		slog.Error("Recipe creation failed - " + err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create recipe")
	}
}

// Update Recipe Entry
func UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	var data []model.Recipe

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err == nil {

		var recipe model.Recipe
		recipe.Id = data[0].Id
		recipe.CompoundCode = data[0].CompoundCode
		recipe.CompoundName = data[0].CompoundName
		recipe.ModifiedOn = time.Now()
		recipe.ModifiedBy = data[0].ModifiedBy
		recipe.BaseQty = data[0].BaseQty
		err := dao.CreateNewOrUpdateExistingRecipe(&recipe)
		if err != nil {
			slog.Error("Recipe Updation failed - " + err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Recipe already exist")
			return
		}

		ProdLineToRecipe, _ := dao.GetProductionLineToRecipe(data[0].ProdLineToRecipe)
		if ProdLineToRecipe.Id != 0 {
			ProdLineToRecipe.ProdLineId = data[0].ProdLineId
			ProdLineToRecipe.ModifiedOn = time.Now()
			ProdLineToRecipe.ModifiedBy, _ = strconv.Atoi(data[0].ModifiedBy)
			err := dao.CreateNewOrUpdateExistingProdLineToRecipe(&ProdLineToRecipe)
			if err != nil {
				slog.Error("Recipe Updation failed - " + err.Error())
				utils.SetResponse(w, http.StatusInternalServerError, "Recipe not exist")
				return
			}

		}

		recipes, err := dao.GetRecipeByParam("Compound_Code", recipe.CompoundCode)
		if len(recipes) > 0 {
			for _, rts := range data {

				reciptTostage, _ := dao.GetRecipeToStageByRecipeAndStageId(recipes[0].Id, int(rts.StageId))
				var recipeToStage model.RecipeToStage
				if len(reciptTostage) > 0 {
					recipeToStage.Id = reciptTostage[0].Id
					recipeToStage.RecipeId = reciptTostage[0].RecipeId
					recipeToStage.StageId = reciptTostage[0].StageId
					recipeToStage.Data = rts.Data
					recipeToStage.ModifiedBy = rts.ModifiedBy
					err := dao.UpdateRecipeToStage(&recipeToStage)
					if err != nil {
						slog.Error("Recipe to stage updation failed - " + err.Error())
						utils.SetResponse(w, http.StatusInternalServerError, "Recipe to stage failed")
						return
					}
				} else {
					recipeToStage.RecipeId = recipes[0].Id
					recipeToStage.StageId = int(rts.StageId)
					recipeToStage.Data = rts.Data
					recipeToStage.CreatedBy = rts.CreatedBy
					err := dao.CreateRecipeToStage(&recipeToStage)
					if err != nil {
						slog.Error("Recipe to stage creation failed - " + err.Error())
						utils.SetResponse(w, http.StatusInternalServerError, "Recipe to stage failed")
						return
					}
				}

			}
			utils.SetResponse(w, http.StatusOK, "Success: successfully updated record")
		} else {
			slog.Error("Recipe updation failed - " + err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Recipe not found")
			return
		}
	} else {
		slog.Error("Recipe creation failed - " + err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update recipe")
	}
}

// GetAllRecipe retrives all recipe
func GetAllRecipe(w http.ResponseWriter, r *http.Request) {

	recipes, err := dao.GetAllRecipe()
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(recipes); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetRecipeByDataKey retrives  recipe based on the datakey value
func GetRecipeByDataKey(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	recipes, err := dao.GetRecipeByDataKey(key)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(recipes); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetRecipeByDataKey retrives  recipe based on the datakey value
func GetRecipeByDataValue(w http.ResponseWriter, r *http.Request) {
	value := r.URL.Query().Get("value")
	recipes, err := dao.GetRecipeByDataValue(value)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(recipes); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetRecipeByDataKey retrives  recipe based on the datakey value
func GetRecipeByDataKeyAndValue(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	recipes, err := dao.GetRecipeByDataKeyAndValue(key, value)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(recipes); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Delete recipe
func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe model.Recipe
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&recipe); err == nil {

		err := dao.DeleteRecipeById(recipe.Id)
		if err != nil {
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to delete record")
			slog.Error("Record deletion failed", "id", recipe.Id, "error", err.Error())
			return
		}

		utils.SetResponse(w, http.StatusOK, "Success: successfully deleted record")
	} else {
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to delete record")
		slog.Error("Record deletion failed", "error", err.Error())
	}
}

// Get Recipe by Parameter
func GetRecipeByParam(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	stage, err := dao.GetRecipeByParam(key, value)
	if err != nil {
		slog.Error("Recipe not found", "key", key, "value", value, "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to find record")
		return
	}
	stages, _ := json.Marshal(stage)
	utils.SetResponse(w, http.StatusOK, string(stages))
}
