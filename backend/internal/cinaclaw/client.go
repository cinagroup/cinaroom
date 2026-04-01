package cinaclaw

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	pb "github.com/cinagroup/cinaseek/backend/proto/cinaclaw"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// CinaClawClient wraps gRPC communication with the cinaclawd daemon.
// It connects via Unix domain socket and provides high-level Go-friendly APIs
// over the low-level bidirectional streaming gRPC protocol.
type CinaClawClient struct {
	conn       *grpc.ClientConn
	client     pb.RpcClient
	socketPath string
}

// NewClient creates a new gRPC client connected to the CinaClaw daemon via Unix socket.
func NewClient(socketPath string) (*CinaClawClient, error) {
	conn, err := grpc.NewClient(
		"unix://"+socketPath,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			return net.DialTimeout("unix", socketPath, 5*time.Second)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cinaclawd at %s: %w", socketPath, err)
	}

	return &CinaClawClient{
		conn:       conn,
		client:     pb.NewRpcClient(conn),
		socketPath: socketPath,
	}, nil
}

// Close closes the gRPC connection.
func (c *CinaClawClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Ping checks connectivity to the CinaClaw daemon.
func (c *CinaClawClient) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, &pb.PingRequest{})
	return err
}

// Version returns the CinaClaw daemon version.
func (c *CinaClawClient) Version(ctx context.Context) (string, error) {
	stream, err := c.client.Version(ctx)
	if err != nil {
		return "", err
	}
	// Send a request
	if err := stream.Send(&pb.VersionRequest{}); err != nil {
		return "", fmt.Errorf("send version request: %w", err)
	}
	stream.CloseSend()

	// Read the reply
	reply, err := stream.Recv()
	if err != nil {
		return "", fmt.Errorf("recv version reply: %w", err)
	}
	return reply.GetVersion(), nil
}

// CreateVM creates (launches) a new virtual machine instance.
func (c *CinaClawClient) CreateVM(ctx context.Context, req *CreateVMRequest) (*VMInfo, error) {
	stream, err := c.client.Create(ctx)
	if err != nil {
		return nil, fmt.Errorf("create stream: %w", err)
	}

	// Send launch request
	protoReq := &pb.LaunchRequest{
		InstanceName: req.Name,
		Image:        req.Image,
		NumCores:     int32(req.CPU),
		MemSize:      req.Memory,
		DiskSpace:    req.Disk,
	}
	if err := stream.Send(protoReq); err != nil {
		return nil, fmt.Errorf("send create request: %w", err)
	}
	stream.CloseSend()

	// Collect replies - wait for the instance name or error
	var vmName string
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("recv create reply: %w", err)
		}
		if name := reply.GetVmInstanceName(); name != "" {
			vmName = name
		}
	}

	if vmName == "" {
		vmName = req.Name
	}

	// Fetch full info for the newly created VM
	return c.GetVMInfo(ctx, vmName)
}

// ListVMs returns all virtual machine instances.
func (c *CinaClawClient) ListVMs(ctx context.Context) ([]*VMInfo, error) {
	stream, err := c.client.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list stream: %w", err)
	}

	if err := stream.Send(&pb.ListRequest{RequestIpv4: true}); err != nil {
		return nil, fmt.Errorf("send list request: %w", err)
	}
	stream.CloseSend()

	var vms []*VMInfo
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("recv list reply: %w", err)
		}

		if il := reply.GetInstanceList(); il != nil {
			for _, inst := range il.GetInstances() {
				ip := ""
				if len(inst.GetIpv4()) > 0 {
					ip = inst.GetIpv4()[0]
				}
				status := statusProtoToString(inst.GetInstanceStatus().GetStatus())
				vms = append(vms, &VMInfo{
					Name:    inst.GetName(),
					Status:  status,
					IP:      ip,
					OS:      inst.GetOs(),
					Release: inst.GetCurrentRelease(),
				})
			}
		}
	}

	return vms, nil
}

// GetVMInfo returns detailed information about a specific VM.
func (c *CinaClawClient) GetVMInfo(ctx context.Context, name string) (*VMInfo, error) {
	stream, err := c.client.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("info stream: %w", err)
	}

	if err := stream.Send(&pb.InfoRequest{
		InstanceSnapshotPairs: []*pb.InstanceSnapshotPair{
			{InstanceName: name},
		},
	}); err != nil {
		return nil, fmt.Errorf("send info request: %w", err)
	}
	stream.CloseSend()

	reply, err := stream.Recv()
	if err != nil {
		return nil, fmt.Errorf("recv info reply: %w", err)
	}

	details := reply.GetDetails()
	if len(details) == 0 {
		return nil, fmt.Errorf("VM %q not found", name)
	}

	d := details[0]
	info := &VMInfo{
		Name:   d.GetName(),
		Status: statusProtoToString(d.GetInstanceStatus().GetStatus()),
		CPU:    int(parseCPUCount(d.GetCpuCount())),
		Memory: d.GetMemoryTotal(),
		Disk:   d.GetDiskTotal(),
	}

	if ii := d.GetInstanceInfo(); ii != nil {
		info.OS = ii.GetOs()
		info.Release = ii.GetImageRelease()
		if len(ii.GetIpv4()) > 0 {
			info.IP = ii.GetIpv4()[0]
		}
		if ts := ii.GetCreationTimestamp(); ts != nil {
			info.CreatedAt = ts.AsTime()
		}
	}

	return info, nil
}

// StartVM starts a stopped VM instance.
func (c *CinaClawClient) StartVM(ctx context.Context, name string) error {
	stream, err := c.client.Start(ctx)
	if err != nil {
		return fmt.Errorf("start stream: %w", err)
	}

	if err := stream.Send(&pb.StartRequest{
		InstanceNames: &pb.InstanceNames{InstanceName: []string{name}},
	}); err != nil {
		return fmt.Errorf("send start request: %w", err)
	}
	stream.CloseSend()

	// Consume replies until done
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("recv start reply: %w", err)
		}
	}

	return nil
}

// StopVM stops a running VM instance.
func (c *CinaClawClient) StopVM(ctx context.Context, name string) error {
	stream, err := c.client.Stop(ctx)
	if err != nil {
		return fmt.Errorf("stop stream: %w", err)
	}

	if err := stream.Send(&pb.StopRequest{
		InstanceNames: &pb.InstanceNames{InstanceName: []string{name}},
	}); err != nil {
		return fmt.Errorf("send stop request: %w", err)
	}
	stream.CloseSend()

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("recv stop reply: %w", err)
		}
	}

	return nil
}

// RestartVM restarts a VM instance.
func (c *CinaClawClient) RestartVM(ctx context.Context, name string) error {
	stream, err := c.client.Restart(ctx)
	if err != nil {
		return fmt.Errorf("restart stream: %w", err)
	}

	if err := stream.Send(&pb.RestartRequest{
		InstanceNames: &pb.InstanceNames{InstanceName: []string{name}},
	}); err != nil {
		return fmt.Errorf("send restart request: %w", err)
	}
	stream.CloseSend()

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("recv restart reply: %w", err)
		}
	}

	return nil
}

// DeleteVM deletes a VM instance.
func (c *CinaClawClient) DeleteVM(ctx context.Context, name string) error {
	stream, err := c.client.Delet(ctx)
	if err != nil {
		return fmt.Errorf("delete stream: %w", err)
	}

	if err := stream.Send(&pb.DeleteRequest{
		InstanceSnapshotPairs: []*pb.InstanceSnapshotPair{
			{InstanceName: name},
		},
		Purge: true,
	}); err != nil {
		return fmt.Errorf("send delete request: %w", err)
	}
	stream.CloseSend()

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("recv delete reply: %w", err)
		}
	}

	return nil
}

// SuspendVM suspends a running VM instance.
func (c *CinaClawClient) SuspendVM(ctx context.Context, name string) error {
	stream, err := c.client.Suspend(ctx)
	if err != nil {
		return fmt.Errorf("suspend stream: %w", err)
	}

	if err := stream.Send(&pb.SuspendRequest{
		InstanceNames: &pb.InstanceNames{InstanceName: []string{name}},
	}); err != nil {
		return fmt.Errorf("send suspend request: %w", err)
	}
	stream.CloseSend()

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("recv suspend reply: %w", err)
		}
	}

	return nil
}

// SnapshotVM creates a snapshot of a VM.
func (c *CinaClawClient) SnapshotVM(ctx context.Context, name, snapshotName string) error {
	stream, err := c.client.Snapshot(ctx)
	if err != nil {
		return fmt.Errorf("snapshot stream: %w", err)
	}

	if err := stream.Send(&pb.SnapshotRequest{
		Instance: name,
		Snapshot: snapshotName,
	}); err != nil {
		return fmt.Errorf("send snapshot request: %w", err)
	}
	stream.CloseSend()

	reply, err := stream.Recv()
	if err != nil && err != io.EOF {
		return fmt.Errorf("recv snapshot reply: %w", err)
	}

	// If no snapshot name was provided, the daemon generates one
	if reply != nil && reply.GetSnapshot() != "" {
		_ = reply.GetSnapshot() // Caller can use the returned name if needed
	}

	return nil
}

// RestoreVM restores a VM from a snapshot.
func (c *CinaClawClient) RestoreVM(ctx context.Context, name, snapshotName string) error {
	stream, err := c.client.Restore(ctx)
	if err != nil {
		return fmt.Errorf("restore stream: %w", err)
	}

	if err := stream.Send(&pb.RestoreRequest{
		Instance:   name,
		Snapshot:   snapshotName,
		Destructive: false,
	}); err != nil {
		return fmt.Errorf("send restore request: %w", err)
	}
	stream.CloseSend()

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("recv restore reply: %w", err)
		}
	}

	return nil
}

// DeleteSnapshot deletes a snapshot of a VM.
func (c *CinaClawClient) DeleteSnapshot(ctx context.Context, name, snapshotName string) error {
	stream, err := c.client.Delet(ctx)
	if err != nil {
		return fmt.Errorf("delete stream: %w", err)
	}

	if err := stream.Send(&pb.DeleteRequest{
		InstanceSnapshotPairs: []*pb.InstanceSnapshotPair{
			{
				InstanceName: name,
				SnapshotName: &snapshotName,
			},
		},
		PurgeSnapshots: true,
	}); err != nil {
		return fmt.Errorf("send delete snapshot request: %w", err)
	}
	stream.CloseSend()

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("recv delete snapshot reply: %w", err)
		}
	}

	return nil
}

// CloneVM clones an existing VM instance.
func (c *CinaClawClient) CloneVM(ctx context.Context, sourceName, destName string) error {
	stream, err := c.client.Clone(ctx)
	if err != nil {
		return fmt.Errorf("clone stream: %w", err)
	}

	req := &pb.CloneRequest{
		SourceName: sourceName,
	}
	if destName != "" {
		req.DestinationName = &destName
	}

	if err := stream.Send(req); err != nil {
		return fmt.Errorf("send clone request: %w", err)
	}
	stream.CloseSend()

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("recv clone reply: %w", err)
		}
	}

	return nil
}

// Mount mounts a host directory into a VM.
func (c *CinaClawClient) Mount(ctx context.Context, hostPath string, targets []MountTarget) error {
	stream, err := c.client.Mount(ctx)
	if err != nil {
		return fmt.Errorf("mount stream: %w", err)
	}

	pbTargets := make([]*pb.TargetPathInfo, len(targets))
	for i, t := range targets {
		pbTargets[i] = &pb.TargetPathInfo{
			InstanceName: t.InstanceName,
			TargetPath:   t.TargetPath,
		}
	}

	if err := stream.Send(&pb.MountRequest{
		SourcePath:  hostPath,
		TargetPaths: pbTargets,
	}); err != nil {
		return fmt.Errorf("send mount request: %w", err)
	}
	stream.CloseSend()

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("recv mount reply: %w", err)
		}
	}

	return nil
}

// Unmount unmounts a directory from a VM.
func (c *CinaClawClient) Unmount(ctx context.Context, targets []MountTarget) error {
	stream, err := c.client.Umount(ctx)
	if err != nil {
		return fmt.Errorf("umount stream: %w", err)
	}

	pbTargets := make([]*pb.TargetPathInfo, len(targets))
	for i, t := range targets {
		pbTargets[i] = &pb.TargetPathInfo{
			InstanceName: t.InstanceName,
			TargetPath:   t.TargetPath,
		}
	}

	if err := stream.Send(&pb.UmountRequest{
		TargetPaths: pbTargets,
	}); err != nil {
		return fmt.Errorf("send umount request: %w", err)
	}
	stream.CloseSend()

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("recv umount reply: %w", err)
		}
	}

	return nil
}

// GetMetrics retrieves resource usage metrics for a VM.
func (c *CinaClawClient) GetMetrics(ctx context.Context, name string) (*VMMetrics, error) {
	stream, err := c.client.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("info stream: %w", err)
	}

	if err := stream.Send(&pb.InfoRequest{
		InstanceSnapshotPairs: []*pb.InstanceSnapshotPair{
			{InstanceName: name},
		},
	}); err != nil {
		return nil, fmt.Errorf("send info request: %w", err)
	}
	stream.CloseSend()

	reply, err := stream.Recv()
	if err != nil {
		return nil, fmt.Errorf("recv info reply: %w", err)
	}

	details := reply.GetDetails()
	if len(details) == 0 {
		return nil, fmt.Errorf("VM %q not found", name)
	}

	d := details[0]
	metrics := &VMMetrics{}

	if ii := d.GetInstanceInfo(); ii != nil {
		metrics.CPUUsage = parseFloat(ii.GetCpuTimes())
		metrics.MemoryUsage = parseFloat(ii.GetMemoryUsage())
		metrics.DiskUsage = parseFloat(ii.GetDiskUsage())
	}

	return metrics, nil
}

// GetSSHInfo returns SSH connection information for a VM.
func (c *CinaClawClient) GetSSHInfo(ctx context.Context, names []string) (map[string]*SSHInfo, error) {
	stream, err := c.client.SshInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("ssh_info stream: %w", err)
	}

	if err := stream.Send(&pb.SSHInfoRequest{InstanceName: names}); err != nil {
		return nil, fmt.Errorf("send ssh_info request: %w", err)
	}
	stream.CloseSend()

	reply, err := stream.Recv()
	if err != nil {
		return nil, fmt.Errorf("recv ssh_info reply: %w", err)
	}

	result := make(map[string]*SSHInfo)
	for name, info := range reply.GetSshInfo() {
		result[name] = &SSHInfo{
			Port:     int(info.GetPort()),
			Host:     info.GetHost(),
			Username: info.GetUsername(),
			Key:      info.GetPrivKeyBase64(),
		}
	}

	return result, nil
}

// FindImages searches for available VM images.
func (c *CinaClawClient) FindImages(ctx context.Context, searchString string) ([]*ImageInfo, error) {
	stream, err := c.client.Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("find stream: %w", err)
	}

	if err := stream.Send(&pb.FindRequest{SearchString: searchString}); err != nil {
		return nil, fmt.Errorf("send find request: %w", err)
	}
	stream.CloseSend()

	reply, err := stream.Recv()
	if err != nil {
		return nil, fmt.Errorf("recv find reply: %w", err)
	}

	var images []*ImageInfo
	for _, img := range reply.GetImagesInfo() {
		images = append(images, &ImageInfo{
			OS:       img.GetOs(),
			Release:  img.GetRelease(),
			Version:  img.GetVersion(),
			Aliases:  img.GetAliases(),
			Codename: img.GetCodename(),
			Remote:   img.GetRemoteName(),
		})
	}

	return images, nil
}

// Recover attempts to recover VMs that are in an error state.
func (c *CinaClawClient) Recover(ctx context.Context, names []string) error {
	stream, err := c.client.Recover(ctx)
	if err != nil {
		return fmt.Errorf("recover stream: %w", err)
	}

	if err := stream.Send(&pb.RecoverRequest{
		InstanceNames: &pb.InstanceNames{InstanceName: names},
	}); err != nil {
		return fmt.Errorf("send recover request: %w", err)
	}
	stream.CloseSend()

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("recv recover reply: %w", err)
		}
	}

	return nil
}

// WaitReady waits until the CinaClaw daemon is ready to accept requests.
func (c *CinaClawClient) WaitReady(ctx context.Context) error {
	stream, err := c.client.WaitReady(ctx)
	if err != nil {
		return fmt.Errorf("wait_ready stream: %w", err)
	}

	if err := stream.Send(&pb.WaitReadyRequest{}); err != nil {
		return fmt.Errorf("send wait_ready request: %w", err)
	}
	stream.CloseSend()

	_, err = stream.Recv()
	if err != nil && err != io.EOF {
		return fmt.Errorf("recv wait_ready reply: %w", err)
	}

	return nil
}

// Helper functions

// statusProtoToString converts a protobuf InstanceStatus_Status to a human-readable string.
func statusProtoToString(status pb.InstanceStatus_Status) string {
	switch status {
	case pb.InstanceStatus_RUNNING:
		return "RUNNING"
	case pb.InstanceStatus_STARTING:
		return "STARTING"
	case pb.InstanceStatus_RESTARTING:
		return "RESTARTING"
	case pb.InstanceStatus_STOPPED:
		return "STOPPED"
	case pb.InstanceStatus_DELETED:
		return "DELETED"
	case pb.InstanceStatus_DELAYED_SHUTDOWN:
		return "DELAYED_SHUTDOWN"
	case pb.InstanceStatus_SUSPENDING:
		return "SUSPENDING"
	case pb.InstanceStatus_SUSPENDED:
		return "SUSPENDED"
	default:
		return "UNKNOWN"
	}
}

// parseCPUCount converts a CPU count string from the proto response to an integer.
func parseCPUCount(s string) int32 {
	var v int32
	fmt.Sscanf(s, "%d", &v)
	if v == 0 {
		v = 1
	}
	return v
}

// parseFloat parses a percentage/usage string from the daemon.
func parseFloat(s string) float64 {
	var v float64
	fmt.Sscanf(s, "%f", &v)
	return v
}
