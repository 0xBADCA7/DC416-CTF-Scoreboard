exports.commands = [
  {
    name: 'register',
    run: require('./register').run,
  },
  {
    name: 'disqualify',
    run: require('./disqualify').run,
  },
  {
    name: 'message',
    run: require('./message').run,
  },
  {
    name: 'dashboard',
    run: require('./dashboard').run,
  },
]