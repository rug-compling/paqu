# PaQu #

## Te doen ##

Zie ook: TODO in diverse bestanden

  - Adminhandleiding
    - Alpino Treebank als alternatief voor Lassy Klein
  - pqserve:
    - Bij tellingen die uit meerdere onderdelen bestaan, gebruiker laten
      kiezen welke onderdelen hij wil zien. Al getoonde tellingen
      inklapbaar maken.
    - Commentaren en labels in invoer van doorlopende tekst negeren. Documenteren.
      - LET OP: ook pqserve/invoer.go aanpassen
    - Expert-opties bij invoer nieuw corpus, met eventueel:
      - Escapen van invoer die al getokeniseerd is: none/half/full
      - Maximaal aantal tokens per zin
      - Kleinere timeout
      - Alternatieve parser voor corpus dat voornamelijk uit vragen bestaat
        - Optie `application_type=qa` vóór de optie `-parse`
      - Door gebruiker andere labels laten kiezen voor `paqu.path1`, `paqu.path2`, etc
