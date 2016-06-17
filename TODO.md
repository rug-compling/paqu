# PaQu #

## Te doen ##

Zie ook: TODO in diverse bestanden

  - PaQu splitsen in twee delen?
    Probleem: extra gedoe voor mensen die data willen uploaden en gelijk
    willen gebruiken in PaQu
     1. deel voor uploaden data en downloaden geparste data
     2. deel voor querying, alleen geparste data uploaden
     3. nieuw: deel voor omzetten data naar PaQu-formaat
  - Alpino-server
    - Komt die er nog? Lijkt er niet op.
    - Aanpassen aan API van de server (huidige server is te oud)
      - zie ook TODOs in `work.go`
  - Adminhandleiding
    - Alpino Treebank als alternatief voor Lassy Klein
  - pqserve:
    - check versie DbXML, uitschakelen als het een gebroken versie is
    - https zonder http? de combinatie maakt het erg ingewikkeld
    - XPath
      - zoeken: POST i.p.v. GET?
        - anders url misschien te lang
        - maar dan werken links naar volgende pagina niet meer
        - pagina updaten i.p.v. opnieuw laden bij uitvoeren van zoeken?
    - weergave op mobiel of tablet
      - weergave van attributen in bomen
    - pqserve stopt niet altijd direct bij sigterm, en na sigkill zijn logs niet compleet
    - Samenvoegen van corpora (na gewoon zoeken, zoeken met xpath, selectie op metadata)
      - Probleem: Samenvoegen van corpora met incompatibele metadata
        (text vs int vs float vs date vs datetime).
      - Als dingen botsen krijgt de gebruiker vanzelf een foutmelding.
        Maar als één subcorpus metadata heeft die een ander subcorpus
        niet heeft, dan de gebruiker waarschuwen?
    - Bij tellingen die uit meerdere onderdelen bestaan, gebruiker laten
      kiezen welke onderdelen hij wil zien?
  - pqalpino verwijderen, inbouwen in pqserve
  - pqtexter verwijderen, inbouwen in pqserve met gebruik van flexgo
    i.p.v. flex
  - pqinit verwijderen, en inbouwen in pqbuild?
  - DbXml
    - zonodig bestanden in verkeerde DbXml-versie omzetten naar goede versie?

## Verzoeken ##

In overweging:

  - pqserve:
    - Invoer nieuw corpus
      - Expert-opties voor Alpino
        - Kleinere timeout
        - Alternatieve parser voor corpus dat voornamelijk uit vragen bestaat
          - Optie `application_type=qa` vóór de optie `-parse`
      - Door gebruiker andere labels laten kiezen voor `paqu.path1`, `paqu.path2`, etc
