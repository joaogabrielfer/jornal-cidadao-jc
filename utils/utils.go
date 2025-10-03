package utils

import(
	"os"
	"log"
	"path/filepath"

	"github/jornal-cidadao-jc/internal/model"
)

func Get_charges_object(charges_dir string) (charges []model.Charge, err error){
	files, err := os.ReadDir(charges_dir)
	if err != nil {
		return nil, err
	}


	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(charges_dir, file.Name())
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Println("Erro obtendo informação da charge: ", file.Name(), err)
			continue
		}

		charges = append(charges, model.Charge{
			Filename: file.Name(),
			URL:      filepath.Join("/static/images/charges", file.Name()),
			Date:     model.FormattedTime(fileInfo.ModTime()),
			Title:	  "",
		})
	}

	total_charges := len(charges)
	for i := range charges {
		charges[i].ID = total_charges - i
	}


	return charges, err
}

