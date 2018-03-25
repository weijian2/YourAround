package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
	"strconv"
	"reflect"
	"gopkg.in/olivere/elastic.v3"
	"github.com/pborman/uuid"
	"strings"
	"context"
	"cloud.google.com/go/storage"
	"io"
	"cloud.google.com/go/bigtable"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/go-redis/redis"
	"time"
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
	Url    string `json:"url"`
}

const (
	INDEX = "around"
	TYPE = "post" // table name in ES to store post data
	DISTANCE = "200km"
	PROJECT_ID = "youraround-cmu"
	BT_INSTANCE = "around-post"
	ES_URL = "http://35.226.169.43:9200/"
	BUCKET_NAME = "post-images-youraround-cmu"
	ENABLE_MEMCACHE = true
	ENABLE_BIGTABLE = false
	REDIS_URL       = "redis-14642.c1.us-central1-2.gce.cloud.redislabs.com:14642"
)
// private key used to sign token
var mySigningKey = []byte("mySecretKey")


func main() {

	// Create a client, which means we create a connection to ES. If there is err, return.
	client, err := elastic.NewClient(elastic.SetURL(ES_URL), elastic.SetSniff(false))
	if err != nil {
		panic(err)
		return
	}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(INDEX).Do()
	if err != nil {
		panic(err)
	}
	if !exists {
		// Create a new index.
		// If not, create a new mapping. For other fields (user, message, etc.)
		// no need to have mapping as they are default. For geo location (lat, lon),
		// we need to tell ES that they are geo points instead of two float points
		// such that ES will use Geo-indexing for them (K-D tree)

		mapping := `{
                    "mappings":{
                           "post":{
                                  "properties":{
                                         "location":{
                                                "type":"geo_point"
                                         }
                                  }
                           }
                    }
             }
             `
		_, err := client.CreateIndex(INDEX).Body(mapping).Do() // Create this index
		if err != nil {
			// Handle error
			panic(err)
		}
	}

	fmt.Println("service start")
	//// without Auth
	//http.HandleFunc("/post", handlerPost)
	//http.HandleFunc("/search", handlerSearch)
	//log.Fatal(http.ListenAndServe(":8080", nil))

	// with Auth
	// reference: https://auth0.com/blog/authentication-in-golang/
	// Create a new router on top of the existing http router as we need to check auth.
	r := mux.NewRouter()
	// Create a new JWT middleware with a Option that uses the key ‘mySigningKey’ such that we know this token is
	// from our server. The signing method is the default HS256 algorithm such that data is encrypted.
	var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	// It means we use jwt middleware to manage these endpoints and if they don’t have valid token, we will reject them.
	r.Handle("/post", jwtMiddleware.Handler(http.HandlerFunc(handlerPost))).Methods("POST")
	r.Handle("/search", jwtMiddleware.Handler(http.HandlerFunc(handlerSearch))).Methods("GET")
	// login and signup don't need middleware, since it is the first time of request, we don't have token now.
	r.Handle("/login", http.HandlerFunc(loginHandler)).Methods("POST")
	r.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

// how to parse multipart form in Go?
// https://golang.org/pkg/net/http/#Request.ParseMultipartForm
// https://github.com/golang-samples/http/blob/master/fileupload/main.go
func handlerPost(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

	// get username from context instead of letting user to input username themselves
	user := r.Context().Value("user")
	claims := user.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"]

	// 32 << 20 is the maxMemory param for ParseMultipartForm, equals to 32MB (1MB = 1024 * 1024 bytes = 2^20 bytes)
	// After you call ParseMultipartForm, the file will be saved in the server memory with maxMemory size.
	// If the file size is larger than maxMemory, the rest of the data will be saved in a system temporary file.
	r.ParseMultipartForm(32 << 20)

	// Parse from form data.
	fmt.Printf("Received one post request %s\n", r.FormValue("message"))
	lat, _ := strconv.ParseFloat(r.FormValue("lat"), 64)
	lon, _ := strconv.ParseFloat(r.FormValue("lon"), 64)
	p := &Post{
		User:    username.(string),
		Message: r.FormValue("message"),
		Location: Location{
			Lat: lat,
			Lon: lon,
		},
	}

	id := uuid.New()

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image is not available", http.StatusInternalServerError)
		fmt.Printf("Image is not available %v.\n", err)
		return
	}
	// defer here because if some codes break in the middle of method, file.close() will not be executed
	defer file.Close()

	// fetch authorized user from local config file, must first authorized with Google
	ctx := context.Background()

	// attrs is file
	_, attrs, err := saveToGCS(ctx, file, BUCKET_NAME, id)
	if err != nil {
		http.Error(w, "GCS is not setup", http.StatusInternalServerError)
		fmt.Printf("GCS is not setup %v\n", err)
		return
	}

	// Update the media link after saving to GCS.
	p.Url = attrs.MediaLink

	// Save to ES.
	saveToES(p, id)

	// Save to BigTable as well.
	if ENABLE_BIGTABLE {
		saveToBigTable(ctx, p, id)
	}
}

// http://localhost:8888/search?lat=10.0&lon=20.0
// http://localhost:8888/search?lat=10.0&lon=20.0&range=100
func handlerSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	fmt.Println("Received one request for search")
	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lon, _ := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
	ran := DISTANCE
	if val := r.URL.Query().Get("range"); val != "" {
		ran = val + "km"
	}
	//fmt.Println("range is ", ran)

	fmt.Printf( "Search received: %f %f %s\n", lat, lon, ran)

	// use cache, reference: https://github.com/go-redis/redis
	// use lat+lon+range as key
	key := r.URL.Query().Get("lat") + ":" + r.URL.Query().Get("lon") + ":" + ran
	// First find query with cache.
	if ENABLE_MEMCACHE {
		// build connection with Redis
		rs_client := redis.NewClient(&redis.Options{
			Addr:     REDIS_URL,
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		val, err := rs_client.Get(key).Result()
		if err != nil {
			fmt.Printf("Cache miss! Redis cannot find the key %s as %v.\n", key, err)
		} else { // cache hit
			fmt.Printf("Cache hit! Redis find the key %s.\n", key)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(val))
			return
		}
	}

	// for here, not enable cache or cache miss, we need to search in database(ES).
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
		// Perform filtering based on keywords such as web spam etc.
		if !containsSensitiveWords(&p.Message) {
			ps = append(ps, p)
		}
	}
	js, err := json.Marshal(ps) // Convert the go object to a string
	if err != nil {
		panic(err)
		return
	}

	// for here, we find result from ES, we need to write result into Redis, use TTL(time to live) as
	// caching strategy to avoid result inconsistent
	if ENABLE_MEMCACHE {
		// build connection with Redis
		rs_client := redis.NewClient(&redis.Options{
			Addr:     REDIS_URL,
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		// Set the cache expiration to be 30 seconds
		err := rs_client.Set(key, string(js), time.Second*30).Err()
		if err != nil {
			fmt.Printf("Redis cannot save the key %s as %v.\n", key, err)
		}

	}

	w.Write(js)
}

// private method, Save a post to ElasticSearch
func saveToES(p *Post, id string) {
	// Create a client
	es_client, err := elastic.NewClient(elastic.SetURL(ES_URL), elastic.SetSniff(false))
	if err != nil {
		panic(err)
		return
	}

	// Save it to index, example taken from https://github.com/olivere/elastic
	_, err = es_client.Index().
		Index(INDEX).
		Type(TYPE).
		Id(id).
		BodyJson(p).
		Refresh(true).
		Do()
	if err != nil {
		panic(err)
		return
	}

	fmt.Printf("Post is saved to Index: %s\n", p.Message)
}

// private method, filter sensitive words
func containsSensitiveWords(post *string) bool {
	sensitiveWords := []string {"fuck", "dick"}
	for _, word := range sensitiveWords {
		if strings.Contains(*post, word) {
			return true
		}
	}
	return false
}

// private method, save image to Google Cloud Storage
// Google example of open a client connection to GCS
// https://cloud.google.com/storage/docs/reference/libraries#client-libraries-install-go
// Google example of writing an object to GCS, see write function
// https://github.com/GoogleCloudPlatform/golang-samples/blob/master/storage/objects/main.go
// r is image file, name is id of this image
func saveToGCS(ctx context.Context, r io.Reader, bucketName, name string) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer client.Close()

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)
	// Next check if the bucket exists
	if _, err = bucket.Attrs(ctx); err != nil {
		return nil, nil, err
	}

	obj := bucket.Object(name) // stored file
	w := obj.NewWriter(ctx)
	// copy(write) image file to GCS's bucket
	if _, err := io.Copy(w, r); err != nil {
		return nil, nil, err
	}
	if err := w.Close(); err != nil {
		return nil, nil, err
	}

	// already finish writting, now modify access permission to all users, ACL is access control list
	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return nil, nil, err
	}

	attrs, err := obj.Attrs(ctx)
	fmt.Printf("Post is saved to GCS: %s\n", attrs.MediaLink)
	return obj, attrs, err
}

// private method, Save a post to BigTable
// how to write obj to bigtable reference: https://cloud.google.com/bigtable/docs/samples-go-hello
// type conversion: string to byte array. byteArray := []byte(myString)
func saveToBigTable(ctx context.Context, p *Post, id string) {
	bt_client, err := bigtable.NewClient(ctx, PROJECT_ID, BT_INSTANCE)
	if err != nil {
		panic(err)
		return
	}
	tbl := bt_client.Open("post")
	mut := bigtable.NewMutation()
	t := bigtable.Now()
	// params: columnFamilyName, columnName, timestamp value
	mut.Set("post", "user", t, []byte(p.User))
	mut.Set("post", "message", t, []byte(p.Message))
	mut.Set("location", "lat", t, []byte(strconv.FormatFloat(p.Location.Lat, 'f', -1, 64)))
	mut.Set("location", "lon", t, []byte(strconv.FormatFloat(p.Location.Lon, 'f', -1, 64)))

	err = tbl.Apply(ctx, id, mut)
	if err != nil {
		panic(err)
		return
	}
	fmt.Printf("Post is saved to BigTable: %s\n", p.Message)
}