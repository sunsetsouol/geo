package service

import (
	"geo/server/backend/dao"
	"geo/server/backend/model"
)

func GetAllPrompts() ([]model.Prompt, error) {
	var prompts []model.Prompt
	if err := dao.DB.Find(&prompts).Error; err != nil {
		return nil, err
	}
	return prompts, nil
}

func GetPromptByID(id uint64) (*model.Prompt, error) {
	var prompt model.Prompt
	if err := dao.DB.First(&prompt, id).Error; err != nil {
		return nil, err
	}
	return &prompt, nil
}

func CreatePrompt(prompt *model.Prompt) error {
	return dao.DB.Create(prompt).Error
}

func UpdatePrompt(prompt *model.Prompt) error {
	return dao.DB.Save(prompt).Error
}

func DeletePrompt(id uint64) error {
	return dao.DB.Delete(&model.Prompt{}, id).Error
}
