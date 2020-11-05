package executor

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"sync"
)

type DockerExecutor struct {
	cli         *client.Client
	runnerImage string

	runningContainers     map[string]string
	runningContainersLock sync.Mutex
}

func NewDockerExecutor(runnerImage string) (*DockerExecutor, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, errors.Wrap(err, "creating docker env client")
	}

	return &DockerExecutor{cli, runnerImage, make(map[string]string), sync.Mutex{}}, nil
}

func (d *DockerExecutor) Close() error {
	return errors.Wrap(d.cli.Close(), "closing docker cli")
}

func (d *DockerExecutor) PrepareBuild(ctx context.Context) error {
	_, err := d.cli.ImagePull(ctx, d.runnerImage, types.ImagePullOptions{})
	if err != nil {
		return errors.Wrap(err, "docker pull")
	}

	return nil
}

func (d *DockerExecutor) BuildPackage(ctx context.Context, cfg *Config) error {
	cfgb, err := json.Marshal(cfg)
	if err != nil {
		return errors.Wrap(err, "marshal config")
	}

	resp, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image:        d.runnerImage,
		Tty:          true,
		AttachStdout: true,
		Env: []string{
			"CONFIG=" + string(cfgb),
		},
	}, &container.HostConfig{
		AutoRemove: true,
	}, nil, "AAAAA-"+cfg.Package.Name)
	if err != nil {
		return errors.Wrap(err, "docker container create")
	}

	if err := d.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return errors.Wrap(err, "docker run")
	}

	d.runningContainersLock.Lock()
	d.runningContainers[cfg.Package.Name] = resp.ID
	d.runningContainersLock.Unlock()

	// Remove from map when done
	go func() {
		statusCh, errCh := d.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

		select {
		case err := <-errCh:
			if err != nil {
				log.Warnf("[docker] waiting for container failed: %v", err)
			}
		case <-statusCh:
		}

		// TODO: Check status

		d.runningContainersLock.Lock()
		delete(d.runningContainers, cfg.Package.Name)
		d.runningContainersLock.Unlock()
	}()

	return nil
}
