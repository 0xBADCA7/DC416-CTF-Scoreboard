const allTeamsQ = `
select id, name, token, score
from teams;
`

const findTeamByIdQ = `
select id, name, token, score
from teams
where id = ?;
`

const findTeamByNameQ = `
select id, name, token, score
from teams
where name = ?;
`

const all = (db) => {
  return new Promise((resolve, reject) => {
    db.all(allTeamsQ, (err, rows) => {
      if (err) {
        reject(err)
      } else {
        resolve(rows)
      }
    })
  })
}


const find = (db, lookupBy) => {
  const { query, arg } = lookupBy.id !== undefined
    ? { query: findTeamByIdQ, arg: lookupBy.id }
    : { query: findTeamByNameQ, arg: lookupBy.name }

  return new Promise((resolve, reject) => {
    db.get(query, arg, (err, row) => {
      if (err) {
        reject(err)
      } else {
        resolve(row)
      }
    })
  })
}


exports.all = all
exports.find = find