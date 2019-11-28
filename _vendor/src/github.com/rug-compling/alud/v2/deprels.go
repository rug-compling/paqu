package alud

import (
	"fmt"
)

func addDependencyRelations(q *context) {

	var node *nodeType

	defer func() {
		if r := recover(); r != nil {
			panic(trace(r, "addDependencyRelations", q, node))
		}
	}()

	for _, node = range q.ptnodes {
		q.depth = 0
		node.udRelation = dependencyLabel(node, q)
		q.depth = 0
		node.udHeadPosition = externalHeadPosition(list(node), q)
		if node.udHeadPosition == 0 && node.udRelation != "root" ||
			node.udHeadPosition != 0 && node.udRelation == "root" {
			panic(fmt.Sprintf(
				"Invalid HEAD:DEPREL combination %s:%s",
				number(node.udHeadPosition),
				node.udRelation))
		}
		if node.udHeadPosition == node.End {
			panic(fmt.Sprintf(
				"DEPREL to self %s:%s",
				number(node.udHeadPosition),
				node.udRelation))
		}
	}
}
