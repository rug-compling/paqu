package alud

func addDependencyRelations(q *context) {
	for _, node := range q.ptnodes {
		q.depth = 0
		node.udRelation = dependencyLabel(node, q)
		q.depth = 0
		node.udHeadPosition = externalHeadPosition(list(node), q)
	}
}
