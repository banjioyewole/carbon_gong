package main
import (
    "html/template"
    "net/http"
    "fmt"
    "github.com/stianeikeland/go-rpio"
    "os"
    "time"
)

var (
    // Use mcu pin 8, corresponds to physical pin 17 on the pi
    pin = rpio.Pin(17)
    gonging = false
    tpl *template.Template
)

func init() {
    tpl = template.Must(template.ParseFiles("gong.gohtml"))
}

func main() {

    fmt.Println("Welcome to Gong, Carbon Edition.")

    if err := rpio.Open(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    defer rpio.Close()

    pin.Output()
    
    http.HandleFunc("/", SayHello)
    http.HandleFunc("/gong", HandleGong)

    if err := http.ListenAndServe(":4664", nil); err != nil {
        panic(err)
    }
}

func SayHello(w http.ResponseWriter, r *http.Request) {
    helloText := "You've reached the GongServer."
    w.Write([]byte(helloText))
}

func HandleGong(w http.ResponseWriter, r *http.Request) {
    timeMessage, isOkToGong := isTimeOk()
    if r.Method == http.MethodPost {
        if isOkToGong {
            go gong()
        }
        w.WriteHeader(http.StatusBadRequest) // to throw the bots off
    } else {
        responseText := fmt.Sprintf("You've reached the GongServer. %s", timeMessage)

        templateData := struct {
            Text string
            ShowGongButton bool
        } {
            responseText,
            isOkToGong,
        }
        w.Header().Set("Content-Type", "text/html")
        tpl.ExecuteTemplate(w, "./gong.gohtml", templateData)
    }
}

func isTimeOk() (string, bool) {
    earlyLateBoundary := 4
    start := 9
    end := 23
    currentTime := time.Now()

    if currentTime.Hour() >= end || currentTime.Hour() < earlyLateBoundary {
        return "Sorry, we don't gong this late.", false
    } else if currentTime.Hour() < start {
        return "Sorry, we don't gong this early.", false
    } else {
        return "", true
    }
}

func gong() {
    if !gonging {
        gonging = true
        fmt.Println("gong!")
        pin.High()
        time.Sleep(time.Second / 10)
        pin.Low()
        time.Sleep(time.Second / 10)
        gonging = false
    }
}

