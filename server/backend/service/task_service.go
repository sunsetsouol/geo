package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"geo/server/backend/config"
	"geo/server/backend/dao"
	"geo/server/backend/model"

	"gorm.io/gorm"
)

type ExposureEvaluation struct {
	BrandScore    float64 `json:"brand_score"`
	ExposureCount int     `json:"exposure_count"`
	ExposureRank  int     `json:"exposure_rank"`
	Analysis      string  `json:"analysis"`
}

func EvaluateExposure(responseText string) (*ExposureEvaluation, error) {
	apiKey := config.AppConfig.LLM.DashScopeAPIKey
	if apiKey == "" {
		apiKey = os.Getenv("DASHSCOPE_API_KEY")
	}

	prompt := fmt.Sprintf(`请分析以下文本中的品牌曝光情况。返回JSON格式，包含以下字段：
brand_score: 品牌评分 (0-100)
exposure_count: 品牌提及次数 (整数)
exposure_rank: 品牌在提及的所有品牌中的排名 (1为最高)
analysis: 简短的分析报告

文本内容：
%s`, responseText)

	requestBody := map[string]interface{}{
		"model": "qwen-plus",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a brand analysis assistant. Always respond in valid JSON.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"response_format": map[string]string{"type": "json_object"},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var llmResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(bodyText, &llmResp); err != nil {
		return nil, err
	}

	if len(llmResp.Choices) == 0 {
		return nil, fmt.Errorf("empty response from LLM")
	}

	var eval ExposureEvaluation
	if err := json.Unmarshal([]byte(llmResp.Choices[0].Message.Content), &eval); err != nil {
		return nil, err
	}

	return &eval, nil
}

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
			// 自动进行曝光度评估
			eval, err := EvaluateExposure(req.ResponseText)
			if err != nil {
				log.Printf("Exposure evaluation failed for task %d: %v", taskID, err)
				// 评分失败不影响任务状态更新，但记录错误
			}

			brandScore := req.BrandScore
			exposureCount := 0
			exposureRank := 0
			analysisReport := req.AnalysisReport

			if eval != nil {
				brandScore = eval.BrandScore
				exposureCount = eval.ExposureCount
				exposureRank = eval.ExposureRank
				analysisReport = eval.Analysis
			}

			result := model.Result{
				TaskID:         taskID,
				ResponseText:   req.ResponseText,
				BrandScore:     brandScore,
				ExposureCount:  exposureCount,
				ExposureRank:   exposureRank,
				AnalysisReport: analysisReport,
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
