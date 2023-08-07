// PM2
module.exports = {
    apps: [{
        name: "alfheimr",
        script: "./bin/alfheimr",
        instances: '1'
    }]
}