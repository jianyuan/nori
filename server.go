package nori

import (
	"os"
)

type Server struct {
	Name  string
	Tasks map[string]*Task
}

func NewServer(name string) *Server {
	return &Server{
		Name:  name,
		Tasks: make(map[string]*Task),
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
	go s.run()

	go s.RunManagementServer(":8080")

	return nil
}

func (s *Server) run() error {
	s.printInfo()

	return nil
}
