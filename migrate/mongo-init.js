

// db = db.getSiblingDB('admin').auth(
//     process.env.MONGO_INITDB_ROOT_USERNAME,
//     process.env.MONGO_INITDB_ROOT_PASSWORD
// );
db.auth('root', 'rootpasswd') 


db = db.getSiblingDB('instanse_mongo')
db = new Mongo().getDB("Person");

db.users.createIndex({ "email": 1 }, { unique: true })
db.users.insertMany(
    [
        {
            "email": "hugo@mail.ru",
            "name": "admin", 
            "password": "$2a$10$NUH4URXlMNPS2NdBfE1DP.DAwcnlSu.cPhN9iv/C6TATzWKzvKqhq", // "password"
            "subscription": "admin"
        },{
            "email": "fawn@mail.ru",
            "name": "kiyosaki", 
            "password": "$2a$10$IcJwWr4Yy3dnH.SPPg1sV.MIgui0DGFjbqsSJEdYDHYSKpqxXEFum", // "123"
            "subscription": "premium"
        },{
            "email": "apostate@mail.ru",
            "name": "torvald",
            "password": "$2a$10$F/50BfVQ5A4fu4G8UIrGourrxlG/ngQ6X.gKh2aHlgOUUBFPs5G4a", // "456"
            "subscription": "econom"
        },{
            "email": "vagabund@mail.ru",
            "name": "eva",
            "password": "$2a$10$qj6G0eAe4HrtNG3EP4Po1.3LB9z12tNctxU6v4NmuvUzgoNSuOLAq", // " " - пробел, 
            "subscription": null
        }
    ]
);


