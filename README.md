# To build with Docker
   sudo docker build -t gymapp .
   sudo docker run -p 3000 --net="host" --name name --rm gymapp -e RDS_DB_NAME='gymapp1' -e RDS_DB_USERNAME='rugile' -r RDS_DB_PASSWORD='maironis' -e RDS_HOSTNAME='127.0.0.1'
   sudo docker stop name

To get the IP

   sudo docker inspect --format '{{ .NetworkSettings.IPAddress }}' <<name>>
# To deploy to Elastic Beanstalk
   eb deploy
   eb status
