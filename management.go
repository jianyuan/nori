package nori

import (
	"expvar"
	"net/http"

	"github.com/jianyuan/nori/log"
)

func (s *Server) RunManagementServer(addr string) {
	log.FromContext(s).Infoln("Management server listening on", addr)
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
