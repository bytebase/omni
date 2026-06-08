package parser

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type dockerEndpoint struct {
	network string
	address string
}

func dockerPreflight(timeout time.Duration) error {
	home, _ := os.UserHomeDir()
	endpoints := dockerEndpointCandidates(os.Getenv("DOCKER_HOST"), home)
	if len(endpoints) == 0 {
		return fmt.Errorf("no supported Docker endpoint configured")
	}

	var errs []string
	for _, endpoint := range endpoints {
		if err := pingDockerEndpoint(endpoint, timeout); err != nil {
			errs = append(errs, fmt.Sprintf("%s:%s: %v", endpoint.network, endpoint.address, err))
			continue
		}
		return nil
	}
	return errors.New(strings.Join(errs, "; "))
}

func dockerEndpointCandidates(dockerHost, home string) []dockerEndpoint {
	if endpoint, ok := parseDockerHostEndpoint(dockerHost); ok {
		return []dockerEndpoint{endpoint}
	}

	endpoints := []dockerEndpoint{{network: "unix", address: "/var/run/docker.sock"}}
	if home != "" {
		endpoints = append(endpoints, dockerEndpoint{
			network: "unix",
			address: filepath.Join(home, ".docker", "run", "docker.sock"),
		})
	}
	return endpoints
}

func parseDockerHostEndpoint(dockerHost string) (dockerEndpoint, bool) {
	if dockerHost == "" {
		return dockerEndpoint{}, false
	}

	parsed, err := url.Parse(dockerHost)
	if err != nil {
		return dockerEndpoint{}, false
	}
	switch parsed.Scheme {
	case "unix":
		if parsed.Path == "" {
			return dockerEndpoint{}, false
		}
		return dockerEndpoint{network: "unix", address: parsed.Path}, true
	case "tcp":
		if parsed.Host == "" {
			return dockerEndpoint{}, false
		}
		return dockerEndpoint{network: "tcp", address: parsed.Host}, true
	default:
		return dockerEndpoint{}, false
	}
}

func pingDockerEndpoint(endpoint dockerEndpoint, timeout time.Duration) error {
	conn, err := (&net.Dialer{Timeout: timeout}).Dial(endpoint.network, endpoint.address)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}
	if _, err := io.WriteString(conn, "GET /_ping HTTP/1.1\r\nHost: docker\r\n\r\n"); err != nil {
		return err
	}

	buf := make([]byte, 128)
	n, err := conn.Read(buf)
	if err != nil {
		return err
	}
	if !strings.Contains(string(buf[:n]), "OK") {
		return fmt.Errorf("unexpected ping response: %q", strings.TrimSpace(string(buf[:n])))
	}
	return nil
}
