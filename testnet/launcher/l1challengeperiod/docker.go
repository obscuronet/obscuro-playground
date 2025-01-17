package l1grantsequencers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/ten-protocol/go-ten/go/common/docker"
)

type SetChallengePeriod struct {
	cfg         *Config
	containerID string
}

func NewSetChallengePeriod(cfg *Config) (*SetChallengePeriod, error) {
	return &SetChallengePeriod{
		cfg: cfg,
	}, nil
}

func (s *SetChallengePeriod) Start() error {
	var err error
	cmds := []string{
		"npx",
		"hardhat",
		"run",
		"--network",
		"layer1",
		"scripts/delay/001_set_challenge_period.ts",
	}

	envs := map[string]string{
		"NETWORK_JSON": fmt.Sprintf(`{ 
            "layer1": {
                "url": "%s",
                "live": false,
                "saveDeployments": true,
                "accounts": [ "%s" ]
            }
        }`, s.cfg.l1HTTPURL, s.cfg.privateKey),
		"MGMT_CONTRACT_ADDRESS": s.cfg.mgmtContractAddress,
		"L1_CHALLENGE_PERIOD":   strconv.Itoa(s.cfg.challengePeriod),
	}

	containerID, err := docker.StartNewContainer(
		"set-challenge-period",
		s.cfg.dockerImage,
		cmds,
		nil,
		envs,
		nil,
		nil,
		false,
	)
	if err != nil {
		return err
	}
	s.containerID = containerID
	return nil
}

func (s *SetChallengePeriod) WaitForFinish() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}
	defer cli.Close()

	// make sure the container has finished execution
	err = docker.WaitForContainerToFinish(s.containerID, 15*time.Minute)
	if err != nil {
		fmt.Println("Error waiting for container to finish: ", err)
		s.PrintLogs(cli)
		return err
	}

	return nil
}

func (s *SetChallengePeriod) PrintLogs(cli *client.Client) {
	logsOptions := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	}

	// Read the container logs
	out, err := cli.ContainerLogs(context.Background(), s.containerID, logsOptions)
	if err != nil {
		fmt.Printf("Error printing out container %s logs... %v\n", s.containerID, err)
	}
	defer out.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, out)
	if err != nil {
		fmt.Printf("Error getting logs for container %s\n", s.containerID)
	}
	fmt.Println(buf.String())
}
