package rabbits

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/streadway/amqp"
)

var logger *log.Logger

const (
	exchName   = "exch.test"
	exchType   = "topic"
	routingKey = "*"
	queueName  = "qmain"
)

// Rabbit used for creating new connections (DDoSing rabbit)
type Rabbit struct {
	cs       string
	Num      int
	prepared bool
	// conn is a connection of listener
	conn *amqp.Connection
	// actconn is a connection of sender
	actconn *amqp.Connection
	actchan *amqp.Channel
	acterrc chan *amqp.Error

	// for queue retrieving
	q amqp.Queue
	d <-chan amqp.Delivery
}

// Prepare checks if connection available for rabbit.
func Prepare(connString string) (Rabbit, error) {
	conn, err := amqp.Dial(connString)
	if err != nil {
		return Rabbit{}, err
	}
	_ = conn.Close()
	logger = log.New(os.Stdout, "[rabbits.Rabbit] ", 0)
	return Rabbit{connString, 0, true, nil, nil, nil, nil, amqp.Queue{}, nil}, nil
}

// ClosingWorker tests how rabbit worker works with closing
func (r *Rabbit) ClosingWorker() (err error) {
	l := log.New(os.Stdout, logger.Prefix()+"[ClosingWorker] ", 0)
	if !r.prepared {
		return errors.New("config not prepared")
	}

	defer func() {
		l.Println("defered initClose")
		if err = r.initClose(); err != nil {
			l.Printf("error while closing: %v", err)
		}
	}()

	for tries := 3; tries > 0; tries-- {
		notdone := true

		l.Println("Sleeping 5 seconds before starting")
		time.Sleep(time.Second * 5)
		l.Println("Initiating")
		if err := r.prepareWorker(); err != nil {
			l.Printf("can't prepare worker: %v", err)
			continue
		}

		exitc := time.After(time.Second * 15)

		for notdone {
			select {
			case d := <-r.d:
				l.Printf("got delivery: %s", d.Body)
			case e := <-r.acterrc:
				l.Printf("got error: %s", e)
			case <-exitc:
				l.Println("15 seconds passed")
				notdone = false
				break
			}
		}
		l.Println("Main loop has been closed. Restarting...")
	}
	l.Println("out of tries. Shutting down.")
	return err
}

func (r *Rabbit) prepareExchange() (q amqp.Queue, d <-chan amqp.Delivery, err error) {

	if err = r.actchan.ExchangeDeclare(
		exchName,
		exchType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return
	}

	if q, err = r.actchan.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return
	}

	if err = r.actchan.QueueBind(
		queueName,
		routingKey,
		exchName,
		false,
		nil,
	); err != nil {
		return
	}

	d, err = r.actchan.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	return
}

func (r *Rabbit) prepareWorker() (err error) {
	l := log.New(os.Stdout, logger.Prefix()+"[prepareWorker] ", 0)
	l.Println("initClose")
	if err = r.initClose(); err != nil {
		return
	}
	l.Println("initRun")
	if err = r.initRun(); err != nil {
		return
	}
	l.Println("prepareExchange")
	if r.q, r.d, err = r.prepareExchange(); err != nil {
		return
	}
	l.Println("notifyClose")
	r.acterrc = r.actconn.NotifyClose(make(chan *amqp.Error))
	return
}

// Run the whole thing
func (r *Rabbit) Run() error {
	if !r.prepared {
		return errors.New("config not prepared")
	}

	logger.Println("running MustListen")
	go r.MustListen()
	logger.Println("running MustRun")
	r.MustRun()

	return nil
}

func (r *Rabbit) initRun() error {
	logger.Println("initRun")
	defer logger.Println("initRun finished")
	var err error
	r.actconn, err = amqp.Dial(r.cs)
	if err != nil {
		return fmt.Errorf("can't dial to mq: %v", err)
	}
	r.actchan, err = r.actconn.Channel()
	if err != nil {
		return fmt.Errorf("can't get channel: %v", err)
	}
	if err := r.actchan.Qos(1, 0, false); err != nil {
		return fmt.Errorf("can't set QOS: %v", err)
	}
	return nil
}

func (r *Rabbit) initClose() error {
	logger.Println("initClose")
	defer logger.Println("initClose finished")

	if r.actchan != nil {
		if err := r.actchan.Close(); err != nil {
			return err
		}
		r.actchan = nil
	}

	if r.actconn != nil {
		if err := r.actconn.Close(); err != nil {
			return err
		}
		r.actconn = nil
		r.actchan = nil
	}
	return nil
}

// MustRun connects to the RabbitMQ or Panics
func (r *Rabbit) MustRun() {
	lg := log.New(os.Stdout, logger.Prefix()+"[MustRun] ", 0)
	lg.Println("MustRun")
	defer lg.Println("MustRun closed")

	lg.Println("MustRun MainLoop")
	for {
		lg.Println("Trying to initiate connection for publishing messages")
		if err := r.initRun(); err != nil {
			lg.Printf("can't init connection of run: %v\n", err)
			continue
		}

		q, err := r.actchan.QueueDeclare(
			"test",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			lg.Printf("can't declare queue: %v", err)
			continue
		}

		errc := r.actconn.NotifyClose(make(chan *amqp.Error))
		go func() {

			now := time.Now()
			id := fmt.Sprintf("%010d", rand.Uint32())
			lg.Printf("[id: %s]Waiting for errc event", id)
			for msg := range errc {
				lg.Printf("[id: %s][took %s]errc error %v", id, time.Since(now), msg)
			}
			lg.Printf("[id %s] channel has been closed", id)

		}()
		lg.Println("Sending publish msg")
		err = r.actchan.Publish(
			q.Name,
			"",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("hello"),
			},
		)
		if err != nil {
			lg.Printf("can't publish, reason: %v\n", err)
		}
		lg.Println("Sleeping 4 seconds before making new connection")
		lg.Println()
		time.Sleep(time.Second * 4)

		lg.Println("Closing old connection before opening new one")
		err = r.initClose()
		if err != nil {
			lg.Printf("can't close actconn/actchan: %s\n", err)
		}

	}
}

// MustListen panics if any error acquired
func (r *Rabbit) MustListen() {
	lg := log.New(os.Stdout, logger.Prefix()+"[MustListen] ", 0)
	lg.Println("MustListen")
	defer lg.Println("MustListen finished")
	for {
		var err error
		lg.Println("MustListen initiating new connection")
		r.conn, err = amqp.Dial(r.cs)
		if err != nil {
			log.Fatal("can't dial to mq: ", err)
		}

		errc := r.conn.NotifyClose(make(chan *amqp.Error))

		ch, err := r.conn.Channel()
		if err != nil {
			log.Fatal("can't get channel: ", err)
		}
		if err := ch.Qos(1, 0, false); err != nil {
			log.Fatal("can't set QOS: ", err)
		}
		q, err := ch.QueueDeclare(
			"test",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatal("can't declare: ", err)
		}
		lg.Println("Consuming new messages")
		msgs, err := ch.Consume(
			q.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		for {
			select {
			case err := <-errc:
				lg.Printf("Channel closed: %v", err)
				break
			case msg := <-msgs:
				lg.Printf("New message: %s", msg.Body)
			}
		}
	}
}

