package friff

import (
	"fmt"
	"os"
	"sync"
)

var moveMu sync.Mutex

type Part struct {
	srcPath  string
	destPath string
	src      map[uint]*Shadow
	dest     map[uint]*Shadow
}

type Diff struct {
	SrcPath    string
	DestPath   string
	Deletions  []*Shadow
	Insertions []*Shadow
	Original   map[uint]*Shadow
}

func MergePaths(left, right string) *Part {
	var wg sync.WaitGroup

	var lShad, rShad map[uint]*Shadow
	wg.Add(2)

	chunkifyRoutine := func(blobAt string, deposit *map[uint]*Shadow, wgg *sync.WaitGroup) {
		defer wgg.Done()
		resShad, err := Chunkify(blobAt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "chunkifying %s got %v\n", blobAt, err)
		} else {
			*deposit = resShad
		}
	}

	go chunkifyRoutine(left, &lShad, &wg)
	go chunkifyRoutine(right, &rShad, &wg)

	wg.Wait()

	return &Part{
		srcPath:  left,
		destPath: right,
		src:      lShad,
		dest:     rShad,
	}
}

func Reconstruct(from, to string) {
	// part := MergePaths(from, to)
}

func MergeShow(p *Part) string {
	// l, r := p.left, p.right
	// TODO: Join up missing chunks
	return ""
}

func move(key uint, from, to map[uint]*Shadow) {
	moveMu.Lock()
	defer moveMu.Unlock()

	retr, ok := from[key]
	if !ok {
		return
	}

	to[key] = retr
	delete(from, key)
}

func (pt *Part) Merge() *Diff {
	src := pt.src
	dest := pt.dest

	untouched := make(map[uint]*Shadow)

	var deletions, insertions []*Shadow
	for srcId, srcShad := range src {
		destShad, ok := dest[srcId]
		if !ok {
			deletions = append(deletions, srcShad)
			continue
		}
		if destShad.Checksum == srcShad.Checksum {
			move(srcId, dest, untouched) // delete(dest, srcId)
			continue
		}

		deletions = append(deletions, srcShad)
	}

	for _, destShad := range dest {
		insertions = append(insertions, destShad)
	}

	return &Diff{
		SrcPath:    pt.srcPath,
		DestPath:   pt.destPath,
		Deletions:  deletions,
		Insertions: insertions,
		Original:   untouched,
	}
}
