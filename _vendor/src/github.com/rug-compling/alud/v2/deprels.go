package alud

import (
	"fmt"
)

func addDependencyRelations(q *context) {
	for _, node := range q.ptnodes {
		q.depth = 0
		node.udRelation = dependencyLabel(node, q, []trace{})
		q.depth = 0
		node.udHeadPosition = externalHeadPosition(list(node), q, []trace{})
		if node.udHeadPosition == 0 && node.udRelation != "root" ||
			node.udHeadPosition != 0 && node.udRelation == "root" {
			panic(fmt.Sprintf(
				"Invalid HEAD:DEPREL combination %s:%s for %s:%s",
				number(node.udHeadPosition),
				node.udRelation,
				number(node.End),
				node.Word))
		}
	}
}
