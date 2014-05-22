
Voorbeelden preparen van een corpus:

    find /net/corpora/LassySmall/Treebank -name '*.xml' | ./pqbuild1 loginfile lassyklein 'Lassy Klein'
    
    find /net/corpora/LassyLarge/NLWIKI20110804 -name '*.dact' | ./pqbuild1 loginfile lassygroot 'Lassy Groot'

Namen van databestanden worden ingelezen van stdin, 1 bestand per regel.
Wanneer een path met bestandsnaam relatief is wordt die automatisch
omgezet in een absoluut path.

De volgende bestanden zijn toegestaan:

Bestand met extensie .xml : 1 parse per bestand
Bestand met extensie .dact : een dbxml-bestand, zoals aangemaakt door alpinocorpus
Bestand met extensie .data.dz : een compact corpusbestand, zoals aangemaakt door alpinocorpus

Je kunt data toevoegen aan een set door de optie -a te gebruiken. Let op
dat je geen data invoegt die er al in stond. Het programma controlleert
niet of er zinnen dubbel in zitten.

Je kunt een bestaande dataset overschrijven met de optie -w.

----

Het eerste argument na het commando pqbuild1 is de naam van een bestand
waarin de string staat om in te loggen op een MySQL-server. Voorbeelden
van zo'n string:

    username:geheim@tcp(zardoz.service.rug.nl:3306)/paqu1

Dit logt in op de MySQL-server die draait op zardoz.service.rug.nl,
poort 3306, op de database 'paqu1', van de gebruiker 'username' met password
'geheim'.

    username:geheim@/paqu1

Dit logt in op de default server op de huidige machine, op de database
'paqu1', van de gebruiker 'username' met password 'geheim'.

----

Het tweede argument na het commando pqbuild1 is een id-string voor deze
dataset. Dit moet een woord van alleen kleine letters a tot z zijn.

Het programma maakt in de database tabellen aan met deze id-string als
prefix. Verder maakt het programma tabellen aan met de prefix 'wordrel'.
Gebruik deze niet als id-string.

----

Het laatste argument is de titel van deze dataset. Die titel zie je
terug in het menu van de interface.

