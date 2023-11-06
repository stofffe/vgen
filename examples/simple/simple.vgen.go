// THIS FILE IS GENERATED BY VGEN
// DO NOT EDIT
package main
import (
	"encoding/json"
	"fmt"
)
type EmailVgen struct {
	Title *string `json:"title"`
	Text *string `json:"text"`
	Sender *string `json:"sender"`
}
func (t EmailVgen) Validate() (Email, error) {
	errs := t.InnerValidation()
	if len(errs) > 0 {
		j, _ := json.Marshal(errs)
		return Email{}, fmt.Errorf("%s", j)
	}
	return Email{}, nil
}
func (t EmailVgen) InnerValidation() map[string][]string {
	errs := make(map[string][]string)
	if t.Title != nil {
		_Title := *t.Title
		{
			if !(len(_Title) > 0) {
				errs[fmt.Sprintf("title")] = append(errs[fmt.Sprintf("title")], fmt.Sprintf("can not be empty"))
			}
			if !(len(_Title) < 50) {
				errs[fmt.Sprintf("title")] = append(errs[fmt.Sprintf("title")], fmt.Sprintf("len must be < 50"))
			}
			_ = _Title
		}
	} else {
		errs["title"] = append(errs["title"], fmt.Sprintf("required"))
	}
	if t.Text != nil {
		_Text := *t.Text
		{
			if !(len(_Text) > 0) {
				errs[fmt.Sprintf("text")] = append(errs[fmt.Sprintf("text")], fmt.Sprintf("can not be empty"))
			}
			if !(len(_Text) > 200) {
				errs[fmt.Sprintf("text")] = append(errs[fmt.Sprintf("text")], fmt.Sprintf("len must be > 200"))
			}
			_ = _Text
		}
	} else {
		errs["text"] = append(errs["text"], fmt.Sprintf("required"))
	}
	if t.Sender != nil {
		_Sender := *t.Sender
		{
			if !(len(_Sender) > 0) {
				errs[fmt.Sprintf("sender")] = append(errs[fmt.Sprintf("sender")], fmt.Sprintf("can not be empty"))
			}
			if !(len(_Sender) < 20) {
				errs[fmt.Sprintf("sender")] = append(errs[fmt.Sprintf("sender")], fmt.Sprintf("len must be < 20"))
			}
			_ = _Sender
		}
	} else {
		errs["sender"] = append(errs["sender"], fmt.Sprintf("required"))
	}
	return errs
}
func (t EmailVgen) Convert() Email {
	var res Email
	if t.Title != nil {
		_Title := *t.Title
		res.Title = _Title
	}
	if t.Text != nil {
		_Text := *t.Text
		res.Text = _Text
	}
	if t.Sender != nil {
		_Sender := *t.Sender
		res.Sender = _Sender
	}
	return Email{}
}
func EmailFromJson(bytes []byte) (Email, error) {
	var v EmailVgen
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return Email{}, err
	}
	r, err := v.Validate()
	if err != nil {
		return Email{}, err
	}
	return r, nil
}
