package hashring

import (
	"fmt"
	"slices"
	"testing"
)

func TestNewHashRing(t *testing.T) {
	nodes := []string{"node1", "node2", "node3"}
	replicas := 3

	hr := NewHashRing(nodes, replicas)

	if len(hr.NodeAddrs) != len(nodes) {
		t.Errorf("Expected %d nodes, got %d", len(nodes), len(hr.NodeAddrs))
	}

	if hr.Replicas != replicas {
		t.Errorf("Expected %d replicas, got %d", replicas, hr.Replicas)
	}

	expectedHashCount := len(nodes) * replicas
	if len(hr.NodesHashList) != expectedHashCount {
		t.Errorf("Expected %d hashes, got %d", expectedHashCount, len(hr.NodesHashList))
	}

	if len(hr.NodeHashToAddr) != expectedHashCount {
		t.Errorf("Expected %d mappings, got %d", expectedHashCount, len(hr.NodeHashToAddr))
	}
}

func TestHashRing_AddNode(t *testing.T) {
	hr := NewHashRing([]string{"node1", "node2"}, 3)
	initialHashCount := len(hr.NodesHashList)

	hr.AddNode("node3")

	if len(hr.NodeAddrs) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(hr.NodeAddrs))
	}

	expectedHashCount := initialHashCount + 3
	if len(hr.NodesHashList) != expectedHashCount {
		t.Errorf("Expected %d hashes, got %d", expectedHashCount, len(hr.NodesHashList))
	}

	if len(hr.NodeHashToAddr) != expectedHashCount {
		t.Errorf("Expected %d mappings, got %d", expectedHashCount, len(hr.NodeHashToAddr))
	}
}

func TestHashRing_RemoveNode(t *testing.T) {
	hr := NewHashRing([]string{"node1", "node2", "node3"}, 3)
	initialHashCount := len(hr.NodesHashList)

	hr.RemoveNode("node2")

	if len(hr.NodeAddrs) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(hr.NodeAddrs))
	}

	expectedHashCount := initialHashCount - 3
	if len(hr.NodesHashList) != expectedHashCount {
		t.Errorf("Expected %d hashes, got %d", expectedHashCount, len(hr.NodesHashList))
	}

	if len(hr.NodeHashToAddr) != expectedHashCount {
		t.Errorf("Expected %d mappings, got %d", expectedHashCount, len(hr.NodeHashToAddr))
	}
}

func TestHashRing_findNextBiggestHash(t *testing.T) {
	hr := NewHashRing([]string{"node1", "node2"}, 1)

	tests := []struct {
		name     string
		lookup   uint32
		expected uint32
	}{
		{"Exact match", hr.NodesHashList[0], hr.NodesHashList[1]},
		{"Between hashes", (hr.NodesHashList[0] + hr.NodesHashList[1]) / 2, hr.NodesHashList[1]},
		{"Wrap around", hr.NodesHashList[1] + 1, hr.NodesHashList[0]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hr.findNextBiggestHash(tt.lookup)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestHashRing_GetNodeAddrForKey(t *testing.T) {
	nodes := []string{"node1", "node2", "node3"}
	hr := NewHashRing(nodes, 1)

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key%d", i)
		node := hr.GetNodeAddrForKey(key)
		if !slices.Contains(nodes, node) {
			t.Errorf("Got unexpected node %s for key %s", node, key)
		}
	}
}
