package storage

import (
	"calulator/models"
	"encoding/json"
	"os"
)

const feedbackFile = "data/feedback.json"

func SaveFeedback(feedback models.Feedback) error {
	feedbacks, err := GetFeedbacks()
	if err != nil {
		feedbacks = []models.Feedback{}
	}
	feedbacks = append(feedbacks, feedback)
	file, err := json.MarshalIndent(feedbacks, "", " ")
	if err != nil {
		return nil
	}
	return os.WriteFile(feedbackFile, file, 0644)
}

func GetFeedbacks() ([]models.Feedback, error) {
	if _, err := os.Stat(feedbackFile); os.IsNotExist(err) {
		return []models.Feedback{}, nil
	}
	data, err := os.ReadFile(feedbackFile)
	if err != nil {
		return nil, err
	}
	var feedbacks []models.Feedback
	err = json.Unmarshal(data, &feedbacks)
	return feedbacks, err
}
