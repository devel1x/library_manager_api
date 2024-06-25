### HOW TO RUN THE PROJECT

docker-compose up 


### API ENDPOINTS

#### USERS
1. POST /api/v1/user/login
2. POST /api/v1/user/signup
3. POST /api/v1/user/auth/refresh

#### BOOKS
1. GET /api/v1/book --*get list of books (supports pagination)*
2. GET /api/v1/book/{{isbn}} --*get book by isbn*
3. POST /api/v1/book/ --*create book*
4. PUT /api/v1/book/{{isbn}} --*update book by isbn*
5. DELETE /api/v1/book/{{isbn}} --*delete book by isbn*

#### SWAGGER
1. GET /swagger/index.html # to see the swagger documentation


