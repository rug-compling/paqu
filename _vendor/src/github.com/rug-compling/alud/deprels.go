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
		if node.udHeadPosition == 0 && node.udRelation != "root" {
			panic(fmt.Sprintf(
				"Invalid HEAD:DEPREL combination 0:%s for %s:%s",
				node.udRelation,
				number(node.End),
				node.Word))
		}
	}
}
