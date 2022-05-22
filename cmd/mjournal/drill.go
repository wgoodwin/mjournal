package main

import (
    "database/sql"
    "sort"

    "github.com/lithammer/fuzzysearch/fuzzy"
)

// TODO I should unit test all of these pieces

type DrillList []string

func (dl DrillList) GetDescription(id int, dbh *sql.DB) (string, error) {
    var result string

    err := dbh.QueryRow("SELECT `description` FROM `drill` WHERE `name` = ?", dl[id]).Scan(&result)
    return result, err
}

func (dl DrillList) Search(search string) DrillList {
    return fuzzy.Find(search, dl)
}

func queryDrillList(dbh  *sql.DB) (DrillList, error) {
    var result DrillList

    rows, err := dbh.Query("SELECT `name` FROM `drill`")
    if err != nil {
        return result, err
    }
    defer rows.Close()

    for rows.Next() {
        var name string
        if err = rows.Scan(&name); err != nil  {
            // TODO We should likely do some sort of warning, here but I don't think we should stop processing
            continue
        }
        result  = append(result, name)
    }

    sort.Strings(result)
    return result, nil
}


