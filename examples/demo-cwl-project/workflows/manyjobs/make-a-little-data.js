const crypto = require('crypto')
if (process.argv.length > 2) {
    let arg = process.argv[2]
    console.log("Argument is: " + arg)
    
    let hashes = [arg]
    
    for (let i = 10; i >= 0; i--) {
        console.log(i + " bottles of " + arg + " on the wall...")
        hasher = crypto.createHash('sha512')
        for (let h of hashes) {
            hasher.update(h)
        }
        hashes.push(hasher.digest('hex'))
    }
    
    console.log("After meditating on the nature of " + arg + ", it turns out to be " + hashes[hashes.length - 1])
    
} else {
    console.log("Didn't get an argument")
}
