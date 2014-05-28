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
      - voor elke gebruiker of per gebruiker
      - max aantal tokens (woorden)
      - maximum lengte in tokens per zin, of:
      - maximale verwerktijd van een zin (alleen bij gebruik lokale parser?)
    - Geheugen
      - globaal instellen dmv `ulimit -v`
    - Processortijd
      - Max aantal corpora dat tegelijk verwerkt wordt, met de aanname
        dat hierbij een processor voor 100% gebruikt wordt.
        - Wanneer corpus geparst wordt door een alpino-server, dan gebruikt
          dit misschien maar een paar procent van een processor
          (afhankelijk van snelheid van server). Je zou dan al de
          tweede fase (verwerking voor MySQL) parallel kunnen starten. En
          dan kom je misschien nog niet aan de 100% van een processor.
      - Scheduler: als er meerdere corpora van een enkele gerbuiker in de
        wachtrij staan de verwerking afwisselen met verwerking corpora van
        andere gebruikers.
  - Toegang
    - deny/allow op basis van ip-adressen
    - deny/allow op basis van e-mailpatroon
    - leesrechten: alleen raadplegen van algemene corpora
    - inlogrechten: mag zelf corpora uploaden
    - downloadrechten, alleen wie is ingelogd
      - eigen corpora
      - corpora die door anderen gedeeld zijn en vrijgegeven voor downloaden
  - Alpino-server
    - Aanpassen aan API van de server (huidige server is te oud)
  - Middelen voor beheerder
    - Statistiek
      - pqstatus
      - handler voor interne status (NumGoroutine, MemStat)
    - Gebruiker verwijderen
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
      - uploaden: tekst in zip, dact
      - downloaden: xml in zip, dact, compact...
      - browse: zinnen, stdout, stderr
      - modernere interface
    - Keuze van corpus
      - Menu met submenu's: algemeen; eigen corpora; anderen -> gedeelde corpora
  - Code
    - Organisatie + documentatie
  - Licentie

### Diversen ###

alles:

  - alle tekst in het Nederlands
  - logo

`pqserve`:

  - benaming "Woordrelaties" overal vervangen door "PaQu" of "paqu" of logo
  - config-optie `cookiepath` verwijderen, afleiden uit `url`
  - is het beter de optie "delen met iedereen" uit te schakelen?
  - foutafhandeling als gebruiker submit doet zonder bestand gekozen te hebben
  - fout van shell: niet de errorcode, maar de laatste regel(s) van
    stderr gebruiken als melding aan de gebruiker

`pqstatus`:

  - gebruikers zonder corpora
  - corpora zonder gebruikers

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
