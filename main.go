package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	song "github.com/ademaligolang/Song_Definitions"
	"github.com/gorilla/mux"
)

// The port this webservice is listening on
const port = ":8081"

// Prepare our variables
var graph ItemGraph
var songs song.Songs
var songTitles = make(map[string]*song.Song)

func main() {
	// Using gorialla mux for http routing
	router := mux.NewRouter().StrictSlash(true)

	// Handles
	router.HandleFunc("/AddSong", AddSong)
	router.HandleFunc("/GetSongGroups", GetSongGroups)

	// A simple listen and serve returning errors to the command line
	error := http.ListenAndServe(port, router)
	fmt.Println(error)
}

// Method to add a song into the graph
func AddSong(w http.ResponseWriter, r *http.Request) {
	// Create a new empty song
	var newSong song.Song

	// Read the body of the request
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	// Panic if we've found an error
	if err != nil {
		panic(err)
	}

	// Panic if it errors while closing the request
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// Try to convert the JSON request into a song object
	if err := json.Unmarshal(body, &newSong); err != nil {
		// Prepare our response to the source of the request
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity)

		// Panic if we can't
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	// Append this new song to the song graph type
	songs.Songs = append(songs.Songs, newSong)

	// Now figure out where to place it in the graph
	songKey := GetSongTitle(newSong.Title)
	fmt.Println(newSong)
	fmt.Println(songKey)
	fmt.Println(songTitles)
	if songTitles[songKey] == nil {
		// If the song key doesn't exist then add it fresh
		songTitles[songKey] = &newSong
		graph.AddSong(&newSong)
	} else {
		// Else find the root of this branch and add it
		lastSong := graph.FindDanglingNode(songTitles[songKey])

		graph.AddSong(&newSong)
		graph.AddEdge(lastSong, &newSong)
	}

	// Print some confirmation to the response writer and set http status
	fmt.Fprintf(w, "Got song: %s by %s", newSong.Title, newSong.Artists)
}

// Method to get a display of all of the song groups
func GetSongGroups(w http.ResponseWriter, r *http.Request) {
	graph.DepthFirstTraverse(w)
}

// A very simple title matching method using regex
func GetSongTitle(song string) string {
	// Compile our regular expression to remove parentheses
	re := regexp.MustCompile(`\(.*\)`)

	// Remove the parenthesis text and set all to lowercase
	cleansedTitle := re.ReplaceAllString(song, "")
	cleansedTitle = strings.ToLower(cleansedTitle)

	// Here we could also strip other song variation indicators such as "live", "cover", "remaster", etc.
	cleansedTitle = strings.ReplaceAll(cleansedTitle, "remix", "")
	cleansedTitle = strings.TrimSpace(cleansedTitle)

	// Return the cleansed title
	return cleansedTitle
}
