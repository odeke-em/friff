package merkle

type Part struct {
	src  map[uint]*Shadow
	dest map[uint]*Shadow
}

type Diff struct {
	deletions  chan *Shadow
	insertions chan *Shadow
}

func merge(pt *Part) *Diff {
	src := pt.src
	dest := pt.dest

	deletions := make(chan *Shadow)
	insertions := make(chan *Shadow)

	diff := Diff{
		deletions:  deletions,
		insertions: insertions,
	}

	for srcId, srcShad := range src {
		destShad, ok := dest[srcId]
		if !ok {
			deletions <- srcShad
			continue
		}

		if destShad.checksum == srcShad.checksum {
			delete(dest, srcId)
			continue
		}

		deletions <- srcShad
	}

	for _, destShad := range dest {
		insertions <- destShad
	}

	return &diff
}
