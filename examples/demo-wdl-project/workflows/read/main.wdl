version 1.0
workflow ReadFile {
    input {
        File input_file
    }
    call read_file { input: input_file = input_file }
}

task read_file {
    input {
        File input_file
    }
    String content = read_string(input_file) 

    command {
        echo '~{content}'
    }
    runtime {
        docker: "ubuntu:latest"
        memory: "4G"
        #cromwell retry
        awsBatchRetryAttempts: 3
        #miniwdl retry
        maxRetries: 3
    }

    output { String out = read_string( stdout() ) }
}

