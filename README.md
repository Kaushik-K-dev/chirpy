# chirpy

## User resource
```json
{
    "id":1,
    "email":"walt@breakingbad.com",
    "hashed_pass":"$2a$10$MnUt50CtHcGjrMseVdcpsOheYkxNahSpeuUs3OjwCaLFAfa2TaZVO","refresh_token":"f27d244ab24f5a17740423c6e12d89991ed23f0346ffa0570c6ec5605a50df29",
    "refresh_expiration":"2024-11-18T10:16:48.520305506+05:30",
    "is_chirpy_red":false
}
```

### POST /api/users
Request Body:
```json
{
  "email": "walt@breakingbad.com",
  "password": "123456"
}
```
Response Body:
```json
{
  "email": "walt@breakingbad.com",
  "id": 1,
  "is_chirpy_red": false
}
```

### POST /api/login
Request Body:
```json
{
  "email": "walt@breakingbad.com",
  "password": "123456"
}
```
Response Body:
```json
{
  "email": "walt@breakingbad.com",
  "id": 1,
  "is_chirpy_red": false,
  "refresh_token": "ceb9d6e50512a2618fdedd2d4927871f86743451f22f0edd3bac520ad271fb83",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiIxIiwiZXhwIjoxNzI2NzIyOTc4LCJpYXQiOjE3MjY3MTkzNzh9.0cmnQUe86jCEPHnx6ZK5YWRcL1JMPE3mfVOAO4cKhkc"
}
```

## Chirp resource

``` json
{
    "id":4,
    "author_id":2,
    "body":"Darn that fly, I just wanna cook"
}
```

### GET /api/chirps
Returns an array of chirps with optional query parameters of author_id and sort (sort only has 2 modes "asc"(default) and "desc")
If author_id is not specified, it returns all chirps in the database

### POST api/chirps
Headers:
```json
{
  "Authorization": "Bearer ${token}"
}
```

Request Body:
```json
{
  "body": "I'm the one who knocks!"
}
```