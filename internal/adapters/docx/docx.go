package docxadapter

import (
	"github.com/nguyenthenguyen/docx"
)

func GenerateTicketDocx(templatePath string, outputPath string, data map[string]string) error {
	r, err := docx.ReadDocxFile(templatePath)
	if err != nil {
		return err
	}
	doc := r.Editable()
	for k, v := range data {
		doc.Replace(k, v, -1)
	}
	return doc.WriteToFile(outputPath)
}
