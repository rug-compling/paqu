package main

import (
	"github.com/pebbe/dbxml"

	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

const (
	SPOD_STD = "std"
	SPOD_DX  = "dx"
)

type Spod struct {
	header  string
	xpath   string
	method  string
	lbl     string
	text    string
	special string
}

type spod_writer struct {
	header map[string][]string
	buffer bytes.Buffer
}

var (
	spods = []Spod{
		// hidden objecten moeten aan het begin, omdat anderen ervan afhankelijk zijn
		{
			"",
			`//node [@pos="verb"]`, // spatie i.v.m. mogelijk duplicaat
			SPOD_STD,
			"has_pos_verb",
			"has @pos=verb",
			"hidden1",
		},
		{
			"",
			"//parser [@cats and @skips]", // spatie i.v.m. mogelijk duplicaat
			SPOD_STD,
			"has_parser",
			"has <parser>",
			"hidden1",
		},
		{
			"",
			"//node [@his]", // spatie i.v.m. mogelijk duplicaat
			SPOD_STD,
			"has_his",
			"has @his",
			"hidden1",
		},
		{
			"",
			`//node [@word="?"]`, // spatie i.v.m. mogelijk duplicaat
			SPOD_STD,
			"has_qm",
			"has @word=?",
			"hidden1",
		},
		{
			"",
			`//node [@stype="ynquestion"]`, // spatie i.v.m. mogelijk duplicaat
			SPOD_STD,
			"has_yn",
			"has @stype=ynquestion",
			"hidden1",
		},
		{
			"",
			`//node [@sc]`, // spatie i.v.m. mogelijk duplicaat
			SPOD_STD,
			"has_sc",
			"has @sc",
			"hidden1",
		},
		{
			`Attributen`,
			`//node [@pos]`, // spatie i.v.m. mogelijk duplicaat
			SPOD_STD,
			"pos",
			"pos",
			"attr",
		},
		{
			``,
			`//node [@postag]`, // spatie i.v.m. mogelijk duplicaat
			SPOD_STD,
			"postag",
			"postag",
			"attr",
		},
		{
			``,
			`//node [@pt]`, // spatie i.v.m. mogelijk duplicaat
			SPOD_STD,
			"pt",
			"pt",
			"attr",
		},
		{
			`Hoofdzinnen//Bij deze queries vergelijken we de verschillende typen hoofdzinnen:
mededelende hoofdzinnen (1), vraagzinnen die met een vraag-constituent
beginnen (2), ja/nee vragen (3), en imperatieven (4).
<p>
<table>
<tr><td>(1)<td> Pieter leest een boek.</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(2)<td>Wie leest er een boek?</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(3)<td>Lees jij een boek?</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(4a)<td>Lees dat boek nou toch eens.</td>
<tr><td>( b)<td> Lees jij dat boek nou toch eens.</tr>
</table>
<p>
Het onderscheid tussen ja/nee vragen en imperatieven is in de
treebanks niet altijd nauwkeurig te bepalen - met name voor
de handmatig geannoteerde treebanks die iets minder informatie
bevatten dan de automatisch geannoteerde treebanks. Bijvoorbeeld
imperatieven die wel een onderwerp bevatten (4b) zullen door de
bestaande queries niet worden gevonden als imperatief.`,
			`//node[@cat="smain"]`,
			SPOD_STD,
			"smain",
			"mededelende hoofdzinnen",
			"",
		},
		{
			"",
			`//node[@cat="whq"]`,
			SPOD_STD,
			"whq",
			"vraagzinnen (wh)",
			"",
		},
		{
			"",
			`%PQ_janee_vragen%`,
			SPOD_STD,
			"janee",
			"ja/nee vragen",
			"qm -yn",
		},
		{
			"",
			`%PQ_imperatieven%`,
			SPOD_STD,
			"imp",
			"imperatieven",
			"",
		},
		{
			`Bijzinnen//Deze queries vergelijken de verschillende soorten bijzinnen.
We onderscheiden de volgende bijzinnen: ingebedde vraagzinnen
met een vraag-constituent (1), finiete bijzinnen (2), infiniete
bijzinnen met "om te" (3) infiniete bijzinnen met alleen "te" (4),
en infiniete bijzinnen met een ander voorzetsel (5). Daarnaast
onderscheiden we relatieve bijzinnen (6) en free relatives (7).
<p>
<table>
<tr><td>(1)<td>(ik vroeg) wie dat boek gelezen heeft</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(2a)<td>(ik lees dat boek) omdat het door Elsschot geschreven is</tr>
<tr><td>( b)<td>(ik denk) dat het boek door Elsschot is geschreven</tr>
<tr><td>( c)<td>(ik vroeg) of het boek door Ellschot is geschreven</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(3a)<td>(ik heb geprobeerd) om een boek te lezen</tr>
<tr><td>( b)<td>(ik ga naar de bieb) om een boek te lezen</tr>
<tr><td>( c)<td>(het was niet nodig) om het boek te lezen</tr>
<tr><td>( d)<td>(slangen zijn) om van te gruwen</tr>
<tr><td>( e)<td>(het kind is oud genoeg) om alleen naar school te gaan</tr>
<tr><td>( f)<td>(het boek is te duur) om te kopen</tr>
<tr><td>( g)<td>(hij heeft voldoende invloed) om het boek te verkopen</tr>
<tr><td>( h)<td>(ze was zo genadig) om het boek voor te lezen</tr>
<tr><td>( i)<td>(een boek) om nooit te vergeten</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(4)<td>(ik heb geprobeerd) een boek te lezen</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(5)<td>(ik lees dat boek) zonder mijn oordeel uit te spreken</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(6)<td>(ik lees een boek) dat door Elsschot is geschreven</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(7)<td>wie dit leest (is gek)</tr>
</table>
<p>
Binnen de finiete bijzinnen maken we onderscheid tussen bijzinnen die
ingeleid worden met "dat" (2b), met "of" (2c) dan wel met een ander
voegwoord (2a).
<p>
De infiniete bijzinnen met "om" worden verder onderverdeeld naar
gelang de bijzin optreedt als complement (3a), of als bepaling (3b),
als onderwerp (3c), als predicatief complement (3d), of als
comperatief complement (3efgh). Infiniete bijzinnen die als
bepaling optreden worden verder onderverdeeld, afhankelijk of
de bepaling bij een werkwoord optreedt (3b) of bij een zelfstandig
naamwoord (3i).
`,
			`%PQ_ingebedde_vraagzinnen%`,
			SPOD_STD,
			"whsub",
			"ingebedde vraagzinnen",
			"",
		},
		{
			"",
			`
//node[@cat='cp'
       and
       node[@rel='body' and @cat='ssub']
]`,
			SPOD_STD,
			"ssub",
			"finiete bijzinnen||",
			"",
		},
		{
			"",
			`
//node[@cat='cp'
       and
       node[@rel='cmp' and @lemma='dat']
       and
       node[@rel='body' and @cat='ssub']
]`,
			SPOD_STD,
			"ssubdat",
			"finiete bijzinnen| met \"dat\"",
			"",
		},
		{
			"",
			`
//node[@cat='cp'
       and
       node[@rel='cmp' and @lemma='of']
       and
       node[@rel='body' and @cat='ssub']
]`,
			SPOD_STD,
			"ssubof",
			"finiete bijzinnen| met \"of\"",
			"",
		},
		{
			"",
			`
//node[@cat='cp'
       and
       node[@rel='cmp'
            and
            not(@lemma=('of','dat'))
       ]
       and
       node[@rel='body' and @cat='ssub']
]`,
			SPOD_STD,
			"ssubcmp",
			"finiete bijzinnen| met andere voegwoorden",
			"",
		},
		{
			"",
			`//node[@cat='oti']`,
			SPOD_STD,
			"oti",
			"infiniete bijzinnen met \"om\"||",
			"",
		},
		{
			"",
			`//node[@cat="oti" and @rel="vc"]`,
			SPOD_STD,
			"otivc",
			"infiniete bijzinnen met \"om\"| die als complement optreden",
			"",
		},
		{
			"",
			`//node[@cat="oti" and @rel="mod"]`,
			SPOD_STD,
			"otimod",
			"infiniete bijzinnen met \"om\"| die als bepaling optreden",
			"",
		},
		{
			"",
			`
//node[@cat="oti"
       and
       ../node[@rel="hd" and @pt="ww"]
]`,
			SPOD_STD,
			"otiww",
			"infiniete bijzinnen met \"om\"| die als bepaling bij een werkwoord optreden",
			"",
		},
		{
			"",
			`
//node[@cat="oti"
       and
       ../node[@rel="hd"
               and
               @pt=("n","vnw")
       ]
]`,
			SPOD_STD,
			"otin",
			"infiniete bijzinnen met \"om\"| die als bepaling bij een zelfstandig naamwoord optreden",
			"",
		},
		{
			"",
			`//node[@cat="oti" and @rel="su"]`,
			SPOD_STD,
			"otisu",
			"infiniete bijzinnen met \"om\"| die als onderwerp fungeren",
			"",
		},
		{
			"",
			`//node[@cat="oti" and @rel="predc"]`,
			SPOD_STD,
			"otipred",
			"infiniete bijzinnen met \"om\"| die als predicaat fungeren",
			"",
		},
		{
			"",
			`//node[@cat="oti" and @rel="obcomp"]`,
			SPOD_STD,
			"otiobc",
			"infiniete bijzinnen met \"om\"| die optreden met combinaties zoals \"te ADJ; zo ADJ; genoeg ADJ; voldoende N\"",
			"",
		},
		{
			"",
			`
//node[@cat='ti'
       and
       not(../@cat=('oti','cp'))
]`,
			SPOD_STD,
			"tite",
			"infiniete bijzinnen|| met alleen \"te\"",
			"",
		},
		{
			"",
			`
//node[@cat='cp'
       and
       node[@rel='body' and @cat='ti']
]`,
			SPOD_STD,
			"ti",
			"infiniete bijzinnen| met ander voorzetsel",
			"",
		},
		{
			"",
			`%PQ_relatieve_bijzinnen%`,
			SPOD_STD,
			"relssub",
			"relatieve bijzinnen",
			"",
		},
		{
			"",
			`%PQ_free_relatives%`,
			SPOD_STD,
			"whrel",
			"free relatives",
			"",
		},
		{
			`Correlatieve comparatieven//Correlatieve comparatieven zijn zinnen zoals:
<p>
<table>
<tr><td>(1)<td>Hoe langer, hoe beter</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(2)<td>Hoe langer hij wordt, des te meer hij last van zijn rug krijgt</tr>
</table>`,
			`%PQ_minimal_corr_comp%`,
			SPOD_STD,
			"corc",
			"totaal",
			"",
		},
		{
			`Woorden met een comparatief complement//Comparatieve adjectieven en woorden zo als "zo", "even", "meer", "minder",
"niet", "niets", "ander", "anders" treden soms op met een complement
dat vaak wordt ingeleid met het woord "dan" of "als". Naast het totale
aantal voorkomens van de constructie kijken we naar het aantal voorbeelden
waarbij "zo" (1) of "even" (2) of een comparatief adjective (3) of
"meer" of "minder" (4) of "niet", "niets", "ander", "anders" (5) het hoofd is.
<p>
Bij comparatieve adjectieven make we ook nog een onderverdeling naar
de aard van het complement. Een vergelijkbare onderverdeling
wordt gemaakt voor de andere hoofden.
<p>
<table>
<tr><td>(1)<td>zo groen als gras</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(2)<td>even dik als jij</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(3a)<td>beter dan jij</tr>
<tr><td>( b)<td>beter dan ik dacht</tr>
<tr><td>( c)<td>beter dan bij de buren</tr>
<tr><td>( d)<td>beter dan gisteren</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(4a)<td>meer geluk dan wijsheid</tr>
<tr><td>( b)<td>meer geluk dan ik dacht</tr>
<tr><td>( c)<td>meer geluk dan bij de buren</tr>
<tr><td>( d)<td>meer geluk dan gisteren</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(5a)<td>niets anders dan ellende</tr>
<tr><td>( b)<td>niets anders dan dat hij komt</tr>
<tr><td>( c)<td>niets anders dan bij de buren</tr>
<tr><td>( d)<td>niets anders dan gisteren</tr>
</table>
`,
			`//node[@rel="hd" and ../node[@rel="obcomp"]]`,
			SPOD_STD,
			"cc",
			"totaal",
			"",
		},
		{
			"",
			`//node[@rel="hd" and ../node[@rel="obcomp"] and @lemma="zo"]`,
			SPOD_STD,
			"cczo",
			"met als hoofd \"zo\"",
			"",
		},
		{
			"",
			`//node[@rel="hd" and ../node[@rel="obcomp"] and @lemma="even"]`,
			SPOD_STD,
			"cceven",
			"met als hoofd \"even\"",
			"",
		},
		{
			"",
			`//node[@rel="hd" and ../node[@rel="obcomp"] and @pt="adj" and @graad="comp"]`,
			SPOD_STD,
			"ccca",
			"met als hoofd comparatief adjectief||",
			"",
		},
		{
			"",
			`
//node[@rel="hd"
       and
       ../node[@rel="obcomp"
               and
               node[@rel="body" and %PQ_np%]
       ]
       and
       @pt="adj"
       and
       @graad="comp"
]`,
			SPOD_STD,
			"ccdannp",
			"met als hoofd comparatief adjectief|, gevolgd door NP",
			"",
		},
		{
			"",
			`
//node[@rel="hd"
       and
       ../node[@rel="obcomp"
               and
               node[@rel="body" and (%PQ_s% or @cat=("cp","ssub") or @pt="ww")]
       ]
       and
       @pt="adj"
       and
       @graad="comp"
]`,
			SPOD_STD,
			"ccdanvs",
			"met als hoofd comparatief adjectief|, gevolgd door VP of S",
			"",
		},
		{
			"",
			`
//node[@rel="hd"
       and
       ../node[@rel="obcomp"
               and
               node[@rel="body" and @cat="pp"]
       ]
       and
       @pt="adj"
       and
       @graad="comp"
]`,
			SPOD_STD,
			"ccdanpp",
			"met als hoofd comparatief adjectief|, gevolgd door PP",
			"",
		},
		{
			"",
			`
//node[@rel="hd"
       and
       ../node[@rel="obcomp"
               and
               node[@rel="body" and (@cat=("advp","ap") or @pt=("adj","bw"))]
       ]
       and
       @pt="adj"
       and
       @graad="comp"
]`,
			SPOD_STD,
			"ccdanav",
			"met als hoofd comparatief adjectief|, gevolgd door A of ADV",
			"",
		},
		{
			"",
			`//node[@rel="hd" and ../node[@rel="obcomp"] and @lemma=("veel","minder","weinig")]`,
			SPOD_STD,
			"ccmm",
			"met als hoofd \"meer\", \"minder\"||",
			"",
		},
		{
			"",
			`
//node[@rel="hd"
       and
       ../node[@rel="obcomp"
               and
               node[@rel="body" and %PQ_np%]
       ]
       and
       @lemma=("veel","minder","weinig")
]`,
			SPOD_STD,
			"ccmdnp",
			"met als hoofd \"meer\", \"minder\"|, gevolgd door NP",
			"",
		},
		{
			"",
			`
//node[@rel="hd"
       and
       ../node[@rel="obcomp"
               and
               node[@rel="body" and (%PQ_s% or @cat=("cp","ssub") or @pt="ww")]
       ]
       and
       @lemma=("veel","minder","weinig")
]`,
			SPOD_STD,
			"ccmdvs",
			"met als hoofd \"meer\", \"minder\"|, gevolgd door VP of S",
			"",
		},
		{
			"",
			`
//node[@rel="hd"
       and
       ../node[@rel="obcomp"
               and
               node[@rel="body" and @cat="pp"]
       ]
       and
       @lemma=("veel","minder","weinig")
]`,
			SPOD_STD,
			"ccmdpp",
			"met als hoofd \"meer\", \"minder\"|, gevolgd door PP",
			"",
		},
		{
			"",
			`
//node[@rel="hd"
       and
       ../node[@rel="obcomp"
               and
               node[@rel="body" and (@cat=("advp","ap") or @pt=("adj","bw"))]
       ]
       and
       @lemma=("veel","minder","weinig")
]`,
			SPOD_STD,
			"ccmdav",
			"met als hoofd \"meer\", \"minder\"|, gevolgd door A of ADV",
			"",
		},
		{
			"",
			`
//node[@rel="hd"
       and
       ../node[@rel="obcomp"]
       and
       @lemma=("niet","niets","ander","anders")
]`,
			SPOD_STD,
			"ccnn",
			"met als hoofd \"niet\", \"niets\", \"ander\", \"anders\"||",
			"",
		},
		{
			"",
			`
//node[@rel="hd"
       and
       ../node[@rel="obcomp"
               and
               node[@rel="body" and %PQ_np%]
       ]
       and
       @lemma=("niet","niets","ander","anders")
]`,
			SPOD_STD,
			"ccndnp",
			"met als hoofd \"niet\", \"niets\", \"ander\", \"anders\"|, gevolgd door NP",
			"",
		},
		{
			"",
			`
//node[@rel="hd"
       and
       ../node[@rel="obcomp"
               and
               node[@rel="body" and (%PQ_s% or @cat=("cp","ssub") or @pt="ww")]
       ]
       and @lemma=("niet","niets","ander","anders")
]`,
			SPOD_STD,
			"ccndvs",
			"met als hoofd \"niet\", \"niets\", \"ander\", \"anders\"|, gevolgd door VP of S",
			"",
		},
		{
			"",
			`
//node[@rel="hd"
       and
       ../node[@rel="obcomp"
               and
               node[@rel="body" and @cat="pp"]
       ]
       and @lemma=("niet","niets","ander","anders")
]`,
			SPOD_STD,
			"ccndpp",
			"met als hoofd \"niet\", \"niets\", \"ander\", \"anders\"|, gevolgd door PP",
			"",
		},
		{
			"",
			`
//node[@rel="hd"
       and
       ../node[@rel="obcomp"
               and
               node[@rel="body" and (@cat=("advp","ap") or @pt=("adj","bw"))]
       ]
       and @lemma=("niet","niets","ander","anders")
]`,
			SPOD_STD,
			"ccndav",
			"met als hoofd \"niet\", \"niets\", \"ander\", \"anders\"|, gevolgd door A of ADV",
			"",
		},
		{
			`Nevenschikkingen//Bij nevenschikkingen maken we onderscheid naar gelang het aantal
coordinatoren. Geen (1), één (2), twee (3) of meer (4).
<p>
<table>
<tr><td>(1)<td>Zij wilden inspraak , medezeggenschap , democratisering</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(2a)<td>Appels en peren</tr>
<tr><td>( b)<td>Appels of peren</tr>
<tr><td>( c)<td>De afzijdige maar invloedrijke waarnemer</tr>
<tr><td>( d)<td>Ze helpen eiwitten, vetten, etc. afbreken</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(3)<td>Noch appels noch peren</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(4)<td>U kunt gezond zijn, of ziek of arm of rijk</tr>
</table>
<p>
Indien er precies één coordinator is, maken we onderscheid
tussen "en" (2a), "of" (2b), "maar" (2c) en het speciale
geval waar de conjunctie wordt afgesloten met een woord
zoals "enzovoorts" dat in de annotatie ook als coordinator
wordt weergegeven (2d).
<p>
We maken ook een onderverdeling naar gelang het aantal conjuncten
dat in de coordinatie optreedt. Eén (5), twee (2), drie (6),
vier, vijf, zes, of meer.
<p>
<table>
<tr><td>(5)<td>de milieuverontreiniging, enzovoorts</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(6)<td>cacaoboter , cacaopoeder of palmolie</tr>
</table>
<p>
Ten slotte wordt een onderscheid gemaakt naar de categorie
van de conjuncten: NP (2), PP (7), hoofdzinnen (8),
VP (9) en bijzinnen (10).
<p>
<table>
<tr><td>(7)<td>in Arnhem en in België</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(8)<td>Drie keer hadden ze dat beloofd en drie keer was die
     belofte weer ingetrokken</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(9)<td>.. die sinds 1981 in Duitsland woont en werkt</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(10)<td>(Hij zei) dat het er bij hoorde, en dat ik niet bang hoefde te zijn</tr>
</table>`,
			`//node[@cat="conj"]`,
			SPOD_STD,
			"conj",
			"alle nevenschikkingen",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       not(node[@rel="crd"])
]`,
			SPOD_STD,
			"crd0",
			"nevenschikkingen zonder coördinator",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       count(node[@rel="crd"])=1
]`,
			SPOD_STD,
			"crd1",
			"nevenschikkingen met 1 coördinator||",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       count(node[@rel="crd"])=1
       and
       node[@rel="crd" and @lemma="en"]
]`,
			SPOD_STD,
			"crd1en",
			"nevenschikkingen met 1 coördinator|, en dat is \"en\"",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       count(node[@rel="crd"])=1
       and
       node[@rel="crd" and @lemma="of"]
]`,
			SPOD_STD,
			"crd1of",
			"nevenschikkingen met 1 coördinator|, en dat is \"of\"",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       count(node[@rel="crd"])=1
       and
       node[@rel="crd" and @lemma="maar"]
]`,
			SPOD_STD,
			"crd1maa",
			"nevenschikkingen met 1 coördinator|, en dat is \"maar\"",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       count(node[@rel="crd"])=1
       and
       node[@rel="crd"]/%PQ_e% = %PQ_e%
]`,
			SPOD_STD,
			"crd1enz",
			"nevenschikkingen met 1 coördinator|, en de coordinator sluit de nevenschikking af (\"enzovoorts\")",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       count(node[@rel="crd"])=2
]`,
			SPOD_STD,
			"crd2",
			"nevenschikkingen met 2 coördinatoren",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       count(node[@rel="crd"])>2
]`,
			SPOD_STD,
			"crd2p",
			"nevenschikkingen met meer dan 2 coördinatoren",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       count(node[@rel="cnj"])=1
]`,
			SPOD_STD,
			"cnj1",
			"nevenschikkingen|| met slechts 1 conjunct",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       count(node[@rel="cnj"])=2
]`,
			SPOD_STD,
			"cnj2",
			"nevenschikkingen| met 2 conjuncten",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       count(node[@rel="cnj"])=3
]`,
			SPOD_STD,
			"cnj3",
			"nevenschikkingen| met 3 conjuncten",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       count(node[@rel="cnj"])=4
]`,
			SPOD_STD,
			"cnj4",
			"nevenschikkingen| met 4 conjuncten",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       count(node[@rel="cnj"])=5
]`,
			SPOD_STD,
			"cnj5",
			"nevenschikkingen| met 5 conjuncten",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       count(node[@rel="cnj"])=6
]`,
			SPOD_STD,
			"cnj6",
			"nevenschikkingen| met 6 conjuncten",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       count(node[@rel="cnj"])>6
]`,
			SPOD_STD,
			"cnj6p",
			"nevenschikkingen| met meer dan 6 conjuncten",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       node[@rel="cnj" and %PQ_np%]
       and
       not(node[@rel="cnj"
                and
                not(%PQ_np%)
           ]
       )
]`,
			SPOD_STD,
			"cnjnp",
			"nevenschikking|| van NP's",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       node[@rel="cnj" and @cat="pp"]
       and
       not(node[@rel="cnj"
                and
                not(@cat="pp")
           ]
       )
]`,
			SPOD_STD,
			"cnjpp",
			"nevenschikking| van PP's",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       node[@rel="cnj" and @cat="smain"]
       and
       not(node[@rel="cnj"
                and
                not(@cat="smain")
           ]
       )
]`,
			SPOD_STD,
			"cnjmain",
			"nevenschikking| van hoofdzinnen",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       node[@rel="cnj"
            and
            (@cat=("ssub","ti","ppart","inf") or @pt="ww")
       ]
       and
       not(node[@rel="cnj"
                and
                not(@cat=("ssub","ti","ppart","inf") or @pt="ww")
           ]
       )
]`,
			SPOD_STD,
			"cnjvp",
			"nevenschikking| van VP",
			"",
		},
		{
			"",
			`
//node[@cat="conj"
       and
       node[@rel="cnj" and @cat="cp"]
       and
       not(node[@rel="cnj"
                and
                not(@cat="cp")
           ]
       )
]`,
			SPOD_STD,
			"cnjcp",
			"nevenschikking| van bijzinnen",
			"",
		},
		{
			`Woordgroepen//Voor woordgroepen onderscheiden we hier naamwoordgroepen (1), voorzetselgroepen
(2),
bijvoeglijk-naamwoordgroepen (3) en bijwoordgroepen (4).
<p>
<table>
<tr><td>(1)<td>een grote boom</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(2)<td>op een tak</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(3)<td>heel erg geliefd</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(4)<td>spelenderwijs</tr>
</table>`,
			`//node[%PQ_np%]`,
			SPOD_STD,
			"np",
			"np",
			"",
		},
		{
			"",
			`//node[@cat='pp']`,
			SPOD_STD,
			"pp",
			"pp",
			"",
		},
		{
			"",
			`
//node[@cat='ap'
       or
       @pt='adj'
       and
       not(@rel='hd')
]`,
			SPOD_STD,
			"ap",
			"ap",
			"",
		},
		{
			"",
			`
//node[@cat='advp'
       or
       @pt='bw'
       and
       not(@rel='hd')
]`,
			SPOD_STD,
			"advp",
			"advp",
			"",
		},
		{
			`Voorzetselgroepen//Voorzetselgroepen worden onderscheiden naar grammaticale 
functie en naar interne structuur. Voorzetselgroepen
kunnen optreden als bepaling bij nomina (1), adjectieven (2)
of bij werkwoorden (6). Dit laatste geval wordt in de traditionele
grammatica een bijwoordelijke bepaling genoemd.  Voorzetselgroepen
kunnen ook als complement optreden. Hier onderscheiden we
voorzetselvoorwerpen (3), locatief/directinele complementen (4) en
predicatieve complementen (5).
<p>
<table>
<tr><td>(1) <td>De vrouw van de buurman
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(2a)<td>De door hem beschreven voorvallen
<tr><td>( b)<td>De op zichzelf rationele beslissingen
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(3) <td>We rekenen op zijn steun
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(4a)<td>De ketchup vindt u bij de groenteafdeling
<tr><td>( b)<td>We rijden wel even naar het station
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(5a)<td>Die aanpak is niet zonder risico
<tr><td>( b)<td>Je moet in het bezit zijn van een visum
<tr><td>( c)<td>De haven is van cruciale betekenis
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(6) <td>We gaan vissen bij de brug
</table>
<p>
Voorzetselgroepen worden onderverdeeld aan de hand van hun
interne structuur. We onderscheiden de meest gangbare structuur
waarbij een voorzetsel direct gevolgd wordt door een NP (7)
en de structuur waarbij een R-pronomen direct of indirect gevolgd
wordt door het bijbehorende voorzetsel (8). Ten slotte volgt
een aparte telling voor het aantal voorzetseluitdrukkingen.
Dit zijn versteende combinaties zoals "ten tijde van", "door middel van",
"naar aanleiding van" die zich als voorzetsel gedragen.
<p>
<table>
<tr><td>(7)<td>in de kooi
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(8)<td>Hij houdt daar helemaal niet van
</table>`,
			`//node[@cat="pp" and @rel="mod" and ../@cat="np"]`,
			SPOD_STD,
			"ppnp",
			"grammaticale functie||, bepalingen bij zelfstandige naamwoorden",
			"",
		},
		{
			"",
			`//node[@cat="pp" and
       @rel="mod" and
       (../@cat="ap" or
        ((../@cat="ppart" or ../@cat="ppres") and
         not(../@rel="vc") and
         not(../@rel="cnj" and ../../@rel="vc")))]`,
			SPOD_STD,
			"ppap",
			"grammaticale functie|, bepalingen bij adjectieven (en als adjectief gebruikte deelwoorden)",
			"",
		},
		{
			"",
			`//node[@cat="pp" and @rel="pc"]`,
			SPOD_STD,
			"pppc",
			"grammaticale functie|, voorzetselvoorwerp",
			"",
		},
		{
			"",
			`//node[@cat="pp" and @rel="ld"]`,
			SPOD_STD,
			"ppld",
			"grammaticale functie|, locatief/directioneel complement",
			"",
		},
		{
			"",
			`//node[@cat="pp" and @rel="predc"]`,
			SPOD_STD,
			"pppredc",
			"grammaticale functie|, predicatief complement",
			"",
		},
		{
			"",
			`//node[@cat="pp" and
       @rel="mod" and
       (../@cat=("smain","sv1","whq","ssub","inf") or
        (../@cat="ppart" and
         (../@rel="vc" or
          (../@rel="cnj" and
           ../../@rel="vc"))))]`,
			SPOD_STD,
			"ppbep",
			"grammaticale functie|, bijwoordelijke bepaling",
			"",
		},
		{
			"",
			`//node[@cat="pp" and
       node[@rel="hd"]/%PQ_b%=%PQ_b% and
       node[@rel="obj1"]/%PQ_e%=%PQ_e%]`,
			SPOD_STD,
			"ppinp",
			"interne structuur||, P + NP",
			"",
		},
		{
			"",
			`//node[@cat="pp" and
       node[@rel="hd"]/%PQ_e%=%PQ_e% and
       node[@rel="obj1" and @pt="vnw"]/%PQ_b%=%PQ_b%]`,
			SPOD_STD,
			"ppirp",
			"interne structuur|, +R-pronomen + P",
			"",
		},
		{
			"",
			`//node[@cat="pp" and node[@rel="hd" and @cat="mwu"]]`,
			SPOD_STD,
			"ppimwu",
			"interne structuur|, complex voorzetsel",
			"",
		},
		{
			`Werkwoorden//Onder het kopje werkwoorden worden de vaste werkwoordelijke
uitdrukkingen geteld (1).
<p>
<table>
<tr><td>(1a)<td>aan bod komen</tr>
<tr><td>( b)<td>ter kennis geven</tr>
<tr><td>( c)<td>op prijs stellen</tr>
</table>
<p>
Daarnaast maken we een onderscheid bij werkwoordclusters
tussen de zogenaamde rode en groene volgorde. Bij de
groene werkwoordvolgorde gaat het voltooid deelwoord
in de werkwoordcluster aan het hulpwerkwoord vooraf (2),
en bij de rode volgorde volgt het voltooid deelwoord het
hulpwerkwoord (3).
<p>
<table>
<tr><td>(2a)<td>zijn boek zal morgen gepresenteerd worden</tr>
<tr><td>( b)<td>ze zouden nooit een woord met elkaar gewisseld hebben</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(3a)<td>zijn boek zal morgen worden gepresenteerd</tr>
<tr><td>( b)<td>ze zouden nooit een woord met elkaar hebben gewisseld</tr>
</table>
<p>
Voorbeelden van werkwoordclusters worden gegeven in (2), (3), en (4).
<p>
<table>
<tr><td>(4a)<td>Moet op korte termijn worden opgetreden?</tr>
<tr><td>( b)<td>De zaak moet aanhangig worden gemaakt</tr>
<tr><td>( c)<td>Ze zeiden dat ik dat zou hebben moeten kunnen zien aankomen</tr>
</table>
<p>
De "accusativus cum infinitivo" constructie vertoont de
beroemde 'cross-serial' afhankelijkheden. Voorbeelden zijn (5).
Deze query werkt niet voor het CGN corpus.
<p>
<table>
<tr><td>(5a)<td>Ik heb hem het boek zien lezen</tr>
<tr><td>( b)<td>Ik heb hem de olifanten helpen voeren</tr>
<tr><td>( c)<td>Ik heb hem de oppasser de olifanten zien helpen voeren</tr>
</table>
<p>
De passieve werkwoorden worden gebruikt in de lijdende vorm, als in
(6). Het gaat hierbij om de gevallen waarbij een onderwerp aanwezig
is. Deze query werkt niet voor het CGN corpus.
<p>
<table>
<tr><td>(6a)<td>Hij wordt geslagen</tr>
<tr><td>( b)<td>Het monument zal spoedig te zien zijn in het Amsterdamse Oosterpark .</tr>
</table>
<p>
De niet-persoonlijke passieven identificeert de gevallen van de lijdende
vorm waarbij geen echt onderwerp beschikbaar is. In die gevallen wordt
vaak "er" als plaatsonderwerp gebruikt.
<p>
<table>
<tr><td>(7a)<td>Er werd niet meer over gesproken</tr>
<tr><td>( b)<td>In Amsterdam werd niet gedemonstreerd</tr>
</table>
`,
			`//node[@rel="svp" and @cat]`,
			SPOD_STD,
			"vwuit",
			"vaste werkwoordelijke uitdrukkingen",
			"",
		},
		{
			"",
			`%PQ_groen%`,
			SPOD_STD,
			"groen",
			"groene werkwoordsvolgorde",
			"",
		},
		{
			"",
			`%PQ_rood%`,
			SPOD_STD,
			"rood",
			"rode werkwoordsvolgorde",
			"",
		},
		{
			"",
			`//node[%PQ_dep_node_in_verbcluster%]`,
			SPOD_STD,
			"wwclus",
			"werkwoordsclusters",
			"",
		},
		{
			"",
			`%PQ_cross_serial_verbcluster%`,
			SPOD_STD,
			"accinf",
			"accusativus cum infinitivo",
			"pos",
		},
		{
			"",
			`//node[%PQ_passive%]`,
			SPOD_STD,
			"passive",
			"passief",
			"pos",
		},
		{
			"",
			`//node[%PQ_impersonal_passive%]`,
			SPOD_STD,
			"nppas",
			"niet-persoonlijke passief",
			"pos",
		},
		{
			`Scheidbare werkwoorden//De queries voor scheidbare werkwoorden zijn niet beschikbaar voor de
handmatig geannoteerde corpora.
<p>
Werkwoorden met een scheidbaar partikel worden gegeven in (1) en (2).
<p>
<table>
<tr><td>(1a)<td>Ik heb hem opgebeld</tr>
<tr><td>( b)<td>Ik bel hem op</tr>
<tr><td>( c)<td>Ik heb hem op willen bellen</tr>
</table>
<p>
We maken daarbij een onderscheid tussen gevallen waarbij het partikel
gescheiden is van het werkwoord (1bc) en gevallen waarbij het partikel
in het werkwoord geincorporeerd is (1a).
<p>
Daarnaast bekijken we de gevallen waarbij het werkwoord niet finiet is
(1ac). En binnen die gevallen onderscheiden we opnieuw gevallen
waarbij het partikel gescheiden is van het werkwoord (1c) en gevallen
waarbij het partikel in het werkwoord geincorporeerd is (1a).`,
			`//node[starts-with(@sc,"part_")]`,
			SPOD_STD,
			"vpart",
			"Werkwoorden met een scheidbaar partikel||",
			"sc",
		},
		{
			"",
			`//node[starts-with(@sc,"part_") and
       @rel="hd" and
       ../node[@rel="svp" and starts-with(@frame,"particle")]]`,
			SPOD_STD,
			"vpartex",
			"Werkwoorden met een scheidbaar partikel|, partikel gescheiden van het werkwoord",
			"sc",
		},
		{
			"",
			`//node[starts-with(@sc,"part_") and
       (   not(@rel="hd") or
           @rel="hd" and not(../node[@rel="svp" and starts-with(@frame,"particle")])
       )]`,
			SPOD_STD,
			"vpartin",
			"Werkwoorden met een scheidbaar partikel|, partikel geïncorporeerd in het werkwoord",
			"sc",
		},
		{
			"",
			`//node[starts-with(@sc,"part_") and
       not(@wvorm="pv")]`,
			SPOD_STD,
			"vprtn",
			"Niet-finiete werkwoorden met een scheidbaar partikel||",
			"sc",
		},
		{
			"",
			`//node[starts-with(@sc,"part_") and
       not(@wvorm="pv") and
       @rel="hd" and
       ../node[@rel="svp" and starts-with(@frame,"particle")]]`,
			SPOD_STD,
			"vprtnex",
			"Niet-finiete werkwoorden met een scheidbaar partikel|, partikel gescheiden van het werkwoord",
			"sc",
		},
		{
			"",
			`//node[starts-with(@sc,"part_") and
       not(@wvorm="pv") and
       (   not(@rel="hd") or
           @rel="hd" and not(../node[@rel="svp" and starts-with(@frame,"particle")])
       )]`,
			SPOD_STD,
			"vprtnin",
			"Niet-finiete werkwoorden met een scheidbaar partikel|, partikel geïncorporeerd in het werkwoord",
			"sc",
		},
		{
			`Inbedding//Bij deze queries wordt gekeken naar de complexiteit van de zinnen in
termen van de inbedding van finitie bijzinnen. Een hoofdzin zonder
finiete bijzin geldt dan als "geen inbedding". Een hoofdzin met een
finiete bijzin geldt als "minstens 1 finiete zinsinbedding. Indien de
finiete bijzin zelf ook een finiete bijzin bevat is er sprake van minstens
2 finiete zinsinbeddingen. En zo verder.`,
			`//node[%PQ_finiete_inbedding0%]`,
			SPOD_STD,
			"inb0",
			"geen inbedding",
			"",
		},
		{
			"",
			`//node[%PQ_finiete_inbedding1%]`,
			SPOD_STD,
			"inb1",
			"minstens|| 1 finiete zinsinbedding",
			"",
		},
		{
			"",
			`//node[%PQ_finiete_inbedding2%]`,
			SPOD_STD,
			"inb2",
			"minstens| 2 finiete zinsinbeddingen",
			"",
		},
		{
			"",
			`//node[%PQ_finiete_inbedding3%]`,
			SPOD_STD,
			"inb3",
			"minstens| 3 finiete zinsinbeddingen",
			"",
		},
		{
			"",
			`//node[%PQ_finiete_inbedding4%]`,
			SPOD_STD,
			"inb4",
			"minstens| 4 finiete zinsinbeddingen",
			"",
		},
		{
			"",
			`//node[%PQ_finiete_inbedding5%]`,
			SPOD_STD,
			"inb5",
			"minstens| 5 finiete zinsinbeddingen",
			"",
		},
		{
			"",
			`//node[%PQ_finiete_inbedding6%]`,
			SPOD_STD,
			"inb6",
			"minstens| 6 finiete zinsinbeddingen",
			"",
		},
		{
			"",
			`//node[%PQ_finiete_inbedding7%]`,
			SPOD_STD,
			"inb7",
			"minstens| 7 finiete zinsinbeddingen",
			"",
		},
		{
			"",
			`//node[%PQ_finiete_inbedding8%]`,
			SPOD_STD,
			"inb8",
			"minstens| 8 finiete zinsinbeddingen",
			"",
		},
		{
			`Topicalizatie en Extractie//De eerste query, "np-topic is subject", bekijkt hoe vaak in een
hoofdzin die met een NP begint, deze NP als subject van de
persoonsvorm fungeert (1), terwijl de tweede query die gevallen telt
waarbij zo'n NP een andere rol bekleedt (2). De derde query, "topic is
niet lokaal", identificeert de gevallen waarbij de eerste constituent
van de hoofdzin geen rol speelt in die hoofdzin, maar een rol speelt
in een ingebedde constituent. Bij de vierde en vijfde query wordt
onderzocht hoe vaak in een wh-vraag of relatieve bijzin het
wh-element, dan wel de woordgroep die het relatieve voornaamwoord
bevat een rol speelt in de hoofdzin (lokaal, 4) dan wel in een ingebedde
bijzin (niet lokaal, 5).
<p>
<table>
<tr><td>(1)<td>De kinderen geloven nog in Sinterklaas</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(2)<td>Sinterklaas keurden ze geen blik waardig</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(3)<td>Wie denk je dat ik tegenkwam?</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(4)<td>Welke procedure moet worden gevolgd?</tr>
<tr><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(5)<td>Ze keek alleen naar wat ze dacht dat het Amerikaanse belang was</tr>
</table>
`,
			`%PQ_vorfeld_np_subject%`,
			SPOD_STD,
			"nptsub",
			"np-topic is subject",
			"",
		},
		{
			"",
			`%PQ_vorfeld_np_no_subject%`,
			SPOD_STD,
			"nptnsub",
			"np-topic is geen subject",
			"",
		},
		{
			"",
			`%PQ_vorfeld_non_local%`,
			SPOD_STD,
			"tnonloc",
			"topic is niet lokaal",
			"",
		},
		{
			"",
			`%PQ_local_extraction%`,
			SPOD_STD,
			"locext",
			"lokale extractie",
			"",
		},
		{
			"",
			`%PQ_non_local_extraction%`,
			SPOD_STD,
			"nlocext",
			"niet-lokale extractie",
			"",
		},
		{
			`Parser succes//Dit deel werkt niet voor handmatig geannoteerde corpora.
<p>
De query telt hoevaak de parser een volledige analyse van de
zin heeft kunnen uitvoeren (dus zonder woorden over te slaan of
fragmenten te combineren).`,
			`//parser[@cats="1" and @skips="0"]`,
			SPOD_STD,
			"ok",
			"volledige parse",
			"parser",
		},
		{
			`Parser succes: geparste delen//Dit deel werkt niet voor handmatig geannoteerde corpora.
<p>
Deze afdeling onderzoekt hoe goed de parser in staat was om de
zin te parsen. Indien de parser geen gehele parse kan maken dan
worden de grootste fragmenten waarvoor wel een parse was gevonden
gecombineerd. Ook kan de parser onder bepaalde voorwaarden
woorden overslaan. De queries tellen hoe vaak geen enkel fragment
werd gevonden, dan wel precies één fragment, twee fragmenten,
drie fragmenten, of vier of meer fragmenten. Indien de parse
uit precies één fragment bestaat kan het toch zijn dat er
geen sprake is van een volledige parse indien bijvoorbeeld
een woord in de zin werd overgeslagen om tot dat fragment te
komen.`,
			`//parser[@cats="0"]`,
			SPOD_STD,
			"cats0",
			"geen enkel deel is geparst",
			"parser",
		},
		{
			"",
			`//parser[@cats="1"]`,
			SPOD_STD,
			"cats1",
			"parse bestaat uit|| één deel",
			"parser",
		},
		{
			"",
			`//parser[@cats="2"]`,
			SPOD_STD,
			"cats2",
			"parse bestaat uit| twee losse delen",
			"parser",
		},
		{
			"",
			`//parser[@cats="3"]`,
			SPOD_STD,
			"cats3",
			"parse bestaat uit| drie losse delen",
			"parser",
		},
		{
			"",
			`//parser[number(@cats) > 3]`,
			SPOD_STD,
			"cats4",
			"parse bestaat uit| vier of meer losse delen",
			"parser",
		},
		{
			`Parser success: overgeslagen woorden//Dit deel werkt niet voor handmatig geannoteerde corpora.
<p>
Deze afdeling onderzoekt hoe goed de parser in staat was om de
zin te parsen. Indien de parser geen gehele parse kan maken dan
worden de grootste fragmenten waarvoor wel een parse was gevonden
gecombineerd. Ook kan de parser onder bepaalde voorwaarden
woorden overslaan. Deze queries tellen hoeveel woorden door de
parser werden overgeslagen voor het bereiken van de beste parse.`,
			`//parser[@skips="0"]`,
			SPOD_STD,
			"skips0",
			"geen enkel woord is overgeslagen",
			"parser",
		},
		{
			"",
			`//parser[@skips="1"]`,
			SPOD_STD,
			"skips1",
			"een van de woorden is overgeslagen",
			"parser",
		},
		{
			"",
			`//parser[@skips="2"]`,
			SPOD_STD,
			"skips2",
			"twee van de woorden zijn overgeslagen",
			"parser",
		},
		{
			"",
			`//parser[@skips="3"]`,
			SPOD_STD,
			"skips3",
			"drie van de woorden zijn overgeslagen",
			"parser",
		},
		{
			"",
			`//parser[number(@skips) > 3]`,
			SPOD_STD,
			"skips4",
			"vier of meer van de woorden zijn overgeslagen",
			"parser",
		},
		{
			`Onbekende woorden//Dit deel werkt niet voor handmatig geannoteerde corpora.
<p>
De parser kent een heel scale aan heuristieken om voor woorden die
niet in het woordenboek staan toch een woordsoort toe te kennen.
In deze afdeling wordt onderzocht hoe vaak zulke onbekende woorden
voorkwamen, en op welke manier onbekende woorden alsnog werden
behandeld. Eerst wordt het totaal aantal woorden gegeven,
gevolgd door het aantal woorden dat in het woordenboek stond.
De derde query geeft het aantal onbekende woorden terug.
De drie volgende queries geven vervolgens aan hoevaak een onbekend
woord als samenstelling werd gezien, of als naam (maar niet uit
de namenlijst) of op een nog andere manier werden behandeld.
`,
			`//node[@his]`,
			SPOD_STD,
			"his",
			"\"alle\" woorden (nodes met attribuut @his)",
			"his",
		},
		{
			"",
			`//node[@his="normal"]`,
			SPOD_STD,
			"normal",
			"woorden uit het woordenboek of de namenlijst",
			"his",
		},
		{
			"",
			`//node[@his and not(@his="normal")]`,
			SPOD_STD,
			"onbeken",
			"woorden niet direct uit het woordenboek",
			"his",
		},
		{
			"",
			`//node[@his="compound"]`,
			SPOD_STD,
			"compoun",
			"woorden herkend als samenstelling",
			"his",
		},
		{
			"",
			`//node[@his="name"]`,
			SPOD_STD,
			"name",
			"woorden herkend als naam (maar niet uit namenlijst)",
			"his",
		},
		{
			"",
			`
//node[@his
       and
       not(@his=("normal","compound","name"))
]`,
			SPOD_STD,
			"noun",
			"onbekende woorden die niet als samenstelling of naam werden herkend",
			"his",
		},
	}
	spodMu        sync.Mutex
	spodSemaphore chan bool
	spodWorking   = make(map[string]bool)
	spodRE        = regexp.MustCompile(`([0-9]+)[^0-9]+([0-9]+)`)
	spodREerr     = regexp.MustCompile(`\[err:`)
)

func spod_init() {
	spodSemaphore = make(chan bool, Cfg.Maxspodjob)
	for i, spod := range spods {
		spods[i].xpath = strings.TrimSpace(spod.xpath)
	}
}

func spod_main(q *Context) {

	writeHead(q, "Syntactic profiler of Dutch", 5)

	fmt.Fprintln(q.w, "Syntactic profiler of Dutch<p>")

	fmt.Fprint(q.w, `
<script type="text/javascript" src="jquery.js"></script>
<script type="text/javascript"><!--
  var xpaths = [
`)

	p := ""
	first := -1
	for i, spod := range spods {
		fmt.Fprintf(q.w, "%s['%s', '%s', '%s']", p, url.QueryEscape(spod.xpath), spod.method, spod.lbl)
		p = ",\n"
		if first < 0 && !strings.HasPrefix(spod.special, "hidden") {
			first = i
		}
	}
	fmt.Fprint(q.w, `
];
  var indexen = {
`)
	p = ""
	for i, spod := range spods {
		fmt.Fprintf(q.w, "%s'%s': %d", p, spod.lbl, i)
		p = ",\n"
	}
	fmt.Fprint(q.w, `
};
function vb(i) {
    var e = document.getElementById("spodform").elements
    window.open("xpath?db=" + e[0].value + "&xpath=" + xpaths[i][0] + "&mt=" + xpaths[i][1]);
}
function alles(n, m, o) {
    var e = document.getElementById("spodform").elements
    for (var i = n; i < m; i++) {
        e[i+1].checked = true;
    }
    if (o) {
        var ee = document.getElementsByClassName("spodblockinner");
        for (var i = 0; i < ee.length; i++) {
            if (ee[i].classList.contains('hide')) {
                ee[i].classList.remove('hide');
            }
        }
    }
}
function niets(n, m, o) {
    var e = document.getElementById("spodform").elements
    for (var i = n; i < m; i++) {
        e[i+1].checked = false;
    }
    if (o) {
        var ee = document.getElementsByClassName("spodblockinner");
        for (var i = 0; i < ee.length; i++) {
            if (! ee[i].classList.contains('hide')) {
                ee[i].classList.add('hide');
            }
        }
    }
}
function omkeren(n, m) {
    var e = document.getElementById("spodform").elements
    for (var i = n; i < m; i++) {
        e[i+1].checked = !e[i+1].checked;
    }
}
function hider(id) {
    var e = document.getElementById(id);
    e.classList.toggle("hide");
}
var queryvisible = false;
var queryid;
function info(id) {
    uninfo();
    var e = $(id);
    e.show();
    e.css("zIndex", 9999);
    queryvisible = true;
    queryid = id;
}
function uninfo() {
    if (queryvisible) {
        var e = $(queryid);
        e.hide();
        e.css("zIndex", 1);
        queryvisible = false;
    }
}
$(document).mouseup(
  function(e) {
    uninfo();
  });
function getChoices() {
  var res = [];
  for (var i = `, first, `; i < `, len(spods), `; i++) {
    var e = document.getElementById('i' + i);
    if (e.checked) {
      res.push(xpaths[i][2]);
    }
  }
  return res;
}
function setChoices(c) {
  var ee = document.getElementsByClassName("spodblockinner");
  for (var i = 0; i < ee.length; i++) {
    if (! ee[i].classList.contains('hide')) {
      ee[i].classList.add('hide');
    }
  }
  var e = document.getElementById("spodform").elements;
  for (var i = `, first, `; i < `, len(spods), `; i++) {
    e[i+1].checked = false;
  }
  for (var i = 0; i < c.length; i++) {
    var idx = indexen[c[i]];
    if (idx) {
      var el = e[idx + 1];
      el.checked = true;
      while (! el.classList.contains('spodblockinner')) {
        el = el.parentNode;
      }
      if (el.classList.contains('hide')) {
        el.classList.remove('hide');
      }
    }
  }
}
function optsave(n) {
  var c = getChoices();
  localStorage.setItem(
      "paqu-spod-"+n,
      JSON.stringify(c)
  );
  document.getElementById('op'+n).disabled = (c.length == 0);
}
function optload(n) {
  var storageContent = localStorage.getItem("paqu-spod-" + n);
  if (storageContent !== undefined) {
    setChoices(JSON.parse(storageContent) || []);
  }
}
//--></script>
<form id="spodform" action="spodform" method="get" accept-charset="utf-8" target="_blank">
corpus: <select name="db">
`)

	html_opts(q, q.opt_dbspod, getprefix(q), "corpus")
	fmt.Fprintln(q.w, "</select>")
	if q.auth {
		fmt.Fprintln(q.w, "<a href=\"corpuslijst\">meer/minder</a>")
	}

	fmt.Fprintf(q.w, `
<p>
<a href="spodlist" target="_blank">lijst van queries</a>
<p>
<a href="javascript:alles(%d, %d, true)">alles</a> &mdash;
<a href="javascript:niets(%d, %d, true)">niets</a> &mdash;
<a href="javascript:omkeren(%d, %d)">omkeren</a>
`, first, len(spods), first, len(spods), first, len(spods))

	inTable := false
	blocknum := 0
	for i, spod := range spods {
		if strings.HasPrefix(spod.special, "hidden") {
			fmt.Fprintf(q.w, `
<input type="hidden" name="i%d" value="t">
`,
				i)
			continue
		}
		spodtext := strings.Replace(spod.text, "|", "", -1)
		if spod.header != "" {
			if inTable {
				fmt.Fprintln(q.w, "</table></div></div>")
			} else {
				inTable = true
			}

			var j int
			for j = i + 1; j < len(spods); j++ {
				if spods[j].header != "" {
					break
				}
			}
			blocknum++
			a := strings.SplitN(spod.header, "//", 2)
			extra := ""
			if len(a) == 2 {
				// we kunnen hier geen <button> gebruiken omdat die de nummering in de war schopt
				extra = fmt.Sprintf(`
<a href="javascript:info('#h%d')" class="spodi">?</a>
<div id="h%d" class="submenu a9999">
  <div class="queryhelp">
    <h4>%s</h4>
%s
  </div>
</div>
`, blocknum, blocknum, html.EscapeString(a[0]), strings.TrimSpace(a[1]))
			}
			fmt.Fprintf(q.w, `
<div class="spodblock">
<a href="javascript:hider('spodblock%d')">%s</a>%s
<div class="spodblockinner hide" id="spodblock%d">
`, blocknum, spodEscape(a[0]), extra, blocknum)
			fmt.Fprintf(q.w, `
<a href="javascript:alles(%d, %d)">alles</a> &mdash;
<a href="javascript:niets(%d, %d)">niets</a> &mdash;
<a href="javascript:omkeren(%d, %d)">omkeren</a>
<p>
<table class="breed">`,
				i, j, i, j, i, j)
		}
		fmt.Fprintf(q.w, `<tr>
  <td><input type="checkbox" name="i%d" id="i%d" value="t">
  <td><a href="javascript:vb(%d)">vb</a>
  <td><label for="i%d">%s</label>
`,
			i,
			i,
			i,
			i,
			spodEscape(spodtext))
	}

	fmt.Fprintf(q.w, `</table></div></div>
<div class="spod"></div>
<div id="spodfoot"><div id="spodinfoot">
<div>
hergebruiken<br>
<button type="button" disabled onclick="optload(1)" id="op1">1</button><button type="button" disabled onclick="optload(2)" id="op2">2</button><button type="button" disabled onclick="optload(3)" id="op3">3</button><br>
<button type="button" disabled onclick="optload(4)" id="op4">4</button><button type="button" disabled onclick="optload(5)" id="op5">5</button><button type="button" disabled onclick="optload(6)" id="op6">6</button>
</div>
<div>
bewaren<br>
<button type="button" onclick="optsave(1)">1</button><button type="button" onclick="optsave(2)">2</button><button type="button" onclick="optsave(3)">3</button><br>
<button type="button" onclick="optsave(4)">4</button><button type="button" onclick="optsave(5)">5</button><button type="button" onclick="optsave(6)">6</button>
</div>
uitvoer: <select name="out">
<option value="html">HTML</option>
<option value="text">Teksttabel</option>
</select>
<p>
<input type="submit" value="verzenden">
</div></div>
</form>
<script type="text/javascript"><!--
for (var n = 1; n < 7; n++) {
  var storageContent = localStorage.getItem("paqu-spod-" + n);
  if (storageContent !== undefined) {
    var d = JSON.parse(storageContent);
    if (d && d.length > 0) {
      document.getElementById('op'+n).disabled = false;
    }
  }
}
//--></script>
`)

	fmt.Fprintln(q.w, "</body></html>")

}

func spod_form(q *Context) {

	seen := ""

	doHtml := first(q.r, "out") == "html"

	if doHtml {
		writeHead(q, "SPOD -- Resultaat", 0)
		defer func() {
			fmt.Fprintf(q.w, "</body>\n</html>\n")
		}()
	} else {
		contentType(q, "text/plain; charset=utf-8")
		nocache(q)
	}

	db := first(q.r, "db")
	if !q.spodprefixes[db] {
		fmt.Fprintln(q.w, "Ongeldig corpus:", db)
		return
	}

	available := map[string]bool{
		"":     true,
		"attr": true,
	}

	dbase, err := dbopen()
	if sysErr(err) {
		return
	}
	defer dbase.Close()
	rows, err := dbase.Query("SELECT `nline` from `" + Cfg.Prefix + "_info` WHERE `id` = \"" + db + "\"")
	if sysErr(err) {
		fmt.Fprintln(q.w, err)
		return
	}
	nlines := 0
	for rows.Next() {
		err = rows.Scan(&nlines)
		rows.Close()
		if sysErr(err) {
			fmt.Fprintln(q.w, err)
			return
		}
	}
	if nlines == 0 {
		err = fmt.Errorf("nlines not found")
		sysErr(err)
		fmt.Fprintln(q.w, err)
		return
	}

	mt := first(q.r, "mt")
	if mt != "std" && mt != "dx" {
		mt = "std"
	}

	var allDone bool
	if doHtml {
		fmt.Fprintf(q.w, "Corpus: %s\n<p>\n", spodEscape(q.desc[db]))
		allDone = spod_stats(q, db, true)
		if !allDone {
			fmt.Fprintln(q.w, "???<br>")
		}
	} else {
		fmt.Fprintln(q.w, "# corpus:", q.desc[db])
		fmt.Fprintln(q.w, "# waarde\t         \tlabel\tomschrijving")
		fmt.Fprintln(q.w, "## Stats")
		allDone = spod_stats(q, db, false)
		if !allDone {
			fmt.Fprintln(q.w, "???")
		}
	}

	inTable := false
	inAttr := false
	inData := false

	if doHtml {
		fmt.Fprintln(q.w, `
<style>
.max100 {
    max-width: 100%;
    overflow: auto;
}
/* The Modal (background) */
.modal {
    display: none; /* Hidden by default */
    position: fixed; /* Stay in place */
    z-index: 1; /* Sit on top */
    left: 0;
    top: 0;
    max-width: 100%;
    width: 100%; /* Full width */
    height: 100%; /* Full height */
    background-color: rgb(0,0,0); /* Fallback color */
    background-color: rgba(0,0,0,0.4); /* Black w/ opacity */
}
.modal-header {
    padding: 1em 2em;
    background-color: #315b7d;
    color: #b5cde1;
}
.modal-header h2 {
    font-size: x-large;
    font-weight: normal;
    text-align: center;
}
/* Modal Body */
.modal-body {
    width: 100%;
    overflow: auto;
    padding: 1em 2em;
    text-align: center;
}
/* Modal Footer */
.modal-footer {
    padding: 1em 2em;
    background-color: #315b7d;
    color: #b5cde1;
    text-align: center;
}
/* Modal Content */
.modal-content {
    position: relative;
    background-color: #fefefe;
    margin: auto;
    padding: 0;
    border: 1px solid #b5cde1;
    width: auto;
    max-width: 94%;
    box-shadow: 0 4px 8px 0 rgba(0,0,0,0.2),0 6px 20px 0 rgba(0,0,0,0.19);
    -webkit-animation-name: animatetop;
    -webkit-animation-duration: 0.4s;
    -moz-animation-name: animatetop;
    -moz-animation-duration: 0.4s;
    -o-animation-name: animatetop;
    -o-animation-duration: 0.4s;
    animation-name: animatetop;
    animation-duration: 0.4s
}

/* Add Animation */
@-webkit-keyframes animatetop {
    from {top: -300px; opacity: 0}
    to {top: 0; opacity: 1}
}
@-moz-keyframes animatetop {
    from {top: -300px; opacity: 0}
    to {top: 0; opacity: 1}
}
@-o-keyframes animatetop {
    from {top: -300px; opacity: 0}
    to {top: 0; opacity: 1}
}

@keyframes animatetop {
    from {top: -300px; opacity: 0}
    to {top: 0; opacity: 1}
}

/* The Close Button */
.close {
    color: #b5cde1;
    float: right;
    font-size: 28px;
    font-weight: bold;
}

.close:hover,
.close:focus {
    color: black;
    text-decoration: none;
    cursor: pointer;
}
.modal-body rect {
  fill: #4682b4;
}

.modal-body text {
  font-family: sans-serif;
  font-size: small;
}
</style>
<div id="myModal" class="modal">
  <!-- Modal content -->
  <div class="modal-content">
    <div class="modal-header">
      <span class="close">&times;</span>
      <h2 id="innerTitle"></h2>
    </div>
    <div class="modal-body"><svg width="960" height="500"></svg></div>
    <div class="modal-footer">`+
			spodEscape(q.desc[db])+`
    </div>
  </div>
</div>
<script src="https://d3js.org/d3.v4.min.js"></script>
<script>
function vb(i) {
    window.open("xpath?db=`+db+`&xpath=" + xpaths[i][0] + "&mt=" + xpaths[i][1]);
}
function vb2(i, a, v) {
    window.open("xpath?db=`+db+`&xpath=%2f%2fnode%5b%40" + a + "%3d%22" + v + "%22%5d&mt=" + xpaths[i][1]);
}

var modal = document.getElementById('myModal');
var span = document.getElementsByClassName("close")[0];
function wg(idx) {

    var d0 = data[idx];
    var prev = 0;
    var s = "woorden,frequentie";
    for (var i = 0; i < d0.length; i++) {
        while (prev < d0[i][0] - 1) {
            prev++;
            s += "\n" + prev + ",0";
        }
        s += "\n" + d0[i][0] + "," + d0[i][1];
        prev = d0[i][0];
    }
    var dat = d3.csvParse(s);
    for (var i = 0; i < dat.length; i++) {
        if (dat[i].frequentie) {
            dat[i].frequentie = + dat[i].frequentie;
        }
    }

    var ticks = [];
    var step = Math.round(dat.length / 20);
    if (step < 1) {
        step = 1;
    }
    for (var i = step; i <= dat.length; i += step) {
        ticks.push(""+i);
    }

    d3.selectAll("svg > *").remove()

	var svg = d3.select("svg"),
	    margin = {top: 20, right: 20, bottom: 40, left: 80},
	    width = +svg.attr("width") - margin.left - margin.right,
	    height = +svg.attr("height") - margin.top - margin.bottom;

	var x = d3.scaleBand().rangeRound([0, width]).padding(0.1),
	    y = d3.scaleLinear().rangeRound([height, 0]);

	var g = svg.append("g")
	    .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

    var ymax = d3.max(dat, function(d) { return d.frequentie; });

    x.domain(dat.map(function(d) { return d.woorden; }));
    y.domain([0, ymax]);

    g.append("g")
        .attr("class", "axis axis--x")
        .attr("transform", "translate(0," + height + ")")
        .call(d3.axisBottom(x).tickValues(ticks))
        .append("text")
        .attr("x", width / 2)
        .attr("y", 0)
        .attr("dy", "37px")
        .attr("text-anchor", "center")
        .attr("fill", "black")
        .text("aantal woorden");

    g.append("g")
        .attr("class", "axis axis--y")
        .call(ymax < 10 ? d3.axisLeft(y).ticks(ymax) : d3.axisLeft(y).tickFormat(d3.format("d")))
        .append("text")
        .attr("transform", "rotate(-90)")
        .attr("y", 6)
        .attr("dy", "0.71em")
        .attr("text-anchor", "end")
        .attr("fill", "black")
        .text("frequentie");

    g.selectAll(".bar")
      .data(dat)
      .enter().append("rect")
      .attr("class", "bar")
      .attr("x", function(d) { return x(d.woorden); })
      .attr("y", function(d) { return y(d.frequentie); })
      .attr("width", x.bandwidth())
      .attr("height", function(d) { return height - y(d.frequentie); });

    d3.selectAll("#innerTitle").text(titles[idx]);
    modal.style.display = "block";
}
span.onclick = function() {
    modal.style.display = "none";
}
window.onclick = function(event) {
    if (event.target == modal) {
        modal.style.display = "none";
    }
}
</script>
`)
	}
	worddata := make([]string, 0)
	wordtitles := make([]string, 0)
	header := ""
	spod_in_use := make(map[string]bool)

	for idx, spod := range spods {
		spod_in_use[spod_fingerprint(idx)] = true

		if strings.HasPrefix(spod.special, "hidden") {
			lines, _, _, done, err := spod_get(q, db, idx)
			if err == nil && done {
				available[strings.Split(spod.lbl, "_")[1]] = lines > 0
			}
			continue
		}

		if spod.header != "" {
			a := strings.SplitN(spod.header, "//", 2)
			header = strings.TrimSpace(a[0])
		}
		spodtext := spod.text
		if strings.Contains(spodtext, "||") || !strings.Contains(spodtext, "|") {
			spodtext = strings.Replace(spodtext, "||", "|", 1)
			seen = ""
		}
		if first(q.r, fmt.Sprintf("i%d", idx)) == "t" {
			if !inTable {
				inTable = true
				if doHtml {
					fmt.Fprintln(q.w, `<div class="max100"><table class="spod">`)
				}
			}
			if spod.special == "attr" {
				if !inAttr {
					inAttr = true
					if doHtml {
						fmt.Fprintln(q.w, `<tr><th colspan="2" class="r b">attributen<th colspan="2" class="r"><th colspan="2"></tr>
`)
					} else {
						fmt.Fprintln(q.w, "# attributen att/totaal\t\tlabel\tomschrijving\t")
					}
				}
				if doHtml {
					fmt.Fprintf(q.w, "<tr><th colspan=\"2\" class=\"r\"><th colspan=\"2\" class=\"r\"><th colspan=\"2\" class=\"left\">%s</tr>\n", spodEscape(spod.text))
				}
			} else {
				if inAttr {
					inAttr = false
				}
				if !inData {
					inData = true
					if doHtml {
						fmt.Fprintln(q.w, `<tr><th colspan="2" class="r b">zinnen<th class="r b">items<th class="r b">woorden<th><th></tr>
`)
					} else {
						fmt.Fprintln(q.w, "# zinnen zinnen/totaal\titems\tlabel\tomschrijving\twoordtelling")
					}
				}
			}
			if doHtml {
				a := strings.SplitN(spodtext, "|", 2)
				if len(a) == 2 {
					if a[0] == seen {
						spodtext = "— " + a[1]
					}
					seen = a[0]
				}
			} else {
				spodtext = strings.Replace(spodtext, "|", "", 1)
			}
			if header != "" && spod.special != "attr" {
				if doHtml {
					fmt.Fprintf(q.w, "<tr><th colspan=\"2\" class=\"r\"><th class=\"r\"><th class=\"r\"><th colspan=\"2\" class=\"left\">%s</tr>\n", spodEscape(header))
				} else {
					fmt.Fprintln(q.w, "##", header)
				}
				header = ""
			}
			avail, ok := available[spod.special]
			if specials := strings.Fields(spod.special); !ok && len(specials) > 1 {
				allok := true
				someavail := false
				for _, spec := range specials {
					neg := false
					if spec[0] == '-' {
						spec = spec[1:]
						neg = true
					}
					av, ok := available[spec]
					if !ok {
						allok = false
						continue
					}
					if neg {
						av = !av
					}
					if av {
						someavail = true
						break
					}
				}
				if someavail {
					available[spod.special] = true
					avail = true
					ok = true
				} else if allok {
					available[spod.special] = false
					avail = false
					ok = true
				}
			}
			if ok && !avail {
				if doHtml {
					fmt.Fprintf(q.w, "<tr><td colspan=\"4\" class=\"center r\"><em>niet voor dit corpus</em><td><td>%s</tr>\n", spodEscape(spodtext))
				} else {
					fmt.Fprintf(q.w, "### niet voor dit corpus\t%s\n", spodtext)
				}
				continue
			}
			var lines, items int
			var wcount string
			var done bool
			var err error
			if avail {
				lines, items, wcount, done, err = spod_get(q, db, idx)
			}
			if err != nil {
				if doHtml {
					fmt.Fprintf(q.w, "<tr><td colspan=\"3\"><b>%s</b>", spodEscape(err.Error()))
				} else {
					fmt.Fprint(q.w, "#", err)
				}
				allDone = false
			} else if done && available[spod.special] {
				if doHtml {
					if spod.special == "attr" {
						for _, line := range strings.Split(wcount, "\n") {
							aa := strings.Fields(line)
							if len(aa) == 4 {
								fmt.Fprintf(q.w, `<tr><td class="right">%s<td class="right r">%s<td colspan="2" class="r">
<td><a href="javascript:vb2(%d, '%s','%s')">vb</a><td>%s</tr>
`,
									aa[0], aa[2], idx, spod.lbl, url.QueryEscape(aa[3]), html.EscapeString(aa[3]))
							}
						}
					} else if spod.special == "parser" {
						fmt.Fprintf(q.w, "<tr><td class=\"right\">%d<td class=\"right r\">%.2f%%",
							lines, float64(lines)/float64(nlines)*100.0)
					} else if spod.special == "his" {
						fmt.Fprintf(q.w, "<tr><td><td class=\"r\"><td class=\"right r\">%d",
							items)
					} else {
						fmt.Fprintf(q.w, "<tr><td class=\"right\">%d<td class=\"right r\">%.2f%%<td class=\"right r\">%d",
							lines, float64(lines)/float64(nlines)*100.0, items)
					}
				} else {
					if spod.special == "attr" {
						for _, line := range strings.Split(wcount, "\n") {
							aa := strings.Fields(line)
							if len(aa) == 4 {
								fmt.Fprintf(q.w, "%s\t%-15s\t\t%s.%s\tattribuut %s: %s\t\n",
									aa[0], aa[1], spod.lbl, aa[3], spod.lbl, aa[3])
							}
						}
					} else {
						v := fmt.Sprintf("%.3g", float64(lines)/float64(nlines))
						fmt.Fprintf(q.w, "%d\t%-15s\t%d", lines, v, items)
					}
				}
			} else {
				wcount = "???"
				if doHtml {
					if spod.special == "attr" {
						fmt.Fprintf(q.w, "<tr><td class=\"right\">???<td class=\"right r\">???<td class=\"right r\" colspan=\"2\"><td><td>???</tr>")
					} else if spod.special == "parser" {
						fmt.Fprintf(q.w, "<tr><td class=\"right\">???<td class=\"right r\">???")
					} else if spod.special == "his" {
						fmt.Fprintf(q.w, "<tr><td><td class=\"r\"><td class=\"right r\">???")
					} else {
						fmt.Fprintf(q.w, "<tr><td class=\"right\">???<td class=\"right r\">???<td class=\"right r\">???")
					}
				} else {
					if spod.special == "attr" {
						fmt.Fprintf(q.w, "???\t???       \t\t%s.???\tattribuut %s: ???\n", spod.lbl, spod.lbl)
						// TODO
					} else {
						fmt.Fprint(q.w, "???\t???       \t???")
					}
				}
				allDone = false
			}
			if doHtml {
				if spod.special == "parser" {
					fmt.Fprintf(
						q.w,
						"<td class=\"r\"><td class=\"r\"><td><a href=\"javascript:vb(%d)\">vb</a><td>%s\n",
						idx, spodEscape(spodtext))
				} else if spod.special == "attr" {
					// niks
				} else {
					counts := strings.Split(wcount, ",")
					sum := 0
					n := 0
					for _, count := range counts {
						a := strings.Split(count, ":")
						if len(a) == 2 {
							i, _ := strconv.Atoi(a[0])
							j, _ := strconv.Atoi(a[1])
							sum += i * j
							n += j
						}
					}
					if sum == 0 {
						t := "&mdash;"
						if wcount == "???" {
							t = "???"
						}
						fmt.Fprintf(
							q.w,
							"<td class=\"right r\">%s<td><a href=\"javascript:vb(%d)\">vb</a><td>%s\n",
							t, idx, spodEscape(spodtext))
					} else {
						fmt.Fprintf(
							q.w,
							"<td class=\"right r\"><a href=\"javascript:wg(%d)\">%.1f</a><td>\n",
							len(worddata), float64(sum)/float64(n))

						fmt.Fprintf(
							q.w,
							"<a href=\"javascript:vb(%d)\">vb</a><td>%s\n",
							idx, spodEscape(spodtext))

						worddata = append(worddata, "[["+
							strings.Replace(
								strings.Replace(wcount, ",", "],[", -1),
								":", ",", -1)+"]]")
						wordtitles = append(wordtitles, jsstringsEscape(spod.text))
					}
				}
			} else {
				if spod.special == "attr" {
					// niks
				} else {
					fmt.Fprintf(q.w, "\t%s\t%s\t%s\n", spod.lbl, spodtext, wcount)
				}
			}
		}
	}
	if doHtml {
		if inTable {
			fmt.Fprintln(q.w, "</table></div>")
		}
		fmt.Fprintln(q.w, "<script>var data = [")
		fmt.Fprintln(q.w, strings.Join(worddata, ",\n"))
		fmt.Fprintln(q.w, "];\nvar titles = [")
		fmt.Fprintln(q.w, strings.Join(wordtitles, ",\n"))
		fmt.Fprintln(q.w, "];\nvar xpaths = [")
		p := ""
		for _, spod := range spods {
			fmt.Fprintf(q.w, "%s['%s', '%s']", p, url.QueryEscape(spod.xpath), spod.method)
			p = ",\n"
		}
		fmt.Fprintln(q.w, "];</script>")
	}
	if !allDone {
		if doHtml {
			fmt.Fprintln(q.w, "<div class=\"footspace\"></div><div id=\"spodinfo\" class=\"info\"><button type=\"button\" onClick=\"location.reload(true)\">Herladen</button> &mdash; Herlaad de pagina om de ontbrekende resultaten te krijgen</div>")
		} else {
			fmt.Fprintln(q.w, "#\n#  -->  HERLAAD DE PAGINA OM DE ONTBREKENDE RESULTATEN TE KRIJGEN  <--\n#")
		}
	}

	// oude varianten verwijderen
	dirpath := filepath.Join(paqudatadir, "data", db, "spod")
	files, err := ioutil.ReadDir(dirpath)
	if sysErr(err) {
		return
	}
	spod_in_use["stats"] = true
	for _, file := range files {
		filename := file.Name()
		if !spod_in_use[filename] {
			err := os.Remove(filepath.Join(dirpath, filename))
			sysErr(err)
		}
	}
}

func spod_stats(q *Context, db string, doHtml bool) bool {
	spodMu.Lock()
	defer spodMu.Unlock()

	key := db
	if spodWorking[key] {
		return false
	}

	dirpath := filepath.Join(paqudatadir, "data", db, "spod")
	os.MkdirAll(dirpath, 0700)
	filename := filepath.Join(dirpath, "stats")

	data, err := ioutil.ReadFile(filename)
	if err == nil {
		if doHtml {
			fmt.Fprintln(q.w, "<table class=\"compact\">")
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				a := strings.Split(line, "\t")
				if len(a) == 5 {
					fmt.Fprintf(q.w, "<tr><td class=\"right\">%s<td>%s\n", a[0], spodEscape(a[4]))
				}
			}
			fmt.Fprintln(q.w, "</table><p>")
		} else {
			q.w.Write(data)
		}
		return true
	}

	spodWorking[key] = true
	go func() {
		go spod_stats_work(q, db, filename)
		spodMu.Lock()
		delete(spodWorking, key)
		spodMu.Unlock()
	}()

	return false
}

func spod_get(q *Context, db string, item int) (lines int, items int, wcount string, done bool, err error) {
	spodMu.Lock()
	defer spodMu.Unlock()

	key := fmt.Sprintf("%s\t%d", db, item)
	if spodWorking[key] {
		return 0, 0, "", false, nil
	}

	dirpath := filepath.Join(paqudatadir, "data", db, "spod")
	os.MkdirAll(dirpath, 0700)

	fingerprint := spod_fingerprint(item)
	filename := filepath.Join(dirpath, fingerprint)
	data, err := ioutil.ReadFile(filename)
	if err == nil && len(data) > 0 {
		if spods[item].special == "attr" {
			return 0, 0, string(data), true, nil
		} else {
			a := strings.Fields(string(data))
			if len(a) == 3 {
				a = append(a, "1:0")
			}
			if len(a) == 4 {
				if a[0] != spods[item].lbl {
					os.Remove(filename)
					return 0, 0, "", false, fmt.Errorf("ERROR: invalid label %q", a[0])
				}
				lines, err := strconv.Atoi(a[1])
				if err != nil {
					return 0, 0, "", false, err
				}
				items, err := strconv.Atoi(a[2])
				if err != nil {
					return 0, 0, "", false, err
				}
				return lines, items, a[3], true, nil
			} else {
				os.Remove(filename)
				return 0, 0, "", false, fmt.Errorf("ERROR: invalid data %q", string(data))
			}
		}
	}

	spodWorking[key] = true
	go func() {
		spod_work(q, key, filename, db, item)
		spodMu.Lock()
		delete(spodWorking, key)
		spodMu.Unlock()
	}()

	return 0, 0, "", false, nil
}

func spod_work(q *Context, key string, filename string, db string, item int) {
	spodSemaphore <- true
	defer func() {
		<-spodSemaphore
	}()

	var u string
	onlyone := spods[item].special == "hidden1"
	attr := spods[item].special == "attr"
	if onlyone {
		u = fmt.Sprintf("http://localhost/?db=%s&xpath=%s&mt=%s&xn=1",
			db,
			url.QueryEscape(spods[item].xpath),
			spods[item].method)
	} else if attr {
		u = fmt.Sprintf("http://localhost/?db=%s&xpath=%s&mt=%s&attr1=%s&d=1",
			db,
			url.QueryEscape(spods[item].xpath),
			spods[item].method,
			spods[item].lbl)
	} else {
		u = fmt.Sprintf("http://localhost/?db=%s&xpath=%s&mt=%s&attr1=word_is_&d=1",
			db,
			url.QueryEscape(spods[item].xpath),
			spods[item].method)
	}

	r, err := http.NewRequest("GET", u, nil)
	if sysErr(err) {
		return
	}
	if sysErr(r.ParseForm()) {
		return
	}
	w := spod_writer{header: make(map[string][]string)}
	dbase, err := dbopen()
	if sysErr(err) {
		return
	}
	defer dbase.Close()
	myQ := Context{
		r:            r,
		w:            &w,
		user:         q.user,
		auth:         q.auth,
		sec:          q.sec,
		quotum:       q.quotum,
		db:           dbase,
		opt_db:       q.opt_db,
		opt_dbmeta:   q.opt_dbmeta,
		ignore:       q.ignore,
		prefixes:     q.prefixes,
		myprefixes:   q.myprefixes,
		spodprefixes: q.spodprefixes,
		protected:    q.protected,
		hasmeta:      q.hasmeta,
		desc:         q.desc,
		lines:        q.lines,
		shared:       q.shared,
		params:       q.params,
		form:         nil,
	}

	fp, err := os.Create(filename)
	if sysErr(err) {
		return
	}

	if onlyone {
		xpath(&myQ)
		if strings.Contains(w.buffer.String(), "<!--NOMATCH-->") {
			fmt.Fprintf(fp, "%s\t0\t0\t\n", spods[item].lbl)
		} else {
			fmt.Fprintf(fp, "%s\t1\t1\t1:1\n", spods[item].lbl)
		}
		fp.Close()
		return
	}

	if attr {
		xpathstats(&myQ)
		scanner := bufio.NewScanner(&w.buffer)
		scanner.Scan()
		scanner.Scan()
		at := make([]string, 0)
		sum := 0
		for scanner.Scan() {
			line := strings.Fields(scanner.Text())
			at = append(at, line[2]+"\t"+line[0])
			a, _ := strconv.Atoi(line[0])
			sum += a
		}
		sort.Strings(at)
		for _, line := range at {
			aa := strings.Fields(line)
			n, _ := strconv.Atoi(aa[1])
			fmt.Fprintf(fp, "%d\t%.3f\t%.2f%%\t%s\n", n, float64(n)/float64(sum), float64(n)/float64(sum)*100.0, aa[0])
		}
		fp.Close()
		return
	}

	xpathstats(&myQ)

	scanner := bufio.NewScanner(&w.buffer)

	scanner.Scan()
	s := scanner.Text()

	if spodREerr.MatchString(s) {
		fmt.Fprintf(fp, "ERROR: %s\n", s)
		chLog <- "ERROR spodREerr.MatchString: " + s
	} else {
		match := spodRE.FindStringSubmatch(s)
		if len(match) == 3 {
			fmt.Fprintf(fp, "%s\t%s\t%s\t", spods[item].lbl, match[1], match[2])
			// skip regel
			scanner.Scan()
			counts := make(map[int]int)
			var n int
			var err error
			for scanner.Scan() {
				s := scanner.Text()
				a := strings.Fields(s)
				n, err = strconv.Atoi(a[0])
				if err != nil {
					fmt.Fprintf(fp, "ERROR: no number found in %q\n", s)
					chLog <- "ERROR Atoi: " + err.Error()
					break
				}
				m := 0
				for _, i := range a[2:] {
					if i == "+" {
						m++
					}
				}
				counts[m] = counts[m] + n
			}
			if err == nil {
				keys := make([]int, 0)
				for key := range counts {
					keys = append(keys, key)
				}
				sort.Ints(keys)
				p := ""
				for _, key := range keys {
					fmt.Fprintf(fp, "%s%d:%d", p, key, counts[key])
					p = ","
				}
			}
			fmt.Fprintln(fp)
		} else {
			fmt.Fprintf(fp, "ERROR: no match found in %q\n", s)
			chLog <- "ERROR spodREerr.FindStringSubmatch: " + s
		}
	}
	fp.Close()
}

func (s *spod_writer) Header() http.Header {
	return http.Header(s.header)
}

func (s *spod_writer) Write(b []byte) (int, error) {
	return s.buffer.Write(b)
}

func (s *spod_writer) WriteHeader(i int) {
}

func spod_stats_work(q *Context, dbname string, outname string) {
	spodSemaphore <- true
	defer func() {
		<-spodSemaphore
	}()

	xx := func(err error) bool {
		if err == nil {
			return false
		}
		fp, _ := os.Create(outname)
		fmt.Fprintln(fp, "ERROR:", err)
		fp.Close()
		return true
	}

	lineCount := 0
	wordCount := 0
	runeCount := 0
	tokens := make(map[string]bool)

	db, err := dbopen()
	if xx(err) {
		return
	}
	defer func() {
		if db != nil {
			db.Close()
		}
	}()

	archnames := make([]string, 0)
	rows, err := db.Query(fmt.Sprintf("SELECT `arch` FROM `%s_c_%s_arch`", Cfg.Prefix, dbname))
	if xx(err) {
		return
	}
	for rows.Next() {
		var s string
		if xx(rows.Scan(&s)) {
			rows.Close()
			return
		}
		archnames = append(archnames, s)
	}
	if xx(rows.Err()) {
		return
	}

	filenames := make([]string, 0)
	if len(archnames) == 0 {
		rows, err := db.Query(fmt.Sprintf("SELECT `file` FROM `%s_c_%s_file`", Cfg.Prefix, dbname))
		if xx(err) {
			return
		}
		for rows.Next() {
			var s string
			if xx(rows.Scan(&s)) {
				rows.Close()
				return
			}
			if strings.HasPrefix(s, "$$/") {
				s = filepath.Join(paqudatadir, "data", dbname, "xml", s[3:])
			}
			filenames = append(filenames, s)
		}
		if xx(rows.Err()) {
			return
		}
	}

	db.Close()
	db = nil

	var processNode func(node *Node)
	processNode = func(node *Node) {
		if node.Word != "" && node.Pt != "let" {
			wordCount++
			runeCount += utf8.RuneCountInString(node.Word)
			tokens[node.Word] = true
		}
		for _, n := range node.NodeList {
			processNode(n)
		}
	}

	processFile := func(data []byte) error {
		lineCount++
		var alpino Alpino_ds
		err := xml.Unmarshal(data, &alpino)
		if err != nil {
			return err
		}
		processNode(alpino.Node0)
		return nil
	}

	if len(archnames) > 0 {
		for _, archname := range archnames {
			db, err := dbxml.OpenRead(archname)
			if xx(err) {
				return
			}
			docs, err := db.All()
			if xx(err) {
				db.Close()
				return
			}
			for docs.Next() {
				if xx(processFile([]byte(docs.Content()))) {
					db.Close()
					return
				}
			}
			db.Close()
		}
	} else {
		for _, filename := range filenames {
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				fp, err := os.Open(filename + ".gz")
				if xx(err) {
					return
				}
				r, err := gzip.NewReader(fp)
				if xx(err) {
					fp.Close()
					return
				}
				data, err = ioutil.ReadAll(r)
				fp.Close()
				if xx(err) {
					return
				}
			}
			if xx(processFile(data)) {
				return
			}
		}
	}

	fp, _ := os.Create(outname)
	fmt.Fprintf(fp,
		"%8d\t\t\tzinnen\tzinnen\n"+
			"%8d\t\t\twoorden\twoorden\n"+
			"%8.4f\t\t\ttt\ttypes per token\n"+
			"%8.4f\t\t\twz\twoorden per zin\n"+
			"%8.4f\t\t\tlw\tletters per woord\n",
		lineCount,
		wordCount,
		float64(len(tokens))/float64(wordCount),
		float64(wordCount)/float64(lineCount),
		float64(runeCount)/float64(wordCount))
	fp.Close()
}

func spod_list(q *Context) {

	contentType(q, "text/plain; charset=utf-8")
	nocache(q)

	for _, spod := range spods {
		if strings.HasPrefix(spod.special, "hidden") {
			continue
		}
		header := spod.header
		if i := strings.Index(header, "//"); i > 0 {
			header = header[:i]
		}
		spodtext := strings.Replace(spod.text, "|", "", -1)
		if header != "" {
			fmt.Fprint(q.w, "\n\n", header, "\n", strings.Repeat("=", len(header)), "\n\n")
		}
		fmt.Fprint(q.w,
			"\n",
			spod.lbl, ": ", spodtext, "\n",
			strings.Repeat("-", len(spod.lbl)+2+len(spodtext)), "\n\n",
			spod.xpath, "\n\n")
	}
}

func spod_fingerprint(item int) string {
	rules := getMacrosRules(&Context{})
	query := macroKY.ReplaceAllStringFunc(spods[item].xpath, func(s string) string {
		return rules[s[1:len(s)-1]]
	})
	query = strings.Join(strings.Fields(query), " ")
	return fmt.Sprintf("%x", md5.Sum([]byte(query+spods[item].method)))
}

func jsstringsEscape(s string) string {
	return "\"" +
		strings.Replace(
			strings.Replace(
				strings.Replace(s, "\\", "\\\\", -1),
				"\"", "\\\"", -1),
			"|", "", -1) +
		"\""
}

func spodEscape(s string) string {
	s = html.EscapeString(s)
	if strings.Contains(s, "|") {
		a := strings.SplitN(s, "|", 2)
		s = "<u>" + a[0] + "</u>" + a[1]
	}
	return s
}
