# PaQu #

Parse & Query

----------------------------------------------------------------

## Status ##

In ontwikkeling, maar al bruikbaar via http://zardoz.service.rug.nl:8067/

----------------------------------------------------------------

## Overzicht ##

Een webserver waarin gebruikers Nederlandstalige tekstdocumenten kunnen
laten analyseren, en vervolgens onderzoeken op woordrelaties.

### PaQu ###

Dit is een server voor meerdere gebruikers. Iedereen kan er de publieke
corpora ondervragen, maar na het aanmaken van een account kan een
gebruiker ook zelf documenten uploaden en laten verwerken tot een nieuw
corpus. Vervolgens kan de gebruiker dat nieuwe corpus alleen zelf
raadplegen, of delen met andere gebruikers die het dan ook kunnen
raadplegen.

Momenteel draait dit op de volgende site: http://zardoz.service.rug.nl:8067/

### PaQu1 ###

PaQu1 is PaQu voor een enkele gebruiker. Er zijn geen accounts. De
gebruiker moet zelf een account hebben op MySQL, en moet zelf corpora
laten parsen door Alpino. Met `pqbuild1` wordt een corpus verwerkt en
klaargezet. Vervolgens kan met `pqserve1` een server gestart worden om
de corpora te raadplegen.

Het database-formaat in MySQL van PaQu1 is niet compatibel met dat van PaQu.

----------------------------------------------------------------

## Te doen ##

### Hoofdtaken ###

  - Quota
    - Corpora
      - ✔ default voor elke gebruiker in setup, door admin aan te passen
        per gebruiker
      - ✔ max aantal tokens (woorden)
      - maximale verwerktijd van een zin (alleen bij gebruik lokale
        parser): in setup
    - Geheugen
      - globaal instellen dmv `ulimit -v`
    - Processortijd
      - ✔  Max aantal corpora dat tegelijk verwerkt wordt, met de aanname
        dat hierbij een processor voor 100% gebruikt wordt.
      - Wanneer corpus geparst wordt door een alpino-server: verdere
        verwerking in parallel met verwerking door de alpinoserver
      - Pas als blijkt dat dit nodig is: Scheduler: als er meerdere
        corpora van een enkele gerbuiker in de wachtrij staan de
        verwerking afwisselen met verwerking corpora van andere
        gebruikers.
  - Toegang
    - ✔ leesrechten: deny/allow op basis van ip-adressen
      - IPv6?
    - ✔ inlogrechten: deny/allow op basis van regexp mail-adressen
    - ✔ downloadrechten, alleen wie is ingelogd, alleen eigen corpora
  - Alpino-server
    - Aanpassen aan API van de server (huidige server is te oud)
  - Middelen voor beheerder
    - Statistiek
      - pqstatus
      - handler voor interne status (NumGoroutine, MemStat)
      - package `expvar`
    - ✔ Gebruiker verwijderen
    - Quotum voor specifieke gebruiker veranderen
    - E-mail aan beheerder bij problemen
  - Adminhandleiding
    - installatie
    - (account-)beheer
  - Gebruikershandleiding
    - Introtekst op hoofdpagina
    - Tekst achter Meer info...
  - Interface
    - Algemene opmaak
    - Beheer van corpora
      - uploaden:
        - meerdere tekstdocumenten in een zipbestand
        - dact
        - regels die al een label hebben
      - downloaden:
        - ✔ xml in zip
        - dact
      - ✔ browse: zinnen, stdout, stderr
      - ✔ modernere interface: verwijderen, opties
      - modernere interface voor afhandeling van opties
    - Keuze van corpus
      - Lijst van beschikbare (gedeelde) corpora die
        toegevoegd/verwijderd kunnen worden in menu
  - Code
    - Organisatie + documentatie
    - Ontwikkeling: `make`, `go fmt`, `go vet`, `golint`
  - Dependencies
    - externe go-pakketten opnemen onder directory `_vendor`
  - Voor later
    - Misschien: subboompjes met statistiek en zoeklink, vergt meer info in database
  - pqserve1 en pqbuild1 verplaatsen? naar ander repository?
  - Licentie


### Diversen ###

alles:

  - alle tekst in het Nederlands
  - logo + balk wijzigen (pqserve, pqserve1)
  - zie TODO in diverse bestanden

`pqserve`:

  - is het beter de optie "delen met iedereen" uit te schakelen?
  - delen beperken: alleen corpora met een minimum grootte?
  - foutafhandeling als gebruiker een leeg bestand uploadt
  - fout van shell: niet de errorcode, maar de laatste regel(s) van
    stderr gebruiken als melding aan de gebruiker

`pqstatus`:

  - gebruikers zonder corpora
  - corpora zonder gebruikers
  - compacte weergave

----------------------------------------------------------------

## Installatie ##

### Systeemeisen ###

Een Linux-machine, veel geheugen, veel processoren...

### Wat je nodig hebt, op Debian ###

pakketten:

 - `mercurial`
 - `git`
 - `graphviz-dev` of `libgraphviz-dev`

overige software:

 - (optioneel, nodig voor verwerking van dact-bestanden)
   [Oracle Berkeley DB XML](http://www.oracle.com/technetwork/database/database-technologies/berkeleydb/downloads):
   libs en headers geïnstalleerd op een standaardlokatie.
 - [Go](http://golang.org/): als met DB XML, dan versie 1.2 of hoger, anders versie 1.0 of hoger
 - [Alpino](http://www.let.rug.nl/vannoord/alp/Alpino/): tokenizer
 - Alpinoparser of toegang tot een alpinoserver
