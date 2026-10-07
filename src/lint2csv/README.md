# README

This file converts the output of `se lint` to a CSV file. It is assumes:

- the input file is `lint.txt`
- the output file will `lint.csv`
- Run: go run lint2csv.go

This code was written by Gemini. See Code Details below.

## Lint Output

Here is an example of such lint output:

```
┏━━━━━━━┳━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃ Code  ┃ Severity      ┃ File                    ┃ Message                    ┃
┡━━━━━━━╇━━━━━━━━━━━━━━━╇━━━━━━━━━━━━━━━━━━━━━━━━━╇━━━━━━━━━━━━━━━━━━━━━━━━━━━━┩
│ f-007 │ Error         │ appendix-1.xhtml        │ File not listed in         │
│       │               │                         │ <spine>.                   │
├───────┼───────────────┼─────────────────────────┼────────────────────────────┤
│ m-045 │ Error         │ appendix-1.xhtml        │ Heading not found in the   │
│       │               │                         │ ToC.                       │
├───────┼───────────────┼─────────────────────────┼────────────────────────────┤
│       │               │                  Line 4 │ XIV                        │
├───────┼───────────────┼─────────────────────────┼────────────────────────────┤
│ s-009 │ Error         │ appendix-1.xhtml        │ <hgroup> element with only │
│       │               │                         │ one child.                 │
├───────┼───────────────┼─────────────────────────┼────────────────────────────┤
│       │               │                 Line 10 │ <hgroup>                   │
│       │               │                         │         <h2>               │
│       │               │                         │                 <span      │
│       │               │                         │ epub:type="se:label">Appe… │
│       │               │                         │                 <span      │
│       │               │                         │ epub:type="z3998:ordinal   │
│       │               │                         │ z3998:roman">I</span>      │
│       │               │                         │         </h2>              │
│       │               │                         │ </hgroup>                  │
├───────┼───────────────┼─────────────────────────┼────────────────────────────┤
│ m-045 │ Error         │ uncopyright.xhtml       │ Heading not found in the   │
│       │               │                         │ ToC.                       │
├───────┼───────────────┼─────────────────────────┼────────────────────────────┤
│       │               │                  Line 4 │ Uncopyright                │
└───────┴───────────────┴─────────────────────────┴────────────────────────────┘
```

## CSV Converted Output

Here is the CSV produced from above. Note that filenames and line numbers are split out into two different columns and empty cells inherit the value from above cell. This makes it more useful as a spreadsheet tool.

```
Code,Severity,File,Line,Message
f-007,Error,appendix-1.xhtml,,File not listed in <spine>.
m-045,Error,appendix-1.xhtml,,Heading not found in the ToC.
m-045,Error,appendix-1.xhtml,Line 4,XIV
s-009,Error,appendix-1.xhtml,,<hgroup> element with only one child.
s-009,Error,appendix-1.xhtml,Line 10,"<hgroup>
         <h2>
                 <span
 epub:type=""se:label"">Appe…
                 <span
 epub:type=""z3998:ordinal
 z3998:roman"">I</span>
         </h2>
 </hgroup>"
m-045,Error,uncopyright.xhtml,,Heading not found in the ToC.
m-045,Error,uncopyright.xhtml,Line 4,Uncopyright
```

## Code Details

This code was written by Gemini and tested by myself. Here are the prompts used:

1. The attached file is tabular using UTF-8 characters to draw the cells. Write a function in Go that converts this to a CSV file so that I can view it with a spreadsheet tool.
2. Let's enhance this code with the following:
    1. the first two columns are sometimes empty. In those cases, copy the content from the cell above into the empty cell.
    2. the third column has both filenames and line numbers. The filenames are left-justified and the line number text is right-justified. Separate these into two different columns: one for the filenames and one for the line numbers. Again copy content from a prior cell to fill in empty cells.
3. There is one bug to fix. The first column in the last row has the UTF-8 cell boundary characters in it. It should have the code from the row above it (in this case "m-045").


