package pointcloud

import (
	"math/rand"
	"strconv"
)

const DefaultMaxPointsPerNode = 50_000
const DefaultMaxDepth = 8

// Node is one octree node. Path is a digit string ("", "0", "03", ...)
// where each digit is the child octant index (0-7) chosen at that depth;
// the root node's Path is "". Points holds only this node's own subsample -
// the remainder of the points in its region live in Children.
type Node struct {
	Path     string
	Depth    int
	Bounds   Bounds
	Points   []Point
	Children [8]*Node
}

func (n *Node) HasChildren() bool {
	for _, c := range n.Children {
		if c != nil {
			return true
		}
	}
	return false
}

// BuildOctree splits points into a Potree-style octree: each node keeps a
// bounded random subsample of the points in its region, and the remainder
// cascades into up to 8 child octants. Every input point ends up in exactly
// one node.
func BuildOctree(points []Point, maxPointsPerNode, maxDepth int) *Node {
	if maxPointsPerNode <= 0 {
		maxPointsPerNode = DefaultMaxPointsPerNode
	}
	if maxDepth <= 0 {
		maxDepth = DefaultMaxDepth
	}
	bounds := BoundsOf(points)
	return buildNode(points, bounds, "", 0, maxPointsPerNode, maxDepth)
}

func buildNode(points []Point, bounds Bounds, path string, depth int, maxPointsPerNode, maxDepth int) *Node {
	node := &Node{Path: path, Depth: depth, Bounds: bounds}

	if len(points) <= maxPointsPerNode || depth >= maxDepth {
		node.Points = points
		return node
	}

	rand.Shuffle(len(points), func(i, j int) {
		points[i], points[j] = points[j], points[i]
	})

	node.Points = points[:maxPointsPerNode]
	remainder := points[maxPointsPerNode:]

	center := bounds.Center()

	var buckets [8][]Point
	for _, p := range remainder {
		idx := octantIndex(p, center)
		buckets[idx] = append(buckets[idx], p)
	}

	for i := 0; i < 8; i++ {
		if len(buckets[i]) == 0 {
			continue
		}
		childBounds := bounds.Octant(i, center)
		node.Children[i] = buildNode(buckets[i], childBounds, path+strconv.Itoa(i), depth+1, maxPointsPerNode, maxDepth)
	}

	return node
}

func octantIndex(p Point, center [3]float64) int {
	idx := 0
	if float64(p.X) >= center[0] {
		idx |= 1
	}
	if float64(p.Y) >= center[1] {
		idx |= 2
	}
	if float64(p.Z) >= center[2] {
		idx |= 4
	}
	return idx
}

// Flatten walks the tree into a flat slice, parents before children, for
// storage.
func Flatten(root *Node) []*Node {
	if root == nil {
		return nil
	}
	result := make([]*Node, 0)
	var walk func(n *Node)
	walk = func(n *Node) {
		if n == nil {
			return
		}
		result = append(result, n)
		for _, c := range n.Children {
			walk(c)
		}
	}
	walk(root)
	return result
}
