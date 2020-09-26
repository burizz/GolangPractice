# Learning GoLang Practice


## gocli

Parse Linux users

## gofte

#### Goals

Mux HTTP API endpoints
MongoDB persistence

1. Seed data; Input from JSON file/s - for now 1 per month(figure out a better way) - seed_data/
2. Updated in MongoDB; if exist, otherwise created
3. CRUD via API endpoints
4. Calculate %
5. Visualize based on - order / utilization / person

TODO:
- Update addToTeam()
- Update deleteFromTeam()
- Update updatePerson()

- Update getPerson(); find by any value - id;firstname;lastname

- DB config file - host/user/db/collection

- SeedDB

- API method for querying based on criteria
  - Range/Amount of hours
  - Month
  - % utilization

Much later
- Cli tool
  - autocomplete
- Frontend

## Go build & run
```
go build
 ~/go/bin/gofte.exe
```

## App tests

HTTP GET :
localhost:10000/team

HTTP POST :
localhost:10000/addToTeam?MemberID=1&Firstname=Boris&Lastname=Yakimov&Hours=108

HTTP PUT :

HTTP DELETE :
## MongoDB container  
```
docker run -d -p 27017-27019:27017-27019 --name mongodb mongo
```  

```
docker exec -it mongodb bash
```  

## Mongo manual stuff reference
```
show dbs
use db_name
show collections
```

```
db.people.find()
{ "_id" : ObjectId("5d70ca754548abb56c4c45f6"), "memberid" : "1", "firstname" : "Boris", "lastname" : "Yakimov", "hours" : 108, "month" : "August" }
{ "_id" : ObjectId("5d70ca844548abb56c4c45f7"), "memberid" : "2", "firstname" : "Vname2", "lastname" : "name2", "hours" : 157, "month" : "August" }
{ "_id" : ObjectId("5d70ca8a4548abb56c4c45f8"), "memberid" : "3", "firstname" : "Kname3", "lastname" : "name3", "hours" : 111, "month" : "August" }
{ "_id" : ObjectId("5d70ca8f4548abb56c4c45f9"), "memberid" : "4", "firstname" : "Yname4", "lastname" : "name4", "hours" : 116, "month" : "August" }
{ "_id" : ObjectId("5d70ca954548abb56c4c45fa"), "memberid" : "5", "firstname" : "Kaname5", "lastname" : "name5", "hours" : 162, "month" : "August" }
{ "_id" : ObjectId("5d70ca9d4548abb56c4c45fb"), "memberid" : "6", "firstname" : "Krname6", "lastname" : "name6", "month" : "August" }
```

```
db.people.remove( {"ID": "7"} )
WriteResult({ "nRemoved" : 1 })
```

```
db.people.findOne( {"memberid": "1"} )
{
        "_id" : ObjectId("5d70ca754548abb56c4c45f6"),
        "memberid" : "1",
        "firstname" : "Boris",
        "lastname" : "Yakimov",
        "hours" : 108,
        "month" : "August"
}
```
