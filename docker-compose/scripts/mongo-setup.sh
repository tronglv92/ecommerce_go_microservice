#!/bin/bash 
sleep 15
m1=host.docker.internal
m2=host.docker.internal
m3=host.docker.internal
port=${PORT:-27017}
# https://agent-hunt.medium.com/mongodb-replica-set-with-docker-40f54d17e10f
# mongodb://root:password@localhost:27021,localhost:27022,localhost:27023/?replicaSet=rs1
# mongodb1=$(getent hosts mongo-1-1 | awk '{ print $1 }')
# mongodb2=$(getent hosts mongo-1-2 | awk '{ print $1 }')
# mongodb3=$(getent hosts mongo-1-3 | awk '{ print $1 }')
# https://kenanabbak.medium.com/how-to-deploy-mongodb-replica-set-on-docker-using-docker-compose-and-shell-script-a28b2fd5506b

# port=${PORT:-27017}
port1=27021
port2=27022
port3=27023



# port=${PORT:-27017}

echo "###### Waiting for ${m1} and ${port1} instance startup.."
until mongosh --host ${m1}:${port1} --eval 'quit(db.runCommand({ ping: 1 }).ok ? 0 : 2)' &>/dev/null; do
  printf '.'
  sleep 1
done
echo "###### Working ${m1} instance found, initiating user setup & initializing rs setup.."

echo "Started.."

echo setup.sh time now: `date +"%T" `
mongosh --host ${m1}:${port1} <<EOF

    var rootUser = 'root';
    var rootPassword = 'password';
    var admin = db.getSiblingDB('admin');
    admin.auth(rootUser, rootPassword);

   var cfg = {
        "_id": "rs1",
        "version": 1,
        "members": [
            {
                "_id": 0,
                "host": "${m1}:${port1}",
                "priority": 3
            },
            {
                "_id": 1,
                "host": "${m2}:${port2}",
                 "priority": 2
            },
            {
                "_id": 2,
                "host": "${m3}:${port3}",
                 "priority": 1
            }
        ]
    };
    rs.initiate(cfg, { force: true });
    rs.status();
EOF

# sleep 10

# mongosh --host ${mongodb1}:${port1}<<EOF
#    use admin;
#    admin = db.getSiblingDB("admin");
#    admin.createUser(
#      {
# 	user: "admin",
#         pwd: "password",
#         roles: [ { role: "root", db: "admin" } ]
#      });
#      db.getSiblingDB("admin").auth("admin", "password");
#      rs.status();
# EOF