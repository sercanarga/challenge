# Challenge
This project is an API service developed in Go. Designed to work under high traffic and used to decrease, increase and query the balance value of existing users.

![last commit](https://badgen.net/github/last-commit/sercanarga/challenge) ![license](https://badgen.net/github/license/sercanarga/challenge)

## Features
- It is developed in Go language.
- Redis database used.
- Messaging is provided with Kafka.
- HTTP server is created with Gin framework.
- API documentation is provided with Swagger.
- It can be easily stand up with Docker.

## Installation
1. Clone the project:
```
git clone https://github.com/sercanarga/challenge.git
```
2. Go to the project directory:
```
cd challenge
```
3. edit the env file:
```env
# General
APP_NAME=challenge
APP_PORT=8080

# Database
DB_DSN=postgres://postgres:123456@postgres:5432/postgres

# Kafka
KAFKA_TOPIC=challenge
KAFKA_BROKER=kafka:9092
```
4. Stand up the project with Docker compose:
```
docker-compose up --build -d
```

## Usage
Go to `localhost:8080/doc` swagger documentation to see the API endpoints and their usage.

_Note: port 8080 is used by default. If you want to use a different port you can change the `APP_PORT` env variable._

## Endpoints
### [GET] /
Shows wallet information of registered users.

#### Request:
```bash
curl -X 'GET' \
'http://localhost:8080/?limit=5&cursor=0' \
-H 'accept: application/json'
```

#### Response:
```json
{
  "wallets": [
    {
      "id": "1",
      "user_id": "1",
      "balances": [
        {
          "currency": "TRY",
          "Amount": 30,
          "LastUpdate": "2024-02-19T22:34:08.125Z"
        }
        ...
      ]
    }
    ...
  ]
}
```
### [POST] /
Increases or decrease the balance of the respective currency in a wallet.

#### Request:
```bash
curl -X 'POST' \
  'http://localhost:8080/' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "events": [
    {
      "app": "string",
      "attributes": {
        "amount": "string",
        "currency": "string"
      },
      "meta": {
        "user": "string"
      },
      "time": "string",
      "type": "string",
      "wallet": "string"
    }
  ]
}'
```
#### Response:
```json
{
  "data": {
   ...
  },
  "result": {
    "errorDetails": "string",
    "statusCode": 0
  }
}
```

## Database Diagram
![database diagram](https://raw.githubusercontent.com/sercanarga/challenge/main/assets/new_database_diagram.jpg?raw=true)

## Performance & Security Tests
AI-powered [deepsource](https://deepsource.com/) was used for security tests. The results are as follows.
![test result](https://raw.githubusercontent.com/sercanarga/challenge/main/assets/deepsource.jpg?raw=true)

vegeta was used for performance testing. The test was performed with 1000 requests per second for 30 seconds. The test results are as follows:
### [GET] /
```bash
echo "GET http://localhost:8080/" | vegeta attack -rate=1000 -duration=30s | tee results.bin | vegeta report
```
```
Requests      [total, rate, throughput]         30000, 1000.04, 1000.01
Duration      [total, attack, wait]             30s, 29.999s, 883µs
Latencies     [min, mean, 50, 90, 95, 99, max]  236.25µs, 652.689µs, 477.751µs, 632.687µs, 721.929µs, 5.119ms, 53.688ms
Bytes In      [total, mean]                     120000, 4.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
```
![test result](https://raw.githubusercontent.com/sercanarga/challenge/main/assets/test_result.jpg?raw=true)

---

### [POST] /
```bash
echo "POST http://localhost:8080/" | vegeta attack -body=body.json -rate=1000 -duration=30s | tee results.bin | vegeta report
```
```
Requests      [total, rate, throughput]         30000, 1000.03, 999.99
Duration      [total, attack, wait]             30s, 29.999s, 1.198ms
Latencies     [min, mean, 50, 90, 95, 99, max]  248.75µs, 1.533ms, 572.906µs, 937.044µs, 1.824ms, 20.116ms, 179.256ms
Bytes In      [total, mean]                     5490000, 183.00
Bytes Out     [total, mean]                     8040000, 268.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      202:30000  
```
![test result 2](https://raw.githubusercontent.com/sercanarga/challenge/main/assets/test_result_2.jpg?raw=true)

## License
Licensed under Apache 2.0. For details, see the [LICENSE](LICENSE) file.