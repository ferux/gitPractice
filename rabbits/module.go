package rabbits

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ferux/gitPractice"

	"github.com/go-kit/kit/log"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var logger *log.Context

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

	l *logrus.Entry
}

// Prepare checks if connection available for rabbit.
func Prepare(connString string) (Rabbit, error) {
	conn, err := amqp.Dial(connString)
	if err != nil {
		return Rabbit{}, err
	}
	_ = conn.Close()
	lg := log.NewLogfmtLogger(os.Stdout)
	logger = log.NewContext(lg).WithPrefix("pkg", "rabbits")
	lr := logrus.New().WithField("pkg", "rabbits")
	if gitPractice.Environment != "develop" {
		logrus.SetLevel(logrus.WarnLevel)
	}
	return Rabbit{connString, 0, true, nil, nil, nil, nil, amqp.Queue{}, nil, lr}, nil
}

// ClosingWorker tests how rabbit worker works with closing
func (r *Rabbit) ClosingWorker() (err error) {
	l := r.l.WithField("fn", "ClosingWorker")
	if !r.prepared {
		return errors.New("config not prepared")
	}

	defer func() {
		l.Info("defered initClose")
		if err = r.initClose(); err != nil {
			l.WithError(err).Warn("can't close")
		}
	}()

	for tries := 3; tries > 0; tries-- {
		notdone := true

		l.Info("Sleeping 3 seconds before starting")
		time.Sleep(time.Second * 3)
		l.Info("Initiating")
		if err := r.prepareWorker(); err != nil {
			l.WithError(err).WithField("tries left", tries).Error("can't prepare worker")
			continue
		}

		exitc := time.After(time.Second * 10)

		for notdone {
			select {
			case d := <-r.d:
				l.WithField("body", string(d.Body)).Print("got new message")
			case e := <-r.acterrc:
				l.WithError(e).Error("got error from acterrc")
			case <-exitc:
				l.Info("10 seconds passed. Exiting.")
				notdone = false
				break
			}
		}
		l.Warn("Main loop has been closed. Restarting...")
	}
	l.Warn("out of tries. Shutting down.")
	return err
}

func msg(l *log.Context, txt string) {
	_ = l.Log("msg", txt)
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
	l := r.l.WithField("fn", "prepareWorker")
	defer func() {
		if err != nil {
			l.WithError(err).Error("got error")
		}
	}()

	l.Info("initClose")
	if err = r.initClose(); err != nil {
		return
	}
	l.Info("initRun")
	if err = r.initRun(); err != nil {
		return
	}
	l.Info("prepareExchange")
	if r.q, r.d, err = r.prepareExchange(); err != nil {
		return
	}
	l.Info("notifyClose")
	r.acterrc = r.actconn.NotifyClose(make(chan *amqp.Error))
	return
}

// Run the whole thing
func (r *Rabbit) Run() error {
	l := logger.WithPrefix("fn", "Run")
	if !r.prepared {
		return errors.New("config not prepared")
	}

	msg(l, "starting MustListen")
	go r.MustListen()
	msg(l, "starting MustRun")
	r.MustRun()

	return nil
}

func (r *Rabbit) initRun() error {
	l := r.l.WithField("fn", "initRun")
	l.Info("start")
	defer l.Info("finish")
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
	l := r.l.WithField("fn", "initClose")
	l.Info("start")
	defer l.Info("finish")

	if r.actchan != nil {
		l.Info("actchan is not nil. Closing")
		if err := r.actchan.Close(); err != nil {
			return err
		}
		r.actchan = nil
	}

	if r.actconn != nil {
		l.Info("actconn is not nil. Closing")
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
	l := logger.WithPrefix("fn", "MustRun")
	msg(l, "start")
	defer msg(l, "finish")

	msg(l, "MustRun MainLoop")
	for {
		msg(l, "Trying to initiate connection for publishing messages")
		if err := r.initRun(); err != nil {
			msg(l.With("err", err), "can't init connection of run")
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
			msg(l.With("err", err), "can't declare queue")
			continue
		}

		errc := r.actconn.NotifyClose(make(chan *amqp.Error))
		go func() {

			now := time.Now()
			id := fmt.Sprintf("%010d", rand.Uint32())
			msg(l.WithPrefix("id", id), "Waiting for errc event")
			for err := range errc {
				msg(l.WithPrefix("id", id).With("time", time.Since(now), "err", err), "errorc")
			}
			msg(l.WithPrefix("id", id), "channel has been closed")

		}()
		msg(l, "Sending publish msg")
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
			msg(l.With("err", err), "can't publish")
		}
		msg(l, "Sleeping 4 seconds before making new connection")
		time.Sleep(time.Second * 4)

		msg(l, "Closing old connection before opening new one")
		err = r.initClose()
		if err != nil {
			msg(l.With("err", err), "can't close actconn/actchan")
		}

	}
}

// MustListen panics if any error acquired
func (r *Rabbit) MustListen() {
	l := logger.WithPrefix("fn", "MustListen")
	msg(l, "start")
	defer msg(l, "finish")
	for {
		var err error
		msg(l, "MustListen initiating new connection")
		r.conn, err = amqp.Dial(r.cs)
		if err != nil {
			msg(l.With("err", err), "can't dial to mq")
			return
		}

		errc := r.conn.NotifyClose(make(chan *amqp.Error))

		ch, err := r.conn.Channel()
		if err != nil {
			msg(l.With("err", err), "can't get channel")
			return
		}
		if err := ch.Qos(1, 0, false); err != nil {
			msg(l.With("err", err), "can't set QOS")
			return
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
			msg(l.With("err", err), "can't declare: ")
		}
		msg(l, "Consuming new messages")
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
				msg(l.With("err", err), "Channel closed with error")
				break
			case msgd := <-msgs:
				msg(l.With("body", string(msgd.Body)), "New message")
			}
		}
	}
}

