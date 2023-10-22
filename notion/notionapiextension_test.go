package notion

import (
	"testing"
	"time"

	"github.com/jomei/notionapi"
)

var now = time.Now()
var notionDate = notionapi.Date(now)

func TestGetStringValueFromProperty(t *testing.T) {
	richText := []notionapi.RichText{
		{
			Type: "rich_text",
			Text: &notionapi.Text{
				Content: "Hello",
			},
		},
		{
			Type: "rich_text",
			Text: &notionapi.Text{
				Content: "world!",
			},
		},
	}
	richTextProperty := &notionapi.RichTextProperty{
		RichText: richText,
		Type:     "rich_text",
	}
	textProperty := &notionapi.TextProperty{
		Text: richText,
		Type: "text",
	}
	titleProperty := &notionapi.TitleProperty{
		Title: richText,
		Type:  "title",
	}
	numberProperty := &notionapi.NumberProperty{
		Number: 42,
		Type:   "number",
	}
	selectProperty := &notionapi.SelectProperty{
		Select: notionapi.Option{
			Name: "Option 1",
		},
		Type: "select",
	}
	multiSelectProperty := &notionapi.MultiSelectProperty{
		MultiSelect: []notionapi.Option{
			{
				Name: "Option 1",
			},
			{
				Name: "Option 2",
			},
		},
		Type: "multi_select",
	}
	dateProperty := &notionapi.DateProperty{
		Date: &notionapi.DateObject{
			Start: &notionDate,
		},
		Type: "date",
	}
	formulaProperty := &notionapi.FormulaProperty{
		Formula: notionapi.Formula{
			String: "Hello world!",
		},
		Type: "formula",
	}
	relationProperty := &notionapi.RelationProperty{
		Relation: []notionapi.Relation{
			{
				ID: notionapi.PageID("abc123"),
			},
			{
				ID: notionapi.PageID("def456"),
			},
		},
		Type: "relation",
	}
	rollupProperty := &notionapi.RollupProperty{
		Rollup: notionapi.Rollup{
			Date: &notionapi.DateObject{
				Start: &notionDate,
				End:   &notionDate,
			},
		},
		Type: "rollup",
	}
	peopleProperty := &notionapi.PeopleProperty{
		People: []notionapi.User{
			{
				Name: "John Doe",
			},
			{
				Name: "Jane Doe",
			},
		},
		Type: "people",
	}
	filesProperty := &notionapi.FilesProperty{
		Files: []notionapi.File{
			{
				Name: "file1.txt",
			},
			{
				Name: "file2.txt",
			},
		},
		Type: "files",
	}
	checkboxProperty := &notionapi.CheckboxProperty{
		Checkbox: true,
		Type:     "checkbox",
	}
	urlProperty := &notionapi.URLProperty{
		URL:  "https://example.com",
		Type: "url",
	}

	tests := []struct {
		name     string
		property notionapi.Property
		want     string
	}{
		{
			name:     "RichTextProperty",
			property: richTextProperty,
			want:     "Hello world!",
		},
		{
			name:     "TextProperty",
			property: textProperty,
			want:     "Hello world!",
		},
		{
			name:     "TitleProperty",
			property: titleProperty,
			want:     "Hello world!",
		},
		{
			name:     "NumberProperty",
			property: numberProperty,
			want:     "42",
		},
		{
			name:     "SelectProperty",
			property: selectProperty,
			want:     "Option 1",
		},
		{
			name:     "MultiSelectProperty",
			property: multiSelectProperty,
			want:     "Option 1 Option 2",
		},
		{
			name:     "DateProperty",
			property: dateProperty,
			want:     time.Now().Format(time.RFC3339),
		},
		{
			name:     "FormulaProperty",
			property: formulaProperty,
			want:     "Hello world!",
		},
		{
			name:     "RelationProperty",
			property: relationProperty,
			want:     "abc123 def456",
		},
		{
			name:     "RollupProperty",
			property: rollupProperty,
			want:     now.Format(time.RFC3339) + " " + now.Format(time.RFC3339),
		},
		{
			name:     "PeopleProperty",
			property: peopleProperty,
			want:     "John Doe Jane Doe",
		},
		{
			name:     "FilesProperty",
			property: filesProperty,
			want:     "file1.txt file2.txt",
		},
		{
			name:     "CheckboxProperty",
			property: checkboxProperty,
			want:     "true",
		},
		{
			name:     "URLProperty",
			property: urlProperty,
			want:     "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStringValueFromProperty(tt.property); got != tt.want {
				t.Errorf("GetStringValueFromProperty() = %v, want %v, property type %v", got, tt.want, tt.property.GetType())
			}
		})
	}
}

func TestUpdateValueFromProp(t *testing.T) {
	now := time.Now()
	richText := []notionapi.RichText{
		{
			Type: "rich_text",
			Text: &notionapi.Text{
				Content: "Hello",
			},
		},
		{
			Type: "rich_text",
			Text: &notionapi.Text{
				Content: "world!",
			},
		},
	}
	richTextProperty := &notionapi.RichTextProperty{
		RichText: richText,
		Type:     "rich_text",
	}
	textProperty := &notionapi.TextProperty{
		Text: richText,
		Type: "text",
	}
	titleProperty := &notionapi.TitleProperty{
		Title: richText,
		Type:  "title",
	}
	numberProperty := &notionapi.NumberProperty{
		Number: 42,
		Type:   "number",
	}
	selectProperty := &notionapi.SelectProperty{
		Select: notionapi.Option{
			Name: "Option 1",
		},
		Type: "select",
	}
	dateProperty := NewDateProperty(now)
	checkboxProperty := &notionapi.CheckboxProperty{
		Checkbox: true,
		Type:     "checkbox",
	}
	urlProperty := &notionapi.URLProperty{
		URL:  "https://example.com",
		Type: "url",
	}

	tests := []struct {
		name         string
		property     notionapi.Property
		newValue     string
		expectedProp notionapi.Property
	}{
		{
			name:         "RichTextProperty",
			property:     richTextProperty,
			newValue:     "New value",
			expectedProp: &notionapi.RichTextProperty{RichText: []notionapi.RichText{{Type: "text", Text: &notionapi.Text{Content: "New value"}}}, Type: "rich_text"},
		},
		{
			name:         "TextProperty",
			property:     textProperty,
			newValue:     "New value",
			expectedProp: &notionapi.TextProperty{Text: []notionapi.RichText{{Type: "text", Text: &notionapi.Text{Content: "New value"}}}, Type: "text"},
		},
		{
			name:         "TitleProperty",
			property:     titleProperty,
			newValue:     "New value",
			expectedProp: &notionapi.TitleProperty{Title: []notionapi.RichText{{Type: "text", Text: &notionapi.Text{Content: "New value"}}}, Type: "title"},
		},
		{
			name:         "NumberProperty",
			property:     numberProperty,
			newValue:     "42.5",
			expectedProp: &notionapi.NumberProperty{Number: 42.5, Type: "number"},
		},
		{
			name:         "SelectProperty",
			property:     selectProperty,
			newValue:     "New value",
			expectedProp: &notionapi.SelectProperty{Select: notionapi.Option{Name: "New value", ID: "", Color: ""}, Type: "select"},
		},
		{
			name:         "DateProperty",
			property:     &dateProperty,
			newValue:     now.Format(time.RFC3339),
			expectedProp: NewDateProperty(now),
		},
		{
			name:         "CheckboxProperty",
			property:     checkboxProperty,
			newValue:     "false",
			expectedProp: &notionapi.CheckboxProperty{Checkbox: false, Type: "checkbox"},
		},
		{
			name:         "URLProperty",
			property:     urlProperty,
			newValue:     "https://newexample.com",
			expectedProp: &notionapi.URLProperty{URL: "https://newexample.com", Type: "url"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prop, _ := UpdateValueFromProp(tt.property, tt.newValue)
			if GetStringValueFromProperty(prop) != GetStringValueFromProperty(tt.property) {
				t.Errorf("UpdateValueFromProp() property = %v, expected %v, type %v", prop, tt.expectedProp, tt.property.GetType())
			}
		})
	}
}
