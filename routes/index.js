var express = require('express');
var Mod = require('../schemas/mod')
var router = express.Router();

/* GET home page. */
router.post('/mod', function(req, res) {
    if (!(req.query.name && req.query.version && req.query.source)) {
        res.status(400)
        res.send('Bad Request')
        return false;
    }
    Mod.findOne({ name: req.query.name.toLowerCase()}, (err, mod) => {
        if (mod){
            mod.name = req.query.name.toLowerCase();
            mod.source = req.query.source;
            mod.version = req.query.version
            mod.save()
            res.send(req.query.name  + ' updated');

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

router.get('/mod', function(req, res) {
    if (!req.query.name) {
        res.status(400)
        res.send('Bad Request')
        return false;
    }
    Mod.findOne({ name: req.query.name}, (err, mod) => {
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

module.exports = router;
