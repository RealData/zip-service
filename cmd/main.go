package main 

import (
	"flag"
	"fmt" 
	"os"
	"runtime"
    r "zip-service/internal/readdir" 
	ce "zip-service/internal/compressextract" 
	"zip-service/internal/filelist"
)

const (
	TOP = 10 
	METHOD = "gzip" 
	PATTERN = "*" 
	FILE = ""
	DIR = "." 
	THREADS = 1 
)

var (
	top int 
	method string 
	pattern string 
	file string 
	dir string 
	compress bool  
	extract bool 
	threads int 
)

func init() { 

	flag.IntVar(&top, "top", TOP, "Select top files to compress/extract") 
	flag.IntVar(&top, "t", TOP, "Select top files to compress/extract") 
	flag.StringVar(&method, "method", METHOD, "Select compression method from gzip, flate, lzw, zlib") 
	flag.StringVar(&method, "m", METHOD, "Select compression method from gzip, flate, lzw, zlib") 
	flag.StringVar(&pattern, "pattern", PATTERN, "Use pattern to filter files") 
	flag.StringVar(&pattern, "p", PATTERN, "Use pattern to filter files") 
	flag.StringVar(&file, "file", FILE, "ZIP file to compress to or extract from") 
	flag.StringVar(&file, "f", FILE, "ZIP file to compress to or extract from") 
	flag.StringVar(&dir, "dir", DIR, "Directory to compress files from or extract files to") 
	flag.StringVar(&dir, "d", DIR, "Directory to compress files from or extract files to") 
	flag.BoolVar(&compress, "compress", false, "Perform compression") 
	flag.BoolVar(&compress, "c", false, "Perform compression") 
	flag.BoolVar(&extract, "extract", false, "Perform extraction") 
	flag.BoolVar(&extract, "e", false, "Perform extraction") 
	flag.IntVar(&threads, "nthreads", THREADS, "Number of threads") 
	flag.IntVar(&threads, "n", THREADS, "Number of threads")

	flag.Parse() 

	if compress == extract {
		fmt.Println("Only one of compress/extract flags can be set") 
		os.Exit(1)
	}

	if len(file) == 0 {
		fmt.Println("Archive file should be provided") 
		os.Exit(1)
	}

	if len(dir) == 0 {
		fmt.Println("Directory to compress/extract should be provided") 
		os.Exit(1)
	} 

	if threads < 1 {
		fmt.Println("Number of threads should be >= 1") 
		os.Exit(1)
	}
}

func main() {

	var compressor ce.Compressor  
	var extractor ce.Extractor  
	switch method {
	case "gzip": 
		compressor = ce.CompressFileGZIP 
		extractor = ce.ExtractFileGZIP 
	case "flate": 
		compressor = ce.CompressFileFlate 
		extractor = ce.ExtractFileFlate 
	case "lzw": 
		compressor = ce.CompressFileLZW 
		extractor = ce.ExtractFileLZW 
	case "zlib": 
		compressor = ce.CompressFileZLIB 
		extractor = ce.ExtractFileZLIB 
	default: 
		fmt.Printf("Method %s is not implemented", method) 
		os.Exit(1)
	}

	if numCPU := runtime.NumCPU(); threads > numCPU {
		threads = numCPU 
	}

	if compress { 
		fmt.Printf("Compressing %d top files filtered with %s pattern from %s directory to %s file using %s method and %d threads \n", top, pattern, dir, file, method, threads)
		files, err := r.ReadDirTopFiles(dir, pattern, top, threads) 
		if err != nil { 
			fmt.Println(err)
			os.Exit(1)
		}
		ce.CompressAndZipFiles(file, dir, files, threads, compressor) 
	}
	
	if extract { 
		fmt.Printf("Extracting %d top files filtered with %s pattern from %s file to %s directory using %s method and %d threads \n", top, pattern, file, dir, method, threads)
		files, err := ce.ExtractFileInfo(file) 
		if err != nil {
			fmt.Println(err) 
			os.Exit(1)
		}
		if pattern != "" {
			files, err = filelist.FilterFiles(files, pattern) 
		}
		if err != nil { 
			fmt.Println(err)
			os.Exit(1)
		}
		files = filelist.ParallelFindTop(files, top, threads) 
		err = ce.UnzipAndExtractFiles(dir, file, files, threads, extractor) 
		if err != nil { 
			fmt.Println(err)
			os.Exit(1)
		}
	}

}
