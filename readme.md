## Base Go     
Source code for Base Go App.

The project uses:

**Golang.     
Gin.  
Postgres.   
Redis.**

We have a script in Makefile which creates necessary keypairs (for creating dev and test JSON Web Tokens), runs   
required docker-compose service for Postgres, migrates our database tables, then shuts down docker-compose.

As an example, we will migrate database changes found in `~/account/migrations` to our postgres-account service   
found in `~/docker-compose.yml`.

If you need to use some commands separately (for example make your own DB migrations),  to distinguish them from    
others, you can check Makefile in the root directory.

## Runing the Application

`.env.dev` files are provided in each application directory for providing development environment variables and     
quickly running this application. 

We will not be adding any critical keys directly to the `.env.dev` file. However, we will eventually  refer to    
access key files in `.env.dev`. Make sure to add these key files to your `.gitignore`.

## From the project root director, run:

```make init```

```docker-compose up``` 

The API will then be available at  **http://127.0.0.1:8000/account**

You can also find all possible API requests/urls when you launch the project in `docker-compose.yml`. 

If you need to make rebuild, you have to use this command:

```docker-compose build```      
After that repeat the command ```docker-compose up``` for launching the project.


### Google Cloud Key

In order to access Google Cloud for storing profile images, you will need to download a service account JSON file   
to your account application folder and call it `serviceAccount.json`.     

This file will be references in ```.env.dev```. You also need to call your Bucket in Google Cloud  as   
`go_base_profile_images` (you can find this information in `.env.dev` file).

Instructions for installing the Google Cloud Storage Client and getting this key are found at:

**https://cloud.google.com/storage/docs/reference/libraries**

## Run

To run this code, you will need docker and docker-compose installed on your machine. In the project root, run:  
```docker-compose up```

## Tests 

You can launch the tests using this command:      

```go test   ./...```

This command lauches all tests in our application (for repository, service and handler layers). If you need to launch     
a separate layer, instead of of using ```...``` you should use a folder name.

Launch tests for a repository layer:        
```go test ./repository```

Launch tests for a service layer:       
```go test ./service```

Launch tests for a handler layer:   
```go test ./handler```
