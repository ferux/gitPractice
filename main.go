package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ferux/gitPractice/abtest"
	"github.com/ferux/gitPractice/hidstruct"
	"github.com/ferux/gitPractice/io"
	"github.com/ferux/gitPractice/rabbits"
	"github.com/ferux/gitPractice/rnd"
	"github.com/ferux/gitPractice/viewer"
)

var name *string
var dur *int

func parseFlags() {
	name = flag.String("name", "dev", "Sets the name of the project")
	dur = flag.Int("d", 10, "Sets the loading speed of the project")
	flag.Parse()
	if *dur < 0 {
		*dur = 10
	}
}

func doSomeWork() {
	fmt.Fprintln(os.Stdout, "🌞 Day time")
	for i := 0; i <= 100; i++ {
		switch i {
		case 0:
			fmt.Fprint(os.Stdout, "[🌑 0%]")
		case 25:
			fmt.Fprint(os.Stdout, "[🌘 25%]")
		case 50:
			fmt.Fprint(os.Stdout, "🌗 50%")
		case 75:
			fmt.Fprint(os.Stdout, "🌖 75%")
		case 100:
			fmt.Fprint(os.Stdout, "🌕 100%\n")
		default:
			if i%5 == 0 {
				fmt.Fprint(os.Stdout, "✈")
			}
		}
		time.Sleep(time.Millisecond * time.Duration(*dur))
	}
}

func doRnd() {
	fmt.Println("in MAIN")
	fmt.Println("Version:", rnd.Version)
	fmt.Println("InitVersion:", rnd.InitVersion)
	rnd.Init("Alex", 19)
	fmt.Println("Version:", rnd.Version)
	fmt.Println("InitVersion:", rnd.InitVersion)
	fmt.Println("stepping in viewer.View()")
	defer fmt.Println("stepping out viewer.View()")
	viewer.View()
}

func doHidden() {
	str := hidstruct.Init("Alex", "qwerty", "ez")
	data, err := json.MarshalIndent(&str, "", " ")
	if err != nil {
		log.Printf("can't marshal: %v", err)
		return
	}
	log.Printf("JSON: %s", data)
}

func doCopy() {
	r := hidstruct.NewRegular("Alex")
	r.Add("password", "qwerty")
	r.Add("login", "Alex")
	r.IDs["Alex"] = 100
	fmt.Printf("%#v\n\n", r)
	newr := r.DeepCopySafe()
	r.Name = "Alex2008"
	r.IDs["Alex"] = 200
	r.Set("login", "Alex2008")
	fmt.Printf("   r: %#v\n\n", r)
	fmt.Printf("newr: %#v", newr)
}

func doIOReadFull() {
	err := io.TryReadFull()
	if err != nil {
		fmt.Printf("can't do tryreadall. reason: %s", err)
	}
}

func doAirbrake() {
	abtest.Run()
}

func doRabbits() {
	r, err := rabbits.Prepare("amqp://localhost:5672")
	if err != nil {
		fmt.Printf("can't prepare rabbit: %v\n", err)
		return
	}
	if err := r.Run(); err != nil {
		fmt.Printf("can't init rabbit: %v", err)
	}
}

func main() {
	// parseFlags()
	// fmt.Fprintf(os.Stdout, "👍 Starting new project %s\n", *name)
	// doSomeWork()
	// fmt.Fprintf(os.Stdout, "👎🏽 Finishing new project %s\n", *name)
	// doRnd()
	// task := prl.Init("Parallels", "v0.0.1")
	// task.Run()
	// doCopy()
	// doAirbrake()
	doRabbits()
}
