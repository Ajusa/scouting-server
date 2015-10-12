package main
import "github.com/gin-gonic/gin"
import "github.com/itsjamie/gin-cors"

import (
  
    "fmt"
    "io/ioutil"
    "bufio"
    "os"
    "strings"
)
func main() {
    gin.SetMode(gin.ReleaseMode)
    r := gin.Default()
    //Creates a new api server

    reader := bufio.NewReader(os.Stdin)
    //Reads from the command line
    fmt.Print("What is the name of the .csv file? (Type it in without the .csv at the end) ")
    //Self explantory
    text2, _ := reader.ReadString('\n')
    //text2 is the file =name, without the .csv
    text3 := strings.TrimSpace(text2)
    //This removes the newline in text2
    text := fmt.Sprint(text3 + ".csv")
    //Makes the name of the csv file so that it can be read        
    temp := fmt.Sprint(text3 + "bak")
    //Makes the backup file
    textBak := fmt.Sprint(temp + ".csv")
    //Makes the bakfile a csv
    firstRead, _ := ioutil.ReadFile(text)
    //Reads the backup file that was typed into the console
    headings := string(firstRead)
    //Makes that reading into a string
    f, _ := os.OpenFile(textBak, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
    //Goes ahead and creates/opens backup file.
    defer f.Close()
    //Needs to exist
    f.WriteString(headings)
    //Writes the headings to the backup file
    r.Use(cors.Middleware(cors.Config{
        Origins:        "*",
        Methods:        "GET, PUT, POST, DELETE",
        RequestHeaders: "Origin, Authorization, Content-Type",
        ExposedHeaders: "",
        Credentials: true,
        ValidateHeaders: false,
	}))

    //^ That entire mess lets us accept requests from anywhere, even off of different domains
    r.POST("/api", func(c *gin.Context) {
    	message := c.PostForm("data")
        //Gets the string of csv from the POST request
        fmt.Println(message)
        //Prints it to the console, so that we feel like we are working
        //^ Because reasons
    	f.WriteString(message)
        //Writes the info recieved into the file
    	c.String(200,message)
        //Tells the application that things (hopefully) did not break
        actualFile, _ := os.Create(text)
        //Goes ahead and opens the actual file
        defer actualFile.Close()
        bakFile, _ := ioutil.ReadFile(textBak)
        //Reads the backup file that was typed into the console
        bakfileInfo := string(bakFile)
        //Strings!
        actualFile.WriteString(bakfileInfo)
    })
    r.GET("/ping",  func(c *gin.Context) {
            c.String(200, "pong")
            //To check for the URL being correct
        })
    r.GET("/data",func(c *gin.Context) {
        csvRead, _ := ioutil.ReadFile(textBak)
        //Reads the backup file that was typed into the console
        csvString := string(csvRead)
        c.String(200, csvString)

        })
    r.Static("/www", "./www")
    r.Run(":80") // listen and serve on 0.0.0.0:80
}
