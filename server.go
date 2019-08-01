package main
import (
	"net/http"
	"strings"
	
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"os"
	"time"
)

var (
	// Use mcu pin 8, corresponds to physical pin 17 on the pi
	pin = rpio.Pin(17)
)

var gonging = false
func gong(){
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

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	helloText := "You've reached the GongServer."
	w.Write([]byte(helloText))
	if(message == "gong"){
		go gong()
	}
}
func main() {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer rpio.Close()

	pin.Output()
	
	http.HandleFunc("/", sayHello)
	if err := http.ListenAndServe(":4664", nil); err != nil {
		panic(err)
	}
}
