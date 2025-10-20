# min-url
Creates a network service that shortens a secure URL, then redirects the shortened URL to the original URL, using only Golang, Docker, and Redis.

# HOW TO USE

1. Copy url to clipboard, then use git clone to clone repository into desired directory.

2. Once you're inside the min-url directory, enter ./min-url to run the Golang HTTP server. You will first see a "PONG <nil>" this indicates that the HTTP server is connected to the redis docker container.

3. While the HTTP server is running, enter in the following command in a separate terminal console:

   curl -X POST -H "Content-Type: application/json" -d '{"url":"https://k-state.edu"}' localhost:8080/shorten

   You will see a link produced that looks like http://localhost:8080/[some arrangement of letters]

   This new link will redirect you to the original url that was sent to be shortened.

   If you want to test other longer urls, just change the part of '{"url":"https://k-state.edu"}' by replacing the example link with any link you want.

   If the link does not contain https, or has a suspicious TLD (Top Level Domain) like .xyz, .co, and so on, instead of .com, or .org, it will reject shortening the url, and will just send an error message to the CLI.
   If your url has 3 periods, then it will say it is an unsafe url, because IP addresses as websites are not considered safe. If your url has more than 100 characters, then you will also get an error message. Your url also must contain a
   "://"

5. Restart the server by using CTRL + C, then running the executable again by using ./min-url. If you restart the server, and use a previous url that you used before, you will see that you get the same shortened link as before, because the program saved your previous input into a redis database that the server connects to whenever it is running.
 
