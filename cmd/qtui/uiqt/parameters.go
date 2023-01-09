package uiqt

import (
	"fmt"
	"github.com/xdumpgo/XDG/modules"
	"github.com/xdumpgo/XDG/qtui"
	"github.com/therecipe/qt/widgets"
	"os"
	"strings"
)

func NewParametersWindow() *qtui.ParametersSettings {

	qtui.ParametersSettingsWindow = qtui.NewParametersSettings(nil)

	// Add Button
	qtui.ParametersSettingsWindow.AddParams.ConnectClicked(func(checked bool) {
		ParameterDialog := qtui.NewAddParamDialog(qtui.ParametersSettingsWindow)

		//Select file button
		ParameterDialog.SelectButton.ConnectClicked(func(bool) {
			wd, _ := os.Getwd()
			FileDialog := widgets.NewQFileDialog2(qtui.ParametersSettingsWindow, "Select filepath for parameter", wd, "Text / List files (*.txt *.lst)")

			FileDialog.ConnectFileSelected(func(file string) {
				ParameterDialog.SelectFileInputField.SetText(file)
			})

			FileDialog.Show()
		})

		//Add Parameter (Ok button)
		ParameterDialog.ButtonBox.ConnectAccepted(func() {
			name := ParameterDialog.NameInputField.Text()
			prefix := ParameterDialog.PrefixInputField.Text()
			filepath := ParameterDialog.SelectFileInputField.Text()

			if len(name) == 0 || len(prefix) == 0 || len(filepath) == 0 {
				qtui.SimpleMB(ParameterDialog, "Please make sure you fill all fields", "Invalid input").Show()
				return
			}

			if _, exist := modules.Generator.Parameters[prefix]; !exist {
				modules.Generator.Parameters[prefix] = &modules.Parameter{
					Name: name,
					FilePath: filepath,
					Prefix: prefix,
				}
				qtui.ParametersSettingsWindow.ParametersList.AddItem(fmt.Sprintf("%s-%s", name, prefix))
			}else {
				qtui.SimpleMB(qtui.ParametersSettingsWindow, "Failed to add parameter", "Error")
			}

		})

		ParameterDialog.Show()
	})

	//Remove Parameter Button
	qtui.ParametersSettingsWindow.RemoveParams.ConnectClicked(func(bool) {
		if index := qtui.ParametersSettingsWindow.ParametersList.CurrentRow(); index > -1 {
			item := qtui.ParametersSettingsWindow.ParametersList.CurrentItem()
			qtui.ParametersSettingsWindow.ParametersList.TakeItem(index)

			if _, exists := modules.Generator.Parameters[strings.Split(item.Text(), "-")[1]]; exists {
				fmt.Println()
				delete(modules.Generator.Parameters, strings.Split(item.Text(), "-")[1])
			}
		}
	})

	//Settings Ok button
	qtui.ParametersSettingsWindow.ButtonBox.ConnectAccepted(func() {
		qtui.Main.ComboBoxParameters.Clear()
		var names []string
		for prefix, parameter := range modules.Generator.Parameters {
			names = append(names, fmt.Sprintf("%s-%s", parameter.Name, prefix))
		}
		qtui.Main.ComboBoxParameters.AddItems(names)

		qtui.ParametersSettingsWindow.Close()
	})

	qtui.ParametersSettingsWindow.ButtonBox.ConnectClicked(func(button *widgets.QAbstractButton) {
		if button.Text() == "Cancel" {
			qtui.ParametersSettingsWindow.Close()
		}
	})

	/*qtui.ParametersSettingsWindow = qtui.NewParametersSettings(nil)

	for prefix, parameter := range modules.Generator.Parameters{
		qtui.ParametersSettingsWindow.GeneratorParams.AddItem(fmt.Sprintf("%s - %s", parameter.Name, prefix))
	}

	qtui.ParametersSettingsWindow.GeneratorAddParams.ConnectClicked(func(checked bool) {
		add := qtui.NewAddParamDialog(qtui.ParametersSettingsWindow)
		add.PushButton.ConnectClicked(func(bool) {
			wd, _ := os.Getwd()
			ofd := widgets.NewQFileDialog2(qtui.ParametersSettingsWindow, "Select filepath for parameter", wd, "Text / List files (*.txt *.lst)")

			ofd.ConnectFileSelected(func(file string) {
				add.LineEdit_3.SetText(file)
			})
		})
		add.ButtonBox.ConnectAccepted(func() {
			name := add.LineEdit.Text()
			prefix := add.LineEdit_2.Text()
			filepath := add.LineEdit_3.Text()

			if len(name) == 0 || len(filepath) == 0 {
				qtui.SimpleMB(add, "Please make sure you fill all fields", "Invalid input").Show()
				return
			}

			if _, ok := modules.Generator.Parameters[prefix]; !ok {
				modules.Generator.Parameters[prefix] = &modules.Parameter{
					Name: name,
					FilePath: filepath,
					Prefix: prefix,
				}
				var params []*modules.Parameter
				if err := viper.UnmarshalKey("generator.Parameters", &params); err != nil {
					log.Fatal(err.Error())
				}
				params = append(params, modules.Generator.Parameters[prefix])
				viper.Set("generator.Parameters", params)
				viper.WriteConfig()
				qtui.ParametersSettingsWindow.GeneratorParams.AddItem(fmt.Sprintf("%s-%s", name, prefix))
			} else {
				qtui.SimpleMB(qtui.ParametersSettingsWindow, "Failed to add parameter", "Error")
			}
		})
		add.Show()
	})

	qtui.ParametersSettingsWindow.GemeratorRemoveParams.ConnectClicked(func(bool) {
		if index := qtui.ParametersSettingsWindow.GeneratorParams.CurrentRow(); index > -1 {
			qtui.ParametersSettingsWindow.GeneratorParams.TakeItem(index)
			if _, ok := modules.Generator.Parameters[strings.Split(qtui.ParametersSettingsWindow.GeneratorParams.CurrentItem().Text(), "-")[1]]; ok {
				delete(modules.Generator.Parameters, qtui.ParametersSettingsWindow.GeneratorParams.CurrentItem().Text())
				var params []*modules.Parameter
				if err := viper.UnmarshalKey("generator.Parameters", &params); err != nil {
					log.Fatal(err.Error())
				}
				var par []*modules.Parameter
				for _, p := range params {
					if p.Prefix == qtui.ParametersSettingsWindow.GeneratorParams.CurrentItem().Text() {
						continue
					}
					par = append(par, p)
				}
				viper.Set("generator.Parameters", par)
			}
		}
	})

	qtui.ParametersSettingsWindow.ParameterbuttonBox.ConnectClicked(func(button *widgets.QAbstractButton) {
		qtui.Main.ComboBoxParameters.Clear()
		var names []string
		for prefix, parameter := range modules.Generator.Parameters {
			names = append(names, fmt.Sprintf("%s-%s", parameter.Name, prefix))
		}
		qtui.Main.ComboBoxParameters.AddItems(names)
	})*/

	return qtui.ParametersSettingsWindow
}
