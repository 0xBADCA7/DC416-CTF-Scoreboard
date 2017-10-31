const sql = require('sqlite3')
const crypto = require('crypto')
const { promisify } = require('util')


const query = `
insert into teams (
  name, score, token
) values (
  $1, 0, $2
);
`

const generateToken = async () => {
  const rand = promisify(crypto.randomBytes)
  const bytes = await rand(16)
  return bytes.toString('hex')
}

const run = async (argv) => {
  const db = new sql.Database(argv.database)
  const token = await generateToken()
  db.run(query, argv.name, token, err => {
    if (err) {
      console.log(`Error registering team: ${err.message}`)
    } else {
      console.log(`Successfully registered the new team.\nName: ${argv.name}\nSubmission token: ${token}`)
    }
  })
}


exports.run = run