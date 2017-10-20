const testTeams = [
  {
    rank: 1,
    name: "Team one",
    score: 155,
    lastSubmission: "5:30 PM"
  },
  {
    rank: 2,
    name: "DC416",
    score: 140,
    lastSubmission: "4:52 PM"
  },
  {
    rank: 3,
    name: "h4xx0rz",
    score: 135,
    lastSubmission: "4:04 PM"
  },
  {
    rank: 4,
    name: "31337",
    score: 135,
    lastSubmission: "4:38 PM"
  },
  {
    rank: 5,
    name: "grep -i flag",
    score: 100,
    lastSubmission: "2:00 PM"
  },
  {
    rank: 6,
    name: "First place",
    score: 25,
    lastSubmission: "1:13 PM"
  },
  {
    rank: 7,
    name: "Lucky #7",
    score: 0,
    lastSubmission: "No submissions yet"
  },
]

const resolvers = {
  Query: {
    teams: () => testTeams,
  }
}

exports.resolvers = resolvers