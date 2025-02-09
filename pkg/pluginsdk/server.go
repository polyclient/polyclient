package pluginsdk

import (
	"context"
	"fmt"
	"log"

	"github.com/polyclient/polyclient/pkg/utils"
	pb "github.com/polyclient/polyclient/proto"
	"google.golang.org/grpc"
)

// pluginServer defines the gRPC server for the plugin.
type pluginServer struct {
	pb.UnimplementedPluginServer
	plugin *Plugin
}

// Register registers the plugin to the gRPC server.
func (s *pluginServer) Register(ctx context.Context, info *pb.PluginInfo) (*pb.PluginInfo, error) {
	return &pb.PluginInfo{
		Name:     s.plugin.Name,
		Version:  s.plugin.Version,
		Actions:  s.plugin.GetActions(),
		Metadata: s.plugin.Metadata,
	}, nil
}

// Execute executes the specified action handler for a plugin.
func (s *pluginServer) Execute(ctx context.Context, req *pb.PluginRequest) (*pb.PluginResponse, error) {
	handler, exists := s.plugin.handlers[req.Action]
	if !exists {
		return &pb.PluginResponse{
			Error: "unknown action: " + req.Action,
		}, nil
	}

	result, err := handler(req.Payload, req.Metadata)
	if err != nil {
		return &pb.PluginResponse{
			Error: err.Error(),
		}, err
	}

	return &pb.PluginResponse{
		Payload: result,
	}, nil
}

// Serve starts a gRPC server for the given Plugin. It determines the appropriate communication
// endpoint based on the operating system. On Unix-like systems, it uses a Unix domain socket,
// which is a special file in the file system used for fast, local inter-process communication (IPC).
// On Windows, it uses a named pipe, which is the Windows equivalent for IPC.
//
// Sockets (Unix domain sockets) and pipes (named pipes on Windows) are used here because they
// provide efficient, local communication without the overhead of TCP/IP networking.
//
// Example usage:
//
//	p := pluginsdk.NewPlugin("myPlugin", "0.0.1")
//	p.RegisterHandler("someAction", func(payload []byte, metadata map[string]string) ([]byte, error) {
//		return []byte("Hello, world!"), nil
//	})
//
//	err := Serve(p)
//	if err != nil {
//		log.Fatalf("Error serving plugin: %v", err)
//	}
func Serve(plugin *Plugin) error {
	socketPath := utils.GetSocketPath(plugin.Name)

	listener, err := createListener(socketPath)
	if err != nil {
		return fmt.Errorf("failed to create listener: %v", err)
	}

	defer func() {
		listener.Close()

		if err := cleanupSocket(socketPath); err != nil {
			log.Printf("Failed to clean up socket: %v", err)
		}
	}()

	server := grpc.NewServer()

	pb.RegisterPluginServer(server, &pluginServer{plugin: plugin})

	log.Printf("Plugin %s v%s starting...", plugin.Name, plugin.Version)

	return server.Serve(listener)
}
