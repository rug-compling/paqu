package spod

const (
	SPOD_STD = "std"
	SPOD_DX  = "dx"
)

type Spod struct {
	Header  string
	Xpath   string
	Method  string
	Lbl     string
	Text    string
	Special string
}

var Spods = []Spod{
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
		`Attributen//Hier worden de woordsoorten geteld. De Lassy part-of-speech labels
worden geteld via het attribuut <b>postag</b>. Er is daarnaast een
attribuut <b>pt</b> dat de vereenvoudigde weergave geeft van
dezelfde informatie: het deel van de postag tot aan het
openingshaakje. De documentatie van postag is te vinden in: Frank van
Eynde,
<a href="http://www.let.rug.nl/vannoord/Lassy/POS_manual.pdf">Part of
speech tagging en lemmatisering van het D-Coi corpus</a>, 2005.
<p>
Vanwege historische redenen wordt hier ook het attribuut <b>pos</b>
opgenomen, maar dit attribuut wordt niet langer onderhouden, is ook
niet gedocumenteerd en wordt ook niet gecorrigeerd in de handmatig
geannoteerde corpora.`,
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
		`//node[%PQ_imperatieven%]`,
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
       count(node[@rel="cnj"])>1
       and
       node[1][@rel="crd"]
       and
       node[3][@rel="crd"]
]`,
		SPOD_STD,
		"crd22",
		"reeksvormers (nevenschikkingen van de vorm crd,cnj,crd,cnj...)",
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
		`Adjectiefgroepen//Adjectiefgroepen worden onderscheiden naar grammaticale functie en naar
interne structuur. Adjectiefgroepen kunnen optreden als bepaling bij nomina (1), (2),
traditioneel &quot;bijvoegelijke bepaling&quot; genoemd, en bij werkwoorden of andere bijvoeglijke
naamwoorden (3), (4), traditioneel bijwoordelijke bepaling genoemd. Adjectiefgroepen kunnen
ook optreden met als relatie &quot;predm&quot;, en dat zijn gevallen die traditionaal bekend staan als
bepaling van gesteldheid (5) (6). Ten slotte kan een adjectiefgroep optreden als predicatief
complement (naamwoordelijk deel van het gezegde), zoals in (7), (8).
<p>
<table>
<tr><td>(1)<td>De <em>aardige</em> buurman</tr>
<tr><td>(2)<td>De <em>nogal vervelende</em> vergadering</tr>
<tr><td><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(3)<td>Hij trapte de bal <em>hard</em> over het doel</tr>
<tr><td>(4)<td>De dienst is <em>volledig</em> kosteloos</tr>
<tr><td><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(5)<td>Hij kwam <em>ziek</em> binnen</tr>
<tr><td>(6)<td>Hij liet het knipsel <em>dol van trots</em> aan iedereen zien</tr>
<tr><td><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(7)<td>Ik vind dat <em>idioot</em></tr>
<tr><td>(8)<td>De dienst is <em>volledig kosteloos</em></tr>
<tr><td><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(9)<td>Het boek is <em>beschikbaar via de uitgever</em></tr>
<tr><td><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(10)<td>Het boek werd <em>tien jaar later</em> gepubliceerd</tr>
<tr><td>(11)<td>De <em>mij onbekende</em> schrijver bleek een natuurtalent</tr>
<tr><td><td>&nbsp;<td>&nbsp;</tr>
<tr><td>(12)<td>De schat is <em>afkomstig uit Japan</em></tr>
<tr><td>(13)<td>Ron is <em>bang dat dit iets te maken heeft met Sirius Zwarts</em></tr>
</table>
<p>
Adjectiefgroepen kunnen ook onderverdeel worden aan de hand van hun interne structuur.
We onderscheiden adjectieven die hun eentje optreden (1), (3), (4), (5), (7); adjectieven die voorkomen
met een bepaling ter linkerzijde (2), (8); adjectieven die voorkomen met een bepaling rechts van
het adjectief (9); adjectieven die voorkomen met een complement, links van het adjectief (10), (11) en
ten slotte de gevallen waarbij het adjectief optreedt met een complement ter rechterzijde (12), (13).`,
		`//node[(@cat="ap" or @pt="adj") and ((@rel="mod" and ../@cat="np") or (@rel="cnj" and ../@rel="mod" and ../../@cat="np"))]`,
		SPOD_STD,
		"ajbepn",
		"grammaticale functie||, bijvoegelijke bepaling (bepaling bij zelfstandig naamwoord)",
		"",
	},
	{
		"",
		`//node[
          (@cat="ap" or @pt="adj")
          and (
                  (@rel="mod" and not(../@cat="np"))
                  or
                  (@rel="cnj" and (../@rel="mod" or ../@index=//node[@rel="mod"]/@index) and not(../../@cat="np"))
                  or
                  //node[
                            (@rel="mod" and not(../@cat="np"))
                            or
                            (@rel="cnj" and (../@rel="mod" or ../@index=//node[@rel="mod"]/@index) and not(../../@cat="np"))
                  ]/@index=@index
          )
]`,
		SPOD_STD,
		"ajbepva",
		"grammaticale functie|, bijwoordelijke bepaling (bepaling bij werkwoord, adjectief, ...)",
		"",
	},
	{
		"",
		`//node[(@cat="ap" or @pt="adj") and ((@rel="predm") or (@rel="cnj" and ../@rel="predm"))]`,
		SPOD_STD,
		"ajpredm",
		"grammaticale functie|, bepaling van gesteldheid etc.",
		"",
	},
	{
		"",
		`//node[
          (@cat="ap" or @pt="adj")
          and
          (
              @rel="predc"
              or
              (@rel="cnj" and (../@rel="predc" or ../@index=//node[@rel="predc"]/@index))
              or
              //node[
                        @rel="predc"
                        or
                        (@rel="cnj" and (../@rel="predc" or ../@index=//node[@rel="predc"]/@index))
              ]/@index=@index
          )
]`,
		SPOD_STD,
		"ajpredc",
		"grammaticale functie|, predicatief complement",
		"",
	},
	{
		"",
		`//node[@pt="adj" and not(@rel="hd")]`,
		SPOD_STD,
		"ajaj",
		"interne structuur||, ADJ",
		"",
	},
	{
		"",
		`//node[@cat="ap" and node[@rel="mod"]/number(@begin) < node[@rel="hd"]/number(@begin)]`,
		SPOD_STD,
		"ajmodaj",
		"interne structuur|, MOD ADJ",
		"",
	},
	{
		"",
		`//node[
          @cat="ap"
          and
          node[
                  (
                      @rel="predc" or @rel="obj1" or @rel="obj2" or @rel="pc" or @rel="me" or @rel="vc"
                  )
                  and
                  not(%PQ_vorfeld%)
          ]/number(@begin)
          <
          node[@rel="hd"]/number(@begin)
]
`,
		SPOD_STD,
		"ajcaj",
		"interne structuur|, COMPL ADJ",
		"",
	},
	{
		"",
		`//node[@cat="ap" and node[@rel="mod"]/number(@begin) > node[@rel="hd"]/number(@begin)]`,
		SPOD_STD,
		"ajajmod",
		"interne structuur|, ADJ MOD",
		"",
	},
	{
		"",
		`//node[
          @cat="ap"
          and
          node[
                  (
                      @rel="predc" or @rel="obj1" or @rel="obj2" or @rel="pc" or @rel="me" or @rel="vc"
                  )
                  and
                  not(%PQ_nachfeld%)
          ]/number(@begin)
          >
          node[@rel="hd"]/number(@begin)
]`,
		SPOD_STD,
		"ajajc",
		"interne structuur|, ADJ COMPL",
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
		"pos_verb",
	},
	{
		"",
		`//node[%PQ_passive%]`,
		SPOD_STD,
		"passive",
		"passief",
		"pos_verb",
	},
	{
		"",
		`//node[%PQ_impersonal_passive%]`,
		SPOD_STD,
		"nppas",
		"niet-persoonlijke passief",
		"pos_verb",
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
		`Werkwoordstijden//We onderscheiden de vier klassieke werkwoordstijden: onvoltooid
tegenwoordige tijd (1), onvoltooid verleden tijd (2), voltooid
tegenwoordige tijd (3) en voltooid verleden tijd (4).
<p>
<table>
<tr><td>(1)<td>De bakker bakt het brood</tr>
<tr><td>(2)<td>De bakker bakte het brood</tr>
<tr><td>(3)<td>De bakker heeft het brood gebakken</td>
<tr><td>(4)<td>De bakker had het brood gebakken</td>
</table>`,
		`//node[%PQ_ott%]`,
		SPOD_STD,
		"ott",
		"Onvoltooid tegenwoordige tijd",
		"",
	},
	{
		"",
		`//node[%PQ_ovt%]`,
		SPOD_STD,
		"ovt",
		"Onvoltooid verleden tijd",
		"",
	},
	{
		"",
		`//node[%PQ_vtt%]`,
		SPOD_STD,
		"vtt",
		"Voltooid tegenwoordige tijd",
		"",
	},
	{
		"",
		`//node[%PQ_vvt%]`,
		SPOD_STD,
		"vvt",
		"Voltooid verleden tijd",
		"",
	},
	{
		`Volgorde van werkwoordstijden//Deze queries onderscheiden finiete zinnen die zelf een finiet
zinscomplement bevatten. Hierbij kijken we naar de combinaties van de
(onvoltooide) werkwoordstijd van de persoonsvorm in de dominerende zin en
de (onvoltooide) werkwoordstijd van de persoonsvorm in de zin die als
complement optreedt.  De combinaties zijn dan tegenwoordige tijd met
tegenwoordige tijd (1), tegenwoordige tijd met verleden tijd (2), verleden
tijd met tegenwoordige tijd (3) en verleden tijd met verleden tijd (4).
<p>
<table>
<tr><td>(1)<td>De minister verklaart dat de economie blijft groeien</tr>
<tr><td>(2)<td>De minister verklaart dat hij het dossier in de trein liet liggen</tr>
<tr><td>(3)<td>De minister verklaarde dat de economie blijft groeien</tr>
<tr><td>(4)<td>De minister verklaarde dat hij het dossier in de trein liet liggen</tr>
</table>`,
		`//node[@rel="hd" and %PQ_ott% and ../node[@rel="vc"]/node[@rel="body"]/node[@rel="hd" and %PQ_ott%]]`,
		SPOD_STD,
		"ottott",
		`ott, ott:  "hij zegt dat hij komt"`,
		"",
	},
	{
		"",
		`//node[@rel="hd" and %PQ_ott% and ../node[@rel="vc"]/node[@rel="body"]/node[@rel="hd" and %PQ_ovt%]]`,
		SPOD_STD,
		"ottovt",
		`ott, ovt:  "hij zegt dat hij kwam"`,
		"",
	},
	{
		"",
		`//node[@rel="hd" and %PQ_ovt% and ../node[@rel="vc"]/node[@rel="body"]/node[@rel="hd" and %PQ_ott%]]`,
		SPOD_STD,
		"ovtott",
		`ovt, ott:  "hij zei dat hij komt"`,
		"",
	},
	{
		"",
		`//node[@rel="hd" and %PQ_ovt% and ../node[@rel="vc"]/node[@rel="body"]/node[@rel="hd" and %PQ_ovt%]]`,
		SPOD_STD,
		"ovtovt",
		`ovt, ovt:  "hij zei dat hij kwam"`,
		"",
	},
	{
		`Inbedding in finiete zinnen//Bij deze queries wordt gekeken naar de complexiteit van de zinnen in
termen van de inbedding van finiete bijzinnen. Een hoofdzin zonder
finiete bijzin geldt dan als "geen inbedding". Een hoofdzin met een
finiete bijzin geldt als "minstens 1 finiete zinsinbedding. Indien de
finiete bijzin zelf ook een finiete bijzin bevat is er sprake van minstens
2 finiete zinsinbeddingen. En zo verder.`,
		`//node[%PQ_finiete_inbedding0%]`,
		SPOD_STD,
		"inb0",
		"finiete zinnen zonder inbedding",
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
		`Topicalisatie en Extractie//De eerste query, "np-topic is subject", bekijkt hoe vaak in een
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
		`//node[@his and not(node[@his])]`,
		SPOD_STD,
		"his",
		"\"alle\" woorden (nodes met attribuut @his)",
		"his",
	},
	{
		"",
		`//node[@his="normal" and not(node[@his])]`,
		SPOD_STD,
		"normal",
		"woorden uit het woordenboek of de namenlijst",
		"his",
	},
	{
		"",
		`//node[@his
       and
       not(@his=("normal","robust_skip","skip"))
       and
       not(node[@his])
]`,
		SPOD_STD,
		"onbeken",
		"woorden niet direct uit het woordenboek",
		"his",
	},
	{
		"",
		`//node[@his="compound" and not(node[@his])]`,
		SPOD_STD,
		"compoun",
		"woorden herkend als samenstelling",
		"his",
	},
	{
		"",
		`//node[@his="name" and not(node[@his])]`,
		SPOD_STD,
		"name",
		"woorden herkend als naam (maar niet uit namenlijst)",
		"his",
	},
	{
		"",
		`//node[@his
       and
       not(@his=("normal","compound","name","robust_skip","skip"))
       and
       not(node[@his])
]`,
		SPOD_STD,
		"noun",
		"onbekende woorden die niet als samenstelling of naam werden herkend",
		"his",
	},
}
