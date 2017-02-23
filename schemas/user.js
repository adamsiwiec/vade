let mongoose = require('mongoose')
mongoose.Promise = global.Promise;

var User = mongoose.Schema({
    username: { type: String, unique: true },
    password: String
})

var user = mongoose.model('User', User);
module.exports = user
