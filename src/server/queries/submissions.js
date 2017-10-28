const createSubmissionQ = `
insert into submissions
  (team_id, flag_id, value, time)
values
  (?, ?, ?, (select date('now')));
`


const create = (db, teamId, flagId, value) => {
  return new Promise((resolve, reject) => {
    db.run(createSubmissionQ, teamId, flagId, value, err => err ? reject(err) : resolve())
  })
}


exports.create = create