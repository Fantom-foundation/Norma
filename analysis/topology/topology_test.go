package topology

import (
	"math/rand"
	"testing"

	"github.com/dominikbraun/graph"
)

func GetDiameter(g graph.Graph[int, int]) int {
	unit := struct{}{}
	edges, err := g.Edges()
	if err != nil {
		panic(err)
	}
	vertices := map[int]struct{}{}
	for _, edge := range edges {
		vertices[edge.Source] = unit
		vertices[edge.Target] = unit
	}

	max := 0
	for from := range vertices {
		for to := range vertices {
			path, err := graph.ShortestPath(g, from, to)
			if err != nil {
				continue
				//panic(err)
			}
			if l := len(path) - 1; l > max {
				max = l
			}
		}
	}

	return max
}

func TestGetDiameter(t *testing.T) {

	g := graph.New(graph.IntHash)

	g.AddVertex(1)
	g.AddVertex(2)
	g.AddVertex(3)
	g.AddVertex(4)

	g.AddEdge(1, 2)
	g.AddEdge(2, 3)
	g.AddEdge(3, 4)
	g.AddEdge(1, 4)

	if want, got := 2, GetDiameter(g); want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestGetDiameterDevNet(t *testing.T) {

	g := graph.New(graph.IntHash)

	for i := 1; i <= 13; i++ {
		g.AddVertex(i)
	}

	g.AddEdge(1, 2)
	g.AddEdge(1, 7)
	g.AddEdge(1, 12)

	g.AddEdge(2, 4)
	g.AddEdge(2, 5)

	g.AddEdge(3, 4)
	g.AddEdge(3, 5)
	g.AddEdge(3, 13)

	g.AddEdge(4, 6)

	g.AddEdge(5, 13)

	g.AddEdge(6, 10)
	g.AddEdge(6, 13)

	g.AddEdge(7, 8)
	g.AddEdge(7, 12)

	g.AddEdge(8, 9)
	g.AddEdge(8, 11)

	g.AddEdge(9, 10)
	g.AddEdge(9, 13)

	g.AddEdge(10, 11)

	g.AddEdge(11, 12)

	g.AddEdge(12, 13)

	t.Fatalf("Diameter of test net is %d\n", GetDiameter(g))
}

func TestGetDiameterRandomNet(t *testing.T) {

	for i := 0; i < 10; i++ {
		g := graph.New(graph.IntHash)

		for i := 1; i <= 13; i++ {
			g.AddVertex(i)
		}

		for i := 1; i <= 13; i++ {
			for j := i + 1; j <= 13; j++ {
				if i == j {
					continue
				}
				if rand.Int31n(4) == 0 {
					g.AddEdge(i, j)
				}
			}
		}

		t.Errorf("Diameter of test net is %d\n", GetDiameter(g))
	}
}
