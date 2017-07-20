package gflow

// NewDag is fuction to generate dag instance from task list
func NewDag(tasks []*Task) Dag {
	m := map[*Task]*vertex{}
	d := Dag{graph: make([]vertex, len(tasks))}
	for i, t := range tasks {
		d.graph[i] = vertex{
			value: t,
			suc:   []*vertex{},
		}
		m[tasks[i]] = &d.graph[i]
	}
	for t, v := range m {
		for _, dt := range t.deps {
			m[dt].suc = append(m[dt].suc, v)
		}
	}
	return d
}

// Dag is struct contains vertex list
type Dag struct {
	graph  []vertex
	index  int
	sorted []*Task
}

// Vertex is struct contains successor list
type vertex struct {
	value   *Task
	suc     []*vertex
	visited bool
}

// Tsort is function to topological sort for DAG
func (d Dag) Tsort() []*Task {
	d.index = len(d.graph) - 1
	d.sorted = make([]*Task, len(d.graph))
	for i := 0; i < len(d.graph); i++ {
		v := &d.graph[i]
		d.traceSuccessor(v)
	}
	return d.sorted
}

func (d *Dag) traceSuccessor(v *vertex) {
	if v.visited {
		return
	}
	for j := 0; j < len(v.suc); j++ {
		d.traceSuccessor(v.suc[j])
	}
	d.visit(v)
}

func (d *Dag) visit(v *vertex) {
	v.visited = true
	d.sorted[d.index] = v.value
	d.index--
}
