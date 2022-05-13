# temptoexcel
Simple application that reads data from tempserver: https://github.com/perolo/tempserver and creates an ExcelSheet with data and graphs.

## How to use
Reads a properties file temptoexcel.properties, override with --prop filename.properties

* start - int: Start index if mode==""
* file - string: Name and path to store data
* template - string: Name and path to template to modify, if empy new file created
* tempserver - string: URL to tempserver
* sensornames - string[]: Name of sensors (comma separated)
* mode - string: Modes of operation
  * last24 - Retrieve the last 24h of data 
  * "" - Retrieve data starting from start id

## Build
`
go build temptoexcel.go
`
## Start
`
temptoexcel
`

## Sheets
As in Template.xlsx

### Sensor Data
Sheet where data is pasted

### Derived Data
Intended for calculations (removed to reduce size)

### Sensor Graphs
X - Y plot of Data
  
### Diff Graphs
X - Y plot of Derived Data (removed to reduce size)
