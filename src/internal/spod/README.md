
De huidige manier van testen van een nieuwe expressie in SPOD verloopt
als volgt:

 1. Neem die nieuwe expressie op in package `spod` en compileer `pqbuild`.
 2. Voer een uitgebreid testcorpus in in PaQu. Alle waardes voor SPOD
    worden berekend door Go.
 3. Verwijder de berekende waardes voor de nieuwe expressie uit
    *paqu.data.dir*`/data/testcorpus/spod` en sla die ergens anders
    op.
 4. Compileer en herstart `pqserve`, en vraag de SPOD-data op (html of
    tekst). Nu worden de verwijderde waardes opnieuw aangemaakt, deze
    keer door DbXML.
 5. Vergelijk de waardes van DbXML met die van Go.

TODO: Bovenstaande manier is omslachtig. Er moet een apart
testprogramma voor komen.

----

Huidige beperkingen van de XPath-parser in Go:

 * Het werkt alleen met de standaardmethode, niet met geëxpandeerde
   index-nodes.
 * Je kunt alleen zoeken op `/node`, niet op `/parser` of `/root`, etc.

----

Herschrijving van XPath-expressies door package `spod` om bepaalde
beperkingen van de parser te omzeilen:

 * `[1]`, `[3]`, `[-1]` → `[first()]`, `[third()]`, `[last()]`
 * `][` → `]/self::node[`
 * `'string'` → `"string"`
 * `=("a","b")` → `="hashcode"`
 * `number(iets)` → ` iets ` --- geen haakjes!

----

Hoewel alle gegevens over elke SPOD-expressie opgenomen zijn in het
bestand `spoddata.go` zijn er toch gevallen die apart behandeld moeten
worden.

In `pqbuild/spod.go` speelt dit voor de **specials** `hidden1`, `attr`,
`parser`, `qm -yn`, `pos_verb`, `sc`, `his`, en voor de **labels** `ok`,
`cats0` t/m `cats4`, `skips0` t/m `skips4`.

In `pqserve/spodtable.go` speelt dit voor de **specials** `hidden1`,
`attr`, `parser`, `his` en voor de **labels** `pos`, `postag`, `pt`,
`cats0` t/m `cats4`, `skips0` t/m `skips4`.

----

De parser voor XPath, `spodmake/spodmake.go`, die XPath-expressies omzet naar Go-code, maakt
gebruik van een gehackte versie van `testXPath` (onderdeel van
`libxml2` versie 2.9.9), waarin de indentatie van de uitvoer met optie
`--tree` niet beperkt is tot 25 niveaus.

Sources: [hier](https://github.com/rug-compling/alud/tree/master/libxml2/2.9.9).

Zie ook: [feature request](https://gitlab.gnome.org/GNOME/libxml2/-/issues/241).
