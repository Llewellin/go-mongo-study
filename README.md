docker setup failed, don't use

not sure wht setup-rs.sh is not working.

run the following to install manually.
```
vagrant ssh mongo-primary-1
mongo --host 44.44.44.11

rs.initiate({ _id : "rs0", members: [ { _id: 0, host: "mongo1:27017", priority: 5 }, { _id: 1, host: "mongo2:27017" },{_id: 2, host: "mongo3:27017" }]})
```

test with vegeta
vegeta attack -targets="mongo.txt" -rate=1 -duration=1s  | vegeta report
vegeta attack -targets="mongo2.txt" -rate=1 -duration=1s  | vegeta report
