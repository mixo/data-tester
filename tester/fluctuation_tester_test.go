package tester

import (
// 	"fmt"
	"github.com/mixo/data-tester/fixture"
	"testing"
	"time"
)

var (
	startDate          = time.Now().AddDate(0, 0, -11)
	endDate            = startDate.AddDate(0, 0, 10)
	rowCountPerDayFrom = 100
	rowCountPerDayTo   = 110
	maxDiff            = 15
	numberDays         = 10
	fluctuations       = []fixture.Fluctuation{
		fixture.Fluctuation{endDate, 1, uint(maxDiff)},
	}
	driver     = "mysql"
	host       = "127.0.0.1"
	port       = "3306"
	user       = "root"
	password   = "root"
	database   = "test"
	tableName  = "datatester_fixture"
	dateColumn = "date"
	columns    = []string{"date", "int_param", "float_param", "group_param"}
    columnsSqls = map[string][]string{
        "mysql": []string{"`date` date", "`int_param` integer", "`float_param` numeric(10, 2)", "`group_param` varchar(255)"},
        "postgres": []string{"\"date\" date", "\"int_param\" integer", "\"float_param\" numeric(10, 2)", "group_param varchar"},
    }

    valueDefinitions = [][]interface{}{
        []interface{}{"int", 300, 330},
        []interface{}{"float", 1000.1, 1100.25},
        []interface{}{"string", "a", "b", "c"},
    }
)

func TestTestYesterdayDataFailed(t *testing.T) {
    test := FunctionalTest{}
    numericColumns := "int_param,float_param"
    groupColumn := "group_param"
    filteredGroups := ""
    dayIndent := 1
    for dbDriver, db := range test.GetDatabases() {
        dataLoader := fixture.DataLoader{tableName, columns, columnsSqls[dbDriver], valueDefinitions}
        dataLoader.Load(db, startDate, endDate, rowCountPerDayFrom, rowCountPerDayTo, fluctuations)
        defer dataLoader.Unload(db)

        var tester DayFluctuationTester
        r := tester.Test(db, tableName, dateColumn, numericColumns, groupColumn, filteredGroups, dayIndent, maxDiff, numberDays)

        r.Show()
        if r.IsOk() {
            t.Error("The test must be failed")
        }
    }
}

// func TestTestYesterdayDataPass(t *testing.T) {
// 	dataLoader := fixture.DataLoader{tableName, columns, columnsSql}
// 	dataLoader.Load(startDate, endDate, rowCountPerDayFrom, rowCountPerDayTo, []fixture.Fluctuation{})
// 	//defer dataLoader.Unload()
//
// 	var tester FluctuationTester
// 	r := tester.TestYesterdayData(driver, host, port, user, password, database, tableName, dateColumn, maxDiff, dayCount)
// 	fmt.Println(r.ToString())
// 	if !r.OK {
// 		t.Error("The test must be passed")
// 	}
// }
