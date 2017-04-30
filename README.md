# To build with Docker
    sudo docker build -t gymapp .
    sudo docker run -p 3000 --net="host" --name name --rm gymapp -e RDS_DB_NAME='gymapp1' -e RDS_DB_USERNAME='rugile' -r RDS_DB_PASSWORD='maironis' -e RDS_HOSTNAME='127.0.0.1'
    sudo docker stop name

To get the IP

    sudo docker inspect --format '{{ .NetworkSettings.IPAddress }}' <<name>>
# To deploy to Elastic Beanstalk
    eb deploy
    eb status

## API documentation
# <<Tokens>>
   **Request a token**
   ----

   * **`POST` /token/request**
   * **Form Params**

     **Required:**
     <_User email and password are required parameters._>
      `email=[string]`
      `password=[string]`

   * **Success Response:**

     <_Returns a token valid for 2 weeks and 200 response._>

     * **Code:** 200
       **Content:**
       ```json
       {
        "error": null,
        "expiration": "2017-05-14T13:03:07+01:00",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InJ1Z2lsZW5hQGdtYWlsLmNvbSIsImV4cCI6IjIwMTctMDUtMTRUMTM6MDM6MDcrMDE6MDAifQ.utL_J900CAIif_FM7UcJiJomRZm_R5VSG8hj6aNfgOA"
        }
       ```

   * **Error Response:**

     * **Code:** 404
       **Content:**
       ```json
       {
         "error": "Invalid email or password."
        }
       ```

   **Renew a token**
   ----

   * **`POST` /token/renew**
   * **Form Params**

     **Required:**
     <_Only required paramenter - the active token which we want to renew._>
      `token=[string]`

   * **Success Response:**

     <_Returns same token valid for 2 weeks from now and a 200 response._>

     * **Code:** 200
       **Content:**
       ```json
       {
        "error": null,
        "expiration": "2017-05-14T13:03:07+01:00",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InJ1Z2lsZW5hQGdtYWlsLmNvbSIsImV4cCI6IjIwMTctMDUtMTRUMTM6MDM6MDcrMDE6MDAifQ.utL_J900CAIif_FM7UcJiJomRZm_R5VSG8hj6aNfgOA"
        }
       ```

   * **Error Response:**

     <_Returns 404 if the token is already expired or does not exist._>

     * **Code:** 404
       **Content:**
       ```json
       {
         "error": "Token not found or already expired."
        }
       ```

   **Destroy a token**
   ----

   * **`POST` /token/destroy**
   * **Form Params**

     **Required:**
     <_Only required paramenter - the active token which we want to renew._>
      `token=[string]`

   * **Success Response:**

     <_Returns a 200 response.._>

     * **Code:** 200
       **Content:**
       ```json
       {
        "error": null,
        }
       ```

   * **Error Response:**

     <_Returns 404 if the token is not found._>

     * **Code:** 404
       **Content:**
       ```json
       {
         "error": "Token not found."
        }
       ```

# <<Users>>
   **Create a user**
   ----

   * **`POST` /user/create**
   * **Form Params**

     **Required:**
     <_User email, password and name. The email must be unique._>
      `email=[string]`
      `password=[string]`
      `name=[string]`

   * **Success Response:**

     <_Returns a 201 response._>

     * **Code:** 201
       **Content:**
       ```json
       {
        "error": null
        }
       ```

   * **Error Response:**

    * **Code:** 400
       **Content:**
       <_If email is in an invalid format, it will not create the user._>
       ```json
       {
         "error": "Invalid email."
        }
       ```
    * **Code:** 400
       **Content:**
       <_If there already is a user with that email, returns an error._>
       ```json
       {
         "error": "User with that email already exists."
        }
       ```
