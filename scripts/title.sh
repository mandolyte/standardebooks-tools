#!/bin/sh

if [ "$3" = "" ] ; then
    echo Usage: title roman-numeral arabic-numeral
    exit
fi
# make the title have title case per SE rules
TITLECASE=`se titlecase "$1"`

# ensure that the roman numeral is uppercased
echo "$text" | tr '[:lower:]' '[:upper:]'
ROMAN=`echo "$2" | tr '[:lower:]' '[:upper:]'`


cat  << EoF | wl-copy
<?xml version="1.0" encoding="utf-8"?>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" epub:prefix="z3998: http://www.daisy.org/z3998/2012/vocab/structure/, se: https://standardebooks.org/vocab/1.0" xml:lang="en-GB">

	<head>
		<title>${ROMAN}: ${TITLECASE}</title>
		<link href="../css/core.css" rel="stylesheet" type="text/css"/>
		<link href="../css/local.css" rel="stylesheet" type="text/css"/>
	</head>
	<body epub:type="bodymatter z3998:fiction">
		<section id="chapter-$3" epub:type="chapter">
			<hgroup>
				<h2 epub:type="z3998:ordinal z3998:roman">${ROMAN}</h2>
				<p epub:type="title">${TITLECASE}</p>
			</hgroup>
EoF
