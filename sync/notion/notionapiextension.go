package notion

import (
	"fmt"

	"github.com/jomei/notionapi"
)

func GetStringValueFromProperty(property notionapi.Property) string {
	switch property.GetType() {
	case notionapi.PropertyTypeRichText:
		richText := property.(*notionapi.RichTextProperty).RichText
		var value string
		for index, richTextItem := range richText {
			value += richTextItem.PlainText
			if index != len(richText)-1 {
				value += " "
			}
		}
		return value
	case notionapi.PropertyTypeText:
		textProperty := property.(*notionapi.TextProperty).Text
		var value string
		for index, item := range textProperty {
			value += item.PlainText
			if index != len(textProperty)-1 {
				value += " "
			}
		}
		return value
	case notionapi.PropertyTypeTitle:
		textProperty := property.(*notionapi.TitleProperty).Title
		var value string
		for index, item := range textProperty {
			value += item.PlainText
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
		return property.(*notionapi.DateProperty).Date.Start.String()
	case notionapi.PropertyTypeFormula:
		return property.(*notionapi.FormulaProperty).Formula.String
	case notionapi.PropertyTypeRelation:
		relation := property.(*notionapi.RelationProperty).Relation
		var value string
		for index, item := range relation {
			value += item.ID.String() + " "
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
			value += item.Name + " "
			if index != len(people)-1 {
				value += " "
			}
		}
		return value
	case notionapi.PropertyTypeFiles:
		files := property.(*notionapi.FilesProperty).Files
		var value string
		for index, item := range files {
			value += item.Name + " "
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
