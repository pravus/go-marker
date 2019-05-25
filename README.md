# marker

It's not magic, but it marks up Markdown files.

## Usage

Marker runs as a stand-alone HTTP service.  Options for the listen address
and syntax highlighter style may be passed on the command line:

```bash
$ marker -http 127.0.0.1:8000 -syntax pygments
```

Once running, `.md` files may be viewed by pointing a browser to `http://localhost:8000/{basename}`.
Marker will try to render the file named `{basename}.md` from the current directory.  For example,
to render this README from the project root one might use the following:

http://localhost:8000/README

## See Also

 * [Go Markdown](https://github.com/gomarkdown/markdown)
 * [Chroma](https://github.com/alecthomas/chroma)
