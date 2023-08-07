// PM2
module.exports = {
    apps: [{
        name: "asgard",
        script: "./bin/asgard",
        instances: '1'
    }]
}