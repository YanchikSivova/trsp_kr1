package main

import (
	"calulator/models"
	"calulator/storage"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	calculate := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var nums models.Nums
		json.NewDecoder(r.Body).Decode(&nums)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int{"result": nums.Num1 + nums.Num2})
	}

	user := models.User{Name: "Сивова Яна", ID: 1}

	getUsers := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}

	postUser := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var user models.UserWithAge
		json.NewDecoder(r.Body).Decode(&user)
		userIsAdult := models.UserIsAdult{
			Name:    user.Name,
			Age:     user.Age,
			IsAdult: user.Age > 18,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(userIsAdult)
	}

	postFeedback := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var feedback models.Feedback
		if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
			sendValidationError(w, []models.ValidationErrorDetail{
				{
					Type:     "json_invalid",
					Location: []string{"body"},
					Msg:      "Invalid JSON",
					Input:    nil,
				},
			}, http.StatusBadRequest)
			return
		}
		validationErrors := validateFeedback(feedback)
		if len(validationErrors) > 0 {
			sendValidationError(w, validationErrors, http.StatusBadRequest)
			return
		}
		err := storage.SaveFeedback(feedback)
		if err != nil {
			http.Error(w, "Feedback could not be saved: "+err.Error(), http.StatusInternalServerError)
			return
		}
		feedbackResponse := fmt.Sprintf("Feedback received! Thank you, %s.", feedback.Name)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": feedbackResponse})

	}

	http.HandleFunc("/calculate", calculate)
	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/user", postUser)
	http.HandleFunc("/feedback", postFeedback)

	http.ListenAndServe(":8080", nil)
}

func validateFeedback(f models.Feedback) []models.ValidationErrorDetail {
	var errors []models.ValidationErrorDetail
	if f.Name == "" {
		errors = append(errors, models.ValidationErrorDetail{
			Type:     "missing",
			Location: []string{"body", "name"},
			Msg:      "field 'name' is required",
			Input:    f.Name,
			Ctx:      map[string]int{"min_length": 2},
		})
	} else if len(f.Name) < 2 {
		errors = append(errors, models.ValidationErrorDetail{
			Type:     "string_too_short",
			Location: []string{"body", "name"},
			Msg:      "String should have at least 2 characters",
			Input:    f.Name,
			Ctx:      map[string]int{"min_length": 2},
		})
	} else if len(f.Name) > 50 {
		errors = append(errors, models.ValidationErrorDetail{
			Type:     "string_too_long",
			Location: []string{"body", "name"},
			Msg:      "String should have at most 50 characters",
			Input:    f.Name,
			Ctx:      map[string]int{"max_length": 50},
		})
	}

	if f.Message == "" {
		errors = append(errors, models.ValidationErrorDetail{
			Type:     "missing",
			Location: []string{"body", "message"},
			Msg:      "field 'message' is required",
			Input:    f.Message,
			Ctx:      map[string]int{"min_length": 10},
		})
	} else if len(f.Message) < 10 {
		errors = append(errors, models.ValidationErrorDetail{
			Type:     "string_too_short",
			Location: []string{"body", "message"},
			Msg:      "String should have at least 10 characters",
			Input:    f.Message,
			Ctx:      map[string]int{"min_length": 10},
		})
	} else if len(f.Message) > 500 {
		errors = append(errors, models.ValidationErrorDetail{
			Type:     "string_too_long",
			Location: []string{"body", "message"},
			Msg:      "String should have at most 500 characters",
			Input:    f.Message,
			Ctx:      map[string]int{"max_length": 500},
		})
	} else {
		lowerMessage := strings.ToLower(f.Message)
		for _, word := range models.ForbiddenWords {
			if strings.Contains(lowerMessage, word) {
				errors = append(errors, models.ValidationErrorDetail{
					Type:     "value_error",
					Location: []string{"body", "message"},
					Msg:      "Value error, Использование недопустимых слов",
					Input:    f.Message,
					Ctx:      map[string]interface{}{"error": map[string]string{"forbidden_word": word}},
				})
				break
			}
		}
	}
	return errors
}

func sendValidationError(w http.ResponseWriter, errors []models.ValidationErrorDetail, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.ValidationError{Detail: errors})
}
