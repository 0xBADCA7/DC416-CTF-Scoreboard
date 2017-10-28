const submissions = require('./submissions')
const flags = require('../../../config/flags.json')

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

const findTeamByTokenQ = `
select id, name, token, score
from teams
where token = ?;
`

const updateScoreQ =`
update teams
set score = ?
where id = ?;
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
  let query = ''
  let arg = ''

  if (lookupBy.id !== undefined) {
    query = findTeamByIdQ
    arg = lookupBy.id
  } else if (lookupBy.name !== undefined) {
    query = findTeamByNameQ
    arg = lookupBy.name
  } else {
    query = findTeamByTokenQ
    arg = lookupBy.token
  }

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

const update = async (db, { id, score }) => {
  return new Promise((reject, resolve) => {
    db.run(updateScoreQ, score, id, err => err ? reject(err) : resolve())
  })
}

const submitFlag = async (db, { token, flag }) => {
  const submitted = flags.find(({ secret }) => secret === flag)

  if (flag === undefined) {
    return Promise.reject(new Error('incorrect flag'))
  }
  const { id, score } = await find(db, { token })
  await submissions.create(db, {
    team: id,
    flag: submitted.id,
    value: submitted.value
  })
  return await update(db, { id, score: score + submitted.value })
}


exports.all = all
exports.find = find
exports.submitFlag = submitFlag