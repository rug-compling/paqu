# PaQu #

Parse & Query

----------------------------------------------------------------

## Status ##

Grotendeels voltooid. Zie TODO.md voor wat er nog moet gebeuren.

----------------------------------------------------------------

## Overzicht ##

PaQu is een webserver waarin gebruikers Nederlandstalige tekstdocumenten
kunnen laten analyseren, en vervolgens onderzoeken op woordrelaties.

Iedereen kan er de publieke corpora ondervragen, maar na het aanmaken
van een account kan een gebruiker ook zelf documenten uploaden en laten
verwerken tot een nieuw corpus. Vervolgens kan de gebruiker dat nieuwe
corpus alleen zelf raadplegen, of delen met andere gebruikers die het
dan ook kunnen raadplegen.

Momenteel draait dit op de volgende site: http://zardoz.service.rug.nl:8067/

----------------------------------------------------------------

## Systeemeisen ##

Een 64-bit Linux-machine (amd64), ruim voldoende geheugen en
schijfruimte, het liefst veel processoren...

Je hebt een account op een MySQL-server nodig.

----------------------------------------------------------------

## Benodigde software ##

### Graphviz: Graph Visualization Software ###

Van Graphviz heb je de library's en de headers nodig. Op Debian worden
die geleverd door het pakket `graphviz-dev` of `libgraphviz-dev`.
Misschien is die versie te oud (TODO: welke versie is nodig?). Dan kun
je het downloaden van http://www.graphviz.org/

### DbXML: Oracle Berkeley DB XML ###

Dit is optioneel, maar wel sterk aanbevolen. Zonder DbXML kun je geen
Dact-bestanden verwerken. Dact-bestanden zijn corpora in het formaat
gebruikt door [Dact](http://rug-compling.github.io/dact/), waarmee
corpora bekeken en geanalyseerd kunnen worden.

Hier kun je DbXML downloaden:
http://www.oracle.com/technetwork/database/database-technologies/berkeleydb/downloads

Let op dat je de software met XML in de naam downloadt.

### Go ###

De meeste programma's zijn geschreven in de programmeertaal Go. Je hebt
dus de Go-compiler nodig. Op Debian is Go beschikbaar in het pakket
`golang`, maar die versie is bij schrijven (juli 2014) te oud.

Je hebt minimaal versie 1.2 nodig. Geef het commando `go version` om dit
te controleren.

Je kunt Go downloaden van http://golang.org/

### Flex: A fast lexical analyzer generator ###

Voor het compileren van één programma heb je het programma flex nodig.
Op Debian is dat beschikbaar in het pakket `flex`.

Dit is dus iets anders dan Adobe Flex of Apache Flex

### Alpino: A dependency parser for Dutch ###

Om PaQu te draaien heb je van
[Alpino](http://www.let.rug.nl/vannoord/alp/Alpino/) de tokenizer en de
parser nodig.

In de nabije toekomst is het mogelijk gebruik te maken van een
Alpino-server, die ergens op een andere website draait. Als je dat
gebruikt heb je zelf van Alpino alleen de tokenizer nodig.

----------------------------------------------------------------

## Installatie ##

PaQu maakt gebruik van een directory voor het opslaan van gegevens.
De default is `$HOME/.paqu`, maar je kunt dit veranderen door de
environment-variabele `PAQU` te zetten. Voor de onderstaande
beschrijving gaan we ervan uit dat je de default gebruikt.

Verder gaan we ervan uit dat je PaQu in `$HOME/paqu` hebt geplaatst
(zonder punt ervoor).

### Compileren ###

Ga naar de directory `~/paqu/src`, kopieer het bestand
`Makefile.cfg.example` naar `Makefile.cfg` en pas het aan. Volg daarvoor
de aanwijzingen in het bestand.

Als je klaar bent, run `make all`. Als alles goed gaat vind je alle
programma's in `~/paqu/bin`. Verplaats de programma's naar een plek in
je PATH, of voeg de directory toe aan PATH.

### Configuratie ###

Voordat je de programma's kunt draaien moet je een configuratiebestand
maken. Kopieer het bestand `~/paqu/setup-example.toml` naar
`~/.paqu/setup.toml` en pas het aan door de instructies in het bestand
op te volgen. Geef het bestand de rechten 600 zodat andere gebruikers op
je computer de wachtwoorden niet kunnen lezen:

    chmod 600 ~/.paqu/setup.toml

Je kunt de configuratie controleren met het commando `pqconfig`. Als je
een typefout maakt in een label in het configuratiebestand, dan wordt
dat item genegeerd. Door het commando `pqconfig` te draaien kun je zien of
alles de juiste waarde heeft. Omdat dit programma ook de wachtwoorden
laat zien heeft het de rechten 700, zodat andere gebruikers op je
computer het programma niet kunnen gebruiken.

### HTTPS ###

Wil je https gebruiken, dan dien je een certificaat aan te maken en die te
laten ondertekenen door een instantie die algemeen erkend wordt.
Doe je dat laatste niet, dan krijgt de gebruiker in z'n webbrowser een
waarschuwing dat de site niet te vertrouwen is, en dat ie de site beter
niet kan bezoeken.

#### HTTPS: Aanmaken certificaat ####

Stel dat je server voor gebruikers toegankelijk is via
`paqu.myserver.nl`, geef dan de volgende commando's:

    cd ~/.paqu
    go run `go env GOROOT`/src/pkg/crypto/tls/generate_cert.go -host=paqu.myserver.nl

Er staan nu twee nieuwe bestanden in `~/.paqu`: `cert.pem` en `key.pem`

#### HTTPS: Ondertekenen van certificaat ####

Voor meer info over het laten ondertekenen van een certificaat, zie
http://nl.wikipedia.org/wiki/Certificaatautoriteit

Voor leveranciers, zie https://www.google.nl/search?q=SSL+certificaten

Als je het certificaat hebt laten ondertekenen dien je die ondertekening
onderaan `cert.pem` toe te voegen.

### Initialiseer de database ###

Om de database te initialiseren, run:

    pqinit

### Lassy Klein ###

Het is zeer aan te bevelen om alvast een corpus te installeren voordat
mensen PaQu gaan gebruiken, omdat ze er dan mee kunnen werken zonder
eerst in te hoeven loggen en hun eigen corpus aan te hoeven maken.

Het is raadzaam om het corpus
[Lassy Klein](http://tst-centrale.org/producten/corpora/lassy-klein-corpus/6-66)
te installeren, omdat alle voorbeelden die de gebruiker in de informatie ziet werken op dat corpus.

Als je dit corpus hebt geïnstalleerd in `~/corpora/LassySmall` kun je
het zo opnemen in PaQu:

    find ~/corpora/LassySmall/Treebank -name '*.xml' | pqbuild -w lassysmall 'Lassy Klein' none 1

Dit duurt wel even, dus wellicht wil je bovenstaande regel in een script
zetten en dat aanroepen met `nohup`.

Nadat het corpus in PaQu is opgenomen moet je de xml-bestanden van Lassy Klein
laten staan waar ze staan. Die xml-bestanden worden zelf niet in PaQu
opgenomen, alleen de analyse van die bestanden. De bestanden zelf zijn
ook nog nodig voor het werken met PaQu.

Wel kun je de xml-bestanden na verwerking comprimeren met `gzip`:

    find ~/corpora/LassySmall/Treebank -name '*.xml' | xargs gzip

----------------------------------------------------------------

## Werken met PaQu ##

### De server starten ###

Je kunt PaQu nu heel simpel starten met `pqserve`. Maar dat is wellicht
niet de beste methode. Beter is het om het script `~/paqu/run.sh` aan te
passen aan je wensen, en regelmatig (bijvoorbeeld elk kwartier) aan te
roepen vanuit `cron`. Dit script test of de server nog reageert
(daarvoor gebruikt het de link `/up` op de server), en start het opnieuw
mocht dit niet zo zijn. Het print dan ook eventuele foutmeldingen van de
vorige run van pqserve, en als je `MAILTO` goed hebt gezet in
`crontab -e`, dan worden die foutmeldingen naar je toegestuurd.

### Log-bestanden ###

Het programma `pqserve` maakt logbestanden aan in directory `~/.paqu`.
Deze bestanden worden automatisch geroteerd als ze te groot worden, dus
je hoeft geen `logrotate` te gebruiken.

### Status ###

Het programma `pqstatus` geeft je een overzicht van gebruikers en
corpora. Wellicht is dit handig om eenmaal per dag aan te roepen vanuit
`cron`, dat dan de resultaten naar je toestuurt.

Info over de interne staat van de server kun je opvragen via deze paden
op de server:

    /debug/vars
    /debug/env

### Beheer ###

Met het programma `pqrmuser` kun je een gebruiker verwijderen, inclusief
alle data die door die gebruiker is opgeslagen, zowel van schijf als uit
de MySQL-database.

Met het programma `pqsetquota` kun je het quotum voor een of meer
gebruikers aanpassen.
