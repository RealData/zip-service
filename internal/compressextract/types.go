package compressextract

type Compressor func(destFilePath string, sourceFilePath string) error
type Extractor func(destFilePath string, sourceFilePath string) error
