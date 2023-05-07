# finance-management
pet finance management project 

### commands:
make postgres - create postgres15 container with postgres db
make postgres-createdb - create app's database finances
make postgres-dropdb - drop db finances
make postgres-migrate-up - execute migration in db finances
make postgres-migrate-down - revert migration in db finances
make sqlc - generate "db" package with query to database

### target functionality:
- create office
- add into office user
- login by user
- setup currency in office
- create account in office
- create transaction between accounts

### DB diagram:
https://dbdiagram.io/d/64576de7dca9fb07c4a412c1