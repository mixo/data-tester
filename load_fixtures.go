package main

import (
	"fmt"
	"github.com/mixo/data-tester/fixture"
	"time"
)

var (
	startDate          = time.Now().AddDate(0, 0, -11)
	endDate            = startDate.AddDate(0, 0, 10)
	rowCountPerDayFrom = 100
	rowCountPerDayTo   = 110
	maxDiff            = uint(15)
	fluctuations       = []fixture.Fluctuation{
		fixture.Fluctuation{endDate, 1, maxDiff},
	}
	tableName   = "datatester_fixture"
	columns     = []string{"date", "int_param", "float_param", "group_param"}
	columnsSql  = []string{"`date` date", "`int_param` integer", "`float_param` numeric(10, 2)", "`group_param` varchar(255)"}
// 	columnsSql  = []string{"\"date\" date", "\"int_param\" integer", "\"float_param\" numeric(10, 2)", "group_param varchar"}
	valueDefinitions = [][]interface{}{
	    []interface{}{"int", 300, 330},
	    []interface{}{"float", 1000.1, 1100.25},
	    []interface{}{"string", "a", "b", "c"},
    }
)

func main() {
	fmt.Println("Load fixtures")
	dataLoader := fixture.DataLoader{tableName, columns, columnsSql, valueDefinitions}
	dataLoader.Load(startDate, endDate, rowCountPerDayFrom, rowCountPerDayTo, fluctuations)
}
