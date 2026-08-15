# rimpact-graph

`rimpact` is the pure core of a selective-build planner. Given a build graph and
the set of packages whose sources changed, it computes the **rebuild-impact
set**: which other packages CI must rebuild.

Edges carry a kind: a `build` edge means the dependent compiles/links against
the dependency; a `runtime` edge means the dependent only loads or reads the
dependency at run time.

```go
impacted := rimpact.ImpactedByChange(pkgs, []string{"base"})
```
