package health

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/longhorn/backing-image-manager/pkg/server"
)

type CheckServer struct {
	pl *server.Manager
}

func NewHealthCheckServer(pl *server.Manager) *CheckServer {
	return &CheckServer{
		pl: pl,
	}
}

func (hc *CheckServer) Check(context.Context, *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	if hc.pl != nil {
		return &healthpb.HealthCheckResponse{
			Status: healthpb.HealthCheckResponse_SERVING,
		}, nil
	}

	return &healthpb.HealthCheckResponse{
		Status: healthpb.HealthCheckResponse_NOT_SERVING,
	}, fmt.Errorf("backing image management service is not running")
}

func (hc *CheckServer) Watch(req *healthpb.HealthCheckRequest, ws healthpb.Health_WatchServer) error {
	for {
		if hc.pl != nil {
			if err := ws.Send(&healthpb.HealthCheckResponse{
				Status: healthpb.HealthCheckResponse_SERVING,
			}); err != nil {
				logrus.Errorf("Failed to send health check result %v for gRPC backing image management server: %v",
					healthpb.HealthCheckResponse_SERVING, err)
			}
		} else {
			if err := ws.Send(&healthpb.HealthCheckResponse{
				Status: healthpb.HealthCheckResponse_NOT_SERVING,
			}); err != nil {
				logrus.Errorf("Failed to send health check result %v for gRPC backing image management server: %v",
					healthpb.HealthCheckResponse_NOT_SERVING, err)
			}

		}
		time.Sleep(time.Second)
	}
}
