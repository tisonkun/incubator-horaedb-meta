// Copyright 2022 CeresDB Project Authors. Licensed under Apache-2.0.

package storage

import "github.com/CeresDB/ceresdbproto/golang/pkg/clusterpb"

type (
	ClusterID    uint32
	SchemaID     uint32
	ShardID      uint32
	TableID      uint64
	ClusterState int
	ShardRole    int
	NodeState    int
)

const (
	ClusterStateEmpty ClusterState = iota + 1
	ClusterStateStable
)

const (
	ShardRoleLeader ShardRole = iota + 1
	ShardRoleFollower
)

const (
	NodeStateOnline NodeState = iota + 1
	NodeStateOffline
)

type ListClustersResult struct {
	Clusters []Cluster
}

type CreateClusterRequest struct {
	Cluster Cluster
}

type CreateClusterViewRequest struct {
	ClusterView ClusterView
}

type GetClusterViewRequest struct {
	ClusterID ClusterID
}

type GetClusterViewResult struct {
	ClusterView ClusterView
}

type UpdateClusterViewRequest struct {
	ClusterID     ClusterID
	ClusterView   ClusterView
	LatestVersion uint64
}

type ListSchemasRequest struct {
	ClusterID ClusterID
}

type ListSchemasResult struct {
	Schemas []Schema
}

type CreateSchemaRequest struct {
	ClusterID ClusterID
	Schema    Schema
}

type CreateTableRequest struct {
	ClusterID ClusterID
	SchemaID  SchemaID
	Table     Table
}

type GetTableRequest struct {
	ClusterID ClusterID
	SchemaID  SchemaID
	TableName string
}

type GetTableResult struct {
	Table  Table
	Exists bool
}

type ListTableRequest struct {
	ClusterID ClusterID
	SchemaID  SchemaID
}

type ListTablesResult struct {
	Tables []Table
}

type DeleteTableRequest struct {
	ClusterID ClusterID
	SchemaID  SchemaID
	TableName string
}

type CreateShardViewsRequest struct {
	ClusterID  ClusterID
	ShardViews []ShardView
}

type ListShardViewsRequest struct {
	ClusterID ClusterID
	ShardIDs  []ShardID
}

type ListShardViewsResult struct {
	ShardViews []ShardView
}

type UpdateShardViewRequest struct {
	ClusterID     ClusterID
	ShardView     ShardView
	LatestVersion uint64
}

type ListNodesRequest struct {
	ClusterID ClusterID
}

type ListNodesResult struct {
	Nodes []Node
}

type CreateOrUpdateNodeRequest struct {
	ClusterID ClusterID
	Node      Node
}

type Cluster struct {
	ID                ClusterID
	Name              string
	MinNodeCount      uint32
	ReplicationFactor uint32
	ShardTotal        uint32
	CreatedAt         uint64
}

type ShardNode struct {
	ID        ShardID
	ShardRole ShardRole
	NodeName  string
}

type ClusterView struct {
	ClusterID  ClusterID
	Version    uint64
	State      ClusterState
	ShardNodes []ShardNode
	CreatedAt  uint64
}

type Schema struct {
	ID        SchemaID
	ClusterID ClusterID
	Name      string
	CreatedAt uint64
}

type Table struct {
	ID        TableID
	Name      string
	SchemaID  SchemaID
	CreatedAt uint64
}

type ShardView struct {
	ShardID   ShardID
	Version   uint64
	TableIDs  []TableID
	CreatedAt uint64
}

type NodeStats struct {
	Lease       uint32
	Zone        string
	NodeVersion string
}

type Node struct {
	Name          string
	NodeStats     NodeStats
	CreatedAt     uint64
	LastTouchTime uint64
	State         NodeState
}

func ConvertShardRolePB(role clusterpb.ShardRole) ShardRole {
	switch role {
	case clusterpb.ShardRole_LEADER:
		return ShardRoleLeader
	case clusterpb.ShardRole_FOLLOWER:
		return ShardRoleFollower
	}
	return ShardRoleFollower
}

func convertNodeToPB(node Node) clusterpb.Node {
	nodeStats := convertNodeStatsToPB(node.NodeStats)
	return clusterpb.Node{
		Name:          node.Name,
		NodeStats:     &nodeStats,
		CreateTime:    node.CreatedAt,
		LastTouchTime: node.LastTouchTime,
		State:         convertNodeStateToPB(node.State),
	}
}

func convertClusterPB(cluster *clusterpb.Cluster) Cluster {
	return Cluster{
		ID:                ClusterID(cluster.Id),
		Name:              cluster.Name,
		MinNodeCount:      cluster.MinNodeCount,
		ReplicationFactor: cluster.ReplicationFactor,
		ShardTotal:        cluster.ShardTotal,
		CreatedAt:         cluster.CreatedAt,
	}
}

func convertClusterToPB(cluster Cluster) clusterpb.Cluster {
	return clusterpb.Cluster{
		Id:                uint32(cluster.ID),
		Name:              cluster.Name,
		MinNodeCount:      cluster.MinNodeCount,
		ReplicationFactor: cluster.ReplicationFactor,
		ShardTotal:        cluster.ShardTotal,
		CreatedAt:         cluster.CreatedAt,
	}
}

func convertClusterStateToPB(state ClusterState) clusterpb.ClusterTopology_ClusterState {
	switch state {
	case ClusterStateEmpty:
		return clusterpb.ClusterTopology_EMPTY
	case ClusterStateStable:
		return clusterpb.ClusterTopology_STABLE
	}
	return clusterpb.ClusterTopology_EMPTY
}

func convertClusterStatePB(state clusterpb.ClusterTopology_ClusterState) ClusterState {
	switch state {
	case clusterpb.ClusterTopology_EMPTY:
		return ClusterStateEmpty
	case clusterpb.ClusterTopology_STABLE:
		return ClusterStateStable
	}
	return ClusterStateEmpty
}

func convertShardRoleToPB(role ShardRole) clusterpb.ShardRole {
	switch role {
	case ShardRoleLeader:
		return clusterpb.ShardRole_LEADER
	case ShardRoleFollower:
		return clusterpb.ShardRole_FOLLOWER
	}
	return clusterpb.ShardRole_FOLLOWER
}

func convertShardNodeToPB(shardNode ShardNode) clusterpb.Shard {
	return clusterpb.Shard{
		Id:        uint32(shardNode.ID),
		ShardRole: convertShardRoleToPB(shardNode.ShardRole),
		Node:      shardNode.NodeName,
	}
}

func convertShardNodePB(shardNode *clusterpb.Shard) ShardNode {
	return ShardNode{
		ID:        ShardID(shardNode.Id),
		ShardRole: ConvertShardRolePB(shardNode.ShardRole),
		NodeName:  shardNode.Node,
	}
}

func convertClusterViewToPB(view ClusterView) clusterpb.ClusterTopology {
	shardViews := make([]*clusterpb.Shard, 0, len(view.ShardNodes))
	for _, shardNode := range view.ShardNodes {
		shardNodePB := convertShardNodeToPB(shardNode)
		shardViews = append(shardViews, &shardNodePB)
	}

	return clusterpb.ClusterTopology{
		ClusterId:        uint32(view.ClusterID),
		Version:          view.Version,
		State:            convertClusterStateToPB(view.State),
		ShardView:        shardViews,
		Cause:            "",
		CreatedAt:        view.CreatedAt,
		ChangeNodeShards: nil,
	}
}

func convertClusterViewPB(view *clusterpb.ClusterTopology) ClusterView {
	shardNodes := make([]ShardNode, 0, len(view.ShardView))
	for _, shardNodePB := range view.ShardView {
		shardNode := convertShardNodePB(shardNodePB)
		shardNodes = append(shardNodes, shardNode)
	}

	return ClusterView{
		ClusterID:  ClusterID(view.ClusterId),
		Version:    view.Version,
		State:      convertClusterStatePB(view.State),
		ShardNodes: shardNodes,
		CreatedAt:  view.CreatedAt,
	}
}

func convertSchemaToPB(schema Schema) clusterpb.Schema {
	return clusterpb.Schema{
		Id:        uint32(schema.ID),
		ClusterId: uint32(schema.ClusterID),
		Name:      schema.Name,
		CreatedAt: schema.CreatedAt,
	}
}

func convertSchemaPB(schema *clusterpb.Schema) Schema {
	return Schema{
		ID:        SchemaID(schema.Id),
		ClusterID: ClusterID(schema.ClusterId),
		Name:      schema.Name,
		CreatedAt: schema.CreatedAt,
	}
}

func convertTableToPB(table Table) clusterpb.Table {
	return clusterpb.Table{
		Id:        uint64(table.ID),
		Name:      table.Name,
		SchemaId:  uint32(table.SchemaID),
		ShardId:   uint32(table.SchemaID),
		Desc:      "",
		CreatedAt: table.CreatedAt,
	}
}

func convertTablePB(table *clusterpb.Table) Table {
	return Table{
		ID:        TableID(table.Id),
		Name:      table.Name,
		SchemaID:  SchemaID(table.SchemaId),
		CreatedAt: table.CreatedAt,
	}
}

func convertShardViewToPB(view ShardView) clusterpb.ShardTopology {
	tableIDs := make([]uint64, 0, len(view.TableIDs))
	for _, id := range view.TableIDs {
		tableIDs = append(tableIDs, uint64(id))
	}

	return clusterpb.ShardTopology{
		ShardId:   uint32(view.ShardID),
		TableIds:  tableIDs,
		Version:   view.Version,
		CreatedAt: view.CreatedAt,
	}
}

func convertShardViewPB(shardTopology *clusterpb.ShardTopology) ShardView {
	tableIDs := make([]TableID, 0, len(shardTopology.TableIds))
	for _, id := range shardTopology.TableIds {
		tableIDs = append(tableIDs, TableID(id))
	}

	return ShardView{
		ShardID:   ShardID(shardTopology.ShardId),
		Version:   shardTopology.Version,
		TableIDs:  tableIDs,
		CreatedAt: shardTopology.CreatedAt,
	}
}

func convertNodeStatsToPB(stats NodeStats) clusterpb.NodeStats {
	return clusterpb.NodeStats{
		Lease:       stats.Lease,
		Zone:        stats.Zone,
		NodeVersion: stats.NodeVersion,
	}
}

func convertNodeStatsPB(stats *clusterpb.NodeStats) NodeStats {
	return NodeStats{
		Lease:       stats.Lease,
		Zone:        stats.Zone,
		NodeVersion: stats.NodeVersion,
	}
}

func convertNodeStateToPB(state NodeState) clusterpb.NodeState {
	switch state {
	case NodeStateOnline:
		return clusterpb.NodeState_ONLINE
	case NodeStateOffline:
		return clusterpb.NodeState_OFFLINE
	}
	return clusterpb.NodeState_OFFLINE
}

func convertNodeStatePB(state clusterpb.NodeState) NodeState {
	switch state {
	case clusterpb.NodeState_ONLINE:
		return NodeStateOnline
	case clusterpb.NodeState_OFFLINE:
		return NodeStateOffline
	}
	return NodeStateOffline
}

func convertNodePB(node *clusterpb.Node) Node {
	nodeStats := convertNodeStatsPB(node.NodeStats)
	return Node{
		Name:          node.Name,
		NodeStats:     nodeStats,
		CreatedAt:     node.CreateTime,
		LastTouchTime: node.LastTouchTime,
		State:         convertNodeStatePB(node.State),
	}
}