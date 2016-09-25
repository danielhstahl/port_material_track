Postgresql information

Setup:
`sudo -u postgres psql`
`CREATE USER porttransaction;`

`ALTER USER porttransaction PASSWORD 'zu3<U@2Cz[uG>Hmf';`
`CREATE DATABASE portServer WITH OWNER=porttransaction;`

`CREATE TABLE PossiblePorts (Port varchar(20) CONSTRAINT portKey PRIMARY KEY);`

`CREATE TABLE PossibleMaterials (Material varchar(20) CONSTRAINT matKey PRIMARY KEY);`

`CREATE TABLE MaterialTransactions (Port varchar(20) not null, Material varchar(20) not null, TransactionDate date not null, Amount int not null);`

`ALTER TABLE MaterialTransactions ADD PRIMARY KEY(Port, Material, TransactionDate);`

`ALTER TABLE MaterialTransactions ADD CONSTRAINT portFK FOREIGN KEY(Port) REFERENCES PossiblePorts(Port);`

`ALTER TABLE MaterialTransactions ADD CONSTRAINT matFK FOREIGN KEY(Material) REFERENCES PossibleMaterials(Material);`



user: `porttransaction`
password: `zu3<U@2Cz[uG>Hmf`
name: `portServer`