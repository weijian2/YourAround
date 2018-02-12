package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
	"strconv"
	"reflect"
	"gopkg.in/olivere/elastic.v3"
)

// type/struct are keywords in GO, struct is similar to class in java, Location is struct name
type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Post struct {
	// `json:"user"` is anotation, for the json parsing of this User field. Otherwise, by default it's 'User'.
	User     string `json:"user"`
	Message  string  `json:"message"`
	Location Location `json:"location"`
}

const (
	INDEX = "around"
	TYPE = "post"
	DISTANCE = "200km"
	// Needs to update
	//PROJECT_ID = "around-xxx"
	//BT_INSTANCE = "around-post"
	// Needs to update this URL if you deploy it to cloud.
	ES_URL = "http://35.196.197.253:9200"
)


func main() {
	fmt.Println("service start")
	http.HandleFunc("/post", handlerPost)
	http.HandleFunc("/search", handlerSearch)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func handlerPost(w http.ResponseWriter, r *http.Request) {
	// Parse from body of request to get a json object.
	fmt.Println("Received one post request")
	decoder := json.NewDecoder(r.Body)
	var p Post
	if err := decoder.Decode(&p); err != nil {
		panic(err)
		return
	}
	fmt.Fprintf(w, "Post received: %s\n", p.Message)
}

// http://localhost:8888/search?lat=10.0&lon=20.0
// http://localhost:8888/search?lat=10.0&lon=20.0&range=100
func handlerSearch(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request for search")
	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lon, _ := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
	ran := DISTANCE
	if val := r.URL.Query().Get("range"); val != "" {
		ran = val + "km"
	}
	//fmt.Println("range is ", ran)

	fmt.Printf( "Search received: %f %f %s\n", lat, lon, ran)

	// Create a client，which means we create a connection to ES. If there is err, return.
	client, err := elastic.NewClient(elastic.SetURL(ES_URL), elastic.SetSniff(false))
	if err != nil {
		panic(err)
		return
	}

	// Define geo distance query as specified in
	// https://www.elastic.co/guide/en/elasticsearch/reference/5.2/query-dsl-geo-distance-query.html
	// Prepare a geo based query to find posts within a geo box.
	q := elastic.NewGeoDistanceQuery("location")
	q = q.Distance(ran).Lat(lat).Lon(lon)

	// Some delay may range from seconds to minutes.
	// Get the results based on Index (similar to dataset) and query (q that we just prepared).
	// Pretty means to format the output.
	searchResult, err := client.Search().
		Index(INDEX).
		Query(q).
		Pretty(true).
		Do()
	if err != nil {
		// Handle error
		panic(err)
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)
	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d post\n", searchResult.TotalHits())

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization.
	var typ Post
	var ps []Post
	// Iterate the result results and if they are type of Post (typ)
	for _, item := range searchResult.Each(reflect.TypeOf(typ)) {
		p := item.(Post) // Cast an item to Post, equals to p = (Post) item in java
		fmt.Printf("Post by %s: %s at lat %v and lon %v\n", p.User, p.Message, p.Location.Lat, p.Location.Lon)
		// TODO: Perform filtering based on keywords such as web spam etc.
		ps = append(ps, p)

	}
	js, err := json.Marshal(ps) // Convert the go object to a string
	if err != nil {
		panic(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)

	//// return fake post
	//p := &Post {
	//	User:"Weijian",
	//	Message:"是中国人就转，不转不是人",
	//	Location: Location {
	//		Lat: lat,
	//		Lon: lon,
	//	},
	//}
	//
	//js, err := json.Marshal(p)
	//if err != nil {
	//	panic(err)
	//	return
	//}
	//
	//w.Header().Set("Content-Type", "application/json")
	//w.Write(js)
	//
	////fmt.Fprintf(w, "Search received: %s %s", lat, lon)
}