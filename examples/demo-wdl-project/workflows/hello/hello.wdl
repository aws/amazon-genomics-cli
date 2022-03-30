version 1.0
workflow hello_agc {
    call hello {}
}
task hello {
    command { echo "Hello Amazon Genomics CLI!" }
    runtime {
        docker: "ubuntu:latest"
        #cromwell retry
        awsBatchRetryAttempts: 3
        #miniwdl retry
        maxRetries: 3
    }
    output { String out = read_string( stdout() ) }
}
