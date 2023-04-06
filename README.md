# PaQu #

Parse & Query

Metadata: https://rug-compling.github.io/paqu/

#### Docker ##

Je kunt PaQu nu ook in Docker draaien. Zie:
https://github.com/rug-compling/paqu-docker

#### Vergelijkbare toepassingen

 * [AlpinoGraph](https://urd2.let.rug.nl/~kleiweg/alpinograph/)
 * [GrETEL](https://gretel.hum.uu.nl/)

----------------------------------------------------------------

## Overzicht ##

PaQu is een webserver waarin gebruikers Nederlandstalige tekstdocumenten
kunnen laten analyseren, en vervolgens onderzoeken op woordrelaties.

Iedereen kan er de publieke corpora ondervragen, maar na het aanmaken
van een account kan een gebruiker ook zelf documenten uploaden en laten
verwerken tot een nieuw corpus. Vervolgens kan de gebruiker dat nieuwe
corpus alleen zelf raadplegen, of delen met andere gebruikers die het
dan ook kunnen raadplegen.

Momenteel draait dit op de volgende site: http://www.let.rug.nl/alfa/paqu

----------------------------------------------------------------

## Systeemeisen ##

Een 64-bit Linux-machine (amd64), ruim voldoende geheugen en
schijfruimte, het liefst veel processoren...

Je hebt een account op een MySQL-server of MariaDB-server nodig.

----------------------------------------------------------------

## Benodigde software ##

### Locale ###

De locale `en_US.utf8` moet geïnstalleerd zijn bij het draaien van PaQu.

### Graphviz: Graph Visualization Software ###

Van Graphviz heb je de library's en de headers nodig. Op Debian worden
die geleverd door het pakket `graphviz-dev` of `libgraphviz-dev`.
Misschien is die versie te oud. Als je een compileerfout krijgt moet je
een nieuwere versie hebben. Die kun je downloaden van
http://www.graphviz.org/

Voor de juiste weergave in de labels van bomen, van letters met een code
boven de U+FFFF (zeldzaam, maar komt voor in Lassy Large, bijv: 𐎠𐎼𐎫𐎧𐏁𐏂𐎠)
heb je minimaal versie 2.36.0 nodig. Vind je dit niet belangrijk, dan
kun je een oudere versie gebruiken.

### DbXML: Oracle Berkeley DB XML ###

Dit is optioneel, maar wel sterk aanbevolen. Zonder DbXML kun je geen
Dact-bestanden verwerken, en kun je geen XPATH-query's uitvoeren.
Dact-bestanden zijn corpora in het formaat gebruikt door
[Dact](http://rug-compling.github.io/dact/), waarmee corpora bekeken en
geanalyseerd kunnen worden.

Hier kun je DbXML downloaden:
http://www.oracle.com/technetwork/database/database-technologies/berkeleydb/downloads

Let op dat je de software met XML in de naam downloadt.

LET OP: Versies 6.x.x t/m 6.1.4 hebben een ernstige bug waardoor
de zoekresultaten met XPATH vaak verre van compleet zijn. (Versie 6.1.4
is op moment van schrijven de meest actuele versie.)

Je kunt versie 2.5.16 hier downloaden: https://github.com/rug-compling/dbxml

### Go: The Go Programming Language ###

De meeste programma's zijn geschreven in de programmeertaal Go. Je hebt
dus de Go-compiler nodig.

Je hebt minimaal Go versie 1.8 nodig. Geef het commando `go version` om
dit te controleren.

Je kunt Go downloaden van http://golang.org/

### Flex: A fast lexical analyzer generator ###

Voor het compileren van één programma heb je het programma flex nodig.
Op Debian is dat beschikbaar in het pakket `flex`.

Dit is dus iets anders dan Adobe Flex of Apache Flex

### Alpino: A dependency parser for Dutch ###

Om PaQu te draaien heb je van
[Alpino](http://www.let.rug.nl/vannoord/alp/Alpino/) de tokenizer en de
parser nodig.

Je kunt ook een Alpino-server gebruiken die voldoet aan de
[Alpino API](https://github.com/rug-compling/alpino-api).
Als je dat doet heb je zelf van Alpino alleen de tokenizer nodig.
Zo'n server kan in principe de data parallel verwerken, en dus veel
sneller zijn dan wanneer je Alpino lokaal gebruikt.

----------------------------------------------------------------

## Installatie ##

PaQu maakt gebruik van een directory voor het opslaan van gegevens, en
een voor het configuratiebestand.

Voor de bepaling van deze directory's wordt eerst gekeken naar de
environment variable `PAQU`. Als die gedefinieerd is:

    data:   $PAQU
    config: $PAQU

Als die niet gedefinieerd is, dan wordt gekeken naar de waarde die bij
het compileren is gedefineerd.

Als deze ook niet is gedefineerd, dan gelden de volgende waardes:

    data:   $XDG_DATA_HOME/paqu
    config: $XDG_CONFIG_HOME/paqu

Als de environment variabelen `XDG_DATA_HOME` of `XDG_CONFIG_HOME` niet
zijn gedefinieerd, dan gelden de volgende waardes:

    data:   $HOME/.local/share/paqu
    config: $HOME/.config/paqu

Hieronder worden deze waardes gesymboliseerd door {data} en {config}.

Verder gaan we ervan uit dat je PaQu in `$HOME/paqu` hebt geplaatst.

### Compileren ###

Ga naar de directory `~/paqu/src`, kopieer het bestand
`Makefile.cfg.example` naar `Makefile.cfg` en pas het aan. Volg daarvoor
de aanwijzingen in het bestand.

Als je klaar bent, run `make all`. Als alles goed gaat vind je alle
programma's in `~/paqu/bin`. Verplaats de programma's naar een plek in
je PATH, of voeg de directory toe aan `PATH`.

### Configuratie ###

Voordat je de programma's kunt draaien moet je een configuratiebestand
maken. Kopieer het bestand `~/paqu/run/setup-example.toml` naar
`{config}/setup.toml` en pas het aan door de instructies in het bestand
op te volgen. Geef het bestand de rechten 600 zodat andere gebruikers op
je computer de wachtwoorden niet kunnen lezen:

    chmod 600 {config}/setup.toml

Dit voorkomt ook dat anderen de programma's kunnen draaien en daarmee
zomaar corpora van willekeurig wie kunnen verwijderen, bijvoorbeeld.

Je kunt de configuratie controleren met het commando `pqconfig`. Als je
een typefout maakt in een label in het configuratiebestand, dan wordt
dat item genegeerd. Door het commando `pqconfig` te draaien kun je zien of
alles de juiste waarde heeft.

### HTTPS ###

Wil je https gebruiken, dan dien je een certificaat aan te maken en die te
laten ondertekenen door een instantie die algemeen erkend wordt.
Doe je dat laatste niet, dan krijgt de gebruiker in z'n webbrowser een
waarschuwing dat de site niet te vertrouwen is, en dat ie de site beter
niet kan bezoeken.

#### HTTPS: Aanmaken certificaat ####

Stel dat je server voor gebruikers toegankelijk is via
`paqu.myserver.nl`, geef dan de volgende commando's:

    cd {config}
    go run `go env GOROOT`/src/crypto/tls/generate_cert.go -host=paqu.myserver.nl

Krijg je een foutmelding `no such file or directory`, doe dan:

    go run `go env GOROOT`/src/pkg/crypto/tls/generate_cert.go -host=paqu.myserver.nl

Er staan nu twee nieuwe bestanden in `{config}`: `cert.pem` en `key.pem`

#### HTTPS: Ondertekenen van certificaat ####

Voor meer info over het laten ondertekenen van een certificaat, zie
http://nl.wikipedia.org/wiki/Certificaatautoriteit

Voor leveranciers, zie https://www.google.nl/search?q=SSL+certificaten

Als je het certificaat hebt laten ondertekenen dien je die ondertekening
onderaan `cert.pem` toe te voegen.

### Initialiseer de database ###

Om de database te initialiseren, run:

    pqinit

Als je een nieuwere versie van PaQu gaat gebruiken, dan moet misschien de
database aangepast worden. Om dat te doen, run:

    pqupgrade

### Lassy Klein ###

Het is zeer aan te bevelen om alvast een corpus te installeren voordat
mensen PaQu gaan gebruiken, omdat ze er dan mee kunnen werken zonder
eerst in te hoeven loggen en hun eigen corpus aan te hoeven maken.

Het is raadzaam om het corpus
[Lassy Klein](http://tst-centrale.org/producten/corpora/lassy-klein-corpus/6-66)
te installeren, omdat alle voorbeelden die de gebruiker in de informatie
ziet werken op dat corpus.

Als je dit corpus hebt geïnstalleerd in `~/corpora/LassySmall` kun je
het zo opnemen in PaQu:

	pqdactx ~/corpora/LassySmall/lassy.dact ~/corpora/LassySmall/lassy.dactx
	echo ~/corpora/LassySmall/lassy.dact | \
	    pqbuild -w -p '.*/corpora/LassySmall/' lassysmall 'Lassy Klein' none 1

Als je geen DbXml beschikbaar hebt kun je het zo opnemen in PaQu:

    find ~/corpora/LassySmall/Treebank -name '*.xml' | sort | \
	    pqbuild -w -p '.*/corpora/LassySmall/Treebank/' lassysmall 'Lassy Klein' none 1

De string achter -p is een reguliere expressie die het deel van elke
bestandsnaam matcht dat wordt verwijderd om van de bestandsnaam een label
te maken.

Dit duurt wel even, dus wellicht wil je bovenstaande regel in een script
zetten en dat aanroepen met `nohup`.

Nadat het corpus in PaQu is opgenomen moet je het dact-bestand (of de
xml-bestanden) van Lassy Klein laten staan waar het staat. Dat bestand
wordt zelf niet in PaQu opgenomen, alleen de analyse ervan. Het bestand
zelf is nog steeds nodig voor het werken met PaQu.

Wel kun je, als je xml-bestanden hebt gebruikt i.p.v. het dact-bestand,
die xml-bestanden na verwerking comprimeren met `gzip`:

    find ~/corpora/LassySmall/Treebank -name '*.xml' | xargs gzip

### Alpino Treebank ###

Kun je Lassy Klein niet installeren, installeer dan het kleinere corpus
Alpino Treebank, te vinden als `cdb.dact` in
[corpora.tar.gz](http://www.let.rug.nl/alfa/docker/paqu/corpora.tar.gz).

----------------------------------------------------------------

## Werken met PaQu ##

### De server starten ###

Je kunt PaQu nu heel simpel starten met `pqserve`. Maar dat is wellicht
niet de beste methode. Beter is het om het script `~/paqu/run/run.sh`
aan te passen aan je wensen, en regelmatig (bijvoorbeeld elk kwartier)
aan te roepen vanuit `cron`. Dit script test of de server nog reageert
(daarvoor gebruikt het de link `/up` op de server), en start het opnieuw
mocht dit niet zo zijn. Het print dan ook eventuele foutmeldingen van de
vorige run van pqserve, en als je `MAILTO` goed hebt gezet in
`crontab -e`, dan worden die foutmeldingen naar je toegestuurd.

### Log-bestanden ###

Het programma `pqserve` maakt logbestanden aan in directory `{data}`.
Deze bestanden worden automatisch geroteerd als ze te groot worden, dus
je hoeft geen `logrotate` te gebruiken.

### Status ###

Het programma `pqstatus` geeft je een overzicht van gebruikers en
corpora. Wellicht is dit handig om eenmaal per dag aan te roepen vanuit
`cron`, dat dan de resultaten naar je toestuurt.

Info over de interne staat van de server kun je opvragen via deze paden
op de server:

    /debug/env 
    /debug/stack
    /debug/vars

### Beheer ###

Met het programma `pqclean` verwijder je automatisch alle gebruikers
zonder corpora die twee maanden niet actief geweest zijn.

Met het programma `pqrmcorpus` kun je een corpus uit de
MySQL/MariaDB-database verwijderen. Is het een corpus dat door een
gebruiker is ge-upload, dan worden de bijbehorende bestanden ook van
schijf verwijderd.

Met het programma `pqrmuser` kun je een gebruiker verwijderen, inclusief
alle data die door die gebruiker is opgeslagen, zowel van schijf als uit
de MySQL/MariaDB-database.

Met het programma `pqsetquota` kun je het quotum voor een of meer
gebruikers aanpassen.
