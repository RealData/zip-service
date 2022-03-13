package main 

import ( 
	"errors"
	"flag"
	"fmt" 
	"os"
	"runtime"
    r "zip-service/internal/readdir" 
	ce "zip-service/internal/compressextract" 
	f "zip-service/internal/filelist"
)

const (
	TOP = 10 
	METHOD = "gzip" 
	PATTERN = "*" 
	FILE = ""
	DIR = "" 
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
	compressor ce.Compressor  
	extractor ce.Extractor  
)

var (
	ceFlagsSetErr = errors.New("Only one of compress/extract flags can be set") 
	ceFlagNotSetErr = errors.New("One of compress/extract flags should be set") 
	zipFileNotProvidedErr = errors.New("Archive file should be provided") 
	dirNotProvidedErr = errors.New("Directory to compress/extract should be provided") 
	numThreadsErr = errors.New("Number of thresds should be >=1") 
	methodNotImplementedErr = errors.New("Method is not implemented")  
)

func parseFlags(args []string) error {

	flags := flag.NewFlagSet(args[0], flag.ExitOnError) 
	
	flags.IntVar(&top, "top", TOP, "Select top files to compress/extract") 
	flags.IntVar(&top, "t", TOP, "Select top files to compress/extract") 
	flags.StringVar(&method, "method", METHOD, "Select compression method from gzip, flate, lzw, zlib") 
	flags.StringVar(&method, "m", METHOD, "Select compression method from gzip, flate, lzw, zlib") 
	flags.StringVar(&pattern, "pattern", PATTERN, "Use pattern to filter files") 
	flags.StringVar(&pattern, "p", PATTERN, "Use pattern to filter files") 
	flags.StringVar(&file, "file", FILE, "ZIP file to compress to or extract from") 
	flags.StringVar(&file, "f", FILE, "ZIP file to compress to or extract from") 
	flags.StringVar(&dir, "dir", DIR, "Directory to compress files from or extract files to") 
	flags.StringVar(&dir, "d", DIR, "Directory to compress files from or extract files to") 
	flags.BoolVar(&compress, "compress", false, "Perform compression") 
	flags.BoolVar(&compress, "c", false, "Perform compression") 
	flags.BoolVar(&extract, "extract", false, "Perform extraction") 
	flags.BoolVar(&extract, "e", false, "Perform extraction") 
	flags.IntVar(&threads, "nthreads", THREADS, "Number of threads") 
	flags.IntVar(&threads, "n", THREADS, "Number of threads")		

	err := flags.Parse(args[1:]) 
	if err != nil {
		return err 
	}

	if compress == extract {
		if compress { 
			return ceFlagsSetErr 
		} else {
			return ceFlagNotSetErr 
		}
	}

	if len(file) == 0 {
		return zipFileNotProvidedErr 
	}

	if len(dir) == 0 {
		return dirNotProvidedErr  
	} 

	if threads < 1 {
		return numThreadsErr  
	}

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
		return methodNotImplementedErr   
	}

	if numCPU := runtime.NumCPU(); threads > numCPU {
		threads = numCPU 
	}

	return nil 

}

func run() error {

	if compress { 
		fmt.Printf("Compressing %d top files filtered with %s pattern from %s directory to %s file using %s method and %d threads \n", top, pattern, dir, file, method, threads)
		files, err := r.ReadDirTopFiles(dir, pattern, top, threads) 
		if err != nil { 
			return err
		}
		ce.CompressAndZipFiles(file, dir, files, threads, compressor) 
	}
	
	if extract { 
		fmt.Printf("Extracting %d top files filtered with %s pattern from %s file to %s directory using %s method and %d threads \n", top, pattern, file, dir, method, threads)
		files, err := ce.ExtractFileInfo(file) 
		if err != nil {
			return err 
		}
		if pattern != "" {
			files, err = f.FilterFiles(files, pattern) 
		}
		if err != nil { 
			return err 
		}
		files = f.ParallelFindTop(files, top, threads) 
		err = ce.UnzipAndExtractFiles(dir, file, files, threads, extractor) 
		if err != nil { 
			return err 
		}
	}

	return nil 

}

func main() {

	err := parseFlags(os.Args) 
	if err != nil {
		fmt.Println(err) 
		os.Exit(1) 
	}
	
	err = run() 
	if err != nil {
		fmt.Println(err) 
		os.Exit(1) 
	}

}
