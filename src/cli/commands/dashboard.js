const sql = require('sqlite3')


const query = `
select S.flag_id, T.name
from submissions S, teams T
where T.id = S.team_id
order by T.score desc;
`

const max = (a, b) => a >= b ? a : b

const display = rows => {
  const uniqueFlagIds = rows.reduce((set, { flag_id }) => set.add(flag_id), new Set())
  const teamNames = rows.map(({ name }) => name)
  const longestTeamName = Array.from(teamNames)
    .map(n => n.length)
    .reduce(max, 0)  // Why doesn't reduce(Math.max, 0) work? WTF?
  const flagIds = Array.from(uniqueFlagIds).sort()
  const longestFlagId = flagIds
    .map(id => id.toString().length)
    .reduce(max, 0)

  // Heading
  const flagHeaders = flagIds
    .map(id => `Flag ${id.toString().padEnd(longestFlagId, ' ')}`)
    .join(' | ')
  const heading = `| ${'Team Name'.padEnd(longestTeamName, ' ')} | ${flagHeaders} |`
  console.log(heading)
  console.log('-'.repeat(heading.length))

  // Rows
  let submissions = teamNames
    .reduce((mapping, name) => ({ ...mapping, [name]: []}), {})
  for (const row of rows) {
    submissions[row.name].push(row.flag_id)
  }
  for (const team in submissions) {
    const flagCols = flagIds
      .map(id => submissions[team].includes(id))
      .map(captured => (captured ? '  X' : '   ').padEnd(longestFlagId + 5, ' '))
      .join(' | ')
    console.log(`| ${team.padEnd(longestTeamName, ' ')} | ${flagCols} |`)
  }
}

const run = argv => {
  const db = new sql.Database(argv.database)
  db.all(query, (err, rows) => {
    if (err) {
      console.log(`Error fetching dashboard information: ${err.message}`)
    } else {
      display(rows)
    }
  })
}


exports.run = run