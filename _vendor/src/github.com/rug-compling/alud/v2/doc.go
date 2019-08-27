/*
Package alud derives Universal Dependencies from sentences parsed with Alpino.

Usually, the input is XML in the alpino_ds format.

The output is in the CoNLL-U format, or the Universal Dependencies can be embedded
into the alpino_ds format (version 1.10), making them available for XPath queries.

It is also possible to embed a user provided file in the CoNLL-U format, and embed this
into the alpino_ds format.

When empty heads are reconstructed (resulting in lines with an ID with a dot), the ID
of the original line is added in the last field of the CoNLL-U format, in the
form CopiedFrom=ID. This information is necessary for correct embedding into the
alpino_ds format.

----

The package is based on a translation of an xquery script written by Gosse Bouma.

See Alpino: https://www.let.rug.nl/vannoord/alp/Alpino/

See Universal Dependencies: https://universaldependencies.org/

See CoNLL-U: https://universaldependencies.org/format.html

See xquery script: https://github.com/gossebouma/lassy2ud
*/
package alud
