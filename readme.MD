# Event Planning BackEnd Api for Klender App

Event Planning App adalah sebuah aplikasi yang digunakan untuk membuat dan management sebuah event.

Membuat API menggunakan bahasa Go dengan menerapkan layered architecture untuk memudahkan pemeliharaan dan peningkatan Fitur

## Acknowledgements

- [Wireframe Example](https://whimsical.com/event-planning-app-RnJLuh3osvawc733NVaQzE)
- [Open API (Referensi untuk Graphql)](https://app.swaggerhub.com/apis-docs/iswanulumam/EventPlanningApp/1.0.0)
- [Gqlgen](https://gqlgen.com)
- [ERD Eventapp API](https://app.diagrams.net/#G1Tgzu7nNEl_V5AEURr1VyVNsI7bc7NeAL)
- [Task Docs](https://docs.google.com/document/d/1SubVXXZ-1KNGTj6s9EyO7rXUd3Pe8mPOMjuhxPy2REQ/edit#)

## Tech Stack

- Go Language
- Labstack Echo
- Mysql
- Docker

## Features

- Bcrypt
- JWT
- GraphQL
- CRUD User
- CRUD Event
- Join/Going Event
- Post and Read Comment


## Queries and Mutations

### User

- How register a new user

```graphql
mutation {
    createUser(
        input: {
            name: String!
            email: String!
            password: String!
            phoneNumber: String
            avatar: String
        }
    ) {
        id: Int
        name: String!
        email: String!
        password: String!
        phoneNumber: String
        avatar: String
    }
}
```

```jsx
// example request body
mutation {
    createUser(
        input: {
            name: "Bagus Brahmantya"
            email: "bagusbpg@gmail.com"
            password: "qwerty1234"
        }
    ) {
        id
        name
    }
}

// example success response
{
    "data": {
        "id": 1
        "name": "Bagus Brahmantya"
    }
}
```

- How registered user login

```graphql
query {
    authLogin(
        email: String!
        password: String!
    ) {
        message: String!
	    id: Int!
	    name: String!
	    email: String!
	    token: String!
    }
}
```

```jsx
// example request body
query {
    authLogin(
        email: "bagusbpg@gmail.com"
        password: "qwerty1234"
    ) {
        id
        name
        token
    }
}

// example success response
{
    "data": {
        "id": 1
        "name": "Bagus Brahmantya"
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NDI1ODE0OTAsImlkIjoxLCJuYW1lIjoiIn0.NW5sOrXy_h9R1LovaarufA-anCJ14QAIwrzYLirkQdQ"
    }
}
```

## Authors

- [@bagusbpg](https://www.github.com/bagusbpg)
- [@zakariyahya](https://www.github.com/zakariyahya)
- [@najibjodiansyah](https://www.github.com/najibjodiansyah)

## Environment Variables

To run this project, you will need to set up the following environment variables to your ./config/config.go file

```bash
defaultConfig.Port
defaultConfig.Database.Driver
defaultConfig.Database.Name
defaultConfig.Database.Address
defaultConfig.Database.Port
defaultConfig.Database.Username
defaultConfig.Database.Password
```

## Run Locally

Clone the project

```bash
  git clone https://github.com/najibjodiansyah/eventapp.git
```

Go to the project directory

```bash
  cd eventapp
```

Install dependencies

```bash
  go mod tidy
```

Start the server

```bash
  go run ./app/main.go
```

## Running Tests

To run tests, run the following command

```bash
  go test ./...
```

## Used By

This project is used by the following companies:

- Group 5 Projet Event Planning App : Klender

## Lessons Learned

What did you learn while building this project? What challenges did you face and how did you overcome them?

## Feedback

If you have any feedback, please reach out to us at najibjodiansyah11@gmail.com

## Support

For support, email najibjodiansyah11@gmail.com or join our Discord channel.
