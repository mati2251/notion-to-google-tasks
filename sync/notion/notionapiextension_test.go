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
	}
	textProperty := &notionapi.TextProperty{
		Text: richText,
	}
	titleProperty := &notionapi.TitleProperty{
		Title: richText,
	}
	numberProperty := &notionapi.NumberProperty{
		Number: 42,
	}
	selectProperty := &notionapi.SelectProperty{
		Select: notionapi.Option{
			Name: "Option 1",
		},
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
	}
	dateProperty := &notionapi.DateProperty{
		Date: &notionapi.DateObject{
			Start: &notionDate,
		},
	}
	formulaProperty := &notionapi.FormulaProperty{
		Formula: notionapi.Formula{
			String: "Hello world!",
		},
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
	}
	rollupProperty := &notionapi.RollupProperty{
		Rollup: notionapi.Rollup{
			Date: &notionapi.DateObject{
				Start: &notionDate,
				End:   &notionDate,
			},
		},
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
	}
	checkboxProperty := &notionapi.CheckboxProperty{
		Checkbox: true,
	}
	urlProperty := &notionapi.URLProperty{
		URL: "https://example.com",
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
			want:     time.Now().Format("2006-01-02 15:04:05 -0700 MST"),
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
			want:     time.Now().Format("2006-01-02 15:04:05 -0700 MST") + " " + time.Now().AddDate(0, 0, 1).Format("2006-01-02 15:04:05 -0700 MST"),
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
				t.Errorf("GetStringValueFromProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}
