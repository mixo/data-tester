package fixture

import (
	"github.com/mixo/godt"
	"github.com/mixo/gosql"
	"time"
	"github.com/urfave/cli/v2"
)

const (
    timeLayout = "2006-01-02"
)

type DataLoader struct {
	TableName   string
	Columns     []string
	ColumnsSql  []string
	ValueDefinitions [][]interface{}
}

func (this DataLoader) GetCliCommand() *cli.Command {
	return &cli.Command{
		Name: "load",
        Aliases: []string{"ld"},
		Description: "Performs the fixtures loading",
		Action: func(c *cli.Context) error {
			tableName := c.String("table-name")
			startDateString := c.String("start-date")
			endDateString := c.String("end-date")
			rowCountPerDayFrom := c.Int("row-count-per-day-from")
			rowCountPerDayTo := c.Int("row-count-per-day-to")
			maxDiff := c.Uint("max-diff")
			if tableName == "" || startDateString == "" || endDateString == "" {// || groupColumn == ""  || dayIndent == 0 || maxDiff == 0 || numberDays == 0 {
				cli.ShowCommandHelp(c, "load")
				return cli.Exit("You must specify the command flags", 1)
			}

			startDate, err := time.Parse(timeLayout, startDateString)
			if err != nil {
			    panic(err)
			}
			endDate, err := time.Parse(timeLayout, endDateString)
			if err != nil {
			    panic(err)
			}

            fluctuations       := []Fluctuation{
                Fluctuation{endDate, 1, maxDiff},
            }

            columns     := []string{"date", "int_param", "float_param", "group_param"}
            columnsSqls := map[string][]string{
                "mysql": []string{"`date` date", "`int_param` integer", "`float_param` numeric(10, 2)", "`group_param` varchar(255)"},
                "postgres": []string{"\"date\" date", "\"int_param\" integer", "\"float_param\" numeric(10, 2)", "group_param varchar"},
            }

            valueDefinitions := [][]interface{}{
                []interface{}{"int", 300, 330},
                []interface{}{"float", 1000.1, 1100.25},
                []interface{}{"string", "a", "b", "c"},
            }

            db := c.App.Metadata["db"].(gosql.DB)
            dataLoader := DataLoader{tableName, columns, columnsSqls[db.Driver], valueDefinitions}
            dataLoader.Load(db, startDate, endDate, rowCountPerDayFrom, rowCountPerDayTo, fluctuations)

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "table-name",
				Aliases: []string{"tn"},
				Usage: "Specify the table name you want to load fixtures",
				Value: "datatester_fixture",
			},
			&cli.StringFlag{
				Name:  "start-date",
				Aliases: []string{"sd"},
				Usage: "Specify the start date",
				Value: time.Now().AddDate(0, 0, -11).Format(timeLayout),
			},
			&cli.StringFlag{
				Name:  "end-date",
				Aliases: []string{"ed"},
				Usage: "Specify the end date",
				Value: time.Now().AddDate(0, 0, -1).Format(timeLayout),
			},
			&cli.IntFlag{
				Name: "row-count-per-day-from",
				Aliases: []string{"cf"},
				Usage: "Specify the minimum number of rows which should be loaded. The actual number of rows will be a random number that is between the minimum and the maximum",
				Value: 100,
			},
			&cli.IntFlag{
				Name: "row-count-per-day-to",
				Aliases: []string{"ct"},
				Usage: "Specify the maximum number of rows which should be loaded. The actual number of rows will be a random number that is between the minimum and the maximum",
				Value: 110,
			},
			&cli.UintFlag{
				Name: "max-diff",
				Aliases: []string{"md"},
				Usage: "Specify the maximum difference in percents between the yesterday row count and the average row count",
				Value: uint(15),
			},
		},
	}
}

func (this DataLoader) Load(db gosql.DB,
	startTime, endTime time.Time,
	rowCountPerDayFrom, rowCountPerDayTo int,
	fluctuations []Fluctuation) {
	var (
		rowCount          int
		rowCountGenerator RowCountGenerator
		rowGenerator      RowGenerator
		rows              [][]interface{}
	)

	db.DropTable(this.TableName)
	db.CreateTable(this.TableName, this.ColumnsSql)

	for _, currentTime := range godt.GetPeriod(startTime, endTime) {
		rowCount = rowCountGenerator.generate(rowCountPerDayFrom, rowCountPerDayTo, currentTime, fluctuations)
		rows = rowGenerator.Generate(rowCount, currentTime, this.ValueDefinitions)
		db.InsertMultiple(this.TableName, rows, this.Columns)
	}
}

func (this DataLoader) Unload(db gosql.DB) {
// 	db := gosql.DB{"mysql", "127.0.0.1", "3306", "root", "", "test"}
	//db := gosql.DB{"postgres", "127.0.0.1", "5432", "test", "test", "test"}
	db.DropTable(this.TableName)
}
