package main

import (

    "net/url"
    "strings"
    "fmt"
)

func IsSafe(rawURL string) bool {

    unsafeTLDs := []string{".xyz", ".buzz", ".shop", ".icu", ".lol", ".sbs", ".cyou", ".cm", ".co", ".ga"}
    
    u, err := url.Parse(rawURL)

    if err != nil {

        fmt.Println("Error parsing raw URL")
        return false

    }

    if len(rawURL) > 100 {

        return false

    }

    //check if raw url has a secure protocol
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

            
        
        


