package dep

// // General purpose topological sort, not specific to the application of
// // library dependencies.  Adapted from Wikipedia pseudo code, one main
// // difference here is that this function does not consume the input graph.
// // WP refers to incoming edges, but does not really need them fully represented.
// // A count of incoming edges, or the in-degree of each node is enough.  Also,
// // WP stops at cycle detection and doesn't output information about the cycle.
// // A little extra code at the end of this function recovers the cyclic nodes.
// func topSortKahn(g graph, in inDegree) (order, cyclic []string) {
// 	var L, S []string
// 	// rem for "remaining edges," this function makes a local copy of the
// 	// in-degrees and consumes that instead of consuming an input.
// 	rem := inDegree{}
// 	for n, d := range in {
// 		if d == 0 {
// 			// accumulate "set of all nodes with no incoming edges"
// 			S = append(S, n)
// 		} else {
// 			// initialize rem from in-degree
// 			rem[n] = d
// 		}
// 	}
// 	for len(S) > 0 {
// 		last := len(S) - 1 // "remove a node n from S"
// 		n := S[last]
// 		S = S[:last]
// 		L = append(L, n) // "add n to tail of L"
// 		for _, m := range g[n] {
// 			// WP pseudo code reads "for each node m..." but it means for each
// 			// node m *remaining in the graph.*  We consume rem rather than
// 			// the graph, so "remaining in the graph" for us means rem[m] > 0.
// 			if rem[m] > 0 {
// 				rem[m]--         // "remove edge from the graph"
// 				if rem[m] == 0 { // if "m has no other incoming edges"
// 					S = append(S, m) // "insert m into S"
// 				}
// 			}
// 		}
// 	}
// 	// "If graph has edges," for us means a value in rem is > 0.
// 	for c, in := range rem {
// 		if in > 0 {
// 			// recover cyclic nodes
// 			for _, nb := range g[c] {
// 				if rem[nb] > 0 {
// 					cyclic = append(cyclic, c)
// 					break
// 				}
// 			}
// 		}
// 	}
// 	if len(cyclic) > 0 {
// 		return nil, cyclic
// 	}
// 	return L, nil
// }
