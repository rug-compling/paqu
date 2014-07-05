
    TODO:

     * Geen curl gebruiken, maar Go-functies

export ALPINO_HOME=/my/opt/Alpino


    Lopende tekst omzetten in genummerde regels:

../../bin/pqtexter -r testdoc-windows.txt | $ALPINO_HOME/Tokenization/tokenize.sh | perl -p -e 's/^/sprintf("%08d|",++$n)/e' > tmp

    Voor tekst die al uit één zin per regel bestaat:

../../bin/pqtexter lines.txt | $ALPINO_HOME/Tokenization/tokenize_nobreak.sh | perl -p -e 's/^/sprintf("%08d|",++$n)/e' > tmp


    Parsen, met server:

../../bin/pqalpino -a $ALPINO_HOME -s http://145.100.57.99/bin/alpino-notk tmp

    Parsen, zonder server:

../../bin/pqalpino -a $ALPINO_HOME tmp

