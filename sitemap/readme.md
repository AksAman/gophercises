Exercise Detail
==============
[link to original exercise](https://github.com/gophercises/sitemap)

A sitemap is basically a map of all of the pages within a specific domain. They are used by search engines and other tools to inform them of all of the pages on your domain.

One way these can be built is by first visiting the root page of the website and making a list of every link on that page that goes to a page on the same domain. For instance, on `calhoun.io` you might find a link to `calhoun.io/hire-me/` along with several other links.

Once you have created the list of links, you could then visit each and add any new links to your list. By repeating this step over and over you would eventually visit every page that on the domain that can be reached by following links from the root page.

In this exercise your goal is to build a sitemap builder like the one described above. The end user will run the program and provide you with a URL (*hint - use a flag or a command line arg for this!*) that you will use to start the process.

Once you have determined all of the pages of a site, your sitemap builder should then output the data in the following XML format:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>http://www.example.com/</loc>
  </url>
  <url>
    <loc>http://www.example.com/dogs</loc>
  </url>
</urlset>
```

*Note: This should be the same as the [standard sitemap protocol](https://www.sitemaps.org/index.html)*

Where each page is listed in its own `<url>` tag and includes the `<loc>` tag inside of it.