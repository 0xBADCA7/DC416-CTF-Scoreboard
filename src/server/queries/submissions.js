const createSubmissionQ = `
insert into submissions
  (team_id, flag_id, value, time)
values
  (?, ?, ?, ?);
`

const findSubmissionsByTeamNameQ = `
select
  team_id as team,
  flag_id as flag,
  time,
  value
from submissions
where team_id = (
  select id
  from teams
  where name = ?
);
`


const create = (db, { team, flag, value }) => {
  return new Promise((resolve, reject) => {
    db.run(createSubmissionQ, team, flag, value, Date.now(), err => err ? reject(err) : resolve())
  })
}

const find = (db, { teamName }) => {
  return new Promise((resolve, reject) => {
    db.all(findSubmissionsByTeamNameQ, teamName, (err, rows) => err ? reject(err) : resolve(rows))
  })
}


exports.create = create
exports.find = find