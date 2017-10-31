const sql = require('sqlite3')


const query = `
insert into messages (
  posted, content
) values (
  $1, $2
);
`

const run = argv => {
  const db = new sql.Database(argv.database)
  const now = new Date().getTime()
  db.run(query, now, argv.msg, err => {
    if (err) {
      console.log(`Error saving message: ${err.message}`)
    } else {
      console.log('Your message has been posted.')
    }
  })
}


exports.run = run