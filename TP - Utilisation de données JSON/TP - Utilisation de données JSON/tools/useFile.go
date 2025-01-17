package tools

import (
	"fmt"
	"os"
)

func ReadFileJson(path string) ([]byte, error) {
	dataFile, errFile := os.ReadFile(path)
	if errFile != nil {
		return nil, fmt.Errorf("ReadFileJson - Erreur de lecture du fichier : %s", errFile.Error())
	}
	return dataFile, nil
}

func WriteFileJson(path string, data []byte) error {
	errFile := os.WriteFile(path, data, 644)
	if errFile != nil {
		return fmt.Errorf("WriteFileJson - Erreur d'écriture du fichier : %s", errFile.Error())
	}
	return nil
}
