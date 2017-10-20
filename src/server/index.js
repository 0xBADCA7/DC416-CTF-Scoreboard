const fs = require('fs')
const { makeExecutableSchema } = require('graphql-tools')
const express = require('express')
const graphqlHTTP = require('express-graphql')
const { resolvers } = require('./resolvers')

const app = express()
const typeDefs = fs.readFileSync('src/server/schema/main.graphql', 'utf8')
const schema = makeExecutableSchema({ typeDefs, resolvers })


app.get('/', (req, res) => {
  const fileContent = fs.readFileSync('app/index.html')
  res.append('Content-Type', 'text/html')
  res.send(fileContent)
})


app.use('/graphql', graphqlHTTP({
  schema,
  graphiql: true,
}))


app.use(express.static('dist'))

app.listen(9001)