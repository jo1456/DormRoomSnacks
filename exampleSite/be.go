//Jeffrey Oberg jo1456
//Project part 2 CRUD APP with front and back end servers
// back end

package main

import (
  "fmt"
  "net"
  "encoding/json"
  "flag"
  "errors"
)

var (
  reportMap map[string]*SurfReport // map holding all reports rather than a DB
  connection net.Conn
  encoder *json.Encoder
  decoder *json.Decoder
)

// general data struct
type SurfReport struct {
  BeachName string
  WaveHeight string
  Period string
  Wind string
}

// struct of the beach name of a specific report for delete and get requests
type Params struct {
  BeachName string
}

// request struct for passing function names
type Request struct {
  FunctionName string
}

// struct of a list of beach names for display on the home page
type ReportsList struct {
  BeachNames []string
}

// struct of status of call
type Status struct {
  Success bool
}

func main() {
        // creates a variable with the passed flag. default value of 8090
        listenPort := flag.String("listen", "8090", "port to listen on")
        flag.Parse()

        // create and load the map
        reportMap = make(map[string]*SurfReport)
        GoodHarbor := &SurfReport{BeachName: "Good Harbor", WaveHeight: "1ft", Period: "7s", Wind: "Glass"}
        LongSands := &SurfReport{BeachName: "Long Sands", WaveHeight: "11ft", Period: "17s", Wind: "Off Shore"}
        Pipeline := &SurfReport{BeachName: "Pipeline", WaveHeight: "15ft", Period: "10s", Wind: "Cross"}
        Mavs := &SurfReport{BeachName: "Mavs", WaveHeight: "52ft", Period: "36s", Wind: "Glass"}
        reportMap[GoodHarbor.BeachName] = GoodHarbor
        reportMap[LongSands.BeachName] = LongSands
        reportMap[Pipeline.BeachName] = Pipeline
        reportMap[Mavs.BeachName] = Mavs

        // infinite loop that accepts connections and proccesses
        // them until the front end disconnects
        for {
          // listen on passed port
          listener, err := net.Listen("tcp", ":" + *listenPort)
          if err != nil {
            fmt.Println(err)
            return
          }
          defer listener.Close()

          // accept a connection
          connection, err := listener.Accept()
          if err != nil {
            fmt.Println(err)
            return
          }

          // assign the global encoder and decoder to the created connection
          encoder = json.NewEncoder(connection)
          decoder = json.NewDecoder(connection)

          closed := false
          // loop to accept and response to requests until it is notified the frontend disconects
          for !closed {
            // status of operation. In case of failure fe can redirect to homepage
            status := Status{Success: true}

            var req Request
            decoder.Decode(&req)

            switch req.FunctionName {
              // list report rpc Function
              // accepts: void
              // returns: list of all beach names, status
            case "ListReports":
              ListReports()

              // view report rpc Function
              // accepts: beach name
              // returns: SurfReport struct for requested beach, status
            case "ViewReport":
              var beachName Params
              decoder.Decode(&beachName)
              err := GetReport(beachName.BeachName)
              if err != nil {
                status.Success = false
              }

              // delete report rpc Function
              // accepts: beach name
              // returns: status
            case "DeleteReport":
              var beachName Params
              decoder.Decode(&beachName)
              delete(reportMap, beachName.BeachName)

              // update report rpc Function
              // accepts: surfReport
              // returns: status
            case "UpdateReport":
              var report SurfReport
              decoder.Decode(&report)
              err := UpdateReport(report)
              if err != nil {
                status.Success = false
              }

              // add report rpc Function
              // accepts: surfReport
              // returns: status
            case "AddReport":
              var report SurfReport
              decoder.Decode(&report)
              err := CreateReport(report)
              if err != nil {
                status.Success = false
              }
              // set the boolean to true so the loop exits and the connection closes
            case "Closed":
              closed = true

            }
            encoder.Encode(status)
          }
          listener.Close() // close the listener so new connections can be made
        }
}

// ListReports encodes and sends a list of all beaches to the requesting front end
// accepts: void
// returns: err
func ListReports() error {
    list := ReportsList{BeachNames: make([]string, 0, len(reportMap))}
    for key, _ := range reportMap {
      list.BeachNames = append(list.BeachNames, key)
    }
    encoder.Encode(list)
    return nil
}

// GetReport encodes and sends a report of the passed beach name to the requesting front end
// accepts: string beach
// returns: err
func GetReport(beach string) error {

    val, ok := reportMap[beach]

    // encode either way as front end is expecting a report even if it does not
    // exist. The error is then used to set status to failed
    encoder.Encode(val)

    if !ok {
      return errors.New("Beach does not exist")
    }

    return nil
}

// UpdateReport accepts a surfReport and updates the corresponding entry in the map
// accepts: SurfReport report
// returns: err
func UpdateReport(report SurfReport) error {
  oldReport, ok := reportMap[report.BeachName]

  // if value is in map then replace old values with new ones
  if ok {
    if report.WaveHeight != "" {
      oldReport.WaveHeight = report.WaveHeight
    }
    if report.Period != "" {
      oldReport.Period = report.Period
    }
    if report.Wind != "" {
      oldReport.Wind = report.Wind
    }
  } else {
    return errors.New("Beach does not exist")
  }
  return nil
}

// CreateReport accepts a surfReport and creates the corresponding entry in the map
// accepts: void
// returns: err
func CreateReport(report SurfReport) error {
  _, valueInMap := reportMap[report.BeachName]

  // only add if there is not already a beach with the same name
  if !valueInMap {
    reportMap[report.BeachName] = &report
  } else {
    return errors.New("Beach already exists")
  }
  return nil
}
