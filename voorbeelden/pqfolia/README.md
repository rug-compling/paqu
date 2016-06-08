# pqfolia #

Met `pqfolia` kun je een corpus in FoLiA-formaat (zie:
https://proycon.github.io/folia/) met interne of externe metadata
omzetten naar het formaat dat je rechtstreeks in PaQu kunt invoeren.
PaQu, de webapp (`pqserve`), gebruikt dit programma, maar je kunt het ook
offline gebruiken. Dat is vooral handig als je een grote hoeveelheid
data hebt, want de hoeveelheid data die je moet uploaden naar PaQu kun
je zo fors reduceren.

Bovendien kun je met `pqfolia` iets wat online niet kan: je kunt
metadata omzetten naar een ander formaat. Je kunt bijvoorbeeld een datum
als `May 12, 2016` omzetten naar `2016-05-12`. De eerste vorm kan door
PaQu niet verwerkt worden als datum, de tweede wel.

Je gebruikt `pqfolia` met een configuratiebestand. Zie het voorbeeld in
`config.toml`, ook voor toelichting bij de opties.
