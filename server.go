package nori

import (
	"os"
	"time"

	"gopkg.in/tomb.v2"

	"github.com/jianyuan/nori/transport"
	"github.com/kr/pretty"
)

type Server struct {
	Name      string
	Tasks     map[string]*Task
	Transport transport.Driver

	tomb *tomb.Tomb
}

func NewServer(name string, transport transport.Driver) *Server {
	return &Server{
		Name:      name,
		Tasks:     make(map[string]*Task),
		Transport: transport,
		tomb:      new(tomb.Tomb),
	}
}

func (s *Server) RegisterTask(t *Task) {
	// TODO: validation
	t.Name = s.Name + "." + t.Name

	if _, existing := s.Tasks[t.Name]; existing {
		log.Panicf("Task %q already registered", t.Name)
	}
	s.Tasks[t.Name] = t
}

func (s *Server) printInfo() {
	if hostname, err := os.Hostname(); err == nil {
		log.Infoln("Hostname:", hostname)
	}

	log.Infoln("Registered tasks:")
	for _, t := range s.Tasks {
		log.Infoln("-", t.Name)
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
		log.Infoln("Connecting")
		select {
		case err := <-s.setupTransport():
			if err != nil {
				log.Errorln("Transport setup error:", err)
				time.Sleep(5 * time.Second)
				continue
			}

			log.Infoln("Connected!")

			s.consumeMessages()

		case <-s.tomb.Dying():
			log.Infoln("Cancelled")
			break

		case <-time.After(5 * time.Second):
			// TODO better retry mechanism
			log.Errorln("Timed out")
			time.Sleep(5 * time.Second)
			continue
		}
	}

	return s.Transport.Close()
}

func (s *Server) setupTransport() <-chan error {
	errChan := make(chan error)
	s.tomb.Go(func() error {
		errChan <- s.Transport.Setup()
		return nil
	})
	return errChan
}

func (s *Server) consumeMessages() {
	reqChan, err := s.Transport.Consume("celery")
	if err != nil {
		log.Errorln("Transport consume error:", err)
		return
	}

	for {
		select {
		case req := <-reqChan:
			pretty.Println("Request:", req)

			if task, ok := s.Tasks[req.TaskName]; ok {
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Errorln("Handler panicked:", r)
						}
					}()
					resp, err := task.Handler(req)
					if err != nil {
						log.Errorln("Task handler errored:", err)
					} else {
						pretty.Println("Response:", resp)
						log.Infoln("Replying...")

						if err := s.Transport.Reply(req, resp); err != nil {
							log.Errorln("Reply errored:", err)
						}
					}
				}()
			} else {
				log.Errorln("Unknown task:", req.TaskName)
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
