package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Activity string

const (
	ActivityCoding Activity = "coding"
)

var toActivity = map[string]Activity{
	"coding": ActivityCoding,
}

func ActivityExist(activity string) bool {
	_, ok := toActivity[activity]
	return ok
}

func (a *Activity) UnmarshalJSON(b []byte) error {
	var value string
	err := json.Unmarshal(b, &value)
	if err != nil {
		return err
	}

	activity, ok := toActivity[value]
	if !ok {
		return errors.New(fmt.Sprintf("%s is not a valid activity", value))
	}

	*a = activity
	return nil
}

func (a Activity) MarshalJSON() ([]byte, error) {
	for key, activity := range toActivity {
		if activity == a {
			return json.Marshal(key)
		}
	}
	return nil, errors.New("activity not found string representation")
}

type Category string

const (
	CategoryProject  Category = "project"
	CategoryGit      Category = "git"
	CategoryFilename Category = "filename"
	CategoryBranch   Category = "branch"
	CategoryLanguage Category = "language"
	CategoryEditor   Category = "editor"
	CategoryKeyCount Category = "keycount"
)

var toCategory = map[string]Category{
	"project":  CategoryProject,
	"git":      CategoryGit,
	"filename": CategoryFilename,
	"branch":   CategoryBranch,
	"language": CategoryLanguage,
	"editor":   CategoryEditor,
	"keycount": CategoryKeyCount,
}

func CategoryExist(category string) bool {
	_, ok := toCategory[category]
	return ok
}

func (c *Category) UnmarshalJSON(b []byte) error {
	var value string
	err := json.Unmarshal(b, &value)
	if err != nil {
		return err
	}

	category, ok := toCategory[value]
	if !ok {
		return errors.New(fmt.Sprintf("%s is not a valid category", value))
	}

	*c = category
	return nil
}

func (c Category) MarshalJSON() ([]byte, error) {
	for key, category := range toCategory {
		if category == c {
			return json.Marshal(key)
		}
	}
	return nil, errors.New("label category not found string representation")
}

type Label struct {
	Category Category `json:"category"`
	Value    string   `json:"value"`
}

type Record struct {
	Start     time.Time `json:"start"`
	Stop      time.Time `json:"stop"`
	Activity  Activity  `json:"activity"`
	SessionId string    `json:"session_id"`
	Labels    []Label   `json:"labels"`
}
