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

	"github.com/perolo/excellogger"
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
	Mode        string `properties:"mode"`
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

	start := 0
	if cfg.Mode == "last24" {
		start, _ = getSensorDataLast()
		start = start - 17280 //TODO Calculate from timestamps and number of sensors

	} else {
		start = cfg.Start
	}

	if cfg.Template != "" {
		err := excellogger.OpenFile(cfg.Template)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		opt := excellogger.Options{SheetName: "Sensor Data"}
		excellogger.NewFile(&opt)
	}
	err := excellogger.SelectSheet("Info")
	if err != nil {
		err = excellogger.NewSheet("Info")
		if err != nil {
			log.Fatal(err)
		}
	}
	excellogger.SetCellFontHeader()
	excellogger.WiteCellln("Introduction")
	excellogger.WiteCellln("Please Do not edit this page!")
	t := time.Now()
	application := os.Args[0]

	excellogger.WiteCellln("Created by: " + application + " : " + t.Format(time.RFC3339))
	fmt.Printf("Application : %s started\n", application)
	excellogger.WiteCellln("")
	if len(os.Args) > 1 {
		for _, arg := range os.Args {
			excellogger.WiteCellln("Arg: " + arg)
			fmt.Printf("Arg : %s \n", arg)
		}
	}
	excellogger.WiteCellln("")
	excellogger.SetCellFontHeader2()
	excellogger.WiteCellln("Properties")

	v := reflect.ValueOf(cfg)

	//		values := make([]interface{}, v.NumField())
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		str := fmt.Sprintf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
		excellogger.WiteCellln(str)
		fmt.Printf(str)
	}
	excellogger.SetAutoColWidth()
	//excellogger.SetCellFontHeader2()
	//excellogger.WiteCellln("Sensor and Temp")
	//excellogger.NextLine()
	err = excellogger.SelectSheet("Sensor Data")
	if err != nil {
		log.Fatal(err)
	}

	excellogger.AutoFilterStart()
	excellogger.SetTableHeader()
	excellogger.WiteCell("Time")
	excellogger.SetTableHeader()
	excellogger.NextCol()
	excellogger.WiteCell("ID")
	excellogger.SetTableHeader()
	excellogger.NextCol()

	for _, name := range sensors {
		excellogger.WiteCell(name)
		excellogger.SetTableHeader()
		excellogger.NextCol()
	}
	excellogger.NextLine()

	cont := true
	limit := 50
	//	ind := 0
	skipfirst := true
	lastmeas := 0
	row := 2
	for cont {
		v := getSensorDataStart(start)

		for _, meas := range v {
			// skip until first sensor found
			if skipfirst {
				if meas.Sensor == 0 {
					skipfirst = false
					lastmeas = -1
				}
			}
			if !skipfirst {
				if lastmeas > meas.Sensor {
					//					excellogger.NextLine()
					lastmeas = -1
					row++
				}
				if lastmeas < meas.Sensor {
					//					excellogger.SetCell(meas.TimeStamp.Format("2006-01-02 15:04:05"), 1, row)
					excellogger.SetCell(meas.TimeStamp, 1, row)
					excellogger.SetCell(meas.ID, 2, row)
					lastTimestamp = meas.TimeStamp
					lastmeas = -1
				}
				if lastmeas < meas.Sensor {
					//					strr := strconv.FormatFloat(meas.Temperature, 'f', 2, 64)
					excellogger.SetCell(meas.Temperature, meas.Sensor+3, row)
					lastmeas = meas.Sensor
				} else {
					fmt.Printf("Que %s row: %v Sensor: %v\n", meas.TimeStamp, row, meas.Sensor)
				}
			}
		}

		if len(v) < limit {

			cont = false
		}
		start = start + limit
	}

	excellogger.AutoFilterEnd()
	excellogger.SetAutoColWidth()

	// Save xlsx file by the given path.
	timestr := lastTimestamp.Format("2006-01-02_15-04-05")
	name := fmt.Sprintf(cfg.File, timestr)
	excellogger.SaveAs(name)

}

func getSensorDataLast() (int, int) {
	var v SensorDataType
	resp, err := http.Get(fmt.Sprintf(cfg.TempServer + "/last"))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		log.Fatal(err)
	}
	lastid := v[len(v)-1].ID
	return lastid, 1000
}

func getSensorDataStart(start int) SensorDataType {
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
	return v
}
