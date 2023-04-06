#!/usr/bin/env python3

import sys
from rdflib import Graph

g = Graph()
g.parse(sys.argv[1])
g.print(out=sys.stdout)
