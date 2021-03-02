package diff

import "github.com/getkin/kin-openapi/openapi3"

// ServersDiff is a diff between two sets of encoding objects: https://swagger.io/specification/#server-object
type ServersDiff struct {
	Added    StringList      `json:"added,omitempty"`
	Deleted  StringList      `json:"deleted,omitempty"`
	Modified ModifiedServers `json:"modified,omitempty"`
}

// ModifiedServers is map of server names to their respective diffs
type ModifiedServers map[string]ServerDiff

func (diff *ServersDiff) empty() bool {
	return len(diff.Added) == 0 &&
		len(diff.Deleted) == 0 &&
		len(diff.Modified) == 0
}

func newServersDiff() *ServersDiff {
	return &ServersDiff{
		Added:    StringList{},
		Deleted:  StringList{},
		Modified: ModifiedServers{},
	}
}

func getServersDiff(config *Config, pServers1, pServers2 *openapi3.Servers) *ServersDiff {

	result := newServersDiff()

	servers1 := derefServers(pServers1)
	servers2 := derefServers(pServers2)

	for _, server1 := range servers1 {
		if server2 := findServer(server1, servers2); server2 != nil {
			if diff := getServerDiff(config, server1, server2); !diff.empty() {
				result.Modified[server1.URL] = diff
			}
		} else {
			result.Deleted = append(result.Deleted, server1.URL)
		}
	}

	for _, server2 := range servers2 {
		if server1 := findServer(server2, servers1); server1 == nil {
			result.Added = append(result.Added, server2.URL)
		}
	}

	return result
}

func derefServers(servers *openapi3.Servers) openapi3.Servers {
	if servers == nil {
		return openapi3.Servers{}
	}

	return *servers
}

func findServer(server1 *openapi3.Server, servers2 openapi3.Servers) *openapi3.Server {
	// TODO: optimize with a map
	for _, server2 := range servers2 {
		if server2.URL == server1.URL {
			return server2
		}
	}

	return nil
}

func (diff *ServersDiff) getSummary() *SummaryDetails {
	return &SummaryDetails{
		Added:    len(diff.Added),
		Deleted:  len(diff.Deleted),
		Modified: len(diff.Modified),
	}
}
