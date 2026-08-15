// Package rimpact computes the rebuild-impact set for a build graph: given the
// set of packages whose sources changed, it decides which OTHER packages CI
// must rebuild. It is the pure core a selective-build planner uses to avoid
// rebuilding the entire monorepo on every change.
package rimpact

// EdgeKind is how one package depends on another.
type EdgeKind string

const (
	// BuildEdge means the dependent compiles/links against the dependency, so
	// a change to the dependency can change the dependent's build output.
	BuildEdge EdgeKind = "build"
	// RuntimeEdge means the dependent only loads/reads the dependency at run
	// time (a data file, a plugin loaded by name, a tool invoked as a
	// subprocess); the dependent's own build output does not embed it.
	RuntimeEdge EdgeKind = "runtime"
)

// Dep is one outgoing dependency edge of a package.
type Dep struct {
	On   string   // the name of the package depended upon
	Kind EdgeKind // how it is depended upon
}

// Package is a node in the build graph.
type Package struct {
	Name string
	Deps []Dep // this package's outgoing dependencies
}

type back struct {
	name string
	kind EdgeKind
}

// ImpactedByChange returns the set of package names that must rebuild when the
// packages named in `changed` have their sources modified. The changed
// packages themselves are always included.
func ImpactedByChange(pkgs []Package, changed []string) map[string]bool {
	// Build the reverse adjacency: for each package X, who depends on X, and
	// by what kind of edge.
	dependents := make(map[string][]back)
	for _, p := range pkgs {
		for _, d := range p.Deps {
			dependents[d.On] = append(dependents[d.On], back{name: p.Name, kind: d.Kind})
		}
	}

	impacted := make(map[string]bool, len(changed))
	queue := make([]string, 0, len(changed))
	for _, c := range changed {
		impacted[c] = true
		queue = append(queue, c)
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, b := range dependents[cur] {
			if b.kind == RuntimeEdge {
				if !impacted[b.name] {
					impacted[b.name] = true
					queue = append(queue, b.name)
				}
				continue
			}
			impacted[b.name] = true
		}
	}

	return impacted
}
