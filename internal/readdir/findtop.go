package readdir 

import ( 
	"math"
	"sort" 
	"sync"
)

// findTop finds top files by their size in a list of FileInfo. 
// Two algorithms are implemented: linear search of top files and simple sort with successive selection of top files. 
// Linear search is asymptotically more effective when `top < log(L) - log(threads)`, although experiments are needed to find the coefficients for this condition.  
// This condition is approximately equivalent to `top < log(L)`
func findTop(files []FileInfo, top int) []FileInfo { 

	l := len(files)
	if l < top {
		top = l  
	}

	if top < 1 {
		return []FileInfo{}
	}

	// Use sort algorithm for small file lists or large `top` parameter values 
	if float64(top) > math.Log(float64(l)) {  
		sort.Slice(files, func(i int, j int) bool { return files[i].Size > files[j].Size }) 
		return files[:top]
	}

	// Otherwise use linear search for top files 
	for i := 0; i < top; i++ {
		maxj := i   
		maxv := files[maxj].Size 
		for j := i; j < len(files); j++ {
			if size := files[j].Size; size > maxv {
				maxj = j 
				maxv = size 
			}
		}
		files[i], files[maxj] = files[maxj], files[i]
	}

	return files[:top]

}

// parallelFindTop finds top fies by their size in a list of FileInfo. It splits the list into `threads` chunks and concurrently calls `findTop` for each of the chunks, the top files from the chunks are then merged, sorted, and `top` files are selected 
func parallelFindTop(files []FileInfo, top int, threads int) []FileInfo {

	if top < 1 {
		return []FileInfo{} 
	}

	var wg sync.WaitGroup 
	l := len(files) 

	if threads > l {
		threads = l 
	} 

	if top > l {
		top = l 
	}

	res := make([]FileInfo, 0, threads * top)
	resChan := make(chan []FileInfo, threads) 

	for i := 0; i < threads; i++ {
		wg.Add(1) 
		go func(files []FileInfo, top int) { 
			defer wg.Done() 
			resChan<- findTop(files, top)   
		}(files[i*l/threads:(i+1)*l/threads], top)
	} 

	go func() {
		wg.Wait() 
		close(resChan) 
	}() 

	for s := range resChan {
		res = append(res, s...)
	}

	sort.Slice(res, func(i int, j int) bool { return res[i].Size > res[j].Size }) 

	return res[:top]  

} 

