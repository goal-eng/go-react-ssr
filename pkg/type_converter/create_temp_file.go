package type_converter

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/natewong1313/go-react-ssr/internal/utils"
	"github.com/natewong1313/go-react-ssr/pkg/config"
)

func createCacheFolder() (string, error) {
	osCacheDir, _ := os.UserCacheDir()
	cacheFolderPath := filepath.Join(osCacheDir, "gossr")
	os.RemoveAll(cacheFolderPath)
	err := os.MkdirAll(cacheFolderPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return cacheFolderPath, nil
}

// https://github.com/tkrajina/typescriptify-golang-structs/blob/master/tscriptify/main.go#L139
func createTemporaryFile(cfg config.Config, folderPath string, structNames []string) (string, error) {
	temporaryFilePath := filepath.Join(folderPath, "generator.go")
	file, err := os.Create(temporaryFilePath)
	if err != nil {
		return temporaryFilePath, err
	}
	defer file.Close()

	t := template.Must(template.New("").Parse(TEMPLATE))

	structsArr := make([]string, 0)
	for _, structName := range structNames {
		structName = strings.TrimSpace(structName)
		if len(structName) > 0 {
			structsArr = append(structsArr, "m."+structName)
		}
	}

	var params TemplateParams
	params.Structs = structsArr

	params.ModuleName, err = getModuleName(cfg.PropsStructsPath)
	if err != nil {
		return temporaryFilePath, err
	}
	params.Interface = true
	params.TargetFile = utils.GetFullFilePath(cfg.GeneratedTypesPath)

	err = t.Execute(file, params)
	if err != nil {
		return temporaryFilePath, err
	}

	return temporaryFilePath, nil
}

func getModuleName(propsStructsPath string) (string, error) {
	dir := filepath.Dir(utils.GetFullFilePath(propsStructsPath))
	cmd := exec.Command("go", "list")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
