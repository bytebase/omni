package parser

import "testing"

func TestRedshiftParserDockerEndpointCandidatesPreferExplicitHost(t *testing.T) {
	endpoints := dockerEndpointCandidates("unix:///tmp/docker.sock", "/home/tester")

	if len(endpoints) != 1 {
		t.Fatalf("expected one endpoint, got %d", len(endpoints))
	}
	if endpoints[0] != (dockerEndpoint{network: "unix", address: "/tmp/docker.sock"}) {
		t.Fatalf("unexpected endpoint: %#v", endpoints[0])
	}
}

func TestRedshiftParserDockerEndpointCandidatesIncludeDefaultSockets(t *testing.T) {
	endpoints := dockerEndpointCandidates("", "/home/tester")

	want := []dockerEndpoint{
		{network: "unix", address: "/var/run/docker.sock"},
		{network: "unix", address: "/home/tester/.docker/run/docker.sock"},
	}
	if len(endpoints) != len(want) {
		t.Fatalf("expected %d endpoints, got %d: %#v", len(want), len(endpoints), endpoints)
	}
	for i := range want {
		if endpoints[i] != want[i] {
			t.Fatalf("endpoint %d: got %#v, want %#v", i, endpoints[i], want[i])
		}
	}
}

func TestRedshiftParserParseDockerHostEndpoint(t *testing.T) {
	tests := []struct {
		name string
		host string
		want dockerEndpoint
		ok   bool
	}{
		{
			name: "unix",
			host: "unix:///Users/test/.docker/run/docker.sock",
			want: dockerEndpoint{network: "unix", address: "/Users/test/.docker/run/docker.sock"},
			ok:   true,
		},
		{
			name: "tcp",
			host: "tcp://127.0.0.1:2375",
			want: dockerEndpoint{network: "tcp", address: "127.0.0.1:2375"},
			ok:   true,
		},
		{
			name: "unsupported",
			host: "ssh://docker.example.com",
			ok:   false,
		},
		{
			name: "empty",
			host: "",
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseDockerHostEndpoint(tt.host)
			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("endpoint = %#v, want %#v", got, tt.want)
			}
		})
	}
}
