package config

import "crypto/sha1"

// TODO config file with chunk size, checksum size, streamstore prefix, download/store paths

// ChecksumSize is the size of checksum block
const ChecksumSize = sha1.Size

// ChecksumStringSize is the character count of checksums as strings
const ChecksumStringSize = ChecksumSize * 2

// ChunkSize is the size of content chunk block
const ChunkSize = 128 << 10

// StoreFolderName is the name of the folder that houses stored
const StoreFolderName = "store"

// DownloadFolderName is the name of the folder that houses in-progress downloads
const DownloadFolderName = "download"

// TempDownloadPrefix is the prefix of temp files stored in the download folder
const TempDownloadPrefix = "streamstore"
