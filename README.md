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

## Te doen, ideeën ##

alles:

  - alle tekst in het Nederlands
  - code documenteren
  - installatiehandleiding
  - logo

`pqserve`:

  - is het beter de optie "delen met iedereen" uit te schakelen?
  - foutafhandeling als gebruiker submit doet zonder bestand gekozen te hebben
  - fout van shell: niet de errorcode, maar de laatste regel(s) van
    stderr gebruiken als melding aan de gebruiker
  - optie: gebruik van een Alpinoserver, is gedaan, maar moet nog
       aangepast worden aan uitvoerformaat van de server (huidige is te oud)
  - keuze van corpus: sticky
    * query zonder db → default uit setup
    * geen query en geen db → uit cookie
  - https
  - alpino draaien met debug=1
  - uploaden van documenten in andere formaten, zoals zipbestand, dactbestanden
  - download corpus als dactbestand
  - browse zinnen, stdout, stderr
  - config: deny/allow op basis van e-mailpatroon
   - lees/query-rechten
   - schrijfrechten (wie een account mag aanmaken)
   - downloadrechten (alleen eigenaar van een corpus)
  - handler voor internse status (NumGoroutine, MemStat, data van pqstatus)
  - limieten
   - geheugen (totaal)
   - processorgebruik (totaal)
   - schijfruimte (totaal, per gebruiker)
   - ruimte in database (totaal, per gebruiker)
  - beheerdersfuncties:
   - gebruik van resources bekijken
   - accounts beheren
   - e-mail aan beheerder bij problemen

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

 - Oracle Berkeley DB XML: libs en headers geïnstalleerd op een standaardlokatie
 - Go versie 1.2 of hoger
 - Alpino: tokenizer
 - Alpinoparser of toegang tot een alpinoserver
