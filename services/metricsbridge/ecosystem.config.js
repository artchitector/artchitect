// PM2
module.exports = {
    apps: [{
        name: "metricsbridge",
        script: "./bin/metricsbridge",
        instances: '1'
    }]
}