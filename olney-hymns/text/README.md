# README

## Links
- Go package for Roman Numerals: https://pkg.go.dev/github.com/brandenc40/romannumeral




## Specs for script

1. File head: book-head.xmlfrag
2. File tail: book-end.xmlfrag
3. Appendix A for input from manager on book file and per hymn semantics/elements
4. Appendix B for how I have setup each hymn for input to script. Here are the salient points:
- `<h>` begins a hymn
- first line is ordinal representing number of hymn
- second line is title and scripture reference, which are split to become a title and a bridgehead
- then you have all the stanzas with all lines left justified
- need to indent per the meter of the hymn



# Appendix A

```xml

<section id="hymn-1-3" epub:type="z3998:hymn">
	<header>
		<hgroup>
			<h3 epub:type="ordinal">3</h3>
			<p epub:type="title">Walking with God</p>
		</hgroup>
        <p epub:type="z3998:author">William Cowper</p>
		<p epub:type="bridgehead">Gen. v,24</p>
	</header>
	<section id="stanza-1-3-1">
		<header>
			<p>I</p>
		</header>
			<p>
				<span>...</span>
			</p>
	</section>
	<section id="stanza-1-3-2">
		<header>
			<p>II</p>
		</header>
			<p>
				<span>...</span>
			</p>
	</section>
	...
```

# Appendix B
```
<h>
1.
Adam.—Gen. iii.

On man, in his own image made,
How much did God bestow!
The whole creation homage paid,
And own'd him lord below.

He dwelt in Eden's garden, stored
With sweets for every sense;
And there, with his descending Lord,
He walk'd in confidence.

But, oh by sin how quickly changed—
His honour forfeited—
His heart from God and truth estranged—
His conscience fill'd with dread!

But when by faith the sinner sees
A pardon bought with blood,
Then he forsakes his foolish pleas
And gladly turns to God.
</h>
```

# AI Promps

Write a golang command line progam that does the following:
    1. Takes a single argument on the command line which is a file name
    2. Read the file into an array of strings
    3. The file will have at least 3800 lines
    4. After the file is read into the array, then have a loop that outputs each element in the array to standard ouput.
    5. At completion, end the program.
    
I have a large chunk of XML text that I want to output first. Change the program to output this chunk of text first. Use this for the XML chunk:
<?xml version="1.0" encoding="utf-8"?>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" epub:prefix="z3998: http://www.daisy.org/z3998/2012/vocab/structure/, se: https://standardebooks.org/vocab/1.0" xml:lang="en-US">
<head>
	<title>Chapter 1</title>
	<link href="../css/core.css" rel="stylesheet" type="text/css"/>
	<link href="../css/local.css" rel="stylesheet" type="text/css"/>
</head>
<body epub:type="bodymatter z3998:fiction">
<section id="book-1" epub:type="part">
    