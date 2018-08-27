package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/ferux/gitPractice"
	"github.com/ferux/gitPractice/abtest"
	"github.com/ferux/gitPractice/hidstruct"
	"github.com/ferux/gitPractice/io"
	"github.com/ferux/gitPractice/rabbits"
	"github.com/ferux/gitPractice/rnd"
	"github.com/ferux/gitPractice/viewer"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/levels"
)

// var defines flags of the application
var (
	dur        = flag.Int("d", 10, "Sets the loading speed of the project")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	info       = flag.Bool("version", false, "shows current version of the application")
)

// Logger is the app-wide logger.
var Logger *levels.Levels
var kitlogger kitlog.Logger

func parseFlags() {
	flag.Parse()

	if *info {
		fmt.Fprintf(
			os.Stdout, "Version: %s\nRevision: %s\nEnvironment: %s\n",
			gitPractice.Version, gitPractice.Revision, gitPractice.Environment,
		)
		// because os.Exit() ignores defer functions in other places
		runtime.Goexit()
	}

	if *dur < 0 {
		*dur = 10
	}

	kitlogger = kitlog.NewLogfmtLogger(os.Stdout)
	loggerContext := kitlog.NewContext(kitlogger)

	log := levels.New(loggerContext)
	Logger = &log
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

func doSomeMagicRabbits() {
	r, err := rabbits.Prepare("amqp://localhost:5672")
	if err != nil {
		panic(err)
	}
	if err := r.ClosingWorker(); err != nil {
		panic(err)
	}
}

func doSomeLogs() {
	m := make(map[string]interface{})
	m["test"] = "hello"
	m["test2"] = 10
	m["another"] = "header"
	logger := Logger.With("fn", "doSomeLogs")
	data, err := json.Marshal(m)
	if err != nil {
		_ = logger.Error().Log("json.Marshal error", err)
	}
	if err := logger.Debug().Log("hello", "hey", "map", string(data)); err != nil {
		log.Println(err)
	}
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(gitPractice.Version + " ")
	log.Print("Application gitPractice started.")
	defer func() {
		log.Print("Application gitPractice finished.")
		os.Exit(0) // if you call runtime.Goexit it finishes current goroutine but not the other ones.
	}()
	parseFlags()

	// doSomeLogs()
	doSomeMagicRabbits()
	// if *cpuprofile != "" {
	// 	f, err := os.Create(*cpuprofile)
	// 	if err != nil {
	// 		log.Fatal("could not create CPU profile: ", err)
	// 	}
	// 	if err := pprof.StartCPUProfile(f); err != nil {
	// 		log.Fatal("could not start CPU profile: ", err)
	// 	}
	// 	defer pprof.StopCPUProfile()
	// }

	// _ = benchs.ConcatAppend(100000)
	// _ = benchs.ConcatCopy(100000)
	// _ = benchs.ConcatBuilderPreGrow(100000)

	// if *memprofile != "" {
	// 	f, err := os.Create(*memprofile)
	// 	if err != nil {
	// 		log.Fatal("could not create memory profile: ", err)
	// 	}
	// 	runtime.GC() // get up-to-date statistics
	// 	if err := pprof.WriteHeapProfile(f); err != nil {
	// 		log.Fatal("could not write memory profile: ", err)
	// 	}
	// 	f.Close()
	// }
	// fmt.Fprintf(os.Stdout, "👍 Starting new project %s\n", *name)
	// doSomeWork()
	// fmt.Fprintf(os.Stdout, "👎🏽 Finishing new project %s\n", *name)
	// doRnd()
	// task := prl.Init("Parallels", "v0.0.1")
	// task.Run()
	// doCopy()
	// doAirbrake()
	// doRabbits()
}
