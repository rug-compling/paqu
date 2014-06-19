
    TODO:

     * Geen curl gebruiken, maar Go-functies
     * 'make' zonder gebruik van 'any2any'
     * Testmateriaal verplaatsen of verwijderen


export ALPINO_HOME=/my/opt/Alpino


    Lopende tekst omzetten in genummerde regels:

./prepare -r testdoc-windows.txt | $ALPINO_HOME/Tokenization/tokenize.sh | perl -p -e 's/^/sprintf("%08d|",++$n)/e' > tmp

    Voor tekst die al uit één zin per regel bestaat:

./prepare lines.txt | $ALPINO_HOME/Tokenization/tokenize_nobreak.sh | perl -p -e 's/^/sprintf("%08d|",++$n)/e' > tmp


    Parsen, met server:

./alpino -a $ALPINO_HOME -s http://145.100.57.99/bin/alpino-notk tmp

    Parsen, zonder server:

./alpino -a $ALPINO_HOME tmp

