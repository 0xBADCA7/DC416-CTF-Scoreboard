const queries = require('../queries')


const testMessages = [
  {
    posted: 'October 15 at 9:30 AM',
    content: 'Registrations are open!',
  },
  {
    posted: 'October 18 at 4:30 PM',
    content: 'We will be posting more information to help you prepare for the CTF soon.',
  },
  {
    posted: 'October 21 at 5:00 PM',
    content: 'The challenges are all available now! Good luck, everyone!'
  },
]


const teamScoreCompare = (t1, t2) => {
  if (t1.score != t2.score) {
    return t1.score - t2.score
  }
  const time1 = new Date(t1.lastSubmission)
  const time2 = new Date(t2.lastSubmission)
  return time1.getTime() - time2.getTime()
}


const resolvers = {
  Query: {
    teams: (_, __, { db }) => {
      return queries.teams.all(db)
        .then(teams => teams.map(team => ({
          rank: 0,
          name: team.name,
          score: team.score,
        })))
      },
    messages: () => testMessages,
  },
  Mutation: {
    submitFlag: (_, args) => {
      return [] 
    }
  },
  Team: {
    lastSubmission: ({ name }, _, { db }) => {
      return queries.teams.find(db, { name })
        .then(team => team.lastSubmission || 'No submissions yet.')
      }
  }
}

exports.resolvers = resolvers