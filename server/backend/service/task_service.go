package service

import (
	"log"
	"strconv"
	"time"

	"geo/server/backend/dao"
	"geo/server/backend/model"

	"gorm.io/gorm"
)

func GenerateDailyTasks() {
	log.Println("Starting daily task generation...")
	var prompts []model.Prompt
	if err := dao.DB.Find(&prompts).Error; err != nil {
		log.Printf("failed to fetch prompts: %v", err)
		return
	}

	for _, p := range prompts {
		task := model.Task{
			PromptID: p.ID,
			Status:   "pending",
		}
		if err := dao.DB.Create(&task).Error; err != nil {
			log.Printf("failed to create task for prompt %d: %v", p.ID, err)
		}
	}
	log.Printf("Successfully generated %d tasks", len(prompts))
}

func GetPendingTasks() ([]model.Task, error) {
	var tasks []model.Task
	err := dao.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("status = ?", "pending").Find(&tasks).Error; err != nil {
			return err
		}
		if len(tasks) > 0 {
			ids := make([]uint64, len(tasks))
			for i, t := range tasks {
				ids[i] = t.ID
			}
			if err := tx.Model(&model.Task{}).Where("id IN ?", ids).Update("status", "processing").Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	for i := range tasks {
		dao.DB.Model(&tasks[i]).Association("Prompt").Find(&tasks[i].Prompt)
	}

	return tasks, nil
}

type UpdateTaskResultReq struct {
	Status         string  `json:"status" binding:"required"`
	ResponseText   string  `json:"response_text"`
	BrandScore     float64 `json:"brand_score"`
	AnalysisReport string  `json:"analysis_report"`
	Citations      []struct {
		URL   string `json:"url"`
		Title string `json:"title"`
	} `json:"citations"`
}

func UpdateTaskResult(taskIDStr string, req UpdateTaskResultReq) error {
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		return err
	}

	return dao.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Task{}).Where("id = ?", taskID).Updates(map[string]interface{}{
			"status":   req.Status,
			"last_run": time.Now(),
		}).Error; err != nil {
			return err
		}

		if req.Status == "completed" {
			result := model.Result{
				TaskID:         taskID,
				ResponseText:   req.ResponseText,
				BrandScore:     req.BrandScore,
				AnalysisReport: req.AnalysisReport,
			}
			if err := tx.Create(&result).Error; err != nil {
				return err
			}

			for _, cit := range req.Citations {
				citation := model.Citation{
					TaskID: taskID,
					URL:    cit.URL,
					Title:  cit.Title,
				}
				if err := tx.Create(&citation).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
