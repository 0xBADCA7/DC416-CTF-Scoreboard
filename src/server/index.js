const fs = require('fs')
const { makeExecutableSchema } = require('graphql-tools')
const express = require('express')
const graphqlHTTP = require('express-graphql')
const sql = require('sqlite3')
const { resolvers } = require('./resolvers')
const { initDB } = require('./queries')

const app = express()
const db = new sql.Database('scoreboard.db')
const typeDefs = fs.readFileSync('src/server/schema/main.graphql', 'utf8')
const schema = makeExecutableSchema({ typeDefs, resolvers })


const start = async () => {
  await initDB(db)

  app.get('/', (req, res) => {
    const fileContent = fs.readFileSync('app/index.html')
    res.append('Content-Type', 'text/html')
    res.send(fileContent)
  })

  app.use('/graphql', graphqlHTTP({
    schema,
    graphiql: true,
    context: { db },
  }))

  app.use(express.static('dist'))

  app.listen(9001)
}


start()