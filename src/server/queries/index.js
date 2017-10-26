const teams = require('./teams')

const createTeamTableQ = `
create table if not exists teams (
  id integer primary key,
  name varchar(128) unique not null,
  token varchar(32) unique not null,
  score integer
);
`

const createMessageTableQ = `
create table if not exists messages (
  id integer primary key,
  posted timestamp not null,
  content text not null
);
`

const createSubmissionTableQ = `
create table if not exists submissions (
  team_id integer,
  time timestamp not null,
  flag_id integer,
  value integer,
  foreign key (team_id) references teams (id)
  primary key (team_id, flag_id)
);
`

const initDB = (db) => {
  db.run(createTeamTableQ)
  db.run(createMessageTableQ)
  db.run(createSubmissionTableQ)
}

exports.initDB = initDB
exports.teams = teams