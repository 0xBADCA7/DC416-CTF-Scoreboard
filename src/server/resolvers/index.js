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
    messages: (_, __, { db }) => {
      return queries.messages.all(db)
        .then(msgs => msgs.map(msg => ({
          posted: msg.posted,
          content: msg.content,
        })))
    }
  },
  Mutation: {
    submitFlag: async (_, { submissionToken, flag }, { db }) => {
      queries.teams.submitFlag(db, { flag, token: submissionToken })
      return {
        correct: true,
        newScore: 0,
        scoreboard: [],
      }
    }
  },
  Team: {
    lastSubmission: async ({ name }, _, { db }) => {
      const submissions = await queries.submissions.find(db, { teamName: name })
      submissions.sort((s1, s2) => s1.time - s2.time)
      if (submissions.length === 0) {
        return 'No submissions yet.'
      }
      const date = new Date(submissions[submissions.length - 1].time)
      return date.toLocaleString()
    }
  }
}

exports.resolvers = resolvers