// Package notinchain provides a temporary sollution for missing data in the blockchain.
// Remove the THOR node REST URL from the configuration once removed.
package notinchain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// BaseURL defines the REST root.
var BaseURL string

var Client http.Client

// TODO(kashif) we can merge this in future into a better caching layer
// not sure at this point if its necessary or not
var cacheDuration time.Duration = 5 * time.Second
var nodesCache []*NodeAccount //For nodeaccounts
var nodesCachedAt time.Time

type NodeCache struct {
	Node     *NodeAccount
	CachedAt time.Time
}

var nodeCache map[string]*NodeCache = make(map[string]*NodeCache) //For nodeaccount

// Return a cached version of nodeaccounts to reduce load on thorchain nodes
func CachedNodeAccountsLookup() ([]*NodeAccount, error) {
	if nodesCache != nil && time.Now().Before(nodesCachedAt.Add(cacheDuration)) {
		return nodesCache, nil
	}

	newNodes, err := NodeAccountsLookup()
	if err != nil {
		return nil, err
	}

	nodesCache = newNodes
	nodesCachedAt = time.Now()
	return newNodes, err
}

// Return a cached version of nodeaccount to reduce load on thorchain nodes
func CachedNodeAccountLookup(address string) (*NodeAccount, error) {
	c, _ := nodeCache[address]
	if c != nil && time.Now().Before(c.CachedAt.Add(cacheDuration)) {
		return c.Node, nil
	}

	newNode, err := NodeAccountLookup(address)
	if err != nil {
		return nil, err
	}

	nodeCache[address] = &NodeCache{
		Node:     newNode,
		CachedAt: time.Now(),
	}
	return newNode, err
}

type JailInfo struct {
	NodeAddr      string `json:"node_address"`
	ReleaseHeight int64  `json:"release_height,string"`
	Reason        string `json:"reason"`
}

type PublicKeys struct {
	Secp256k1 string `json:"secp256k1"`
	Ed25519   string `json:"ed25519"`
}

type NodeAccount struct {
	NodeAddr         string     `json:"node_address"`
	Status           string     `json:"status"`
	Bond             int64      `json:"bond,string"`
	PublicKeys       PublicKeys `json:"pub_key_set"`
	RequestedToLeave bool       `json:"requested_to_leave"`
	ForcedToLeave    bool       `json:"forced_to_leave"`
	LeaveHeight      int64      `json:"leave_height,string"`
	IpAddress        string     `json:"ip_address"`
	Version          string     `json:"version"`
	SlashPoints      int64      `json:"slash_points,string"`
	Jail             JailInfo   `json:"jail"`
	CurrentAward     int64      `json:"current_award,string"`
}

// Get all nodes from the thorchain api
func NodeAccountsLookup() ([]*NodeAccount, error) {
	resp, err := Client.Get(BaseURL + "/nodeaccounts")
	if err != nil {
		return nil, fmt.Errorf("node accounts unavailable from REST on %w", err)
	}
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("node accounts REST HTTP status %q, want 2xx", resp.Status)
	}
	var accounts []*NodeAccount
	if err := json.NewDecoder(resp.Body).Decode(&accounts); err != nil {
		return nil, fmt.Errorf("node accounts irresolvable from REST on %w", err)
	}
	return accounts, nil
}

// Get node details by address from the thorchain api
func NodeAccountLookup(addr string) (*NodeAccount, error) {
	resp, err := Client.Get(BaseURL + "/nodeaccount/" + addr)
	if err != nil {
		return nil, fmt.Errorf("node account unavailable from REST on %w", err)
	}
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("node account REST HTTP status %q, want 2xx", resp.Status)
	}
	var account *NodeAccount
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, fmt.Errorf("node account irresolvable from REST on %w", err)
	}
	return account, nil
}

type VaultData struct {
	TotalReserve int64 `json:"total_reserve,string"`
}

// Get vault data from the thorchain api
func VaultDataLookup() (*VaultData, error) {
	resp, err := Client.Get(BaseURL + "/vault")
	if err != nil {
		return nil, fmt.Errorf("vault data unavailable from REST on %w", err)
	}
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("vault data REST HTTP status %q, want 2xx", resp.Status)
	}
	var data *VaultData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("vault data irresolvable from REST on %w", err)
	}
	return data, nil
}

type Constants struct {
	Int64Values map[string]int64 `json:"int_64_values"`
}

var constants *Constants
var constantsMu sync.Mutex

// Looks up thorchain constants, query is run once then cached in memory
func ConstantsLookup() (*Constants, error) {
	if constants == nil {
		resp, err := Client.Get(BaseURL + "/constants")
		if err != nil {
			return nil, fmt.Errorf("constants unavailable from REST on %w", err)
		}
		if resp.StatusCode/100 != 2 {
			return nil, fmt.Errorf("constants REST HTTP status %q, want 2xx", resp.Status)
		}
		constantsMu.Lock()
		defer constantsMu.Unlock()
		if err := json.NewDecoder(resp.Body).Decode(&constants); err != nil {
			return nil, fmt.Errorf("constants irresolvable from REST on %w", err)
		}
	}

	return constants, nil
}
