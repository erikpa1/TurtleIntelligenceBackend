package tools

import (
	"github.com/RyanCarrier/dijkstra/v2"
	"turtle/lg"
)

type PathResult struct {
	Path     []string
	Distance float64
}

type Graph struct {
	djikstraGraph *dijkstra.Graph
	incIndex      int
	indexToUid    map[int]string
	uidToIndex    map[string]int
}

func NewGraph() *Graph {
	tmp := &Graph{}
	tmp.djikstraGraph = &dijkstra.Graph{}
	tmp.indexToUid = map[int]string{}
	tmp.uidToIndex = map[string]int{}
	tmp.incIndex = 0

	return tmp
}

func (self *Graph) AddEmptyVertex(uid string) {
	self.uidToIndex[uid] = self.incIndex
	self.indexToUid[self.incIndex] = uid
	self.djikstraGraph.AddEmptyVertex(self.incIndex)
	self.incIndex += 1
}

func (self *Graph) AddArc(a, b string, distance float64) {

	a_i, a_ok := self.uidToIndex[a]
	b_i, b_ok := self.uidToIndex[b]

	if a_ok && b_ok {
		err := self.djikstraGraph.AddArc(a_i, b_i, uint64(distance*100))
		if err != nil {
			lg.LogE("Unable to find connection between")
		}
	} else {
		if a_ok == false {
			lg.LogE("Unable to convert A", a)
		} else {
			lg.LogE("Unable to convert B", b)
		}

	}
}
func (self *Graph) Shortest(a, b string) PathResult {

	if a == b {
		lg.LogE("The same routes searching: [", a, "] , [", b, "]")
	}

	a_i, a_ok := self.uidToIndex[a]
	b_i, b_ok := self.uidToIndex[b]

	if a_ok && b_ok {
		shortest, err := self.djikstraGraph.Shortest(a_i, b_i)

		if err == nil {
			var uidPath = make([]string, len(shortest.Path))

			for _, val := range shortest.Path {
				uidPath = append(uidPath, self.indexToUid[val])
			}
			return PathResult{
				Distance: float64(shortest.Distance) / 100,
				Path:     uidPath,
			}
		} else {
			lg.LogE("Error between: [", a, "] or [", b, "]")
			lg.LogE(err)
		}
	} else {
		lg.LogE("Unable to find: [", a, "] or [", b, "]")
	}

	return PathResult{}
}
