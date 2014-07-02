# PaQu #

## Te doen ##

Zie ook: TODO in diverse bestanden

  - Alpino-server
    - Aanpassen aan API van de server (huidige server is te oud)
      - zie ook TODOs in work.go
  - Gebruikershandleiding
    - Introtekst op hoofdpagina
    - Tekst achter Meer info...
  - Interface
    - https zonder http? de combinatie maakt het erg ingewikkeld
    - Logo + balk wijzigen
    - Pop-up help, waar nodig
    - Beheer van corpora
      - uploaden:
        - menu, met tussenkopjes: wel of geen default?
        - detectie van soort data
          - ZIP: 50 4B 03 04
          - DbXML: zie onder
          - tijdens verwerking, in veld "Opmerkingen" soort vermelden
        - xml-bestanden in zipbestand
        - regels
          - met/zonder label
          - wel/niet getokeniseerd
            - met/zonder metanotatie (geen automatische herkenning)
      - Voordat er iets gedownload wordt controleren of er nog
        onge'gzip'te bestanden zijn
  - Code
    - Organisatie + documentatie
    - Ontwikkeling: `make`, `go fmt`, `go vet`, `golint`

  - pqbuild
    - invoer: *.xml.gz

  - alle programma's zonder pq hernoemen

  - niet meer dan 10.000 bestanden in een directory

```
# Berkeley DB
#
# Ian Darwin's file /etc/magic files: big/little-endian version.
#
# Hash 1.85/1.86 databases store metadata in network byte order.
# Btree 1.85/1.86 databases store the metadata in host byte order.
# Hash and Btree 2.X and later databases store the metadata in host byte order.

0   long    0x00061561  Berkeley DB
>8  belong  4321
>>4 belong  >2      1.86
>>4 belong  <3      1.85
>>4 belong  >0      (Hash, version %d, native byte-order)
>8  belong  1234
>>4 belong  >2      1.86
>>4 belong  <3      1.85
>>4 belong  >0      (Hash, version %d, little-endian)

0   belong  0x00061561  Berkeley DB
>8  belong  4321
>>4 belong  >2      1.86
>>4 belong  <3      1.85
>>4 belong  >0      (Hash, version %d, big-endian)
>8  belong  1234
>>4 belong  >2      1.86
>>4 belong  <3      1.85
>>4 belong  >0      (Hash, version %d, native byte-order)

0   long    0x00053162  Berkeley DB 1.85/1.86
>4  long    >0      (Btree, version %d, native byte-order)
0   belong  0x00053162  Berkeley DB 1.85/1.86
>4  belong  >0      (Btree, version %d, big-endian)
0   lelong  0x00053162  Berkeley DB 1.85/1.86
>4  lelong  >0      (Btree, version %d, little-endian)

12  long    0x00061561  Berkeley DB
>16 long    >0      (Hash, version %d, native byte-order)
12  belong  0x00061561  Berkeley DB
>16 belong  >0      (Hash, version %d, big-endian)
12  lelong  0x00061561  Berkeley DB
>16 lelong  >0      (Hash, version %d, little-endian)

12  long    0x00053162  Berkeley DB
>16 long    >0      (Btree, version %d, native byte-order)
12  belong  0x00053162  Berkeley DB
>16 belong  >0      (Btree, version %d, big-endian)
12  lelong  0x00053162  Berkeley DB
>16 lelong  >0      (Btree, version %d, little-endian)

12  long    0x00042253  Berkeley DB
>16 long    >0      (Queue, version %d, native byte-order)
12  belong  0x00042253  Berkeley DB
>16 belong  >0      (Queue, version %d, big-endian)
12  lelong  0x00042253  Berkeley DB
>16 lelong  >0      (Queue, version %d, little-endian)
```
