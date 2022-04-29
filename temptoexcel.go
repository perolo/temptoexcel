package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/magiconair/properties"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	excelutils "github.com/perolo/excel-utils"
)

type SensorDataType []struct {
	ID          int       `json:"id"`
	Sensor      int       `json:"sensor"`
	TimeStamp   time.Time `json:"timeStamp"`
	Temperature float64   `json:"temperature"`
}

var cfg Config

type Config struct {
	Start       int    `properties:"start"`
	File        string `properties:"file"`
	Template    string `properties:"template"`
	TempServer  string `properties:"tempserver"`
	SensorNames string `properties:"sensornames"`
}

func main() { //nolint:funlen
	var lastTimestamp time.Time
	propPtr := flag.String("prop", "temptoexcel.properties", "a string")
	flag.Parse()

	p := properties.MustLoadFile(*propPtr, properties.ISO_8859_1)

	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	sensors := strings.Split(cfg.SensorNames, ",")

	if cfg.Template != "" {
		err := excelutils.OpenFile(cfg.Template)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		opt := excelutils.Options{SheetName: "Sensor Data"}
		excelutils.NewFile(&opt)
	}
	err := excelutils.SelectSheet("Info")
	if err != nil {
		err = excelutils.NewSheet("Info")
		if err != nil {
			log.Fatal(err)
		}
	}
	excelutils.SetCellFontHeader()
	excelutils.WiteCellln("Introduction")
	excelutils.WiteCellln("Please Do not edit this page!")
	t := time.Now()
	application := os.Args[0]

	excelutils.WiteCellln("Created by: " + application + " : " + t.Format(time.RFC3339))
	excelutils.WiteCellln("")
	if len(os.Args) > 1 {
		for _, arg := range os.Args {
			excelutils.WiteCellln("Arg: " + arg)
		}
	}
	excelutils.WiteCellln("")
	excelutils.SetCellFontHeader2()
	excelutils.WiteCellln("Properties")

	v := reflect.ValueOf(cfg)

	//		values := make([]interface{}, v.NumField())
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		str := fmt.Sprintf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
		excelutils.WiteCellln(str)

	}
	excelutils.SetAutoColWidth()
	//excelutils.SetCellFontHeader2()
	//excelutils.WiteCellln("Sensor and Temp")
	//excelutils.NextLine()
	err = excelutils.SelectSheet("Sensor Data")
	if err != nil {
		log.Fatal(err)
	}

	excelutils.AutoFilterStart()
	excelutils.SetTableHeader()
	excelutils.WiteCell("Time")
	excelutils.SetTableHeader()
	excelutils.NextCol()
	excelutils.WiteCell("ID")
	excelutils.SetTableHeader()
	excelutils.NextCol()

	for _, name := range sensors {
		excelutils.WiteCell(name)
		excelutils.SetTableHeader()
		excelutils.NextCol()
	}
	excelutils.NextLine()

	start := cfg.Start
	cont := true
	limit := 50
	//	ind := 0
	skipfirst := true
	colcount := 0
	for cont {
		var v SensorDataType
		resp, err := http.Get(fmt.Sprintf(cfg.TempServer+"/start/%v", start))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&v)
		if err != nil {
			log.Fatal(err)
		}

		for _, meas := range v {
			// skip until first sensor found
			if skipfirst {
				if meas.Sensor == 0 {
					skipfirst = false
				}
			}
			if !skipfirst {
				if meas.Sensor == 0 || colcount == 0 {
					excelutils.WiteCellnc(meas.TimeStamp)
					excelutils.WiteCellnc(meas.ID)
					for colcount < meas.Sensor {
						excelutils.WiteCellnc("")
						colcount++
					}
				}
				if colcount == meas.Sensor {
					excelutils.WiteCellnc(meas.Temperature)
					colcount++
				} else {
					fmt.Printf("Que %s colcount: %v Sensor: %v\n", meas.TimeStamp, colcount, meas.Sensor)
					colcount = meas.Sensor
				}
				if meas.ID == 44190 {
					fmt.Printf("Que %s colcount: %v Sensor: %v\n", meas.TimeStamp, colcount, meas.Sensor)

				}
				if meas.Sensor == 5 {
					excelutils.NextLine()
					colcount = 0
					lastTimestamp = meas.TimeStamp
				}
			}
		}

		if len(v) < limit {

			cont = false
		}
		start = start + limit
	}

	excelutils.AutoFilterEnd()
	excelutils.SetAutoColWidth()

	// Save xlsx file by the given path.
	timestr := lastTimestamp.Format("2006-01-02_15-04-05")
	name := fmt.Sprintf(cfg.File, timestr)
	excelutils.SaveAs(name)

}
