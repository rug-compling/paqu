# PaQu #

## Te doen ##

Zie ook: TODO in diverse bestanden

  - Alpino-server
    - Aanpassen aan API van de server (huidige server is te oud)
      - zie ook TODOs in `work.go`
  - Gebruikershandleiding
    - Introtekst op hoofdpagina
    - Tekst achter *Meer info...*
    - Helptekst voor pagina *XPath*
    - Helptekst voor pagina *Mijn corpora*
  - Interface
    - https zonder http? de combinatie maakt het erg ingewikkeld
    - Logo + balk wijzigen
  - Code
    - pqserve: gebruik/codering van labels is een rommeltje, dat moet
      beter, zie beneden
      - benaming: label/bestand
      - uitleg op tab browse: `%` en `_` in label
    - pqserve: alle javascript onderbrengen in één apart bestand
    - Organisatie + documentatie
  - pqserve:
    - XPath
      - macro's
    - weergave op mobiel of tablet


## Gebruik/codering van labels ##

Weergave labels in:

  - Tab Corpura
    - overzicht
	  - fouten
	  - alle zinnen
	- download zinnen
	- download xml
	- download dact
  - Onder boomweergave

Bronnen:

  - lokale bestanden (voor globale corpora)
    - weergave van labels zonder compleet path
  - geüpload door gebruikers, door PaQu geparst, zonder labels
    - labels kunnen ongewijzigd worden gebruikt
  - geüpload door gebruikers, door PaQu geparst, met labels
    - weergave van originele labels
	- bij downloaden (zip of dact) moeten originele labels worden
      gebruikt (met corpusnaam als prefix???)
  - geüpload door gebruikers, door server geparst, zonder labels
    - ???
  - geüpload door gebruikers, door server geparst, met labels
    - ???
  - geüpload als dact-bestand
    - weergave van originele labels
	- bij downloaden (zip of dact) moeten originele labels worden
      gebruikt (met corpusnaam als prefix???)
  - geüpload als zip-bestand van xml-bestanden
    - weergave van originele labels
	- bij downloaden (zip of dact) moeten originele labels worden
      gebruikt (met corpusnaam als prefix???)
  - geüpload als treebank
    - ???  
