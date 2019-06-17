# Song_WebService

This webservice holds the list of songs and graph data types in memory. It has two URL's of interest.

http://localhost:8081/AddSong - This takes a json string for a single song and adds it into the graph at an appropriate place

http://localhost:8081/GetSongGroups - This does the job of traversing the graph and returning the text descriptions of the song groups found in the graph.
