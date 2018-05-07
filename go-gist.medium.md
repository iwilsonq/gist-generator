There is a lot of busywork that comes with writing blogs and posting them to different platforms.

For example, I usually start writing an article by creating markdown file within my code editor. Then, as I write the post, I have to insert code snippets and images.

I am grateful for the ability to decorate my posts with code and images in a way that makes sense for the kind of writing that I do.

However, problems arise when I need to post this same markdown to [Medium](https://medium.com/@iwilsonq). Although their text editor is pretty, it is not perfectly safe to copy and paste my markdown files there.

Here are some issues:

*   code words such as `fields` or `someFunction` wrapped in backticks are not formatted properly on paste
*   larger code blocks do not currently support syntax highlighting

```javascript
function doSomething() {
	return function(name) {
		return 'Hello ' + name
	}
}
```

*   words that are **bold** or _italicized_ maintain their styling, but are not separated from the \*s that wrap them in my markdown file

And here is a screenshot of those 3 bullet points pasted:

![copypastemedium](unsafe_copy_paste.png)

Call me lazy, but after writing a 2500+ word post, I must:

*   go through each line and make sure my grammar makes sense.
*   assure my code examples work since some of my articles are tutorials.
*   make my code look pretty through [Github Gists](https://gist.github.com), thus bringing along with it a whole new issue of tabs vs spaces.
*   clean up every markdown-style hyperlink and remove excess asterisks or backticks from highlighted words and phrases.

This becomes draining. It becomes a lot of work to try to produce content that can be cross-posted to multiple sites.

Why dont I just cave in and post exclusively to [dev.to](https://dev.to/iwilsonq) where you literally just need to paste your markdown file?

Because that's not my nindo. I must become omnipresent throughout the internet.

## Cleaning Code Snippets

Let's start by grabbing the lowest hanging fruit, that is, clean up those code snippets! How can we do this?

Well, the brute force way would be to copy and paste as normal but, remove the backticks denoting code blocks, press "CMD+option+6" to create a Medium Codeblock&trade;, and then struggle with the awkward newlines and tabs within the Medium Codeblock&trade;.

Nope. I'm not gonna do that again.

The next option is to manually create Github Gists, which delegates the code formatting to Github, so we get syntax highlighting out of it.

However, in my most recent toils, the tab size seems to default to 8, which is absolutely ridiculous. Nobody uses size-8 tabs, _Github_. It seemed that the only way to fix this issue, is by replacing all instances of tabs with spaces (2 in my case).

The benefit is we get prettier code snippets. The drawback is that every tab must be replaced by spaces. My previous article had 18 JavaScript code snippets. It was a miserable editing experience.

But fear not, I have a plan that should help us to save time and, by association, money!

# Introducing the Gist Generator

\*_I hope I come up with a cooler name sometime_

I decided to take advantage of my dev chops to utilize the Github API in order to create these Gists. This script, written in Go, will go as follows:

1.  Read in the markdown file
2.  Parse designated language snippets like "go" or "javascript"
3.  For each snippet, create a gist
4.  Keep a struct with references to the snippet and the gist
5.  Replace the snippets with their associated gist URLs
6.  Write a new markdown file, with code snippets replaced by gists for Medium

If you interested in checking out the source code, you can check it out the repo [here](https://github.com/iwilsonq/gist-generator). I'd like to put together some tests before I decide to sell anybody on it, but it'll be open source in any case. Currently the usage goes like this:

```
./gist-generator -f example.md -lang javascript -token <GITHUB_ACCESS_TOKEN>
```

It would probably be easier if the user didn't need to specify the language of their snippets, maybe they want to use multiple different languages in their article? Add that to the todo list.

If I make a web client from which this process would be run, the user would sign in through OAuth with their Github account. This would eliminate the step where they have to manually generate their personal access token and paste it into the CLI. Again, that's something that would be nice to have but isn't in my crosshairs just yet.

Let's check out some more pressing issues.

## Problems Encountered

While solving the problems I cited earlier in this post was my main objective, there were a couple of places in writing this script where I ran into trouble.

One such problem was handling the formatting of the gist that would be created. What looks good in my markdown file may not necessarily look good in the gist.

In my first pass, I had excess whitespace due to newlines and oversized tabs. By iterating over the bytes in the code snippet, I could remove unnecessary newlines and replaces tabs with spaces. Here's an example:

https://gist.github.com/2e15a9f4eb4f5b812fb311b06712abab

Now out of all of the Go code I wrote for this script, why did I show the trimming of the newlines and tabs first? Why didn't I show the structs I created to represent gists and snippets?

Because the whitespace issue was the key to solving one of my biggest problems stated in the beginning! It was the easiest problem that yielded the biggest results.

There were also some lower level issues that dealt with how to best parse files or build strings using byte buffers. I'll tackle this in a future post since it deserves its own.

## Lessons Learned and Future Projects

I loved this project for several reasons:

*   It was small enough that it wasn't too intimidating
*   The end goal was clear and practical: to improve the quality of my blog posts
*   I had the opportunity to practice developing a Go command line tool
*   It can easily fit into an ecosystem of similar blog-productivity tools

By that last point, I mean that the gist generation step may be one of several.

I could write a similar script that would eliminate the excess markdown characters from the blog text. I discovered one way to partially do this step via [this post](https://medium.com/@andymcfee/how-to-import-markdown-into-medium-c06dc981bd96).

Medium allows one to import a markdown file if its pasted in a gist (since it requires a URL). This causes it to handle backtick and asterisk wrapped words properly. It even handles the hyperlinks within your text.

To handle grammatical mistakes I could add a step that perhaps calls a simple Grammarly*ish* API. If an open source solution doesn't already exist, that might be another future project idea.

I could even create a dashboard that would allow me to easily delete gists if a mistake in my script cause faults in any of them. The current [official gist dashboard](https://gist.github.com) is rather inefficent for editing or deleting gists.

Finally, I've been floating around the idea of a cloud-based markdown editor that would be accessible from mobile. You know, when I'm waiting in line at Disneyland and want to edit one of my articles. Times like these.

# Wrapping Up

So much to do, so little time it seems. I suppose my next step is to cover those other aesthetic issues regarding importing stories from markdown.

At this point, any project that reduces the friction involved in writing technical content is valuable. This project arose from my desire to eliminate some of the pain associated with tons of redundant editing. Some pains are cause by lack of familarity with platforms, others are caused by an actual limitation in the platform.

Writing articles like this one and embarking on similar projects solves both problems: you learn the capababilites of that platform all the while improving the parts where it falls short.

A great learning experience, no doubt, in the name of productivity.

Again, [here](https://github.com/iwilsonq/gist-generator) is the link to the project's repo.

Curious for more posts or witty remarks? Give me some likes and follow me on [Medium](https://medium.com/@iwilsonq), [Github](https://github.com/iwilsonq) and [Twitter](https://twitter.com/iwilsonq)!
