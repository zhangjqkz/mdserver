package main

const defaultTemplate = `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{.FileName}}</title>
    <link rel="stylesheet" href="/markdown-lib/github-markdown.css">
    <link rel="stylesheet" href="/markdown-lib/default.css">
    <style>
        .markdown-body {
            box-sizing: border-box;
            min-width: 200px;
            max-width: 980px;
            margin: 0 auto;
            padding: 45px;
        }
        @media (max-width: 767px) {
            .markdown-body {
                padding: 15px;
            }
        }
    </style>
</head>
<body>
    <article class="markdown-body">{{.Content}}</article>
    <script src="/markdown-lib/highlight.pack.js"></script>
    <script>hljs.initHighlightingOnLoad();</script>
</body>
</html>
`