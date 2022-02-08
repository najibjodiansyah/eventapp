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

- **How a new user register**

```graphql
mutation {
    createUser (
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
    createUser (
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
        "createUser": {
            "id": 1
            "name": "Bagus Brahmantya"
        }
    }
}
```

- **How a registered user login**

```graphql
query {
    authLogin (
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
    authLogin (
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
        "authLogin": {
            "id": 1
            "name": "Bagus Brahmantya"
            "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NDI1ODE0OTAsImlkIjoxLCJuYW1lIjoiIn0.NW5sOrXy_h9R1LovaarufA-anCJ14QAIwrzYLirkQdQ"
        }
    }
}
```

- **How anyone get list of all users**

```graphql
query {
    users {
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
query {
    users {
        id
	    name
    }
}

// example success response
{
    "data": {
        "users": [
            {
                "id": 1
                "name": "Bagus Brahmantya"
            },
            {
                "id": 2
                "name": "Najib Jodiansyah"
            },
            {
                "id": 3
                "name": "Yahya Zakariya"
            }
        ]
    }
}
```

- **How anyone get user by id**

```graphql
query {
    userById (
        id: Int!
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
query {
    user (
        id: 1
    ) {
	    name
        email
    }
}

// example success response
{
    "data": {
        "userById": {
            "name": "Bagus Brahmantya"
            "email": "bagusbpg@gmail.com"
        }
    }
}
```

- **How registered user update his/her own profile**

```graphql
mutation {
    updateUser (
        id: Int!
        set: {
            name: String
	        email: String
	        password: String
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
    updateUser (
        id: 1
        set: {
	        email: "bagusbrahmantya@gmail.com"
	        password: "yametkudasi"
            phoneNumber: "+6281234567890"
            avatar: "http://imgur.com/afywt07sav"
        }
    ) {
	    id
	    name
	    email
	    password
        phoneNumber
        avatar
    }
}

// example success response
{
    "data": {
        "updateUser": {
            "id": 1
            "name": "Bagus Brahmantya"
            "email": "bagusbrahmantya@gmail.com"
            "password":
            "phoneNumber": "+6281234567890"
            "avatar": "http://imgur.com/afywt07sav"
        }
    }
}
```

- **How registered user delete his/her own account**

```graphql
mutation {
    deleteUser (
        id: Int!
    ) {
        code: Int!
	    message: String!
    }
}
```

```jsx
// example request body
mutation {
    deleteUser (
        id: 1
    ) {
	    code
        message
    }
}

// example success response
{
    "data": {
        "deleteUser": {
            "code": 200
            "message": "success delete user"
        }
    }
}
```

### Event

- **How anyone get list of all events (pagination mode: on)**

```graphql
query {
    events (
        page: Int!
    ) {
        totalPage
        event{
            id: Int
            name: String!
            username: String!
            host: String!
            description: String!
            datetime: String!
            location: String!
            category: String!
            photo: String!
        }
    }
}
```

```jsx
// example request body
query {
    events(page:1)
    {
    totalPage
    event{  id
            name
            host
            category
        }
    }
}

// example success response
{
    "data": {
        "totalPage": 3
        "events": [
            {
                "id": 1
                "name": "Git 101"
                "host": "Alterra Academy"
                "category": "Education"
            },
            {
                "id": 2
                "name": "Unit Testing for Beginner"
                "host": "Alterra Academy"
                "category": "Education"
            },
            {
                "id": 3
                "name": "Introduction to TypeScript"
                "host": "Alterra Academy"
                "category": "Education"
            }
        ]
    }
}
```

- **How anyone get list of all events hosted by particular user (pagination mode: off)**

```graphql
query {
    eventByHostId (
        userId: Int!
    ) {
        id: Int
	    name: String!
	    username: String!
	    host: String!
	    description: String!
	    datetime: String!
	    location: String!
	    category: String!
	    photo: String!
    }
}
```

```jsx
// example request body
query {
    eventByHostId (
        userId: 1
    ) {
        id
	    name
        username
    }
}

// example success response
{
    "data": {
        "eventByHostId": [
            {
                "id": 1
                "name": "Git 101"
                "username": "Bagus Brahmantya"
            },
            {
                "id": 2
                "name": "Unit Testing for Beginner"
                "username": "Bagus Brahmantya"
            },
            {
                "id": 3
                "name": "Introduction to TypeScript"
                "username": "Bagus Brahmantya"
            }
        ]
    }
}
```

- **How anyone get list of all events held in certain location (pagination mode: on)**

```graphql
query {
    eventByLocation (
        location: String!
        page: Int!
    ) {
        totalPage
        event{
            id: Int
            name: String!
            username: String!
            host: String!
            description: String!
            datetime: String!
            location: String!
            category: String!
            photo: String!
    }
}
```

```jsx
// example request body
query {
    eventByLocation (
        location: "Malang"
        page: 1
    ) {
    totalPage
    event{  id
            name
            location
    }
}

// example success response
{
    "data": {
        "totalPage": 3
        "eventByLocation": [
            {
                "id": 1
                "name": "Git 101"
                "location": "Malang"
            },
            {
                "id": 2
                "name": "Unit Testing for Beginner"
                "location": "Malang"
            }
        ]
    }
}
```

- **How anyone get list of all events based on search keyword in event name (pagination mode: on)**

```graphql
query {
    eventByKeyword (
        keyword: String!
        page: Int!
    ) {
        totalPage
        event{
            id: Int
            name: String!
            username: String!
            host: String!
            description: String!
            datetime: String!
            location: String!
            category: String!
            photo: String!
    }
}
```

```jsx
// example request body
query {
    eventByKeyword (
        keyword: "beginner"
        page: 1
    ) {
    totalPage
    event{  id
            name
    }
    }
}

// example success response
{
    "data": {
        "totalPage": 1
        "eventByKeyword": [
            {
                "id": 2
                "name": "Unit Testing for Beginner"
            }
        ]
    }
}
```

- **How anyone get list of all events based on category (pagination mode: on)**

```graphql
query {
    eventByCategory (
        category: String!
        page: Int!
    ) {
        totalPage
        event{
            id: Int
            name: String!
            username: String!
            host: String!
            description: String!
            datetime: String!
            location: String!
            category: String!
            photo: String!
    }
}
```

```jsx
// example request body
query {
    eventByCategory (
        category: "Education"
        page: 1
    ) {
	totalPage
    event{  name
            category
    }
}

// example success response
{
    "data": {
        "totalPage": 1
        "eventByCategory": [
            {
                "name": "Git 101"
                "category": "Education"
            },
            {
                "name": "Unit Testing for Beginner"
                "category": "Education"
            },
            {
                "name": "Introduction to TypeScript"
                "category": "Education"
            }
        ]
    }
}
```

- **How registered user get list of all events he/she participates (pagination mode: off)**

```graphql
query {
    eventByParticipantId (
        userId: Int!
    ) {
        id: Int
	    name: String!
	    username: String!
	    host: String!
	    description: String!
	    datetime: String!
	    location: String!
	    category: String!
	    photo: String!
    }
}
```

```jsx
// example request body
query {
    eventByParticipantId (
        userId: 1
    ) {
	    name
        username
        host
        location
        category
    }
}

// example success response
{
    "data": {
        "eventByParticipantId": {
            "name": "Git 101"
            "username": "Bagus Brahmantya"
            "host": "Alterra Academy"
            "location": "Malang"
            "category": "Education"
        }
    }
}
```

- **How anyone retrieve particular event by id**

```graphql
query {
    eventById (
        id: Int!
    ) {
        id: Int
	    name: String!
	    username: String!
	    host: String!
	    description: String!
	    datetime: String!
	    location: String!
	    category: String!
	    photo: String!
    }
}
```

```jsx
// example request body
query {
    eventById (
        id: 1
    ) {
	    name

    }
}

// example success response
{
    "data": {
        "eventByParticipantId": [
            {
                "name": "Gaining Profit from Cyrptocurrency"
            },
            {
                "name": "Creating Profitable NFT"
            },
            {
                "name": "Life is a Joke"
            }
        ]
    }
}
```

- **How registered user create an event**

```graphql
mutation {
    createEvent (
        input: {
            name: String!
	        userid: Int
	        host: String!
	        description: String!
	        datetime: String!
	        location: String!
	        category: String!
	        photo: String!
        }
    ) {
        id: Int
	    name: String!
	    username: String!
	    host: String!
	    description: String!
	    datetime: String!
	    location: String!
	    category: String!
	    photo: String!
    }
}
```

```jsx
// example request body
mutation {
    createEvent (
        input: {
            name: "Arabica Coffee Speciality"
	        host: "Starbucks"
	        description: "Apa dan bagaimana"
	        datetime: "2022-02-04 15:00:00"
	        location: "IKN"
	        category: "Hobby"
	        photo: "http://imgur.com/t70asvafhde"
        }
    ) {
	    id
	    name
	    username
	    host
    }
}

// example success response
{
    "data": {
        "createEvent": {
            "id": 18
            "name": "Arabica Coffee Speciality"
            "username": "Najib Jodiansyah"
            "host": "Starbucks"
        }
    }
}
```

- **How registered user update an event he/she hosts**

```graphql
mutation {
    updateEvent (
        id: Int!
        input: {
            name: String
	        host: String
	        description: String
	        datetime: String
	        location: String
	        category: String
	        photo: String
        }
    ) {
        id: Int
	    name: String!
	    username: String!
	    host: String!
	    description: String!
	    datetime: String!
	    location: String!
	    category: String!
	    photo: String!
    }
}
```

```jsx
// example request body
mutation {
    updateEvent (
        id: 18
        input: {
            name: "Arabica Coffee Speciality for Beginner"
        }
    ) {
	    name
    }
}

// example success response
{
    "data": {
        "updateEvent": {
            "name": "Arabica Coffee Speciality for Beginner"
        }
    }
}
```

- **How registered user delete an event he/she hosts**

```graphql
mutation {
    deleteEvent (
        id: Int!
    ) {
        code: Int!
        message: String!
    }
}
```

```jsx
// example request body
mutation {
    deleteEvent (
        id: 18
    ) {
	    code
        message
    }
}

// example success response
{
    "data": {
        "deleteEvent": {
            "code": 200
            "message": "success delete event"
        }
    }
}
```

### Participants

- **How anyone get list of all participants in a particular event (pagination mode: off)**

```graphql
query {
    participants (
        eventId: Int!
    ) {
	    name: String!
	    avatar: String!
    }
}
```

```jsx
// example request body
query {
    participants (
        eventId: 1
    ) {
	    name
        avatar
    }
}

// example success response
{
    "data": {
        "participants": [
            {
                "name": "Bagus Brahmantya"
                "avatar": "http://imgur.com/bistayo"
            },
            {
                "name": "Najib Jodiansyah"
                "avatar": "http://imgur.com/squarepants"
            },
            {
                "name": "Yahya Zakariya"
                "avatar": "http://imgur.com/tuankrab"
            }
        ]
    }
}
```

- **How registered user join an event**

```graphql
mutation {
    joinEvent (
        eventid: Int!
    ) {
        code: Int!
        message: String!
    }
}
```

```jsx
// example request body
mutation {
    joinEvent (
        eventid: 1
    ) {
	    code
        message
    }
}

// example success response
{
    "data": {
        "joinEvent": {
            "code": 200
            "message": "success join event"
        }
    }
}
```

- **How registered user unjoin an event he/she previosly joined**

```graphql
mutation {
    unjoinEvent (
        eventid: Int!
    ) {
        code: Int!
        message: String!
    }
}
```

```jsx
// example request body
mutation {
    unjoinEvent (
        eventid: 1
    ) {
	    code
        message
    }
}

// example success response
{
    "data": {
        "joinEvent": {
            "code": 200
            "message": "success unjoin event"
        }
    }
}
```

### Comments

- **How anyone get list of all comments in a particular event (pagination mode: off)**

```graphql
query {
    comments (
        eventId: Int!
    ) {
        id: Int
	    userId: Int!
	    name: String!
	    avatar: String!
	    content: String!
	    createdAt: String!
    }
}
```

```jsx
// example request body
query {
    comments (
        eventId: 1
    ) {
	    name
        content
    }
}

// example success response
{
    "data": {
        "comments": [
            {
                "name": "Bagus Brahmantya"
                "content": "Dikasih snack ga boss?"
            },
            {
                "name": "Najib Jodiansyah"
                "content": "Pembicaranya ganti aja dong."
            },
            {
                "name": "Yahya Zakariya"
                "avatar": "Ini bukannya pas tanggal merah ya? Panitianya ngaco nih."
            }
        ]
    }
}
```

- **How registered user make a comment on an event**

```graphql
mutation {
    createComment (
        eventid: Int!
        input: String!
    ) {
        id: Int
	    userId: Int!
	    name: String!
	    avatar: String!
	    content: String!
	    createdAt: String!
    }
}
```

```jsx
// example request body
mutation {
    createComment (
        eventid: 1
        input: "Ini wajib on camera ya pas webinarnya? kameraku rusak bos"
    ) {
	    id
	    userId
	    name
	    avatar
	    content
	    createdAt
    }
}

// example success response
{
    "data": {
        "createComment": {
            "id": 1
	        "userId": 1
	        "name": "Bagus Brahmantya"
	        "avatar": "http://imgur.com/afywt07sav"
	        "content": "Ini wajib on camera ya pas webinarnya? kameraku rusak bos"
	        "createdAt": "2022-02-04 06:35:40"
        }
    }
}
```

- **How registered user delete a comment he/she made in an event**

```graphql
mutation {
    deleteComment (
        commentid: Int!
    ) {
        code: Int!
        message: String!
    }
}
```

```jsx
// example request body
mutation {
    deleteComment (
        commentid: 1
    ) {
	    code
        message
    }
}

// example success response
{
    "data": {
        "deleteComment": {
            "code": 200
            "message": "success delete comment"
        }
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
  go run ./app/server.go
```

## Running Tests

To run tests, run the following command

```bash
  go test ./...
```

## Used By

This project is used by the following companies:

- Group 3 Projet Event Planning App : Klender

## Lessons Learned

What did you learn while building this project? What challenges did you face and how did you overcome them?

## Feedback

If you have any feedback, please reach out to us at najibjodiansyah11@gmail.com

## Support

For support, email najibjodiansyah11@gmail.com or join our Discord channel.
