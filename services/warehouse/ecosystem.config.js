// PM2
module.exports = {
    apps: [{
        name: "warehouse",
        script: "./bin/warehouse",
        instances: '1'
    }]
}