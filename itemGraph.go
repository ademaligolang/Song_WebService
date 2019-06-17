package main

import (
	"fmt"
	"net/http"
	"sync"

	song "github.com/ademaligolang/Song_Definitions"
)

// The ItemGraph struct handles songs like a directed graph data type
type ItemGraph struct {
	nodes       []*song.Song
	edges       map[*song.Song][]*song.Song
	visitedList map[*song.Song]bool
	lock        sync.RWMutex
}

// Add a song to the nodes slice
func (g *ItemGraph) AddSong(n *song.Song) {
	g.lock.Lock()
	g.nodes = append(g.nodes, n)
	g.lock.Unlock()
}

// Add an edge between two nodes of songs
func (g *ItemGraph) AddEdge(n1, n2 *song.Song) {
	g.lock.Lock()
	if g.edges == nil {
		g.edges = make(map[*song.Song][]*song.Song)
	}

	// Add the edge entry to the map
	g.edges[n1] = append(g.edges[n1], n2)
	// If we wanted a unidirectional graph we would uncomment the line below
	//g.edges[n2] = append(g.edges[n2], n1)
	g.lock.Unlock()
}

// This performs a depth first traverse to collect the song group information
func (g *ItemGraph) DepthFirstTraverse(w http.ResponseWriter) {

	// Reset the visited list
	g.visitedList = make(map[*song.Song]bool)

	// Prepare our variables
	numGroups := 0

	// Iterate over all of the nodes
	for _, song := range g.nodes {

		// If we haven't visited this node
		if !g.visitedList[song] {
			// Increase the groups count
			numGroups++

			// Print this data to the writer
			fmt.Fprintf(w, "Song Group: %d\nGroup Title: %s\n\n", numGroups, song.Title)

			// Recursion lets us travel depth first
			g.DFS(w, song)

			// Line breaks to make it look nicer
			fmt.Fprintf(w, "\n----------------------------------------------------------\n")
		}
	}
}

// DFS is a part of the depth first search
func (g *ItemGraph) DFS(w http.ResponseWriter, song *song.Song) {

	// We are now visiting this node
	g.visitedList[song] = true

	// Print it's song information
	fmt.Fprintf(w, "%s\n", song.String())

	// Iterate over all edges connected to this song node
	for _, childSong := range g.edges[song] {

		// If we haven't already visited then let's go deeper
		if !g.visitedList[childSong] {
			g.DFS(w, childSong)
		}
	}
}

// Breadth First Travel, unused but still implemented for fun
func (g *ItemGraph) BreadthFirstTraverse() {

	// Reset the visited list and create a queue
	g.visitedList = make(map[*song.Song]bool)
	queue := []*song.Song{}

	// Iterate over every song node in the graph
	for _, s := range g.nodes {

		// If we haven't visited then add it to the queue
		if !g.visitedList[s] {
			queue = append(queue, s)
		}

		// If the queue is empty then break out
		if len(queue) == 0 {
			break
		}

		// Grab the song information
		song := queue[0]

		// Remove this song from the queue
		queue = queue[1:len(queue)]

		// Declare we've visited it in the list
		g.visitedList[song] = true

		// For all of the edges children
		for _, childSong := range g.edges[song] {

			// If we haven't visited then add it to the queue and move on
			if !g.visitedList[childSong] {
				queue = append(queue, childSong)
				g.visitedList[childSong] = true
			}
		}

		// Finally write the song information
		fmt.Println(song)
	}
}

// Recursive function used to find a dangling node
// A dangling node is a node with no direction
// and it's where we want to add any new songs that match the group
func (g *ItemGraph) FindDanglingNode(song *song.Song) *song.Song {

	// If this song node has edges
	if len(g.edges[song]) > 0 {

		// Iterate through all of them
		for _, childSong := range g.edges[song] {

			// and recursively call this method to go deeper
			return g.FindDanglingNode(childSong)
		}
	}

	// If we get here then this song is dangling
	return song
}
