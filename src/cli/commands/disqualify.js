const sql = require('sqlite3')


const query = `delete from teams where name = $1;`

const run = argv => {
  const db = new sql.Database(argv.database)
  db.run(query, argv.name, err => {
    if (err) {
      console.log(`Error registering team: ${err.message}`)
    } else {
      console.log(`Success! "${argv.name}" will no longer be listed on the scoreboard or allowed to submit flags.`)
    }
  })
}


exports.run = run