package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	excelutils "github.com/perolo/excel-utils"
)

type SensorDataType []struct {
	ID          int       `json:"id"`
	Sensor      int       `json:"sensor"`
	TimeStamp   time.Time `json:"timeStamp"`
	Temperature float64   `json:"temperature"`
}

type AllDataType struct {
	ID int `json:"id"`
	//	Sensor      int        `json:"sensor"`
	TimeStamp time.Time  `json:"timeStamp"`
	Temp      [6]float64 `json:"temperature"`
}

func main() { //nolint:funlen

	allData := make(map[int]AllDataType)
	start := 43500
	cont := true
	limit := 50
	ind := 0
	for cont {
		var v SensorDataType
		resp, err := http.Get(fmt.Sprintf("http://192.168.50.152:8081/start/%v", start))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&v)
		if err != nil {
			log.Fatal(err)
		}

		for _, meas := range v {
			if _, ok := allData[ind]; ok {
				r := allData[ind]
				r.Temp[meas.Sensor] = meas.Temperature
				allData[ind] = r
			} else {
				n := AllDataType{}
				n.ID = meas.ID
				n.TimeStamp = meas.TimeStamp
				n.Temp[meas.Sensor] = meas.Temperature
				allData[ind] = n
			}
			if meas.Sensor == 5 {
				ind++
			}
			//allData[meas.TimeStamp] = meas
		}

		//allData = append(allData, v...)
		if len(v) < limit {
			cont = false
		}
		start = start + limit
	}

	excelutils.NewFile()
	/*
		excelutils.SetCellFontHeader()
		excelutils.WiteCellln("Introduction")
		excelutils.WiteCellln("Please Do not edit this page!")
		t := time.Now()
		excelutils.WiteCellln("Created by: " + "ConfUser" + " : " + t.Format(time.RFC3339))
		excelutils.WiteCellln("")

		excelutils.SetCellFontHeader2()
		excelutils.WiteCellln("Sensor and Temp")
		excelutils.NextLine()
	*/
	excelutils.AutoFilterStart()
	excelutils.SetTableHeader()
	excelutils.WiteCell("Time")
	excelutils.SetTableHeader()
	excelutils.NextCol()
	excelutils.WiteCell("ID")
	excelutils.SetTableHeader()
	excelutils.NextCol()

	excelutils.WiteCell("Varm V Ut")
	excelutils.SetTableHeader()
	excelutils.NextCol()
	excelutils.WiteCell("Värme Framledn")
	excelutils.SetTableHeader()
	excelutils.NextCol()
	excelutils.WiteCell("Ute Temp")
	excelutils.SetTableHeader()
	excelutils.NextCol()
	excelutils.WiteCell("KB In")
	excelutils.SetTableHeader()
	excelutils.NextCol()
	excelutils.WiteCell("Värme Retur")
	excelutils.SetTableHeader()
	excelutils.NextCol()
	excelutils.WiteCell("KB Ut")
	excelutils.SetTableHeader()
	excelutils.NextCol()

	/*
		for i := 1; i < 7; i++ {
			excelutils.SetTableHeader()
			//		excelutils.WiteCellnc(allData[i].Sensor)
			excelutils.WiteCellnc(i)
		}*/
	excelutils.NextLine()

	keys := make([]int, 0, len(allData))
	for k := range allData {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		for i := 0; i < 6; i++ {
			if allData[k].Temp[i] == 0 {
				// attempt to repair
				if (allData[k-1].Temp[i] != 0) && (allData[k+1].Temp[i] != 0) {
					fmt.Printf("Try Repair\n")
					t := allData[k]
					t.Temp[i] = (allData[k-1].Temp[i] + allData[k+1].Temp[i]) / 2
					allData[k] = t
				} else {
					fmt.Printf("Faileed Repair\n")
				}
			}
		}

	}
	//	max := len(allData)
	//	key := 0
	for _, k := range keys {
		//	for key := range allData {
		excelutils.WiteCellnc(allData[k].TimeStamp)
		excelutils.WiteCellnc(allData[k].ID)
		//excelutils.WiteCellnc()
		for i := 0; i < 6; i++ {
			excelutils.WiteCellnc(allData[k].Temp[i])
		}
		//key = key + 6
		excelutils.NextLine()
	}

	excelutils.SetAutoColWidth()
	excelutils.AutoFilterEnd()

	//excelutils.SetColWidth("A", "A", 60)
	// Save xlsx file by the given path.
	excelutils.SaveAs("C:\\Temp\\Temp.xlsx")

}
