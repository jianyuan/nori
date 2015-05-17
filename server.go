package nori

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

type Server struct {
	Name  string
	Tasks []*Task
}

func NewServer(name string) *Server {
	return &Server{
		Name: name,
	}
}

func (s *Server) RegisterTask(t *Task) {
	// TODO: validation
	t.Name = s.Name + "." + t.Name
	s.Tasks = append(s.Tasks, t)
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
	go s.run()

	return nil
}

func (s *Server) run() error {
	s.printInfo()

	return nil
}
