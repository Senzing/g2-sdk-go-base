//go:build windows

package szengine

var expectedExportCsvEntityReport = []string{
	`RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL_CODE,MATCH_KEY,DATA_SOURCE,RECORD_ID`,
	`10003,0,"","","CUSTOMERS","1001"`,
	`10003,0,"RESOLVED","+NAME+DOB+PHONE","CUSTOMERS","1002"`,
	`10003,0,"RESOLVED","+NAME+DOB+EMAIL","CUSTOMERS","1003"`,
}

var expectedExportCsvEntityReportIterator = []string{
	`RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL_CODE,MATCH_KEY,DATA_SOURCE,RECORD_ID`,
	`100013,0,"","","CUSTOMERS","1001"`,
	`100013,0,"RESOLVED","+NAME+DOB+PHONE","CUSTOMERS","1002"`,
	`100013,0,"RESOLVED","+NAME+DOB+EMAIL","CUSTOMERS","1003"`,
}

var expectedExportCsvEntityReportIteratorNilCsvColumnList = []string{
	`RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL_CODE,MATCH_KEY,DATA_SOURCE,RECORD_ID`,
	`100019,0,"","","CUSTOMERS","1001"`,
	`100019,0,"RESOLVED","+NAME+DOB+PHONE","CUSTOMERS","1002"`,
	`100019,0,"RESOLVED","+NAME+DOB+EMAIL","CUSTOMERS","1003"`,
}
