# To build with Docker
   sudo docker build -t gymapp .
   sudo docker run --publish 3000:8080 --name name --rm gymapp

# To deploy to Elastic Beanstalk
   eb deploy
   eb status
