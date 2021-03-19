package config

// ResponseContinue tells the client that the response was understood and the operation can continue
const ResponseContinue = byte(1)

// ResponseSuccess tells the client that the response was successful
const ResponseSuccess = byte(2)

// ResponseInternalError tells the client an internal error occurred
const ResponseInternalError = byte(3)

// ResponseNotFound tells the client the entity could not be found
const ResponseNotFound = byte(4)

// ResponseInvalidOperation tells the client the operation requested does not exist
const ResponseInvalidOperation = byte(5)

// ResponseInvalidChecksum tells the client the checksum is invalid
const ResponseInvalidChecksum = byte(6)

// ResponseChecksumError tells the client a checksum validation failed
const ResponseChecksumError = byte(7)

// ResponseMalformedRequest tells the client the request was malformed
const ResponseMalformedRequest = byte(8)

// WriteErrorContent is a fixed value sent to the http body indicating internal server error during writes
var ResponseWriteError = []byte{60, 13, 24, 62, 62, 171, 74, 212, 189, 94, 66, 244, 236, 231, 230, 207}

// RequestUpload asks to upload an entity
const RequestUpload = byte(1)

// RequestDownload asks to download an entity
const RequestDownload = byte(2)

// RequestDelete asks to delete an entity
const RequestDelete = byte(3)
