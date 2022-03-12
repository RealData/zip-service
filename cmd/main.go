package main 

import (
	"flag"
	"fmt" 
	"os"
    r "zip-service/internal/readdir" 
	ce "zip-service/internal/compressextract" 
)

const (
	TOP = 10 
	METHOD = "gzip" 
	PATTERN = "" 
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
	flag.StringVar(&method, "method", METHOD, "Select compression method") 
	flag.StringVar(&method, "m", METHOD, "Select compression method") 
	flag.StringVar(&pattern, "pattern", PATTERN, "Use pattern to filter files") 
	flag.StringVar(&pattern, "p", PATTERN, "Use pattern to filter files") 
	flag.StringVar(&file, "file", FILE, "Archive file") 
	flag.StringVar(&file, "f", FILE, "Archive file") 
	flag.StringVar(&dir, "dir", DIR, "Directory to compress/extract files") 
	flag.StringVar(&dir, "d", DIR, "Directory to compress/extract files") 
	flag.BoolVar(&compress, "compress", false, "Compress") 
	flag.BoolVar(&compress, "c", false, "Compress") 
	flag.BoolVar(&extract, "extract", false, "Extract") 
	flag.BoolVar(&extract, "e", false, "Extract") 
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


	if  compress { 
		files, err := r.ReadDirTopFiles(dir, pattern, top, threads) 
		if err != nil { 
			os.Exit(1)
		}
		ce.CompressAndZipFiles(file, dir, files, threads, compressor) 
	}
	
	if extract {
		ce.UnzipAndExtractFiles(dir, file, threads, extractor) 
	}

}
