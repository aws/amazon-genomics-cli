version 1.0
workflow hello_agc {
    call hello {}
}
task hello {
    command { echo "Hello Amazon Genomics CLI!" }
    runtime {
        docker: "ubuntu:latest"
    }
    output { String out = read_string( stdout() ) }
}
