# PaQu #

## Te doen ##

Zie ook: TODO in diverse bestanden

  - Alpino-server
    - Aanpassen aan API van de server (huidige server is te oud)
      - zie ook TODOs in `work.go`
  - Gebruikershandleiding
    - Tekst voor pagina *Info*: onder kopje *Metadata*
    - Tekst voor pagina *Info*: onder kopje *Corpora*
  - Adminhandleiding
    - optie `loginurl` in `setup-example.toml`
	- gebruik van `$PAQU/contact.html`
  - Metadata
    - Hoe voer je dat in?
  - pqserve:
    - https zonder http? de combinatie maakt het erg ingewikkeld
    - gewoon zoeken
	  - Nu alleen link naar boom achter gevonden zinnen. Toevoegen: link
        naar boom achter elke woord/tag-combinatie onder elke zin.
    - XPath
      - zoeken: POST i.p.v. GET?
	    - anders url misschien te lang
		- maar dan werken links naar volgende pagina niet meer
		- pagina updaten i.p.v. opnieuw laden bij uitvoeren van zoeken?
    - weergave op mobiel of tablet
	  - weergave van attributen in bomen
	- pqserve stopt niet altijd direct bij sigterm, en na sigkill zijn logs niet compleet
    - Onder *Metadata*: subcorpus maken op basis van metadata
     - Vanuit meerdere corpora, dan moeten die corpora dezelfde soorten metadata hebben.
     - setHigh/unHigh
    - Samenvoegen van corpora (na gewoon zoeken, zoeken met xpath, selectie op metadata)
     - Probleem: Samenvoegen van corpora met incompatibele metadata
       (text vs int vs float vs date vs datetime).
      - Als dingen botsen krijgt de gebruiker vanzelf een foutmelding.
        Maar als één subcorpus metadata heeft die een ander subcorpus
        niet heeft, dan de gebruiker waarschuwen?
  - programma tags.go toevoegen
    - kopiëren uit /net/aistaff/alfa/lassy
    - aanpassen
	- documentatie
  - DbXml
    - zonodig bestanden in verkeerde DbXml-versie omzetten naar goede versie?
