package code

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func RunInDocker(translator, code string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return "", fmt.Errorf("cli docker: %s", err.Error())
	}
	defer cli.Close()

	filename, run := getFileNameRun(translator)
	if filename == "" {
		return "", fmt.Errorf("translator not found: %s", translator)
	}

	resp, err := cli.ContainerCreate(ctx,
		client.ContainerCreateOptions{
			Name: "temp",
			Config: &container.Config{
				Image:      "code-processor:latest",
				Cmd:        []string{"sh", "-c", run},
				WorkingDir: "/code",
			},
			HostConfig: &container.HostConfig{
				AutoRemove: true,
			},
			NetworkingConfig: nil,
			Platform:         nil,
		})
	if err != nil {
		return "", fmt.Errorf("container create: %s", err.Error())
	}
	defer func() {
		_, _ = cli.ContainerRemove(ctx, resp.ID, client.ContainerRemoveOptions{Force: true})
	}()

	tarContent, err := makeTarArchive(filename, code)
	if err != nil {
		return "", fmt.Errorf("make tar: %w", err)
	}

	if err = copyCodeToContainer(ctx, cli, resp.ID, tarContent); err != nil {
		return "", fmt.Errorf("copy to container: %s, %s", resp.ID, err.Error())
	}

	if _, err = cli.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("start to container: %s, %s", resp.ID, err.Error())
	}

	statusCh := cli.ContainerWait(ctx, resp.ID, client.ContainerWaitOptions{Condition: container.WaitConditionNotRunning})

	select {
	case err = <-statusCh.Error:
		return "", fmt.Errorf("error waiting: %s", err.Error())
	case <-statusCh.Result:
	case <-ctx.Done():
		return "", fmt.Errorf("timeout context done")
	}

	output, err := cli.ContainerLogs(ctx, resp.ID, client.ContainerLogsOptions{ShowStderr: true, ShowStdout: true})
	if err != nil {
		return "", fmt.Errorf("output: %s", err.Error())
	}
	defer output.Close()

	var stdout, stderr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdout, &stderr, output)
	if err != nil {
		return "", fmt.Errorf("parse output: %w", err)
	}

	if stderr.Len() > 0 {
		return stdout.String(), fmt.Errorf("stderr: %s", stderr.String())
	}
	return stdout.String(), nil
}

func getFileNameRun(translator string) (string, string) {
	switch translator {
	case "python3":
		return "user.py", "python3 /code/user.py"
	case "gcc":
		return "user.c", "gcc /code/user.c -o /code/user && /code/user"
	case "clang":
		return "user.c", "clang /code/user.c -o /code/user && /code/user"
	default:
		return "", ""
	}
}

func copyCodeToContainer(ctx context.Context, cli *client.Client, containerID string, content io.Reader) error {
	_, err := cli.CopyToContainer(
		ctx,
		containerID,
		client.CopyToContainerOptions{
			DestinationPath:           "/code",
			Content:                   content,
			AllowOverwriteDirWithFile: true,
			CopyUIDGID:                false,
		},
	)
	return err
}

func makeTarArchive(filename, code string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	hdr := &tar.Header{
		Name: filename,
		Mode: 0o644,
		Size: int64(len(code)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return nil, err
	}
	if _, err := tw.Write([]byte(code)); err != nil {
		return nil, err
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	return bytes.NewReader(buf.Bytes()), nil
}
