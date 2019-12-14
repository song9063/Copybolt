package main

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var strSource string = ""
var strDest string = ""

func main(){
	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(500,350)
	window.SetWindowTitle("Busang Inc.")

	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())
	window.SetCentralWidget(widget)

	setting := core.NewQSettings4(".boltcopyconf", core.QSettings__IniFormat, widget )


	titleLabel := widgets.NewQLabel(widget, 0)
	titleLabel.SetText("GO Project Duplicator")
	titleLabel.SetStyleSheet("QLabel {font-weight:bold;}")
	titleLabel.SetMinimumSize2(500,10)
	widget.Layout().AddWidget(titleLabel)


	// 원본 코드
	inputPath := widgets.NewQLabel(widget, 0)
	inputPath.SetText("-")
	widget.Layout().AddWidget(inputPath)
	btSetInPath := widgets.NewQPushButton2("Project Root Path", widget)
	btSetInPath.ConnectClicked(func(bool){
		strSource = widgets.QFileDialog_GetExistingDirectory(widget, "Project Root", "~", widgets.QFileDialog__ShowDirsOnly)
		if len(strSource) > 0{
			inputPath.SetText(strSource)
			setting.SetValue("src",core.NewQVariant12(strSource))
		}

	})
	widget.Layout().AddWidget(btSetInPath)

	// 패키지명
	packNameWrap := widgets.NewQWidget(widget, 0)
	packNameWrap.SetLayout(widgets.NewQHBoxLayout())

	lblPackName := widgets.NewQLabel(widget, 0)
	lblPackName.SetText("busangweb.com/")
	packNameWrap.Layout().AddWidget(lblPackName)

	packName := widgets.NewQLineEdit(widget)
	packName.SetPlaceholderText("packagename")
	packNameWrap.Layout().AddWidget(packName)
	widget.Layout().AddWidget(packNameWrap)


	// 출력 경로
	outputPath := widgets.NewQLabel(widget, 0)
	outputPath.SetText("-")
	widget.Layout().AddWidget(outputPath)
	btSetPath := widgets.NewQPushButton2("Output Path", widget)
	btSetPath.ConnectClicked(func(bool){
		strDest = widgets.QFileDialog_GetExistingDirectory(widget, "Output Path", "~", widgets.QFileDialog__ShowDirsOnly)
		if len(strDest) > 0{
			outputPath.SetText(strDest)
			setting.SetValue("dest",core.NewQVariant12(strDest))
		}

	})
	widget.Layout().AddWidget(btSetPath)

	// 복사버튼
	btCopy := widgets.NewQPushButton2("Copy", widget)
	btCopy.ConnectClicked(func(bool){
		strPackageName := packName.Text()

		if len(strPackageName) < 1 {
			msg := widgets.NewQMessageBox(widget)
			msgResp := msg.Information(widget, "Hello", "Input your package name.", widgets.QMessageBox__Ok, widgets.QMessageBox__NoButton)
			if msgResp == widgets.QMessageBox__No {
				fmt.Println("Yes")
			}
			return
		}

		if len(strSource) < 1 || len(strDest) < 1 {
			msg := widgets.NewQMessageBox(widget)
			msgResp := msg.Information(widget, "Hello", "Select path.", widgets.QMessageBox__Ok, widgets.QMessageBox__NoButton)
			if msgResp == widgets.QMessageBox__No {
				fmt.Println("Yes")
			}
			return
		}
		btCopy.SetEnabled(false)

		fmt.Print("Source:")
		fmt.Println(strSource)
		fmt.Print("Dest:")
		fmt.Println(strDest + "/" + strPackageName)


		// 파일 복사 시작
		strFilesTo := strDest + "/" + strPackageName

		if _, err := os.Stat(strFilesTo); os.IsNotExist(err) {
			os.Mkdir(strFilesTo, os.ModePerm)
			err = copyDir(strSource, strFilesTo)
		}

		// 패키지 수정 시작
		err := filepath.Walk(strFilesTo, func(path string, info os.FileInfo, ferr error) error {
			if ferr != nil {
				return ferr
			}

			name := info.Name()
			if strings.HasSuffix(name, ".go") {
				fmt.Println("GOGOGO " + name)

				if cerr := changePackageName(path, strPackageName); cerr != nil {
					fmt.Println("Error while change package name")
					fmt.Println(cerr)
				}

			}


			return nil
		})
		if err != nil {
			fmt.Println(err)
		}


		btCopy.SetEnabled(true)

		msg := widgets.NewQMessageBox(widget)
		msgResp := msg.Information(widget, "Hello", "Success!!", widgets.QMessageBox__Ok, widgets.QMessageBox__NoButton)
		if msgResp == widgets.QMessageBox__No {
			fmt.Println("Yes")
		}
	})
	widget.Layout().AddWidget(btCopy)


	strOldSrc := setting.Value("src", core.NewQVariant12("")).ToString()
	strOldDest := setting.Value("dest", core.NewQVariant12("")).ToString()
	strSource = strOldSrc
	strDest = strOldDest
	inputPath.SetText( strOldSrc)
	outputPath.SetText(strOldDest)
	fmt.Println(strOldSrc)




	window.Show()
	app.Exec()
}
func changePackageName(fileName, newPackageName string ) error {
	input, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		if strings.Contains(line, "busangweb.com/goboltwebbase/") {
			lines[i] = strings.ReplaceAll(line, "busangweb.com/goboltwebbase/", "busangweb.com/" + newPackageName + "/")
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(fileName, []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
func copyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}
func copyDir(src, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}

	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}else{
			if err = copyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}