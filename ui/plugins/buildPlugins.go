package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// 1. 将 demo.svg 放入 public/plugins/external/panel 下
// 2. 新增 plugin.json 在 public/plugins/external/panel 下
// 3. 将 demo 代码文件夹移动到 src/views/dashboard/plugins/panel 下
// 4. 自动生成 src/views/dashboard/plugins/externalPlugins 文件

func main() {
	// Walk("./panel", true, true, walker)
	cmd := exec.Command("bash", "-c", "mkdir -p ../public/plugins/external/panel")
	cmd.CombinedOutput()

	cmd = exec.Command("bash", "-c", "mkdir -p ../public/plugins/external/datasource")
	cmd.CombinedOutput()

	// generate panel plugins
	panels, err := os.ReadDir("./panel")
	if err != nil {
		log.Fatal("read panel direrror", err)
	}

	pluginsList := make([]map[string]string, 0)

	externalPluginsStr1 := `
import { PanelPluginComponents } from "types/plugins/plugin"
`
	externalPluginsStr2 := `
export const panelPlugins: Record<string,PanelPluginComponents> = {`
	for _, panel := range panels {
		panelType := panel.Name()
		// cp .svg to public/plugins/external/panel
		cmdStr := fmt.Sprintf("cp ./panel/%s/%s.svg ../public/plugins/external/panel", panelType, panelType)
		cmd := exec.Command("bash", "-c", cmdStr)
		if _, err := cmd.CombinedOutput(); err != nil {
			log.Println("copy plugin .svg  error: ", err, ", panel: ", panelType)
			continue
		}

		// cp panel codes into src/views/dashboard/plugins/panel
		cmdStr = fmt.Sprintf("cp -r ./panel/%s ../src/views/dashboard/plugins/panel", panelType)
		cmd = exec.Command("bash", "-c", cmdStr)
		if _, err := cmd.CombinedOutput(); err != nil {
			log.Println("copy plugin code dir  error: ", err, ", panel: ", panelType)
			continue
		}

		pluginsList = append(pluginsList, map[string]string{
			"type": panelType,
		})

		componentStr := strings.Title(panelType) + "Components"
		externalPluginsStr1 += fmt.Sprintf("\nimport %s from \"./panel/%s\"", componentStr, panelType)

		externalPluginsStr2 += fmt.Sprintf("\n\t\"%s\": %s,", panelType, componentStr)
	}

	externalPluginsStr2 += "\n}"

	externalPluginFile := externalPluginsStr1 + externalPluginsStr2

	// generate externalPlugins.ts file
	err = os.WriteFile("../src/views/dashboard/plugins/externalPlugins.ts", []byte(externalPluginFile), 0666)
	if err != nil {
		log.Fatal("write plugin.json error", err)
	}

	// generate plugin.json
	pluginsJson, _ := json.Marshal(pluginsList)
	err = os.WriteFile("../public/plugins/external/panel/plugins.json", pluginsJson, 0666)
	if err != nil {
		log.Fatal("write plugin.json error", err)
	}

	log.Println("Generate panel plugins file successfully!")
}
