package main

import (
    "sort"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
)

func TestQueryDrillList(t *testing.T) {
    db, mock, _ := sqlmock.New()
    defer db.Close()

    mock.ExpectQuery("SELECT `name` FROM `drill`").
        WillReturnRows(sqlmock.NewRows([]string{"name"}).
        AddRow("ADrill").
        AddRow("BDrill").
        AddRow("DDrill").
        AddRow("CDrill"))

    dList, err := queryDrillList(db)
    if assert.NoError(t, err) {
        // Expecting the list to be ordered
        assert.Equal(t, "ADrill", dList[0])
        assert.Equal(t, "BDrill", dList[1])
        assert.Equal(t, "CDrill", dList[2])
        assert.Equal(t, "DDrill", dList[3])
    }
}

func TestDrillListGetDescription(t *testing.T) {
    db, mock, _ := sqlmock.New()
    defer db.Close()

    mock.ExpectQuery("SELECT `description` FROM `drill` WHERE (.*)").
        WillReturnRows(sqlmock.NewRows([]string{"description"}).
        AddRow("DescriptionA"))

    var dList DrillList
    dList = append(dList, "ADrill")
    result,  err := dList.GetDescription(0, db)
    if assert.NoError(t, err) {
        assert.Equal(t, "DescriptionA", result)
    }
}

func TestDrillListSearch(t *testing.T) {
    var dList DrillList
    dList = []string{
        "Foo",
        "Bar",
        "Boo",
        "Baz",
        "Woo",
        "Wow",
    }
    sort.Strings(dList)

    result := dList.Search("oo")
    assert.Equal(t, 3, len(result))
    assert.Equal(t, "Boo", result[0])
    assert.Equal(t, "Foo", result[1])
    assert.Equal(t, "Woo", result[2])

    result = dList.Search("o")
    assert.Equal(t, 4, len(result))
    assert.Equal(t, "Boo", result[0])
    assert.Equal(t, "Foo", result[1])
    assert.Equal(t, "Woo", result[2])
    assert.Equal(t, "Wow", result[3])

    result = dList.Search("Ba")
    assert.Equal(t, 2, len(result))
    assert.Equal(t, "Bar", result[0])
    assert.Equal(t, "Baz", result[1])
}
