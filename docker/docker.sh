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
      { _id : 0, host : "localhost:27017" },
    ]
  }
)


docker run --name mongo1 -p 27018:27017  -v /Users/tingsheng.lee/Desktop/projects/go-mongo/concurrency/mongo/db1:/data/db -d  docker.io/mongo  --replSet "rs0"
docker run --name mongo2 -p 27019:27017  -v /Users/tingsheng.lee/Desktop/projects/go-mongo/concurrency/mongo/db2:/data/db -d  docker.io/mongo  --replSet "rs0"
docker run --name mongo3 -p 27020:27017  -v /Users/tingsheng.lee/Desktop/projects/go-mongo/concurrency/mongo/db3:/data/db -d  docker.io/mongo --replSet "rs0"

rs.initiate({ _id : "rs0", members: [ { _id: 0, host: "172.17.0.2:27017" }, { _id: 1, host: "172.17.0.3:27017" },{_id: 2, host: "172.17.0.4:27017" }]})
