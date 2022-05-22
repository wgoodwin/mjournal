package main

import (
    _ "embed"
    "fmt"
    "os"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"
)
//go:embed version.txt
var Version string

type dataType struct  {
    Name    string
    Value   string
}
var data = []dataType{
    dataType{"Test1", "This is test1."},
    dataType{"Test2", "This is test2."},
    dataType{"Test3", "This is test3 with an absurdly long string......sdjhgflkjasgdflkjadshgkdashgkjdsahgkljdshgjksdahgkjdlashgkdjashgaklsdjhg."},
}


func main() {
    myApp := app.New()

    // TODO This seems to get cut off in the title bar, I should really find out why
    myWindow := myApp.NewWindow("Martial Journal - " + Version)
    myWindow.Resize(fyne.NewSize(600,600))

    appPath, err := GetEnvironmentPath()
    if err != nil {
        showClosingErrorDialog(myWindow, err)
    }

    if err = os.MkdirAll(appPath, os.ModePerm); err != nil {
        showClosingErrorDialog(myWindow, err)
    }


    // TODO  The list needs to go into a container with with a search bar and an add button
    splitContainer := container.NewHSplit(container.NewVBox(), container.NewVBox())
    listContainer := NewSearchListContainer(splitContainer)
    splitContainer.Leading = listContainer

    myWindow.SetContent(splitContainer)
    myWindow.ShowAndRun()
}

func NewSearchListContainer(parent *container.Split) *fyne.Container {
    searchBar := widget.NewEntry()
    searchBar.SetPlaceHolder("Search Here")
    searchBar.OnChanged = func(s string) {
        fmt.Println(searchBar.Text)
    }
    addButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func () {
        fmt.Println("Add Element")
    })

    topContainer := container.NewGridWithColumns(2, searchBar, addButton)
    list := widget.NewList(
        func() int {
            return len(data)
        },
        func() fyne.CanvasObject {
            return widget.NewLabel("")
        },
        func(i widget.ListItemID, o fyne.CanvasObject) {
            o.(*widget.Label).SetText(data[i].Name)
        },
    )
    list.OnSelected = func(id widget.ListItemID) {
        l_name := widget.NewLabel(data[id].Name)
        l_name.TextStyle = fyne.TextStyle{Bold: true}
        l_description := widget.NewRichTextFromMarkdown(data[id].Value)
        // TODO This leads into some UI ideas, perhaps having an edit/remove button at the bottom which converts to a save/cancel on edit
        e_button := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func () {
            // TODO This should swap cell with an editable one
        })
        parent.Trailing = container.NewVBox(
            container.NewHBox(l_name, e_button),
            container.NewHBox(l_description, e_button),
        )
        parent.Refresh()
    }
    return container.NewVBox(topContainer, list)
}

func showClosingErrorDialog(parent fyne.Window, err error) {
    errDiag := dialog.NewError(err, parent)
    errDiag.SetOnClosed(
        func() {
            os.Exit(1)
        },
    )
    errDiag.Show()
}
