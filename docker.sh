docker exec -it localmongo1 mongo

rs.initiate(
  {
    _id : 'rs0',
    members: [
      { _id : 0, host : "mongo1:27017" },
      { _id : 1, host : "mongo2:27017" },
      { _id : 2, host : "mongo3:27017", arbiterOnly: true }
    ]
  }
)

exit


rs.initiate(
  {
    _id : 'rs0',
    members: [
      { _id : 0, host : "127.0.0.1:27017" },
    ]
  }
)