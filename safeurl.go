package main

import (

    "net/url"
    "strings"
    "fmt"
)

func IsSafe(rawURL string) bool {

    // A string slice containing many unsafe TLDs
    unsafeTLDs := []string{".xyz", ".buzz", ".shop", ".icu", ".lol", ".sbs", ".cyou", ".cm", ".co", ".ga"}
    
    u, err := url.Parse(rawURL)

    // If there is an error parsing the url, say there is an error and return false
    if err != nil {

        fmt.Println("Error parsing raw URL")
        return false

    }

    // If the length of the raw URL is greater than 100 characters, then return false
    if len(rawURL) > 100 {

        return false

    }

    //check if raw url has the https protocol
    if u.Scheme != "https" {

        return false

    }

    //check if raw url is an IP address
    hostname := u.Hostname()

    count := 0
    
    for i := 0; i < len(hostname); i++{

        if string(hostname[i]) == "." {

            count += 1

            if count == 3 {

                return false 

            }
        }

    }

    if !strings.Contains(rawURL, "://") {

        return false

    }

    //check if raw url has a suspicious TLD
    for _, value := range unsafeTLDs {

       if strings.Contains(rawURL, value) && !strings.Contains(rawURL, ".com") {
            
            return false

       }
		
	} 

    return true

}

            
        
        


