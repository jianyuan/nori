package nori

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/net/context"

	"gopkg.in/tomb.v2"

	"github.com/jianyuan/nori/log"
	"github.com/jianyuan/nori/message"
	"github.com/jianyuan/nori/transport"
	"github.com/kr/pretty"
)

type Server struct {
	context.Context
	Tasks  map[string]*Task
	config *Configuration
	tomb   *tomb.Tomb
}

type Configuration struct {
	Name      string
	Transport transport.Driver
}

func NewServer(ctx context.Context, config *Configuration) (*Server, error) {
	ctx, err := configureLogger(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Logger configuration error: %s", err)
	}

	srv := &Server{
		Context: ctx,
		Tasks:   make(map[string]*Task),
		config:  config,
		tomb:    new(tomb.Tomb),
	}

	log.FromContext(srv).Info("Server set up successful")

	return srv, nil
}

func (s *Server) RegisterTask(t *Task) {
	// TODO: validation
	t.Name = s.config.Name + "." + t.Name

	if _, existing := s.Tasks[t.Name]; existing {
		log.FromContext(s).Panicf("Task %q already registered", t.Name)
	}
	s.Tasks[t.Name] = t
}

func (s *Server) printInfo() {
	if hostname, err := os.Hostname(); err == nil {
		log.FromContext(s).Infoln("Hostname:", hostname)
	}

	log.FromContext(s).Infoln("Registered tasks:")
	for _, t := range s.Tasks {
		log.FromContext(s).Infoln("-", t.Name)
	}

}

func (s *Server) Run() error {
	s.tomb.Go(s.run)

	go s.RunManagementServer(":8080")

	return nil
}

func (s *Server) run() error {
	s.printInfo()

	for {
		log.FromContext(s).Infoln("Connecting")
		select {
		case err := <-s.setupTransport():
			if err != nil {
				log.FromContext(s).Errorln("Transport setup error:", err)
				time.Sleep(5 * time.Second)
				continue
			}

			log.FromContext(s).Infoln("Connected!")

			s.consumeMessages()

		case <-s.tomb.Dying():
			log.FromContext(s).Infoln("Cancelled")
			break

		case <-time.After(5 * time.Second):
			// TODO better retry mechanism
			log.FromContext(s).Errorln("Timed out")
			time.Sleep(5 * time.Second)
			continue
		}
	}

	return s.config.Transport.Close()
}

func (s *Server) setupTransport() <-chan error {
	errChan := make(chan error)
	s.tomb.Go(func() error {
		if err := s.config.Transport.Init(s.Context); err != nil {
			errChan <- err
			return nil
		}
		errChan <- s.config.Transport.Setup()
		return nil
	})
	return errChan
}

func (s *Server) consumeMessages() {
	reqChan, err := s.config.Transport.Consume("celery")
	if err != nil {
		log.FromContext(s).Errorln("Transport consume error:", err)
		return
	}

	for {
		select {
		case req := <-reqChan:
			pretty.Println("Request:", req)

			if task, ok := s.Tasks[req.TaskName]; ok {
				resp, err := callTaskHandlerSafely(task.Handler, req)
				if err != nil {
					log.FromContext(s).Errorln("Task handler errored:", err)
				} else {
					pretty.Println("Response:", resp)
					log.FromContext(s).Infoln("Replying...")

					if err := s.config.Transport.Reply(req, resp); err != nil {
						log.FromContext(s).Errorln("Reply errored:", err)
					}
				}
			} else {
				log.FromContext(s).Errorln("Unknown task:", req.TaskName)
			}

		case <-s.tomb.Dying():
			return
		}
	}
}

func (s *Server) Wait() error {
	return s.tomb.Wait()
}

func (s *Server) Stop() {
	s.tomb.Kill(nil)
}

func callTaskHandlerSafely(t TaskHandlerFunc, req *message.Request) (message.Response, error) {
	var resp message.Response
	var err error
	resp, err = func() (message.Response, error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("Handler panicked: %v", r)
			}
		}()
		return t(req)
	}()
	if err != nil {
		return nil, err
	}
	return resp, nil
}
