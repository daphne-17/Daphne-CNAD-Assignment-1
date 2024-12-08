Design consideration of microservices

The microservices is designed to operate independently of the others
auth-service: manages user authentication and generation of the JWT token
user-service: handles user-related data such as modification and viewing of profile, rental history and viewing of membership status
billing-services: handles payment processing and promo codes
reservation-service: manages vehicle reservations and availability

With this seperation it allows us to develop, test and deploy services independently without affecting the rest of the system.

As for scalability, with microservices it allows independent scaling

And lets say if for example vehicle-service faces an issue it wont affect the other services which makes this architecture fault-tolerant

Each services also has its own database to ensure data seperation and reduces coupling:
auth-service will handle authentication data
user-service will handle user and membership details 
billing-service manages billing and promotions and lastly user-service will store vehicle status and reservation records.

And as for security i have used jwt authentication so all microservices are secured using JSON Web tokens to ensure only authorized users are able to access protected routes.
Services also validate tokens before processing requests ensuring secure communication.

And as for error handling, each services is designed to handle errors gracefully and will return meaning HTTP status code

![image](https://github.com/user-attachments/assets/f45999bf-0506-4e71-a980-cd422b12ad9e)

Instructions to run the file

1. Run the sql script and change each link to your database root and passord in each database.go in each service
2. open 4 terminals and cd to each service in each terminal respectively
3. run go run main.go in each terminal
4. run login.html and click sign up
5. register a user and login
6. you are good to go :) can just run the different requirements
