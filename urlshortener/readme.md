The goal of this exercise is to create an http.Handler that will look at the path of any incoming web request and determine if it should redirect the user to a new page, much like URL shortener would.

For instance, if we have a redirect setup for /dogs to https://www.somesite.com/a-story-about-dogs we would look for any incoming web requests with the path /dogs and redirect them.


### **Handlers Added**:
#### These handlers create `http.HandlerFunc` out of data passed
- [x]   MapHandler
        Takes in a map of string to string (src to dest) 
- [x]   JSONFileHandler
        Takes JSON filename
        Data Format:
        
        [
          {
            "path": "/urlshort",
            "url": "/long-url"
          }
        ]
- [x]   YAMLHandler
        Takes in YAML string with format:
        
        - path: /urlshort
          url: /long-url
- [x]   YAMLFileHandler
        Takes YAML filename
        Format:same as YAMLHandler

Will may be add DB support in future