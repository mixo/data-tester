package fixture

import (
	"github.com/mixo/godt"
	"github.com/mixo/gosql"
	"time"
)

type DataLoader struct {
	TableName   string
	Columns     []string
	ColumnsSql  []string
	ValueDefinitions [][]interface{}
}

func (this DataLoader) Load(
	startTime, endTime time.Time,
	rowCountPerDayFrom, rowCountPerDayTo int,
	fluctuations []Fluctuation) {
	var (
		rowCount          int
		rowCountGenerator RowCountGenerator
		rowGenerator      RowGenerator
		rows              [][]interface{}
		db gosql.DB = gosql.DB{"mysql", "127.0.0.1", "3306", "root", "", "test"}
		//db gosql.DB = gosql.DB{"postgres", "127.0.0.1", "5432", "test", "test", "test"}
	)

	db.DropTable(this.TableName)
	db.CreateTable(this.TableName, this.ColumnsSql)

	for _, currentTime := range godt.GetPeriod(startTime, endTime) {
		rowCount = rowCountGenerator.generate(rowCountPerDayFrom, rowCountPerDayTo, currentTime, fluctuations)
		rows = rowGenerator.Generate(rowCount, currentTime, this.ValueDefinitions)
		db.InsertMultiple(this.TableName, rows, this.Columns)
	}
}

func (this DataLoader) Unload() {
	db := gosql.DB{"mysql", "127.0.0.1", "3306", "root", "", "test"}
	//db := gosql.DB{"postgres", "127.0.0.1", "5432", "test", "test", "test"}
	db.DropTable(this.TableName)
}
