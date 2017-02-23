var express = require('express');
var Mod = require('../schemas/mod')
var User = require('../schemas/user')
var router = express.Router();
var jwt = require('jsonwebtoken')
var secret = require('../config.json').secret


router.get('/mod', function(req, res) {
    if (!req.query.name) {
        res.status(400)
        res.send('Bad Request')
        return false;
    }
    Mod.findOne({
        name: req.query.name
    }, (err, mod) => {
        if (mod) {
            res.json({
                name: mod.name,
                version: mod.version,
                source: mod.source

            })
        } else {
            res.json({
                none: true
            })
        }
    })
});

router.post('/login', (req, res) => {
    User.findOne({
        username: req.body.username
    }, (err, user) => {
        console.log(user.password)
        if (err) {
            console.log(err)
        } else if (req.body.password == user.password) {
            var token = jwt.sign({
                username: user.username,
                password: user.password
            }, secret, {
                expiresIn: "1y"
            });

            // return the information including token as JSON
            res.json({
                success: true,
                message: 'Enjoy your token!',
                token: token
            });
        } else {
            res.json({
                success: false,
                message: "Wrong password"
            })

        }
    })
})


router.post('/create', (req, res) => {
    var newUser = new User({
        username: req.body.username,
        password: req.body.password
    })

    newUser.save((err) => {
        if (err) {
            res.json({
                success: false,
                message: "A user with username " + req.body.username + " already exists"
            })
        } else {

            res.json({
                success: true
            })
        }
    })
})

router.use((req, res, next) => {
    if (!req.user) {
        res.json({
            success: false,
            message: 'Authentication failed. User not found.'
        });
    } else if (user) {
        return next
    }

    var token = jwt.sign(user, 'secret', {
        expiresInMinutes: 1440 // expires in 24 hours
    });

    // return the information including token as JSON
    res.json({
        success: true,
        message: 'Enjoy your token!',
        token: token
    });


});

router.post('/mod', function(req, res) {
    if (!(req.query.name && req.query.version && req.query.source)) {
        res.status(400)
        res.send('Bad Request')
        return false;
    }
    Mod.findOne({
        name: req.query.name.toLowerCase()
    }, (err, mod) => {
        if (mod) {
            mod.name = req.query.name.toLowerCase();
            mod.source = req.query.source;
            mod.version = req.query.version
            mod.save()
            res.send(req.query.name + ' updated');

        } else {
            var module = new Mod({
                name: req.query.name.toLowerCase(),
                version: req.query.version,
                source: req.query.source
            })
            module.save()
            res.send(req.query.name.toLowerCase() + ' created')

        }


    })

});

module.exports = router;
