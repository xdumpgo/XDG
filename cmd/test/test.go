package main

import (
"os"

"github.com/therecipe/qt/core"
"github.com/therecipe/qt/qml"
"github.com/therecipe/qt/widgets"
)

func main() {
	widgets.NewQApplication(len(os.Args), os.Args)
	var engine = qml.NewQQmlApplicationEngine(nil)
	engine.Load(core.NewQUrl3("qrc:///qml/main.qml", 0))
	widgets.QApplication_Exec()
}
