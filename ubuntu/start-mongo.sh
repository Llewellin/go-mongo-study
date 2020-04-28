export HOME=/home/vagrant
echo $HOME

cat >> $HOME/mongod.conf <<EOF
# mongod.conf

# for documentation of all options, see:
#   http://docs.mongodb.org/manual/reference/configuration-options/

# Where and how to store data.
storage:
  dbPath: /var/lib/mongodb
  journal:
    enabled: true
#  engine:
#  mmapv1:
#  wiredTiger:

# where to write logging data.
systemLog:
  destination: file
  logAppend: true
  path: /var/log/mongodb/mongod.log

# network interfaces
net:
  port: 27017
#  bindIp: 127.0.0.1
  bindIp: 44.44.44.11


# how the process runs
processManagement:
  timeZoneInfo: /usr/share/zoneinfo

#security:

#operationProfiling:

replication:
  replSetName: rs0
#sharding:

## Enterprise-Only Options:

#auditLog:

#snmp:
EOF

cat >> /etc/systemd/system/mongod.service <<EOF
[Unit]
Description=mongod
Documentation=https://github.com/coreos

[Service]
ExecStart=/usr/bin/mongod \\
  --config=/home/vagrant/mongod.conf
EOF

sudo systemctl start mongod
sudo systemctl status mongod
sudo systemctl enable mongod