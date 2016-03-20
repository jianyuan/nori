package nori

import (
	"expvar"
	"net/http"
)

func (s *Server) RunManagementServer(addr string) {
	log.Infoln("Management server listening on", addr)
	s.setupExpvar()
	go http.ListenAndServe(addr, nil)
}

func (s *Server) setupExpvar() {
	parentMap := expvar.NewMap("nori")

	parentMap.Set("Tasks", expvar.Func(func() interface{} {
		tasks := make([]map[string]interface{}, 0, len(s.Tasks))
		for _, task := range s.Tasks {
			tasks = append(tasks, map[string]interface{}{
				"Name": task.Name,
			})
		}
		return tasks
	}))
}
