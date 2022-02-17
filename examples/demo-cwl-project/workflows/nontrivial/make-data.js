if (process.argv.length > 2) {
    let arg = process.argv[2]
    console.log("Argument is: " + arg);
    
    for (let i = 1000; i >= 0; i--) {
        console.log(i + " bottles of " + arg + " on the wall...");
    }
    
} else {
    console.log("Didn't get an argument");
}
