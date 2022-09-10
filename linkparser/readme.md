Exercise Detail
==============
[link to original exercise](https://github.com/gophercises/link)


In this exercise your goal is create a package that makes it easy to parse an HTML file and extract all of the links (`<a href="">...</a>` tags). For each extracted link you should return a data structure that includes both the `href`, as well as the text inside the link. Any HTML inside of the link can be stripped out, along with any extra whitespace including newlines, back-to-back spaces, etc.

Links will be nested in different HTML elements, and it is very possible that you will have to deal with HTML similar to code below.

```html
<a href="/dog">
  <span>Something in a span</span>
  Text not in a span
  <b>Bold text!</b>
</a>
```

In situations like these we want to get output that looks roughly like:

```go
Link{
  Href: "/dog",
  Text: "Something in a span Text not in a span Bold text!",
}
```
