package tester

import (
	"fmt"
	"github.com/mixo/gosql"
	"github.com/urfave/cli/v2"
	"math"
	"strings"
	"time"
)

type DayFluctuationTester struct{}

func (this DayFluctuationTester) GetCliCommand() *cli.Command {
	return &cli.Command{
		Name:        "day-fluctuation",
		Aliases:     []string{"df"},
		Description: "Checks that a day number of rows doesn't differ from an average number of rows per day",
		Action: func(c *cli.Context) error {
			tableName := c.String("table-name")
			dateColumn := c.String("date-column")
			numericColumns := c.String("numeric-columns")
			groupColumn := c.String("group-column")
			filteredGroups := c.String("filtered-groups")
			dayIndent := c.Int("day-indent")
			maxDiff := c.Int("max-diff")
			numberDays := c.Int("number-days")
			if tableName == "" || dateColumn == "" || groupColumn == "" || dayIndent == 0 || maxDiff == 0 || numberDays == 0 {
				cli.ShowCommandHelp(c, "day-fluctuation")
				return cli.Exit("You must specify the command flags", 1)
			}

			r := (DayFluctuationTester{}).Test(c.App.Metadata["db"].(gosql.DB), tableName, dateColumn, numericColumns, groupColumn, filteredGroups, dayIndent, maxDiff, numberDays)

			r.Show()
			if !r.IsOk() {
				return cli.Exit("", 1)
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "table-name",
				Aliases: []string{"tn"},
				Usage:   "Specify the table name you want to test",
			},
			&cli.StringFlag{
				Name:    "date-column",
				Aliases: []string{"dc"},
				Usage:   "Specify the date column that will be used to count rows",
			},
			&cli.StringFlag{
				Name:     "numeric-columns",
				Aliases:  []string{"nc"},
				Usage:    "Specify comma-separated numeric columns which should be checked. For example, amount,price,profit",
				Required: false,
			},
			&cli.StringFlag{
				Name:    "group-column",
				Aliases: []string{"gc"},
				Usage:   "Specify the group column that will be used to count and sum columns within groups",
			},
			&cli.StringFlag{
				Name:    "filtered-groups",
				Aliases: []string{"fg"},
				Usage:   "Specify comma-separated groups that should be checked",
			},
			&cli.IntFlag{
				Name:    "day-indent",
				Aliases: []string{"di"},
				Usage:   "For which day, check a row count. For example, 1 - tomorrow, 0 - today, -1 - yesterday, -2 the day before yesterday and so on",
			},
			&cli.IntFlag{
				Name:    "max-diff",
				Aliases: []string{"md"},
				Usage:   "Specify the maximum difference in percents between the yesterday row count and the average row count",
			},
			&cli.IntFlag{
				Name:    "number-days",
				Aliases: []string{"nd"},
				Usage:   "Specify the number of days for which to calculate the average row count",
			},
		},
	}
}

func (this DayFluctuationTester) Test(db gosql.DB,
	tableName, dateColumn, numericColumnsList, groupColumn, filteredGroupsString string,
	dayIndent, maxDiff, numberDays int) TestResult {

	quantityColumn := "_quantity"
	day := time.Now().AddDate(0, 0, dayIndent)
	startDate := day.AddDate(0, 0, -numberDays)
	endDate := day.AddDate(0, 0, -1)
	var numericColumns []string
	if numericColumnsList != "" {
		numericColumns = strings.Split(numericColumnsList, ",")
	}
	filteredGroupStrings := strings.Split(filteredGroupsString, ",")
	filteredGroups := make([]interface{}, 0)
	for _, filteredGroupString := range filteredGroupStrings {
		if filteredGroupString != "" {
			filteredGroups = append(filteredGroups, filteredGroupString)
		}
	}

	avgParams := db.GetAvgRowParamsPerDay(tableName, dateColumn, startDate, endDate, numberDays, quantityColumn, numericColumns, groupColumn, filteredGroups)
	dayParams := db.GetRowParamsOnDate(tableName, dateColumn, day, quantityColumn, numericColumns, groupColumn, filteredGroups)
	groupsAvgParams := this.groupMaps(avgParams, groupColumn)
	groupsDayParams := this.groupMaps(dayParams, groupColumn)

	columns := append(numericColumns, quantityColumn)
	groupsDiffs, ok := this.compareGroupsParams(columns, maxDiff, groupsAvgParams, groupsDayParams)

	return DayFluctuationTesterResult{ok, groupsAvgParams, groupsDayParams, groupsDiffs, maxDiff, startDate, endDate, day, quantityColumn, groupColumn}
}

func (this DayFluctuationTester) groupMaps(maps []map[string]interface{},
	groupColumn string) map[string]map[string]interface{} {
	var group string

	groupMaps := make(map[string]map[string]interface{}, 0)
	for _, singleMap := range maps {
		switch singleMap[groupColumn].(type) {
		case int:
			group = fmt.Sprintf("%d", singleMap[groupColumn].(int))
		case string:
			group = singleMap[groupColumn].(string)
		default:
			panic(fmt.Sprintf("Undefined type %T", singleMap[groupColumn]))
		}
		groupMaps[group] = singleMap
	}

	return groupMaps
}

func (this DayFluctuationTester) compareGroupsParams(columns []string, maxDiff int,
	groupsAvgParams, groupsDayParams map[string]map[string]interface{}) (map[string]map[string]float64, bool) {
	ok := true
	groupsDiffs := make(map[string]map[string]float64, 0)
	for group, groupAvgParams := range groupsAvgParams {
		groupsDiffs[group] = make(map[string]float64, 0)
		for _, param := range columns {
			var dayValue float64
			if groupsDayParams[group] == nil {
				groupsDayParams[group] = make(map[string]interface{}, 0)
			}
			if groupsDayParams[group][param] == nil {
				groupsDayParams[group][param] = 0.0
			}

			switch groupsDayParams[group][param].(type) {
			case int:
				dayValue = float64(groupsDayParams[group][param].(int))
			case int8:
				dayValue = float64(groupsDayParams[group][param].(int8))
			case int16:
				dayValue = float64(groupsDayParams[group][param].(int16))
			case int32:
				dayValue = float64(groupsDayParams[group][param].(int32))
			case int64:
				dayValue = float64(groupsDayParams[group][param].(int64))
			default:
				dayValue = groupsDayParams[group][param].(float64)
			}

			diff := this.calculateDiffInPercents(dayValue, groupAvgParams[param].(float64))
			groupsDiffs[group][param] = diff

			ok = ok && int(math.Abs(diff)) <= maxDiff
		}
	}

	return groupsDiffs, ok
}

func (this DayFluctuationTester) calculateDiffInPercents(dayRowCount, avgRowCountPerDay float64) float64 {
	return math.Round((dayRowCount * 100 / avgRowCountPerDay) - 100)
}
