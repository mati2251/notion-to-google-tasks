package notion

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jomei/notionapi"
)

func GetStringValueFromProperty(property notionapi.Property) string {
	switch property.GetType() {
	case notionapi.PropertyTypeRichText:
		richText := property.(*notionapi.RichTextProperty).RichText
		var value string
		for index, richTextItem := range richText {
			value += richTextItem.Text.Content
			if index != len(richText)-1 {
				value += " "
			}
		}
		return value
	case notionapi.PropertyTypeText:
		textProperty := property.(*notionapi.TextProperty).Text
		var value string
		for index, item := range textProperty {
			value += item.Text.Content
			if index != len(textProperty)-1 {
				value += " "
			}
		}
		return value
	case notionapi.PropertyTypeTitle:
		textProperty := property.(*notionapi.TitleProperty).Title
		var value string
		for index, item := range textProperty {
			value += item.Text.Content
			if index != len(textProperty)-1 {
				value += " "
			}
		}
		return value
	case notionapi.PropertyTypeNumber:
		return fmt.Sprintf("%v", property.(*notionapi.NumberProperty).Number)
	case notionapi.PropertyTypeSelect:
		return property.(*notionapi.SelectProperty).Select.Name
	case notionapi.PropertyTypeMultiSelect:
		multiSelect := property.(*notionapi.MultiSelectProperty).MultiSelect
		var value string
		for index, item := range multiSelect {
			value += item.Name
			if index != len(multiSelect)-1 {
				value += " "
			}
		}
		return value
	case notionapi.PropertyTypeDate:
		if property.(*notionapi.DateProperty).Date == nil {
			return ""
		}
		if property.(*notionapi.DateProperty).Date.End == nil {
			if property.(*notionapi.DateProperty).Date.Start == nil {
				return ""
			}
			return property.(*notionapi.DateProperty).Date.Start.String()
		}
		return property.(*notionapi.DateProperty).Date.End.String()
	case notionapi.PropertyTypeFormula:
		return property.(*notionapi.FormulaProperty).Formula.String
	case notionapi.PropertyTypeRelation:
		relation := property.(*notionapi.RelationProperty).Relation
		var value string
		for index, item := range relation {
			value += item.ID.String()
			if index != len(relation)-1 {
				value += " "
			}
		}
		return value
	case notionapi.PropertyTypeRollup:
		rollup := property.(*notionapi.RollupProperty).Rollup
		return rollup.Date.End.String() + " " + rollup.Date.Start.String()
	case notionapi.PropertyTypePeople:
		people := property.(*notionapi.PeopleProperty).People
		var value string
		for index, item := range people {
			value += item.Name
			if index != len(people)-1 {
				value += " "
			}
		}
		return value
	case notionapi.PropertyTypeFiles:
		files := property.(*notionapi.FilesProperty).Files
		var value string
		for index, item := range files {
			value += item.Name
			if index != len(files)-1 {
				value += " "
			}
		}
		return value
	case notionapi.PropertyTypeCheckbox:
		return fmt.Sprintf("%v", property.(*notionapi.CheckboxProperty).Checkbox)
	case notionapi.PropertyTypeURL:
		return property.(*notionapi.URLProperty).URL
	case notionapi.PropertyTypeEmail:
		return property.(*notionapi.EmailProperty).Email
	case notionapi.PropertyTypePhoneNumber:
		return property.(*notionapi.PhoneNumberProperty).PhoneNumber
	case notionapi.PropertyTypeCreatedTime:
		return property.(*notionapi.CreatedTimeProperty).CreatedTime.String()
	case notionapi.PropertyTypeCreatedBy:
		return property.(*notionapi.CreatedByProperty).CreatedBy.Name
	case notionapi.PropertyTypeLastEditedTime:
		return property.(*notionapi.LastEditedTimeProperty).LastEditedTime.String()
	case notionapi.PropertyTypeLastEditedBy:
		return property.(*notionapi.LastEditedByProperty).LastEditedBy.Name
	case notionapi.PropertyTypeStatus:
		return property.(*notionapi.StatusProperty).Status.Name
	case notionapi.PropertyTypeUniqueID:
		return property.(*notionapi.UniqueIDProperty).UniqueID.String()
	case notionapi.PropertyTypeVerification:
		return property.(*notionapi.VerificationProperty).Verification.VerifiedBy.Name
	default:
		return ""
	}
}

func UpdateValueFromProp(property notionapi.Property, newValue string) (notionapi.Property, error) {
	var obj notionapi.Property = property
	if property == nil {
		return nil, errors.New("property is nil")
	}
	switch property.GetType() {
	case notionapi.PropertyTypeRichText:
		obj.(*notionapi.RichTextProperty).RichText = []notionapi.RichText{
			{
				Type: "text",
				Text: &notionapi.Text{
					Content: newValue,
				},
			}}
	case notionapi.PropertyTypeText:
		obj.(*notionapi.TextProperty).Text = []notionapi.RichText{
			{
				Type: "text",
				Text: &notionapi.Text{
					Content: newValue,
				},
			}}
	case notionapi.PropertyTypeTitle:
		obj.(*notionapi.TitleProperty).Title = []notionapi.RichText{
			{
				Type: "text",
				Text: &notionapi.Text{
					Content: newValue,
				},
			}}
	case notionapi.PropertyTypeNumber:
		num, err := strconv.ParseFloat(newValue, 64)
		if err != nil {
			return nil, errors.Join(err, errors.New("error parsing newValue to float"))
		}
		obj.(*notionapi.NumberProperty).Number = num
	case notionapi.PropertyTypeSelect:
		obj.(*notionapi.SelectProperty).Select.Name = newValue
		obj.(*notionapi.SelectProperty).Select.ID = ""
		obj.(*notionapi.SelectProperty).Select.Color = ""
	case notionapi.PropertyTypeDate:
		date, err := time.Parse(time.RFC3339, newValue)
		notionDate := notionapi.Date(date)
		if err != nil {
			return nil, errors.Join(err, errors.New("error parsing newValue to time"))
		}
		dateObject := notionapi.DateObject{
			Start: &notionDate,
		}
		obj.(*notionapi.DateProperty).Date = &dateObject
	case notionapi.PropertyTypeCheckbox:
		checked, err := strconv.ParseBool(newValue)
		if err != nil {
			return nil, errors.Join(err, errors.New("error parsing newValue to bool"))
		}
		obj.(*notionapi.CheckboxProperty).Checkbox = checked
	case notionapi.PropertyTypeURL:
		obj.(*notionapi.URLProperty).URL = newValue
	case notionapi.PropertyTypeEmail:
		obj.(*notionapi.EmailProperty).Email = newValue
	case notionapi.PropertyTypePhoneNumber:
		obj.(*notionapi.PhoneNumberProperty).PhoneNumber = newValue
	case notionapi.PropertyTypeCreatedTime:
		date, err := time.Parse(time.RFC3339, newValue)
		if err != nil {
			return nil, errors.Join(err, errors.New("error parsing newValue to time"))
		}
		obj.(*notionapi.CreatedTimeProperty).CreatedTime = date
	case notionapi.PropertyTypeLastEditedTime:
		date, err := time.Parse(time.RFC3339, newValue)
		if err != nil {
			return nil, errors.Join(err, errors.New("error parsing newValue to time"))
		}
		obj.(*notionapi.LastEditedTimeProperty).LastEditedTime = date
	case notionapi.PropertyTypeStatus:
		obj.(*notionapi.StatusProperty).Status.Name = newValue
		obj.(*notionapi.StatusProperty).Status.ID = ""
		obj.(*notionapi.StatusProperty).Status.Color = ""
	default:
		return nil, errors.New("property type not supported")
	}
	return obj, nil
}

func NewRichText(content string) notionapi.RichText {
	return notionapi.RichText{
		Type: "text",
		Text: &notionapi.Text{
			Content: content,
		},
	}
}

func NewDateProperty(date time.Time) notionapi.DateProperty {
	notionDate := notionapi.Date(date)
	return notionapi.DateProperty{
		Type: "date",
		Date: &notionapi.DateObject{
			Start: &notionDate,
		},
	}
}

func NewRichTextProperty(content string) *notionapi.RichTextProperty {
	return &notionapi.RichTextProperty{
		Type: "rich_text",
		RichText: []notionapi.RichText{
			NewRichText(content),
		},
	}
}
