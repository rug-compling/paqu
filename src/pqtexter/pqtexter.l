/*
 * File: texter.l
 *
 * (c) Peter Kleiweg
 *     Wed May  7 09:58:58 2014
 *
 * Compile:
 *
 *    flex pqtexter.l
 *    gcc -s -Wall -o pqtexter pqtexter.c -lfl
 *    rm pqtexter.c
 *
 */

%{

#define my_VERSION "0.01"

#define __NO_MATH_INLINES

#ifdef __WIN32__
#  define my_PATH_SEP '\\'
#else
#  define my_PATH_SEP '/'
#endif

#ifdef __MSDOS__
#  ifndef __COMPACT__
#    error Memory model COMPACT required
#  endif  /* __COMPACT__  */
#  include <dir.h>
#endif  /* __MSDOS__  */
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

/* veelvoud van 4 */
#define BUFSIZE 4096

unsigned char    
    buffer [BUFSIZE];
char
    *programname;
int
    charset = 0,
    run = 0,
    nl = 0;

void
    out (long unsigned),
    out1 (unsigned),
    get_programname (char const *argv0),
    errit (char const *format, ...),
    syntax (void);
int
    get_charset ();


%}

%option 8bit batch
%option noyywrap nounput noinput
%option outfile="pqtexter.c"

SPC	([ \t\n\r\f]+)

%Start _u0001 _u1000 _u01 _u10 _u8 _win _b4u0001 _b4u1000 _b2u01 _b2u10 _b3u8

%%

<INITIAL>{

[\x00-\x7F]                     { /* kan us-ascii of utf-8 zijn */ }
[\xC0-\xDF][\x80-\xBF]          { /* kan utf-8 zijn */ }
[\xE0-\xEF][\x80-\xBF]{2}       { /* kan utf-8 zijn */ }
[\xF0-\xF7][\x80-\xBF]{3}       { /* kan utf-8 zijn */ }
[\xF8-\xFB][\x80-\xBF]{4}       { /* kan utf-8 zijn */ }
[\xFC-\xFD][\x80-\xBF]{5}       { /* kan utf-8 zijn */ }
.                               {
                                  /* geen utf-8, dan gaan we uit van windows-1252 */
                                  charset = _win;
                                  yyterminate();
                                }
}

<_b4u0001>{

(.|\n){4}                       { /* skip BOM */ BEGIN _u0001; }

}

<_b4u1000>{

(.|\n){4}                       { /* skip BOM */  BEGIN _u1000; }

}

<_b2u01>{

(.|\n){2}                       { /* skip BOM */  BEGIN _u01; }

}

<_b2u10>{

(.|\n){2}                       { /* skip BOM */  BEGIN _u10; }

}

<_b3u8>{

(.|\n){3}                       { /* skip BOM */  BEGIN _u8; }

}

<_u1000>{

\n\0\0\0\r\0\0\0|\r\0\0\0\n\0\0\0|\r\0\0\0|\n\0\0\0     { out1 (10); }
....                            { out ((long unsigned)(unsigned char)yytext[0] + 
				       256 * (long unsigned)(unsigned char)yytext[1] +
				       256 * 256 * (long unsigned)(unsigned char)yytext[2] +
				       256 * 256 * 256 * (long unsigned)(unsigned char)yytext[3]); }

}

<_u0001>{

\0\0\0\n\0\0\0\r|\0\0\0\r\0\0\0\n|\0\0\0\r|\0\0\0\n     { out1 (10); }
....                            { out ((long unsigned)(unsigned char)yytext[3] +
				       256 * (long unsigned)(unsigned char)yytext[2] +
				       256 * 256 * (long unsigned)(unsigned char)yytext[1] +
				       256 * 256 * 256 * (long unsigned)(unsigned char)yytext[0]); }

}

<_u10>{

\n\0\r\0|\r\0\n\0|\r\0|\n\0     { out1 (10); }
..                              { out ((long unsigned)(unsigned char)yytext[0] + 256 * (long unsigned)(unsigned char)yytext[1]); }

}

<_u01>{

\0\n\0\r|\0\r\0\n|\0\r|\0\n     { out1 (10); }
..                              { out ((long unsigned)(unsigned char)yytext[1] + 256 * (long unsigned)(unsigned char)yytext[0]); }

}

<_u8>{

\n\r|\r\n|\r|\n                 { out1 (10); }
.                               { out1 ((unsigned char)yytext[0]); }

}

<_win>{

\n\r|\r\n|\r|\n { out1 (10); }

\x80            { out (0x20ac); }
\x81            { out (0xfffe); }
\x82            { out (0x201a); }
\x83            { out (0x0192); }
\x84            { out (0x201e); }
\x85            { out (0x2026); }
\x86            { out (0x2020); }
\x87            { out (0x2021); }
\x88            { out (0x02c6); }
\x89            { out (0x2030); }
\x8a            { out (0x0160); }
\x8b            { out (0x2039); }
\x8c            { out (0x0152); }
\x8d            { out (0xfffe); }
\x8e            { out (0x017d); }
\x8f            { out (0xfffe); }
\x90            { out (0xfffe); }
\x91            { out (0x2018); }
\x92            { out (0x2019); }
\x93            { out (0x201c); }
\x94            { out (0x201d); }
\x95            { out (0x2022); }
\x96            { out (0x2013); }
\x97            { out (0x2014); }
\x98            { out (0x02dc); }
\x99            { out (0x2122); }
\x9a            { out (0x0161); }
\x9b            { out (0x203a); }
\x9c            { out (0x0153); }
\x9d            { out (0xfffe); }
\x9e            { out (0x017e); }
\x9f            { out (0x0178); }

[\x00-\x7f]     { out1 ( yytext[0]); }
.               { out ((unsigned char) yytext[0]); }

}



%%

int main (int argc, char *argv [])
{
    char
	*filename;

    get_programname (argv [0]);

    while (argc > 1 && argv [1][0] == '-') {
	if (argv[1][1] == 'r') {
	    run = 1;
	} else {
	    syntax ();
	    return 1;
	}
	argc--;
	argv++;
    }

    if (argc == 2) {
	filename = argv[1];
    } else {
	syntax();
	return 0;
    }

    yyin = fopen (filename, "r");
    if (! yyin)
	errit ("Opening file \"%s\": %s", filename, strerror (errno));

    charset = get_charset ();         /* test: utf8-bom, utf-16(-bom), utf-32(-bom) */
    if (!charset) {
        charset = _u8;                /* veronderstel: utf-8 */
	fseek (yyin, 0, SEEK_SET);
	yylex ();                     /* test: windows-1252 */
	yyrestart (yyin);
    }

    fseek (yyin, 0, SEEK_SET);
    BEGIN (charset);
    yylex ();                         /* eigenlijke verwerking */

    fclose (yyin);

    if (run)
	putchar (10);

    return 0;
}

void out (long unsigned c) {
    if (c < 0x80) {
	out1 (c);
	return;
    }
    if (c < 0x800) {
	out1 ((c>>6  & 0x1F) | 0xC0);
	out1 ((c>>0  & 0x3F) | 0x80);
	return;
    }
    if (c < 0x10000) {
	out1 ((c>>12 & 0x0F) | 0xE0);
	out1 ((c>>6  & 0x3F) | 0x80);
	out1 ((c>>0  & 0x3F) | 0x80);
	return;
    }
    if (c < 0x200000) {
	out1 ((c>>18 & 0x07) | 0xF0);
	out1 ((c>>12 & 0x3F) | 0x80);
	out1 ((c>>6  & 0x3F) | 0x80);
	out1 ((c>>0  & 0x3F) | 0x80);
	return;
    }

    if (c < 0x4000000) {
	out1 ((c>>24 & 0x03) | 0xF8);
	out1 ((c>>18 & 0x3F) | 0x80);
	out1 ((c>>12 & 0x3F) | 0x80);
	out1 ((c>>6  & 0x3F) | 0x80);
	out1 ((c>>0  & 0x3F) | 0x80);
	return;
    }

    out1 ((c>>32 & 0x01) | 0xFC);
    out1 ((c>>24 & 0x3F) | 0x80);
    out1 ((c>>18 & 0x3F) | 0x80);
    out1 ((c>>12 & 0x3F) | 0x80);
    out1 ((c>>6  & 0x3F) | 0x80);
    out1 ((c>>0  & 0x3F) | 0x80);

}

void out1 (unsigned c) {
    if (! run) {
	putchar (c);
	return;
    }
    if (c == 12) {
	putchar (10);
	nl = 0;
	return;
    }
    if (c == 10) {
	nl++;
	return;
    }
    if (c == 32 || c == 9) {
	putchar (32);
	return;
    }
    if (nl == 1)
	putchar (32);
    else if (nl > 1)
	putchar (10);
    putchar (c);
    nl = 0;
}

/* Test alleen utf-8 (met BOM), utf-16 en utf-32 (met/zonder BOM) */
int get_charset() {
    int
	i,
	u1000,
	u0001,
	u10,
	u01,
	n;

    n = fread (buffer, 1, BUFSIZE, yyin);

    if (n > 3 && buffer[0] == 0xff && buffer[1] == 0xfe && buffer[2] == 0 && buffer[3] == 0)
	return _b4u1000;
    if (n > 3 && buffer[0] == 0 && buffer[1] == 0 && buffer[2] == 0xfe && buffer[3] == 0xff)
	return _b4u1000;
    if (n > 1 && buffer[0] == 0xff && buffer[1] == 0xfe)
	return _b2u10;
    if (n > 1 && buffer[0] == 0xfe && buffer[1] == 0xff)
	return _b2u01;
    if (n > 2 && buffer[0] == 0xef && buffer[1] == 0xbb && buffer[2] == 0xbf)
	return _b3u8;

    u1000 = u0001 = 0;
    for (i = 0; i < n; i += 4) {
	if (buffer[i] == 0 && buffer[i+1] == 0 && buffer[i+2] == 0 && buffer[i+3] != 0)
	    u0001++;
	else if (buffer[i] != 0 && buffer[i+1] == 0 && buffer[i+2] == 0 && buffer[i+3] == 0)
	    u1000++;
    }
    if (u1000 > 1 && u1000 > n/4 * .9)
	return _u1000;
    if (u0001 > 1 && u0001 > n/4 * .9)
	return _u0001;
    
    u10 = u01 = 0;
    for (i = 0; i < n; i += 2) {
	if (buffer[i] == 0 && buffer[i+1] != 0)
	    u01++;
	else if (buffer[i] != 0 && buffer[i+1] == 0)
	    u10++;
    }
    if (u10 > 1 && u10 > n/2 * .9)
	return _u10;
    if (u01 > 1 && u01 > n/2 * .9) {
	return _u01;
    }

    return 0;
}


void errit (char const *format, ...)
{
    va_list
	list;

    fprintf (stderr, "\nError %s: ", programname);

    va_start (list, format);
    vfprintf (stderr, format, list);

    fprintf (stderr, "\n\n");

    exit (1);
}

void get_programname (char const *argv0)
{
#ifdef __MSDOS__
    char
        name [MAXFILE];
    fnsplit (argv0, NULL, NULL, name, NULL);
    programname = strdup (name);
#else
    char
        *p;
    p = strrchr (argv0, my_PATH_SEP);
    if (p)
        programname = strdup (p + 1);
    else
        programname = strdup (argv0);
#endif    
}

void syntax ()
{
    fprintf (
	stderr,
	"\n"
	"Version " my_VERSION "\n"
	"\n"
	"Usage: %s [-r] textfile\n"
	"\n"
	"  -r : reformat running text\n"
	"\n"
	"Dit programma heeft drie functies:\n\n"
	" 1. Zet alle varianten van regeleinde om naar newline\n"
	" 2. Zet invoer om naar utf-8\n"
	" 3. Met optie -r: haal newlines weg, behalve tussen paragrafen\n\n"
	"Bij 2: er wordt verondersteld dat de invoer in een van de volgende tekensets is:\n"
	"utf-8, utf-16, utf-32 (alle met/zonder BOM), us-ascii, iso-8859-1, windows-1252\n\n"
	"Bij 3: een paragraafgrens bestaat uit 2 of meer regeleindes met alleen spaties\n"
	"of tabs ertussen, of een formfeed\n\n",
	programname
    );
    exit (1);
}
