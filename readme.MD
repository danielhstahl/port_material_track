Postgresql information
`sudo apt-get install postgresql`

Setup (note that the PG user and PG password are generic...use your own!)
`sudo -u postgres psql`
`CREATE USER porttransaction;`

`ALTER USER porttransaction PASSWORD 'zu3<U@2Cz[uG>Hmf';`
`CREATE DATABASE portserver WITH OWNER=porttransaction;`
`\connect portserver`
`CREATE SCHEMA main AUTHORIZATION porttransaction;` 

`CREATE TABLE main.possibleports (Port varchar(50) CONSTRAINT portkey PRIMARY KEY);`

`CREATE TABLE main.possiblematerials (material varchar(50) CONSTRAINT matkey PRIMARY KEY);`

`CREATE TABLE main.materialtransactions (port varchar(50) not null, material varchar(50) not null, transactiondate date not null, amount int not null, comment text);`

`ALTER TABLE main.materialtransactions ADD PRIMARY KEY(port, material, transactiondate);`

`ALTER TABLE main.materialtransactions ADD CONSTRAINT portfk FOREIGN KEY(port) REFERENCES portserver.main.possibleports(port);`

`ALTER TABLE main.materialtransactions ADD CONSTRAINT matfk FOREIGN KEY(material) REFERENCES main.possiblematerials(material);`

`GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA main TO porttransaction;`

Put an environment.json file like 

`{
    "PGUsername":"porttransaction",
    "PGPassword":"zu3<U@2Cz[uG>Hmf",
    "Statichtml":"~/Documents/htmlApps/portClient/build"
}`
in the binary folder

user: `porttransaction`
password: `zu3<U@2Cz[uG>Hmf`
name: `portserver`

# install

`go get`
`go build`

If errors occur, try
`go mod tidy`
`go mod vendor`