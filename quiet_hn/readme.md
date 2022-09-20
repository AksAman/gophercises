# Exercise: Quiet Hacker News

Original link : [gophercises/quiet_hn](https://github.com/gophercises/quiet_hn)

## Exercise details

### Problem statement
- Writing a program that 
    - [ ] Creates an HTTP server
    - [ ] Serves a single page that displays the top N stories from Hacker News
    - [ ] In order to get the stories:
        - Use the [Hacker News API](https://github.com/HackerNews/API)
        - Base Endpoint: https://hacker-news.firebaseio.com/v0
    - [ ] Get the stories using goroutines and channels
    - [ ] Stories must retain their original order
    - [ ] Make sure you always get exact N stories, not more, not less (important with concurrency)
    - [ ] Implements Caching with 
        - [ ] In-memory cache 
        - [ ] Redis
    - [ ] Caching should consider race conditions
    - [ ] Implement background cache updating
    - [ ] Implement Rate-Limiting using channels

----

### Hacker News Endpoints

1. Top Stories: [/topstories.json](https://hacker-news.firebaseio.com/v0/topstories.json)
    
    Response: List of story ids as `[]int`
    Example:
    ```javascript
    [
        32918301,
        32916994,
        32916318,
        32911299,
        32913125,
        ....
    ]
    ```


1. Story: [/item/{id}.json](https://hacker-news.firebaseio.com/v0/item/8863.json)

    Stories, comments, jobs, Ask HNs and even polls are just items. They're identified by their ids, which are unique integers, and live under `/v0/item/<id>`.

    All items have some of the following properties, with required properties in bold:

    Field | Description
    ------|------------
    **id** | The item's unique id.
    deleted | `true` if the item is deleted.
    type | The type of item. One of "job", "story", "comment", "poll", or "pollopt".
    by | The username of the item's author.
    time | Creation date of the item, in [Unix Time](http://en.wikipedia.org/wiki/Unix_time).
    text | The comment, story or poll text. HTML.
    dead | `true` if the item is dead.
    parent | The comment's parent: either another comment or the relevant story.
    poll | The pollopt's associated poll.
    kids | The ids of the item's comments, in ranked display order.
    url | The URL of the story.
    score | The story's score, or the votes for a pollopt.
    title | The title of the story, poll or job. HTML.
    parts | A list of related pollopts, in display order.
    descendants | In the case of stories or polls, the total comment count.
    
    For example, a story: https://hacker-news.firebaseio.com/v0/item/8863.json?print=pretty

    ```javascript
    {
        "by" : "dhouston",
        "descendants" : 71,
        "id" : 8863,
        "kids" : [ 8952, 9224, ...],
        "score" : 111,
        "time" : 1175714200,
        "title" : "My YC app: Dropbox - Throw away your USB drive",
        "type" : "story",
        "url" : "http://www.getdropbox.com/u/2/screencast.html"
    }
    ```


---
### Learning Outcomes
- [ ] HTTP in Go
- [ ] Templates in Go
- [ ] Concurrency
- [ ] Channels
- [ ] Goroutines
- [ ] Caching
- [ ] Redis in Go
- [ ] Rate-limiting
