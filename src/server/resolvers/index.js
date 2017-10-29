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


const teamSortCompare = (t1, t2) => {
  if (t2.score != t1.score) {
    return t2.score - t1.score
  }
  const time1 = new Date(t1.lastSubmission)
  const time2 = new Date(t2.lastSubmission)
  return time2.getTime() - time1.getTime()
}


const displayTeams = (teams) => {
  teams.sort(teamSortCompare)
  for (const index in teams) {
    const i = parseInt(index)
    teams[i].rank = i + 1
  }
  return teams
}

const resolvers = {
  Query: {
    teams: async (_, __, { db }) => {
      const teams = await queries.teams.all(db)
      return displayTeams(teams)
    },
    messages: async (_, __, { db }) => {
      const messages = await queries.messages.all(db)
      return messages.map(msg => {
        msg.posted = Math.round(msg.posted / 1000.0)
        return msg
      })
    }
  },
  Mutation: {
    submitFlag: async (_, { submissionToken, flag }, { db }) => {
      const result = await queries.teams.submitFlag(db, { flag, token: submissionToken })
      const teams = await queries.teams.all(db)
      return {
        correct: result,
        teams: displayTeams(teams),
      }
    }
  },
  Team: {
    lastSubmission: async ({ name }, _, { db }) => {
      const submissions = await queries.submissions.find(db, { teamName: name })
      submissions.sort((s1, s2) => s1.time - s2.time)
      if (submissions.length === 0) {
        return null
      }
      return Math.round(submissions[submissions.length - 1].time / 1000.0)
    }
  }
}

exports.resolvers = resolvers