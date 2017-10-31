const yargs = require('yargs')
const { commands } = require('./commands')


const argv = yargs
  .option('database', {
    alias: 'd',
    default: 'scoreboard.db',
  })
  .command('register <name>', 'Register a new team.', yargs => {
    return yargs.positional('name', {
      describe: 'The name of the new team to register. It must be unique. Multiple words must be enclosed in quotes.',
    })
  })
  .command('disqualify <name>', 'Disqualify a team from participating.', yargs => {
    return yargs.positional('name', {
      describe: 'The name of the team to disqualify. Multiple words must be enclosed in quotes.',
    })
  })
  .command('message <msg>', 'Post a message for all participants to see.', yargs => {
    return yargs.positional('msg', {
      describe: 'A message to post to participants. Multiple words must be enclosed in quotes.',
    })
  })
  .command('dashboard', 'View information about teams and submissions.')
  .argv



const main = () => {
  const command = commands.find(({ name }) => name === argv._[0])
  if (command === undefined) {
    console.log(`Unknown command "${argv._[0]}"`)
    return
  }
  command.run(argv)
}


main()