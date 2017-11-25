const sql = require('sqlite3')


const queryDeleteSubmissions = `delete from submissions where team_id = (select id from teams where name = $1 limit 1);`
const queryDeleteTeam = `delete from teams where name = $1;`

const run = argv => {
  const db = new sql.Database(argv.database)
  db.run(queryDeleteSubmissions, argv.name, err => {
    if (err) {
      console.log(`Error registering team: ${err.message}`)
    } else {
      db.run(queryDeleteTeam, argv.name, err => {
        if (err) {
          console.log(`Error registering team: ${err.message}`)
        } else {
          console.log(`${argv.name} has been deregistered successfully!`)
        }
      })
    }
  })
}


exports.run = run