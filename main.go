package main

import (

    "math/rand"
    "context"
    "time"
    "log"
    "io"
    "net/http"
    "encoding/json"
    "fmt"
    "strings"
    "github.com/redis/go-redis/v9"
)

//A struct that is used to encode a http request into JSON
type shortreq struct {

    Url string  `json:"url"`

}

var ctx = context.Background()

//Initialize the redis client
var client = redis.NewClient(&redis.Options{
 
    Addr:        "localhost:6379",
    Password:    "",
    DB:          0,

})

/*
var memoryStore = MemoryStore{
    
    Map:                 make(map[string]string),
    Mutex:               new(sync.Mutex),

}

var m *MemoryStore = &memoryStore
*/

func main(){

    //ping the redis client to see if there is a connection
    pong, err := client.Ping(ctx).Result()
    fmt.Println(pong, err)

    // set up http endpoints
    http.HandleFunc("/", RedirectHandler)
    http.HandleFunc("/healthz", HealthCheckHandler)
    http.HandleFunc("/shorten", ShortenHandler)
    http.ListenAndServe(":8080", nil)
}

//A function that creates a random string for the shortened url
//Can be considered the ID for the min url
func makeRandomString() string {

    const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
    seed := rand.NewSource(time.Now().UnixNano())
    random := rand.New(seed)

    result := make([]byte, 8)

    for i := range result {

        result[i] = charset[random.Intn(len(charset))]

    }

    return string(result)

}

//This handles a request to shorten a url 
func ShortenHandler(w http.ResponseWriter, req *http.Request){

    defer req.Body.Close()

    data, err := io.ReadAll(req.Body)
    
    if err != nil {
        
        http.Error(w, err.Error(), http.StatusBadRequest) 
        return
    }

    shortenreq := shortreq{}

    err = json.Unmarshal(data, &shortenreq)

    if err != nil {

        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if shortenreq.Url == "" {

        w.Write([]byte("Error: Empty request"))
        return
    }

    if !IsSafe(shortenreq.Url) {

        w.Write([]byte("Error: URL is not secure"))
        w.Write([]byte("\n"))
        return
    }
    
    fmt.Printf("The long URL is %s\n", shortenreq.Url)

    w.WriteHeader(http.StatusOK)

    //Check if the long URL -> minID key-value pair is in memory
    id, err := client.Get(ctx, shortenreq.Url).Result()

    // If the miniURL id exists, then output the mini URL to w
    if err == nil {
 
        fmt.Fprintf(w, "http://localhost:8080/%s", id) 
    
    } else {
        
        //If the mini URL does not exist, then create a new one
        
        minID := makeRandomString()            

        fmt.Printf("The min URL ID is %s\n ", minID)
       
        //Set the new key-value pair into memory

        _, err := client.Set(ctx, shortenreq.Url, minID, 0).Result()

        if err != nil {

            log.Fatal(err)

        }

        // Output shortened URl to w
        
        fmt.Fprintf(w, "http://localhost:8080/%s", minID) 

    }

/*
    m.Mutex.Lock()
    
    //check if original url is in map
    id, ok := m.Map[shortenreq.Url]
    
    m.Mutex.Unlock()

    if ok {

        fmt.Fprintf(w, "http://localhost:8080/%s", id) 
    
    } else {

        minID := makeRandomString()            

        fmt.Printf("The min URL ID is %s\n ", minID)
        
        fmt.Fprintf(w, "http://localhost:8080/%s", minID) 

        //Save minID to map
        //m.Map[shortenreq.Url] = minID
    }

*/    

    w.Write([]byte("\n"))
}

//This handles a request to redirect the shortened url to the original url
func RedirectHandler(w http.ResponseWriter, req *http.Request) {

    
    fmt.Printf("The URL path of the request is %s\n", req.URL.Path)
    minID := strings.TrimPrefix(req.URL.Path, "/")

    //Try to find original URL (which is a key) in memory

    keys, err := client.Keys(ctx, "*").Result() 

    if err != nil {

        log.Fatal(err)

    }

    found := false

    var origURL string

    for _, key := range keys {

        id, err := client.Get(ctx, key).Result()
        
        if err != nil {

            log.Fatal(err)

        }


        if id == minID {

            origURL = key
            found = true
            break

        }

    }


    if found {

        http.Redirect(w, req, origURL, http.StatusPermanentRedirect)
    
    } else {

        w.WriteHeader(http.StatusBadRequest)

    }
                        
/*
    found := false

    var origURL string

    /*
    m.Mutex.Lock()
    for key, value := range m.Map {

        if string(value) == (minID) {
            
            origURL = key
            found = true
            break
        }

    }

    m.Mutex.Unlock()
    
    if found {

        http.Redirect(w, req, origURL, http.StatusPermanentRedirect)

    } else {

        w.WriteHeader(http.StatusBadRequest)
    }

*/
}


//This handler checks if the server is healthy
func HealthCheckHandler(w http.ResponseWriter, req *http.Request) {
    
    data := map[string]string{
        
        "health":      "ok",
    }
    
    w.Header().Set("Content-Type", "application/json")

    js, err := json.Marshal(data)

    if err != nil {

        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write(js)
}


        
    
