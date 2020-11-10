//Jeffrey Oberg jo1456
//Project part 2 CRUD APP with front and back end servers
// Front end

package main

import (
  "github.com/kataras/iris/v12"
  "strings"
  "sort"
  "fmt"
  "flag"
  "net"
  "encoding/json"
)

var (
  connection net.Conn
  encoder *json.Encoder
  decoder *json.Decoder
  app *iris.Application
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
// This is mostly for ViewReport. In the case where someone attempts to view
// a report that has been deleted they should be redirected to the home page.
// Other functions follow same error handling for consistency
type Status struct {
  Success bool
}

func main() {

    // creates a variable with the passed flag. default value of 8080
    listenPort := flag.String("listen", "8080", "port to listen on")
    backendHostandPort := flag.String("backend", ":8090", "host name and port of backend. (Format: hostName:port)")
    flag.Parse()

    // attempt to connect to the backend
    connection, err := net.Dial("tcp", *backendHostandPort)
    if err != nil {
      fmt.Println("Connection to passed backend host and port failed. Error:")
      fmt.Println(err)
      return
    }
    defer connection.Close()

    // assign the global encoder and decoder to the created connection
    encoder = json.NewEncoder(connection)
    decoder = json.NewDecoder(connection)

    defer encoder.Encode(Request{FunctionName: "Closed"})

    app = iris.New()

    // add handlers
    app.Get("/", GetHomePage)
    app.Get("/add-report-page", GetAddReportPage)
    app.Get("/update-form/{beachName}", GetUpdateFormPage)
    app.Get("/delete/{beachName}", DeleteReport)
    app.Get("/view-report/{beachName}", ViewReport)
    app.Post("/update", UpdateReport)
    app.Post("/add-report", AddReport)

    // turn on the app
    app.Listen(":"+*listenPort)
}

// This function provides the html for the homepage of the website which will list all reports.
// Provides the html to the passed context
// Gets beaches from backend
// Accepts: iris.context which will be provided by handler that calls the fucntions
// returns: nothing
func GetHomePage(ctx iris.Context) {

  // Request all report names from backend. This function accepts no paramaters
  encoder.Encode(Request{FunctionName: "ListReports"})

  // recieve the list of reports from the backend.
  var list ReportsList
  decoder.Decode(&list)

  var status Status
  decoder.Decode(&status) // status is always returned

  // sort the list for consistent display
  sort.Strings(list.BeachNames)

  // populate the html with recieved data
  var beaches string
  for i := range list.BeachNames {
    name := list.BeachNames[i]
    beaches = beaches + "<tr><td><a href=/view-report/"+strings.Replace(name," ", "_",-1)+">" + name + "</a></td></tr>"
  }

  beforeTable := `<!DOCTYPE html>
                  <html>
                    <body>

                    <h1>Click a beach to view the current surf conditions!</h1>

                    <table>
                      <tr>
                        <th>Select a Beach</th>
                      <tr>`
  afterTable := `   </table>

                    <br>
                    <br>
                    <a href="/add-report-page"> add report </a>

                    </body>
                  </html>
                  `

  ctx.HTML(beforeTable + beaches + afterTable)
}

// This function provides the html for the page to add a new report. The C in CRUD
// Provides the html to the passed context. The form will then hit the /add-report endpoint.
// Accepts: iris.context which will be provided by handler that calls the fucntions
// returns: nothing
func GetAddReportPage(ctx iris.Context) {
  ctx.HTML(`
    <html>

        <form action="/add-report" method="POST">
          <label>Beach:</label>
          <input type="text" name="beachName">

          <label>Height:</label>
          <input type="text" name="waveHeight">

          <label>Period:</label>
          <input type="text" name="period">

          <label>Wind:</label>
          <input type="text" name="wind">

          <input type="submit" value="Submit">
        </form>

      </body>
    </html>
    `)
}

// This function provides the html for the page to update a report. The U in CRUD
// Provides the html to the passed context. The form will hit the update endpoint.
// Accepts: iris.context which will be provided by handler that calls the fucntions
// returns: nothing
func GetUpdateFormPage(ctx iris.Context) {
  var params Params
  err := ctx.ReadParams(&params)
  if err != nil {
      GetHomePage(ctx)
      return
  }

  // replace underscores with spaces to restore beach name
  params.BeachName = strings.Replace(params.BeachName, "_", " ", -1)

  html := "<html><h1>Update Report for " + params.BeachName + "</h1>"
  html = html + `<form action="/update" method="POST">`
  html = html + ` <input type="hidden" name="beachName" value="`+params.BeachName+`">`
  html = html + `
                  <label>Height:</label>
                  <input type="text" name="waveHeight">

                  <label>Period:</label>
                  <input type="text" name="period">

                  <label>Wind:</label>
                  <input type="text" name="wind">

                  <input type="submit" value="Submit">
                </form>

              </body>
            </html>`

  ctx.HTML(html)

}

// This function provides the html for the page to read a report. The R in CRUD
// Provides the html to the passed context.
// Accepts: iris.context which will be provided by handler that calls the fucntions
// returns: nothing
func ViewReport(ctx iris.Context) {
  var params Params
  err := ctx.ReadParams(&params)
  if err != nil {
      GetHomePage(ctx)
      return
  }
  params.BeachName = strings.Replace(params.BeachName, "_", " ", -1)

  // Request a report from the backend by passing the ViewReport function name
  // and the name of the beach to get the report for.
  encoder.Encode(Request{FunctionName: "ViewReport"})
  encoder.Encode(params)

  // recieve the report from the backend
  var report SurfReport
  decoder.Decode(&report)

  // if request fails, return to homepage
  var status Status
  decoder.Decode(&status)

  if !status.Success {
    GetHomePage(ctx)
    return
  }

  // populate the html with the recieved report
  if report.BeachName != "" {
    html :=      "<table><tr><td>Beach: "+report.BeachName + "</td></tr>"
    html = html   + "<tr><td>Wave Height: "+report.WaveHeight + "</td></tr>"
    html = html   + "<tr><td>Period: "+report.Period + "</td></tr>"
    html = html   + "<tr><td>Wind: "+report.Wind + "</td></tr>"
    html = html+ "</table>"
    html = html+ `<a href="/update-form/`+strings.Replace(params.BeachName, " ", "_", -1)+`">Update</a> <br>`
    html = html+` <a href="/delete/`+strings.Replace(params.BeachName, " ", "_", -1)+`"">delete</a> <br>`
    html = html+` <a href="/">Return</a>`
    ctx.HTML(html)
  }
}

// This function deletes a report from the reportMap global variable and then sends the user
// to the home page where the report map is used to reload the reports. The D in CRUD
// Accepts: iris.context which will be provided by handler that calls the fucntions
// returns: nothing
func DeleteReport(ctx iris.Context) {
  var params Params
  err := ctx.ReadParams(&params)
  if err != nil {
      GetHomePage(ctx)
      return
  }
  params.BeachName = strings.Replace(params.BeachName, "_", " ", -1)

  // Request to delete a report by passing the DeleteReport function name
  // and the name of the report to delete
  encoder.Encode(Request{FunctionName: "DeleteReport"})
  encoder.Encode(params)

  var status Status
  decoder.Decode(&status) // status is always returned

  ctx.Redirect("/", iris.StatusPermanentRedirect)
}

// This function updates reports. It accepts new wave height, period, and Wind
// from the form and updates the element in the map.
// Accepts: iris.context which will be provided by handler that calls the fucntions
// returns: nothing
func UpdateReport(ctx iris.Context) {
  var report SurfReport

  err := ctx.ReadForm(&report)
  if err != nil {
      GetHomePage(ctx)
      return
  }

  // request to update a report by passing the UpdateReport function name
  // and the new report struct
  encoder.Encode(Request{FunctionName: "UpdateReport"})
  encoder.Encode(report)

  // if request fails, return to homepage
  var status Status
  decoder.Decode(&status)

  if !status.Success {
    GetHomePage(ctx)
    return
  }

  ctx.StatusCode(iris.StatusCreated)
  GetHomePage(ctx)
}

// This function adds a report. It accepts a new name, wave height, period, and Wind
// from the form and adds the element to the map.
// Accepts: iris.context which will be provided by handler that calls the fucntions
// returns: nothing
func AddReport(ctx iris.Context) {
  var report SurfReport

  err := ctx.ReadForm(&report)
  if err != nil {
      GetHomePage(ctx)
      return
  }

  // request to add a report by passing the UpdateReport function name
  // and the new report struct
  encoder.Encode(Request{FunctionName: "AddReport"})
  encoder.Encode(report)

  // if request fails, return to homepage
  var status Status
  decoder.Decode(&status)

  if !status.Success {
    GetHomePage(ctx)
    return
  }

  ctx.StatusCode(iris.StatusCreated)
  GetHomePage(ctx)
}
