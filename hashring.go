package main

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"slices"
)

func NewHashRing(nodes []string, replicas int) HashRing {
	nodeHashToAddrs := make(map[uint32]string)
	nodeHashes := make([]uint32, len(nodes)*replicas)
	count := 0
	for _, nodeAddr := range nodes {
		for replica := range replicas {
			virtualNodeHash := generateKeyHash(fmt.Sprintf("%s-%d", nodeAddr, replica))
			nodeHashToAddrs[virtualNodeHash] = nodeAddr
			nodeHashes[count] = virtualNodeHash
			count++
		}
	}
	slices.Sort(nodeHashes)

	return HashRing{
		NodeAddrs:      nodes,
		Replicas:       replicas,
		NodeHashToAddr: nodeHashToAddrs,
		NodesHashList:  nodeHashes,
	}
}

func (hr *HashRing) AddNode(node string) {
	hr.NodeAddrs = append(hr.NodeAddrs, node)
	for i := range hr.Replicas {
		nodeHash := generateKeyHash(fmt.Sprintf("%s-%d", node, i))
		hr.NodeHashToAddr[nodeHash] = node
		hr.NodesHashList = append(hr.NodesHashList, nodeHash)
	}
	slices.Sort(hr.NodesHashList)
}

func (hr *HashRing) RemoveNode(node string) {
	for i := range hr.Replicas {
		nodeHash := generateKeyHash(fmt.Sprintf("%s-%d", node, i))
		hr.NodesHashList = slices.DeleteFunc(hr.NodesHashList, func(hash uint32) bool { return hash == nodeHash })
		delete(hr.NodeHashToAddr, nodeHash)
	}
}

func (hr HashRing) findNextBiggestHash(lookupValue uint32) uint32 {
	n, found := slices.BinarySearch(hr.NodesHashList, lookupValue)
	if found {
		n++
	}
	if n > len(hr.NodesHashList)-1 {
		return hr.NodesHashList[0]
	}
	return hr.NodesHashList[n]
}

func (hr HashRing) GetNodeAddrForKey(key string) string {
	keyHash := generateKeyHash(key)
	nodeHash := hr.findNextBiggestHash(keyHash)
	return hr.NodeHashToAddr[nodeHash]
}

type HashRing struct {
	NodeAddrs      []string
	Replicas       int
	NodeHashToAddr map[uint32]string
	NodesHashList  []uint32
}

func generateKeyHash(key string) uint32 {
	hash := md5.Sum([]byte(key))
	hashNumber := binary.LittleEndian.Uint32(hash[:])
	return hashNumber
}
