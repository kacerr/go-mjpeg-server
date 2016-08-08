## Synopsis
Server that allow mjpeg video streaming over HTTP protocol. 
There might be multiple sources of data at some point, but for start it would only stream from files (.jpg) prepared in the directory

## Feature list
- source directory can be set at server startup from command line
- server:port/mjpeg streams the data from source directory
- server:port/devbug/vars shows application statistics (number of currently running streams included)



## Nice to have features
- allow streaming from multiple sources
- allow streaming from multiple source types
- let the user specify framerate