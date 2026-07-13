package routeguide

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"os"
	"sync"

	routeguidepb "github.com/agentio/routeguide-sidecar/genproto/routeguidepb"
	"github.com/agentio/routeguide-sidecar/internal/routeguide/constants"
	"github.com/agentio/sidecar"
)

type Server struct {
	savedFeatures []*routeguidepb.Feature
	mu            sync.Mutex // protects routeNotes
	routeNotes    map[string][]*routeguidepb.RouteNote
}

func NewServer(filepath string) (*Server, error) {
	s := &Server{}
	s.savedFeatures = []*routeguidepb.Feature{}
	s.routeNotes = map[string][]*routeguidepb.RouteNote{}
	err := s.loadFeatures(filepath)
	return s, err
}

//go:embed route_guide_db.json
var testdata []byte

func (s *Server) loadFeatures(filePath string) error {
	data := testdata
	if filePath != "" {
		var err error
		if data, err = os.ReadFile(filePath); err != nil {
			return err
		}
	}
	return json.Unmarshal(data, &s.savedFeatures)
}

func (s *Server) ServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(constants.RouteGuideGetFeatureProcedure, sidecar.HandleUnary(s.getFeature))
	mux.HandleFunc(constants.RouteGuideListFeaturesProcedure, sidecar.HandleServerStreaming(s.listFeatures))
	mux.HandleFunc(constants.RouteGuideRecordRouteProcedure, sidecar.HandleClientStreaming(s.recordRoute))
	mux.HandleFunc(constants.RouteGuideRouteChatProcedure, sidecar.HandleBidiStreaming(s.routeChat))
	return mux
}
