package abtest

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/airbrake/gobrake"
)

// Run the package
func Run() {
	fmt.Println("testing airbrake")
	nt := Gobrake("dev", "v1.0.0", "1", "tcp://127.0.0.1:3000", 1)
	nt.Notify(errors.New("hello"), nil)
	defer nt.NotifyOnPanic()
	panic("test")
}

// Gobrake creates new gobrake client
func Gobrake(env, ver, key, host string, id int64) ErrorNotifier {
	gobrake.SetLogger(log.New(os.Stderr, "ERRBIT ", 0))

	notifier := gobrake.NewNotifier(id, key)
	notifier.SetHost(host)
	notifier.AddFilter(func(n *gobrake.Notice) *gobrake.Notice {
		n.Context["environment"] = env
		n.Context["version"] = ver
		return n
	})

	return notifier
}

// ErrorNotifier for interface
type ErrorNotifier interface {
	Notify(err interface{}, req *http.Request)
	Close() error
	NotifyOnPanic()
}
