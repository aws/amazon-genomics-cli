version 1.0

workflow myWorkflow {
    call myTaskA
    call myTaskB
}

task myTaskA {
    input {
        String addressee
    }

    command {
        echo "hello world ${addressee} from myTaskA"
    }
    output {
        String out = read_string(stdout())
    }
}

task myTaskB {
    input {
        String addressee
    }

    command {
        echo "hello world ${addressee} from myTaskB"
    }
    output {
        String out = read_string(stdout())
    }
}