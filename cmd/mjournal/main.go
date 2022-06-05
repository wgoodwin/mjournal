package main

import (
    "database/sql"
    _ "embed"
    "fmt"
    "path/filepath"
    "os"

    // UI Components
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"

    // SQL Database Pieces
    _ "github.com/mattn/go-sqlite3"
    "github.com/qustavo/dotsql"

)
//go:embed version.txt
var Version string

var drills DrillList
var  dbh *sql.DB

func main() {
    myApp := app.New()

    // TODO This seems to get cut off in the title bar, I should really find out why
    // TODO I did notice the title is a little higher than I would like in i3 (top seems cut off) I should look to see if I can offset that somehow
    myWindow := myApp.NewWindow("Martial Journal - " + Version)
    myWindow.Resize(fyne.NewSize(600,600))

    appPath, err := GetEnvironmentPath()
    if err != nil {
        showClosingErrorDialog(myWindow, err)
    }

    if err = os.MkdirAll(appPath, os.ModePerm); err != nil {
        showClosingErrorDialog(myWindow, err)
    }

    dbh, err = sql.Open("sqlite3", filepath.Join(appPath, "mjournal.db"))
    if err != nil {
        showClosingErrorDialog(myWindow, err)
    }

    // TODO This should likely get pulled from the environment module
    dot, err := dotsql.LoadFromFile("database/v0.0.1.sql")
    if err != nil {
        showClosingErrorDialog(myWindow, err)
    }

    // TODO These should be called  in order by some sort of changeset configuration
    _, err = dot.Exec(dbh, "create-drill-table")
    if err != nil {
        showClosingErrorDialog(myWindow, err)
    }
    
    _, err = dot.Exec(dbh, "create-tag-table")
    if err != nil {
        showClosingErrorDialog(myWindow, err)
    }

    _, err = dot.Exec(dbh, "create-drill-tags")
    if err != nil {
        showClosingErrorDialog(myWindow, err)
    }

    drills, err = queryDrillList(dbh)
    if err != nil {
        showClosingErrorDialog(myWindow, err)
    }

    // TODO  The list needs to go into a container with with a search bar and an add button
    splitContainer := container.NewHSplit(container.NewVBox(), container.NewVBox())
    listContainer := NewSearchListContainer(splitContainer)
    splitContainer.Leading = listContainer

    myWindow.SetContent(splitContainer)
    myWindow.ShowAndRun()
}

// FIXME Right now this is using a hard position without a layout, in the future this should likely have a custom layout
func NewSearchListContainer(parent *container.Split) *fyne.Container {
    list := widget.NewList(
        func() int {
            return len(drills)
        },
        func() fyne.CanvasObject {
            return widget.NewLabel("")
        },
        func(i widget.ListItemID, o fyne.CanvasObject) {
            o.(*widget.Label).SetText(drills[i])
        },
    )

    // TODO This might be better served as a dialog that opens with a save and cancel
    list.OnSelected = func(id widget.ListItemID) {
        l_name := widget.NewLabel(drills[id])
        l_name.TextStyle = fyne.TextStyle{Bold: true}

        descText, err := drills.GetDescription(id, dbh)
        if err != nil { //  TODO This needs a better  handling
            descText = err.Error()
        }
        l_description := widget.NewRichTextFromMarkdown(descText)
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
    list.Move(fyne.NewPos(0, 60))
    list.Resize(fyne.NewSize(390, 600))

    searchBar := widget.NewEntry()
    searchBar.SetPlaceHolder("Search Here")
    searchBar.OnChanged = func(s string) {
        fmt.Println(searchBar.Text)
    }
    searchBar.Move(fyne.NewPos(0, 0))
    searchBar.Resize(fyne.NewSize(300, 40))

    addButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func () {
        drills = append(drills, "Untitled")
        list.Refresh()
        list.Select(len(drills) - 1)
    })
    addButton.Move(fyne.NewPos(350, 0)) // FIXME I would like this button to be the length of the container
    addButton.Resize(fyne.NewSize(40, 40))

    topContainer := container.NewWithoutLayout(searchBar, addButton)
    topContainer.Move(fyne.NewPos(0, 0))
    topContainer.Resize(fyne.NewSize(500, 50))

    return container.NewWithoutLayout(topContainer, list)
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
