package nori

import (
	"expvar"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

var (
	registeredTasks = expvar.NewMap("tasks")
)

func (s *Server) RunManagementServer(addr string) {
	log.Infoln("Management server started on ", addr)

	http.ListenAndServe(addr, nil)
}
