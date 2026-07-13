# Overview

The idea for this program is as follows.

- take the original body.xhtml
- copy from body.xhtml each chapter and paste into a file with a name like chapter-1-orig.xhtml
- edit the original to be enclosed in body element tags just to make valid xml.
- might have to change initial paragraph from div to p tags
- copy that same chapter from my real repo for the book
- edit the real chapter to remove headers, leaving only body/section enclosing elements.
- run against that chapter (see example below)
- fix any false positives (such as changing "some one" to "someone" or "any one" to "any one")
- continue until they match perfectly (as shown in example) or you find a real descrepancy.
- if a real descrepancy is found, fix the real book repo
- continue until done.


Example:

```
$ go run main.go chapter-1-orig.xhtml chapter-1.xhtml 
Mismatch found at paragraph index 47!
------------------------------------------------------------
Filename: chapter-1-orig.xhtml (Index: 47)
Up to mismatch: "true but i do not think he was mad not from what you have told me but let us see what the commotion is some"

Filename: chapter-1.xhtml (Index: 47)
Up to mismatch: "true but i do not think he was mad not from what you have told me but let us see what the commotion is someone"
------------------------------------------------------------
exit status 1
$ go run main.go chapter-1-orig.xhtml chapter-1.xhtml 
Success! Both files have the same number of paragraphs, and each corresponding paragraph has the exact same word count.
$ 
```