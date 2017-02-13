let mongoose = require('mongoose')
mongoose.Promise = global.Promise;

var Mod = mongoose.Schema({
    name: String,
    version: String,
    source: String
})

var mod = mongoose.model('Mod', Mod);
module.exports = mod
