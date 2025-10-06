package utils

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github/jornal-cidadao-jc/internal/model"
)

func Get_charges_object(charges_dir string) ([]model.Charge, error) {
	files, err := os.ReadDir(charges_dir)
	if err != nil {
		return nil, err
	}

	var charges []model.Charge

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
			Title:    "", 
		})
	}

	sort.Slice(charges, func(i, j int) bool {
		return time.Time(charges[i].Date).Before(time.Time(charges[j].Date))
	})

	for i := range charges {
		charges[i].ID = i + 1
	}

	return charges, nil
}
